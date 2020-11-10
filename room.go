package main

import (
    "log"
    "net/http"
    "sync"
    
    "github.com/gorilla/websocket"
)

type room struct {
    broadcast chan Post
    clients sync.Map
    upgrader websocket.Upgrader
    helpList *IdList
    callList *IdList
}

func newRoom() *room {
    return &room {
        broadcast: make(chan Post),
        clients: sync.Map{},
        upgrader: websocket.Upgrader{},
        helpList: NewIdList(),
        callList: NewIdList(),
    }
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    client, err := r.upgrader.Upgrade(w, req, nil)
    if err != nil {
        // log.Fatal("error upgrading GET request to a websocket::", err)
        log.Printf("error upgrading GET request to a websocket: %v", err)
        return;
    }
    defer client.Close()

    r.clients.Store(client, true)
    initialBroadcast(r, client)
    
    for {
        var post Post
        err := client.ReadJSON(&post)
        if err != nil {
            log.Printf("error occurred while reading post: %v", err)
            r.clients.Delete(client)
            break
        }
        
        if post.Method == "help" {
            procHelp(r, &post)
        } else if post.Method == "call" {
            procCall(r, &post)
        } else if post.Method == "deleteHelp" {
            procDeleteHelp(r, &post)
        } else if post.Method == "deleteCall" {
            procDeleteCall(r, &post)
        } else {
            log.Printf("unknown post: %v", post)
        }
    }
}