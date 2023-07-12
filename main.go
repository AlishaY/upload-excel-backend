package main

import (
	"log"
	"upload-excel-backend/controller"
	"upload-excel-backend/database"
	// "upload-excel-backend/model"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	database.ConnectDB()

	// Create the Gin router
	router := gin.Default()

	// Initialize the controller with the database connection
	kpiController := controller.NewKPIController(database.DB)

	// Define the API routes
	api := router.Group("/api")
	{
		api.GET("/kpis", kpiController.GetKPIs)
		api.POST("/kpis", kpiController.CreateKPI)
	}

	// Start the server
	log.Fatal(router.Run(":8080"))
}
