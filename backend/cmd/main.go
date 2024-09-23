package main

import (
	"context"
	"gohelp/cmd/config"
	"gohelp/cmd/handler"
	"gohelp/internal/service/auth"
	"gohelp/internal/service/forum"
	"gohelp/internal/storage"
	"gohelp/internal/storage/mongo"
	"gohelp/internal/storage/postgresql"
	"log"
	"net/http"
	"time"
)
// @title OverflowStack
// @description Community Assistent System

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	config.LoadEnv()
	mongodb := storage.CreateMongoClient(ctx)
	db := storage.InitDB()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("db is not connected")
	}
	forumdb := mongodb.Database("forum")
	forumRepo := mongo.NewForumStorage(forumdb, mongodb)
	userRepo := postgresql.NewUserRepository(db)
	userService := auth.NewUserService(userRepo)
	forumService := forum.NewForumService(forumRepo)
	userHandler := handler.NewHandler(userService, forumService)

	log.Fatal(http.ListenAndServe(":8080", userHandler.InitRoutes()))

}
