package main

import (
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var clients sync.Map

var broadcast = make(chan Post)

var upgrader = websocket.Upgrader{}

var helpList, callList = NewIdList(), NewIdList()

func HandleClients(w http.ResponseWriter, r *http.Request) {
    websocket, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal("error upgrading GET request to a websocket::", err)
    }
    defer websocket.Close()

    clients.Store(websocket, true)
    initialBroadcast(websocket)
    
    for {
        var post Post
        err := websocket.ReadJSON(&post)
        if err != nil {
            log.Printf("error occurred while reading post: %v", err)
            // delete(clients, websocket)
            clients.Delete(websocket)
            break
        }
        
        if post.Method == "help" {
            procHelp(&post)
        } else if post.Method == "call" {
            procCall(&post)
        } else if post.Method == "deleteHelp" {
            procDeleteHelp(&post)
        } else if post.Method == "deleteCall" {
            procDeleteCall(&post)
        } else {
            log.Printf("unknown post: %v", post)
        }
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })
    
    http.HandleFunc("/teacher", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "teacher.html")
    })
    
    http.HandleFunc("/student", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "student.html")
    })
    
    http.HandleFunc("/update", HandleClients)
    go broadcastPostsToClients()
    
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("error starting http server::", err)
        return
    }
}
