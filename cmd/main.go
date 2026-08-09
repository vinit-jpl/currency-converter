package main

import (
	"context"
	"currency-converter/internal/config"
	"currency-converter/internal/handlers"
	"currency-converter/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	converterSvc := service.NewConverterService(cfg)
	converterHandler := handlers.NewConverterHandler(converterSvc)

	router := gin.Default()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	router.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})
	router.GET("/app.js", func(c *gin.Context) {
		c.File("./public/app.js")
	})
	router.GET("/styles.css", func(c *gin.Context) {
		c.File("./public/styles.css")
	})

	router.GET("/currencies", converterHandler.GetCurrencies)
	router.POST("/convert", converterHandler.ConvertCurrency)

	go func() {
		log.Printf("starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server....")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
