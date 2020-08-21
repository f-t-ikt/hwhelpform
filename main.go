package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
)

type subject struct {
    Id   string
    Name string
}

var templates = make(map[string]*template.Template)

func handleTemplates(name string, data subject) func(http.ResponseWriter, *http.Request) {
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
    templates["index"] = loadTemplate("_index")
    templates["student"] = loadTemplate("student")
    templates["teacher"] = loadTemplate("teacher")
    
    digital := subject {"digital", "Digital"}
    mpu := subject {"mpu", "MPU"}
    psoc := subject {"psoc", "PSoC"}

    http.Handle("/", http.FileServer(http.Dir("./templates")))
    http.Handle("/resources/", http.StripPrefix("/resources/",http.FileServer(http.Dir("./resources"))))
    http.HandleFunc("/digital", handleTemplates("index", digital))
    http.HandleFunc("/mpu", handleTemplates("index", mpu))
    http.HandleFunc("/psoc", handleTemplates("index", psoc))
    http.HandleFunc("/digital/student", handleTemplates("student", digital))
    http.HandleFunc("/mpu/student", handleTemplates("student", mpu))
    http.HandleFunc("/psoc/student", handleTemplates("student", psoc))
    http.HandleFunc("/digital/teacher", handleTemplates("teacher", digital))
    http.HandleFunc("/mpu/teacher", handleTemplates("teacher", mpu))
    http.HandleFunc("/psoc/teacher", handleTemplates("teacher", psoc))
    
    digitalRoom := newRoom()
    mpuRoom := newRoom()
    psocRoom := newRoom()
    http.Handle("/digital/update", digitalRoom)
    http.Handle("/mpu/update", mpuRoom)
    http.Handle("/psoc/update", psocRoom)
    go digitalRoom.broadcastPostsToClients()
    go mpuRoom.broadcastPostsToClients()
    go psocRoom.broadcastPostsToClients()
    
    server := http.Server {
        Addr: ":" + os.Getenv("PORT"),
    }
    server.ListenAndServe()
}
