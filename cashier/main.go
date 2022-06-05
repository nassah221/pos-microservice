package main

import (
	"fmt"
	"log"

	db "pos-microservices/cashier/mongo"
)

func main() {
	cfg, err := db.NewConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Println(cfg)
}
