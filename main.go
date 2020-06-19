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

// クライアントからは JSON 形式で受け取る 
type Post struct {
    // Message string `json:message`
    Method string `json:Method`
    Id     int    `json:Id`
}

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

func initialBroadcast(client *websocket.Conn) {
    helpList.Each(func(v interface{}) bool {
        post := Post {
            Method: "help",
            Id: v.(int),
        }
        err := client.WriteJSON(post)
        if err != nil {
            log.Printf("error occurred while writing post to client: %v", err)
            client.Close()
            clients.Delete(client)
            return false
        }
        return true
    })
    
    if _, ok := clients.Load(client); !ok {
        return
    }
    
    callList.Each(func(v interface{}) bool {
        post := Post {
            Method: "call",
            Id: v.(int),
        }
        err := client.WriteJSON(post)
        if err != nil {
            log.Printf("error occurred while writing post to client: %v", err)
            client.Close()
            clients.Delete(client)
            return false
        }
        return true
    })
}

func procHelp(post *Post) {
    id := post.Id
    if helpList.Contains(id) {
        return
    }
    
    if callList.Contains(id) {
        callList.Remove(id)
        del := Post {
            Method: "deleteCall",
            Id: id,
        }
        broadcast <- del
    }
    
    helpList.Add(id)
    broadcast <- *post
}

func broadcastPostsToClients() {
    for {
        // メッセージ受け取り
        post := <-broadcast
        // クライアントの数だけループ
        // for client := range clients {
        //　書き込む
            // err := client.WriteJSON(post)
            // if err != nil {
                // log.Printf("error occurred while writing post to client: %v", err)
                // client.Close()
                // delete(clients, client)
            // }
        // }
        clients.Range(func(client, stored interface{})bool {
            err := client.(*websocket.Conn).WriteJSON(post)
            if err != nil {
                log.Printf("error occurred while writing post to client: %v", err)
                client.(*websocket.Conn).Close()
                clients.Delete(client)
            }
            return true
        })
    }
}
