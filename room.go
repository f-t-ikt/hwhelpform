package main

import (
    "log"
    "net/http"
    "sync"
    
    "github.com/gorilla/websocket"
)

type Post struct {
    Method string `json:Method`
    Id     int    `json:Id`
}

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

func (r *room) initialBroadcast(client *websocket.Conn) {
    r.helpList.Each(func(v interface{}) bool {
        post := Post {
            Method: "help",
            Id: v.(int),
        }
        err := client.WriteJSON(post)
        if err != nil {
            log.Printf("error occurred while writing post to client: %v", err)
            client.Close()
            r.clients.Delete(client)
            return false
        }
        return true
    })
    
    if _, ok := r.clients.Load(client); !ok {
        return
    }
    
    r.callList.Each(func(v interface{}) bool {
        post := Post {
            Method: "call",
            Id: v.(int),
        }
        err := client.WriteJSON(post)
        if err != nil {
            log.Printf("error occurred while writing post to client: %v", err)
            client.Close()
            r.clients.Delete(client)
            return false
        }
        return true
    })
}

func (r *room) procHelp(post *Post) {
    id := post.Id
    if r.helpList.Contains(id) {
        return
    }
    
    if r.callList.Contains(id) {
        r.callList.Remove(id)
        del := Post {
            Method: "deleteCall",
            Id: id,
        }
        r.broadcast <- del
    }
    
    r.helpList.Add(id)
    r.broadcast <- *post
}

func (r *room) procCall(post *Post) {
    id := post.Id
    if r.callList.Contains(id) {
        return
    }
    
    if r.helpList.Contains(id) {
        r.helpList.Remove(id)
        del := Post {
            Method: "deleteHelp",
            Id: id,
        }
        r.broadcast <- del
    }
    
    r.callList.Add(id)
    r.broadcast <- *post
}

func (r *room) procDeleteHelp(post *Post) {
    r.helpList.Remove(post.Id)
    r.broadcast <- *post
}

func (r *room) procDeleteCall(post *Post) {
    r.callList.Remove(post.Id)
    r.broadcast <- *post
}

func (r *room) broadcastPostsToClients() {
    for {
        post := <- r.broadcast
        r.clients.Range(func(client, stored interface{})bool {
            err := client.(*websocket.Conn).WriteJSON(post)
            if err != nil {
                log.Printf("error occurred while writing post to client: %v", err)
                client.(*websocket.Conn).Close()
                r.clients.Delete(client)
            }
            return true
        })
    }
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    websocket, err := r.upgrader.Upgrade(w, req, nil)
    if err != nil {
        log.Fatal("error upgrading GET request to a websocket::", err)
        log.Printf("error upgrading GET request to a websocket: %v", err)
    }
    defer websocket.Close()

    r.clients.Store(websocket, true)
    r.initialBroadcast(websocket)
    
    for {
        var post Post
        err := websocket.ReadJSON(&post)
        if err != nil {
            log.Printf("error occurred while reading post: %v", err)
            r.clients.Delete(websocket)
            break
        }
        
        if post.Method == "help" {
            r.procHelp(&post)
        } else if post.Method == "call" {
            r.procCall(&post)
        } else if post.Method == "deleteHelp" {
            r.procDeleteHelp(&post)
        } else if post.Method == "deleteCall" {
            r.procDeleteCall(&post)
        } else {
            log.Printf("unknown post: %v", post)
        }
    }
}