package main

import (
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/satori/go.uuid"
)

func wsPage(res http.ResponseWriter, req *http.Request) {
    conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
    if err != nil {
        http.NotFound(res, req)
        return
    }
    client := &Client{
        id: uuid.NewV4().String(),
        socket: conn,
        send: make(chan []byte),
    }

    manager.register <- client

    go client.read()
    go client.write()
}

func main() {
    println("hello, chaty")
    go manager.start()
    http.HandleFunc("/ws", wsPage)
    http.ListenAndServe(":12345", nil)
}
