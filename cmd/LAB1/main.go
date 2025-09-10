package main

import (
	"LAB1/internal/api"
	"log"
)

func main() {
	log.Println("Aplication start")
	api.StartServer()
	log.Println("Aplication terminated")
}
