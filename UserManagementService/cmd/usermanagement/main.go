package main

import (
	"advancedProg2Final/UserManagementService/pkg/api"
	"advancedProg2Final/UserManagementService/pkg/repository"
	"advancedProg2Final/UserManagementService/pkg/service"
	"log"
	"net/http"
)

func main() {
	database, err := repository.NewConnection("172.17.0.2", 5432, "postgres", "123456", "postgres")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	router := api.NewRouter(userService)

	log.Fatal(http.ListenAndServe(":8080", router))
}
