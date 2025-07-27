package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vvikash157/url_shortener/config"
	"github.com/vvikash157/url_shortener/controller"
	"github.com/vvikash157/url_shortener/repository"
	"github.com/vvikash157/url_shortener/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	godotenv.Load()

	db := config.InitDB()
	defer db.Close()

	redisClient := config.InitRedis()
	defer redisClient.Close()

	urlRepo := repository.NewPostgresUrlRepository(db)

	urlCache := repository.NewRedisCacheClient(redisClient)

	urlService := services.NewURLService(urlRepo, urlCache, "www.abc.com")

	controller := controller.NewUrlController(urlService)

	router.HandleFunc("/shortner", controller.ShortUrlHandler).Methods("POST")
	router.HandleFunc("/{code}", controller.RedirectOnURL).Methods("GET")

	c := cors.AllowAll()
	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default for local dev
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("shutting down server....")

		err := gracefulShutDown(server, 25*time.Second)
		if err != nil {
			log.Printf("getting error from while shutting Down %s", err.Error())
		}

		os.Exit(0)

	}()

	fmt.Printf("server started at %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error %s", err)
	}
	log.Fatal()
}

func gracefulShutDown(server *http.Server, maximumTime time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), maximumTime)
	defer cancel()
	fmt.Println("Server gracefully Shutdown")
	return server.Shutdown(ctx)
}
