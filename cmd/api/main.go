package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const (
	webPort = 4321
)

func main() {
	fmt.Println("This is logging service")
	// connect to MongoDB

	h := NewHandler()

	err := http.ListenAndServe(fmt.Sprintf(":%d", webPort), h.router)
	if err != nil {
		log.Panic(err)
	}
}
