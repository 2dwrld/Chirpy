package main

import (
	"encoding/json"
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
    w.Header().Add("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `<html>
                    <body>
                        <h1>Welcome, Chirpy Admin</h1>
                        <p>Chirpy has been visited %d times!</p>
                    </body>
                    </html>`, cfg.fileserverHits)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits = 0
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Hits reset to 0")
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, http.StatusText(http.StatusOK))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
    type Chirp struct {
        Body string `json:"body"`
    }
    type ErrorResponse struct {
        Error string `json:"error"`
    }
    type SuccessResponse struct {
        Valid bool `json:"valid"`
    }

    var chirp Chirp
    err := json.NewDecoder(r.Body).Decode(&chirp)
    if err != nil {
        fmt.Printf("error decoding %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{Error: "Something went wrong"})
        return
    }

    if len(chirp.Body) > 140 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"})
        return
    }
   
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(SuccessResponse{Valid: true})
}

func main() {
	mux := http.NewServeMux()
    apiCfg := apiConfig {
        fileserverHits: 0,
    }

    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

    mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
    mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
    mux.HandleFunc("GET /api/healthz", handlerHealthz)
    
    server := &http.Server {
        Addr: ":8080",
        Handler: mux,
    }
    log.Fatal(server.ListenAndServe())
}
