package main

import (
    // "log"
    "github.com/gorilla/websocket"
)

type Post struct {
    Method string `json:Method`
    Id     int    `json:Id`
}

func initialBroadcast(client *websocket.Conn) {
    helpList.Each(func(v interface{}) bool {
        post := Post {
            Method: "help",
            Id: v.(int),
        }
        err := client.WriteJSON(post)
        if err != nil {
            // log.Printf("error occurred while writing post to client: %v", err)
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
            // log.Printf("error occurred while writing post to client: %v", err)
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

func procCall(post *Post) {
    id := post.Id
    if callList.Contains(id) {
        return
    }
    
    if helpList.Contains(id) {
        helpList.Remove(id)
        del := Post {
            Method: "deleteHelp",
            Id: id,
        }
        broadcast <- del
    }
    
    callList.Add(id)
    broadcast <- *post
}

func procDeleteHelp(post *Post) {
    helpList.Remove(post.Id)
    broadcast <- *post
}

func procDeleteCall(post *Post) {
    callList.Remove(post.Id)
    broadcast <- *post
}

func broadcastPostsToClients() {
    for {
        post := <-broadcast
        clients.Range(func(client, stored interface{})bool {
            err := client.(*websocket.Conn).WriteJSON(post)
            if err != nil {
                // log.Printf("error occurred while writing post to client: %v", err)
                client.(*websocket.Conn).Close()
                clients.Delete(client)
            }
            return true
        })
    }
}
