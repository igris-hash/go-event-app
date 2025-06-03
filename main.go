package main

import (
	"github.com/gin-gonic/gin"
	"github.com/igris-hash/go-event-app/db"
	"github.com/igris-hash/go-event-app/routes"
)

func main() {
	db.InitDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8000")

}
