package main

import (
	"fmt"
	"log"
	"net/http"

	"vibework-chat/backend/internal/config"
)

func main() {
	cfg := config.Load()
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, http.NewServeMux()); err != nil {
		log.Fatal(err)
	}
}
