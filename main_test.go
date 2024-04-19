package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestSocketImpl_broadcast(t *testing.T) {
    type fields struct {
        wsconns map[*websocket.Conn]bool
        closed  bool
    }
    type args struct {
        messageType int
        message     []byte
    }
    tests := []struct {
        name   string
        fields fields
        args   args
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &SocketImpl{
                wsconns: tt.fields.wsconns,
                closed:  tt.fields.closed,
            }
            s.broadcast(tt.args.messageType, tt.args.message)
        })
    }
}

func TestNewSocket(t *testing.T) {
    tests := []struct {
        name string
        want *SocketImpl
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := NewSocket(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewSocket() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_main(t *testing.T) {
    tests := []struct {
        name string
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            main()
        })
    }
}

func TestSocketImpl_wsHandler(t *testing.T) {
    type fields struct {
        wsconns map[*websocket.Conn]bool
        closed  bool
    }
    type args struct {
        w http.ResponseWriter
        r *http.Request
    }
    tests := []struct {
        name   string
        fields fields
        args   args
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &SocketImpl{
                wsconns: tt.fields.wsconns,
                closed:  tt.fields.closed,
            }
            s.wsHandler(tt.args.w, tt.args.r)
        })
    }
}

// MockWebSocketServer creates a mock WebSocket server for testing.
func MockWebSocketServer(t *testing.T) *httptest.Server {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        upgrader := websocket.Upgrader{}
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            t.Fatalf("error upgrading connection to WebSocket: %v", err)
        }
        defer conn.Close()
    }))

    return server
}

func TestWebSocketBroadcastNoConnections(t *testing.T) {
    server := MockWebSocketServer(t)
    fmt.Println(server.URL)
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
    fmt.Println(wsURL)
    defer server.Close()

    // Define the number of WebSocket connections to establish
    // numConnections := 100

    // Channels for synchronization
    // done := make(chan struct{})
    // errors := make(chan error)

    // Connect from multiple WebSocket clients concurrently
    // go func() {
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
    if err != nil {
        // errors <- err
        fmt.Println("error occured")
        fmt.Println(err)
        return
    }
    // defer conn.Close()

    message := []byte("test message")

    if err := conn.WriteMessage(1, message); err != nil {
        // errors <- err
        fmt.Println("error occured in writing message")
        fmt.Println(err)
        return
    }

    time.Sleep(time.Millisecond * 1000)

    _, receivedMessage, err := conn.ReadMessage()
    if err != nil {
        fmt.Println("error in reading mesage")
        fmt.Println(err)
        // errors <- err
        return
    }

    receivedMessageStr := string(receivedMessage)
    expectedMessage := "test message"
    if receivedMessageStr != expectedMessage {
        // errors <- ffmt.Errorf("expected message %q but received %q", expectedMessage, receivedMessageStr)
        // fmt.Errorf("expected message %q but received %q", expectedMessage, receivedMessageStr)
        fmt.Println("expected message something but received something different")
        return
    }

    // done <- struct{}{}
    // }()

    // Wait for all connections to complete
}
// TestWebSocketBroadcast tests the broadcast functionality of the WebSocket server.
func TestWebSocketBroadcast(t *testing.T) {
    server := MockWebSocketServer(t)
    fmt.Println(server.URL)
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
    fmt.Println(wsURL)
    defer server.Close()

    // Define the number of WebSocket connections to establish
    numConnections := 1

    // Channels for synchronization
    done := make(chan struct{})
    errors := make(chan error)

    // Connect from multiple WebSocket clients concurrently
    for i := 0; i < numConnections; i++ {
        go func() {
            conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
            if err != nil {
                errors <- err
                return
            }
            defer conn.Close()

            message := []byte("test message")

            if err := conn.WriteMessage(1, message); err != nil {
                errors <- err
                return
            }

            time.Sleep(time.Millisecond * 100)

            _, receivedMessage, err := conn.ReadMessage()
            if err != nil {
                errors <- err
                return
            }

            receivedMessageStr := string(receivedMessage)
            expectedMessage := "test message"
            if receivedMessageStr != expectedMessage {
                errors <- fmt.Errorf("expected message %q but received %q", expectedMessage, receivedMessageStr)
                return
            }

            done <- struct{}{}
        }()
    }

    // Wait for all connections to complete
    for i := 0; i < numConnections; i++ {
        select {
        case <-done:
        case err := <-errors:
            t.Fatalf("error connecting to WebSocket server: %v", err)
        }
    }
}

// TestWebSocketLatency tests the latency of the WebSocket server.
func TestWebSocketLatency(t *testing.T) {
    server := MockWebSocketServer(t)
    defer server.Close()

    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
    start := time.Now()
    conn, _, err := websocket.DefaultDialer.Dial(wsURL ,nil)
    if err != nil {
        t.Fatalf("error connecting to WebSocket server: %v", err)
    }
    defer conn.Close()

    elapsed := time.Since(start)
    t.Logf("WebSocket connection established in %v", elapsed)
}

func TestMain(m *testing.M) {
    // Run tests
    // print
    code := m.Run()

    // Perform any teardown tasks here

    // Exit with the test result code
    os.Exit(code)
}
