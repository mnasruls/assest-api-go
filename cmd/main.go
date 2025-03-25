package main

import (
	"assets-api-go/internal/config"
	"assets-api-go/internal/server"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {

	env, errs := config.InitAndCheckEnv()
	if errs != nil {
		log.Fatal("Error init env :", errs)
	}

	log.Println("env :", env)

	db, err := config.InitDb(env)
	if err != nil {
		log.Fatal("Error open connection db :", err)
		return
	}

	// run auto migrate
	err = config.AutoMigrate(db)
	if err != nil {
		log.Fatal("Error migrate db :", err)
		return
	}

	// serve API
	api := server.NewRestApi(db)
	if err = api.Serve(":" + env.AppPort); err != nil {
		log.Fatal(err)
	}
}
