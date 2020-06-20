package main

import (
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

// WebSocket サーバーにつなぎにいくクライアント
// var clients = make(map[*websocket.Conn]bool)
var clients sync.Map

// クライアントから受け取るメッセージを格納
var broadcast = make(chan Post)

// WebSocket 更新用
var upgrader = websocket.Upgrader{}

var helpList, callList = NewIdList(), NewIdList()

// クライアントのハンドラ
func HandleClients(w http.ResponseWriter, r *http.Request) {
    // websocket の状態を更新
    websocket, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal("error upgrading GET request to a websocket::", err)
    }
    // websocket を閉じる
    defer websocket.Close()

    // clients[websocket] = true
    clients.Store(websocket, true)
    initialBroadcast(websocket)
    
    for {
        var post Post
        // メッセージ読み込み
        err := websocket.ReadJSON(&post)
        if err != nil {
            log.Printf("error occurred while reading post: %v", err)
            // delete(clients, websocket)
            clients.Delete(websocket)
            break
        }
        
       // メッセージを受け取る
        // broadcast <- post
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
    // localhost:8080 でアクセスした時に index.html を読み込む
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
