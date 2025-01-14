package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type WsServer struct {
	conns  map[*websocket.Conn]bool
	wsPort string
}

func NewWsServer(wsPort string) *WsServer {
	return &WsServer{
		conns:  make(map[*websocket.Conn]bool),
		wsPort: wsPort,
	}
}

func (s *WsServer) Run() {
	http.Handle("/ws", websocket.Handler(s.handleWS))
	http.Handle("/orderbookfeed", websocket.Handler(s.handleWSOrderBook))

	log.Println("Websocket running on port: ", s.wsPort)
	http.ListenAndServe(s.wsPort, nil)
}

func (s *WsServer) handleWS(ws *websocket.Conn) {
	fmt.Println("new ws connection comming from client:", ws.RemoteAddr())

	//protect with mutex
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *WsServer) handleWSOrderBook(ws *websocket.Conn) {
	fmt.Println("new ws connection comming from client to orderbook feed:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

func (s *WsServer) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("ws read error:", err)
			continue
		}
		msg := buf[:n]

		s.broadcast(msg)
	}
}

func (s *WsServer) broadcast(b []byte) {
	for ws := range s.conns {
		func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}
