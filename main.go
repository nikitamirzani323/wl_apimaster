package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nikitamirzani323/wl_api_master/db"
	"github.com/nikitamirzani323/wl_api_master/helpers"
	"github.com/nikitamirzani323/wl_api_master/routers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}

	initRedis := helpers.RedisHealth()

	if !initRedis {
		panic("cannot load redis")
	}

	db.Init()
	app := routers.Init()
	log.Fatal(app.Listen(":1011"))
}
