package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Gateway is alive"))
    })

    log.Println("Gateway running on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}