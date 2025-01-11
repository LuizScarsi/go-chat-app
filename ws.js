let socket = new WebSocket("ws://localhost:3000/ws");
socket.onmessage = (event) => {
  console.log("received from server: ", event.data);
};
socket.send("hello from client");
