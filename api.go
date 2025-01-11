package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type WsServer struct {
	conns  map[*websocket.Conn]bool
	wsPort string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func NewWsServer(wsPort string) *WsServer {
	return &WsServer{
		conns:  make(map[*websocket.Conn]bool),
		wsPort: wsPort,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(s.handleRoot))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleByID))
	// router.HandleFunc("/chat", makeHTTPHandleFunc(s.handleChat))
	// router.HandleFunc("/chat", s.handleChat)
	// go handleMessages()

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *WsServer) Run() {
	http.Handle("/ws", websocket.Handler(s.handleWS))
	log.Println("Websocket running on port: ", s.wsPort)
	http.ListenAndServe(s.wsPort, nil)
}

func (s *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) error {
	t, _ := template.ParseFiles("index.html")
	err := t.Execute(w, "teste")
	return err
}

func (s *WsServer) handleWS(ws *websocket.Conn) {
	fmt.Println("new ws connection comming from client: ", ws.RemoteAddr())

	//protect with mutex
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *WsServer) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("ws read error: ", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("Thank you for the message!"))
	}
}

// var upgrader = websocket.Upgrader{
// 	// implement security checks
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// var clients = make(map[*websocket.Conn]bool)
// var broadcast = make(chan ChatMessageRequest)

// func (s *APIServer) handleChat(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		// return err
// 	}
// 	defer conn.Close()

// 	s.conns[conn] = true

// 	for {
// 		var messageReq ChatMessageRequest
// 		if err := conn.ReadJSON(&messageReq); err != nil {
// 			delete(s.conns, conn)
// 			fmt.Println(err)
// 			// return err
// 		}
// 		broadcast <- messageReq
// 	}
// 	// for {
// 	// 	messageType, p, err := conn.ReadMessage()
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	if err := conn.WriteMessage(messageType, p); err != nil {
// 	// 		return err
// 	// 	}
// 	// }
// }

// func handleMessages() {
// 	for {
// 		msg := <-broadcast
// 		fmt.Println("printing message received from broadcast")
// 		fmt.Println(msg.Message)
// 		for client := range clients {
// 			if err := client.WriteJSON(msg); err != nil {
// 				log.Fatal(err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountByID(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.NickName)

	err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id is given %s", idStr)
	}
	return id, err
}
