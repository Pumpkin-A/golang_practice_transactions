package main

import (
	"myGolangProject/Nastya/api"
	"myGolangProject/Nastya/models"
)

func main() {
	api.NewServer(models.GlobalConfig)
}
