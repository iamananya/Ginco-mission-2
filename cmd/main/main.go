package main

import (
	"log"
	"os"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/iamananya/Ginco-mission-2/pkg/controllers"
	"github.com/iamananya/Ginco-mission-2/pkg/middlewares"
	"github.com/iamananya/Ginco-mission-2/pkg/routes"
)

func main() {
	// if len(os.Args) > 1 && os.Args[1] == "benchmark" {
	// 	// Run the benchmark and exit
	// 	runBenchmark()
	// 	return
	// }

	router := gin.Default()
	store := cookie.NewStore([]byte("secret-key"))
	router.Use(sessions.Sessions("session-movie", store))
	// CORS ISSUE
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}
	config.AllowHeaders = []string{"Content-Type", "access-control-allow-headers", "access-control-allow-methods", "access-control-allow-origin", "session-id", "Session-ID"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	router.POST("/login", controllers.Login)
	// Apply authentication middleware to all routes except login and register
	router.Use(middlewares.AuthMiddleware())
	routes.RegisterTicketRoutes(router)

	log.Fatal(router.Run(":9010"))
}
func runBenchmark() {
	// Run the benchmark tests
	benchmarkResult := testing.Benchmark(controllers.BenchmarkSeatBooking)

	// Print benchmark results
	benchmarkResultString := benchmarkResult.String()
	log.Println(benchmarkResultString)

	// Write benchmark results to a file
	file, err := os.Create("benchmark_results.txt")
	if err != nil {
		log.Fatal("Failed to create benchmark_results.txt:", err)
	}
	defer file.Close()

	file.WriteString(benchmarkResultString)
}
