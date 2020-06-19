package main

import (
    "container/list"
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

// クライアントからは JSON 形式で受け取る 
type Post struct {
    // Message string `json:message`
    Method string `json:Method`
    Id     int    `json:Id`
}

type IdList struct {
    list *list.List
}

func NewIdList() *IdList {
    return &IdList {
        list: list.New(),
    }
}

func (il *IdList) Add(v interface{}) *list.Element {
    il.Lock()
    defer il.Unlock()
    retunr il.list.PushBack(v)
}

func (il *IdList) Remove(v interface{}) interface{} {
    il.Lock()
    defer il.Unlock()
    for e := il.list.Front(); e != nil; e = e.Next() {
        if e.Value == v {
            return il.list.Remove(e)
        }
    }
    return nil
}

func (il *IdList) Contains(v interface{}) bool {
    il.Lock()
    defer il.Unlock()
    for e := il.list.Front(); e != nil; e = e.Next() {
        if e.Value == v {
            return true
        }
    }
    return false
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
        broadcast <- post
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
