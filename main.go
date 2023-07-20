package main

import (
	"log"
	"upload-excel-backend/controller"
	"upload-excel-backend/database"
	// "upload-excel-backend/model"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/cors"
)

func main() {
	// Initialize the database connection
	database.ConnectDB()

	// Create the Gin router
	router := gin.Default()
	// router.Use(cors.Default())

	// Initialize the controller with the database connection
	kpiController := controller.NewKPIController(database.DB)

	// Define the API routes
		// router.GET("/kpis", kpiController.GetKPIs)
		// router.POST("/kpis", kpiController.CreateKPI)
		router.POST("/postfile", kpiController.PostFile)

	// Start the server
	log.Fatal(router.Run(":8080"))
}
