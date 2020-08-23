package main

import (
    "github.com/gorilla/websocket"
)

type Post struct {
    Method string `json:Method`
    Id     int    `json:Id`
    Date   string `json:Date`
}

func initialBroadcast(r *room, client *websocket.Conn) {
    r.helpList.Each(func(v interface{}) bool {
        post := v.(*Post)
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
        post := v.(*Post)
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

func procHelp(r *room, post *Post) {
    id := post.Id
    if r.helpList.ContainsId(post) {
        return
    }
    
    if r.callList.ContainsId(post) {
        r.callList.RemoveById(post)
        del := Post {
            Method: "deleteCall",
            Id: id,
        }
        r.broadcast <- del
    }
    
    r.helpList.Add(post)
    r.broadcast <- *post
}

func procCall(r *room, post *Post) {
    id := post.Id
    if r.callList.ContainsId(post) {
        return
    }
    
    if r.helpList.ContainsId(post) {
        r.helpList.RemoveById(post)
        del := Post {
            Method: "deleteHelp",
            Id: id,
        }
        r.broadcast <- del
    }
    
    r.callList.Add(post)
    r.broadcast <- *post
}

func procDeleteHelp(r *room, post *Post) {
    r.helpList.RemoveById(post)
    r.broadcast <- *post
}

func procDeleteCall(r *room, post *Post) {
    r.callList.RemoveById(post)
    r.broadcast <- *post
}

func broadcastPostsToClients(r *room) {
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
