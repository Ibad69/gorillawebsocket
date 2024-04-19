package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var ClientList = make(map[*websocket.Conn]bool)

type SocketImpl struct {
    wsconns map[*websocket.Conn]bool
    closed bool
}

type User struct {
	ID string `json:"id"`
}

func NewSocket() *SocketImpl {
    return &SocketImpl{wsconns: make(map[*websocket.Conn]bool)}
} 

func (s *SocketImpl) broadcast(messageType int, message []byte){
    fmt.Println(s)
    for ws := range s.wsconns{
        // fmt.Println("emitting to connection")
        // fmt.Println(ws)
        go func(ws *websocket.Conn) {
            if err:= ws.WriteMessage(messageType, message); err != nil{
                fmt.Println("error occured in broadcasting message")
                fmt.Println(err)
            }  
        }(ws)
    }
} 

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


func main (){
    socket := NewSocket()
    http.HandleFunc("/ws", socket.wsHandler)
    err := http.ListenAndServe(":8081", nil)
    if err !=nil {
       fmt.Println("server couldn run")
        fmt.Println(err)
    }
    fmt.Println("server running ")
}

func (s *SocketImpl) wsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("WebSocket handler")
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Error occurred in upgrading the connection to WebSocket")
        fmt.Println(err)
        return
    }
    defer conn.Close()

    fmt.Println("Client connected")
    fmt.Println("Remote address:", conn.RemoteAddr())

    // Add the connection to the WebSocket connections map
    s.wsconns[conn] = true

    // Send a test message to the client
    msg := []byte("Test message")
    if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
        fmt.Println("Error occurred in writing test message")
        fmt.Println(err)
        return
    }

    // Read messages from the client
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Error occurred in reading message")
            fmt.Println(err)

            // Remove the connection from the WebSocket connections map
            delete(s.wsconns, conn)
            return
        }

        fmt.Println("Message received:", string(message))
        fmt.Println("Message type:", messageType)

        // Broadcast the received message to all connected clients
        s.broadcast(websocket.TextMessage, message)
    }
}

// func  (s *SocketImpl) wsHandler(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("websocket handler")
//     conn, err := upgrader.Upgrade(w, r, nil)
//     fmt.Println("client connected")
//     s.wsconns[conn] = true
//     if err != nil {
//         fmt.Println("err occured in upgrading the connection to websocket")
//         fmt.Println(err)
//     }
//     msg := []byte("test message")
//     conn.WriteMessage(1, msg)
//     for {
//         messageType, message, err := conn.ReadMessage()
//         if err != nil{
//             fmt.Println("error occured in reading message")
//             fmt.Println(err)
//         }
//         fmt.Println("message received")
//         fmt.Println(string(message))
//         fmt.Println(messageType)
//         s.broadcast(1, message)
//     }
    
// }
