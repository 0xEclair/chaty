package main

import (
    "encoding/json"
    
    "github.com/gorilla/websocket"
)

func (c *Client) read() {
    defer func() {
        manager.unregister <- c
        c.socket.Close()
    }()

    for {
        _, message, err := c.socket.ReadMessage()
        if err != nil {
            manager.unregister <- c
            c.socket.Close()
            break
        }
        jsonMessage, _ := json.Marshal(&Message{
            Sender: c.id,
            Content: string(message),
        })
        manager.broadcast <- jsonMessage
    }
}

func (c *Client) write() {
    defer func() {
        c.socket.Close()
    }()

    for {
        select {
        case message, ok := <- c.send: {
            if !ok {
                c.socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.socket.WriteMessage(websocket.TextMessage, message)
        }
        }
    }
}