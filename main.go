package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
    fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits++
        next.ServeHTTP(w,r)
    })
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits = 0
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Hits reset to 0")
}


func main() {
	mux := http.NewServeMux()
    apiCfg := apiConfig {
        fileserverHits: 0,
    }

    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    //mux.Handle("/app/assets/", http.StripPrefix("/app/assets", http.FileServer(http.Dir("assets"))))
    mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
    mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
    mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, http.StatusText(http.StatusOK))
    })
    
    server := &http.Server {
        Addr: ":8080",
        Handler: mux,
    }
    log.Fatal(server.ListenAndServe())
}
