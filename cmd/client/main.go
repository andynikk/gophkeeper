// Package client: инициализация и запуск клиента
package main

import (
	"fmt"
	"gophkeeper/internal/client"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

// main запуск клиентского приложения
func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	c := client.NewClient()
	client.InitForms().Run(c)
}
