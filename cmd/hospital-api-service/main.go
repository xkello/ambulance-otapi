package main

import (
	"context"
	"github.com/Serbel97/ambulance-api/api"
	"github.com/Serbel97/ambulance-api/internal/hospital_wl"
	"github.com/Serbel97/ambulance-api/internal/db_service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("HOSPITAL_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("HOSPITAL_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update  middleware
	dbService := db_service.NewMongoService[hospital_wl.Hospital](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	// request routings
	handleFunctions := &hospital_wl.ApiHandleFunctions{
		HospitalRolesAPI:  hospital_wl.NewHospitalRolesApi(),
		HospitalEmployeeListAPI: hospital_wl.NewHospitalEmployeeListApi(),
		HospitalsAPI:           hospital_wl.NewHospitalsApi(),
	}
	hospital_wl.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
