package main

import (
	"fmt"
	"log"
	"net/http"

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
	fmt.Println("server started at 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
