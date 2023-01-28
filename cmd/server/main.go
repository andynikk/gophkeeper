package main

import "gophkeeper/internal/handlers"

// TODO: запуск сервера
func main() {
	srv := handlers.NewByConfig()
	srv.Run()
}
