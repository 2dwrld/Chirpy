package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
    mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
    mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("assets"))))
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        _, err := w.Write([]byte("OK"))
        if err != nil {
            fmt.Println("error writing response: ", err)
        } 
    })
    
    server := http.Server {
        Addr: ":8080",
        Handler: mux,
    }
    log.Fatal(server.ListenAndServe())
}
