package main

import "encoding/json"

func (manager *ClientManager) start() {
    for {
        select {
        case conn := <-manager.register: {
            manager.clients[conn] = true
            jsonMessage, _ := json.Marshal(&Message{
                Content: "/A new socket has connected.",
            })
            manager.send(jsonMessage, conn)
        }
        case conn := <-manager.unregister: {
            if _, ok := manager.clients[conn]; ok {
                close(conn.send)
                delete(manager.clients, conn)
                jsonMessage, _ := json.Marshal(&Message{
                    Content: "/A socket has disconnected.",
                })
                manager.send(jsonMessage, conn)
            }
        }
        case message := <-manager.broadcast: {
            for conn := range manager.clients {
                select {
                case conn.send <- message:
                default: {
                    close(conn.send)
                    delete(manager.clients, conn)
                }
                }
            }
        }
        }
    }
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
    for conn := range manager.clients {
        if conn == ignore {
            continue
        }
        conn.send <- message
    }
}