package main

import (
    "log"
    "net/http"
    "os"
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
        // log.Fatal("error upgrading GET request to a websocket::", err)
        log.Printf("error upgrading GET request to a websocket: %v", err)
    }
    defer websocket.Close()

    clients.Store(websocket, true)
    initialBroadcast(websocket)
    
    for {
        var post Post
        err := websocket.ReadJSON(&post)
        if err != nil {
            // log.Printf("error occurred while reading post: %v", err)
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

func handleIndex(name string, data subject) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := templates[name].Execute(w, data); err != nil {
            log.Printf("failed to execute template: %v", err)
        }
    }
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles("templates/" + name + ".html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("./static")))
    
    http.HandleFunc("/update", HandleClients)
    go broadcastPostsToClients()
    
    server := http.Server {
        Addr: ":" + os.Getenv("PORT"),
    }
    server.ListenAndServe()
}
