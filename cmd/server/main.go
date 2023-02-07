// Package server: инициализация и запуск сервера
package main

import (
	"fmt"
	"gophkeeper/internal/handlers"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

// main запуск сервера
func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	srv := handlers.NewServer()
	srv.Run()
}
