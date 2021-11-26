package main

import (
	"ipkeeper/pkg/service"
	"log"
	"os"

	handler "ipkeeper/pkg/delivery/http"
	repo "ipkeeper/pkg/repository/mongo"

	"github.com/gin-gonic/gin"
)

func main() {
	mgoConnInfo := repo.MongoInfo{
		Host:     "172.16.20.12",
		UserName: os.Getenv("MGO_USER"),
		Password: os.Getenv("MGO_PASS"),
		Port:     27017,
	}

	mgoConnect := &repo.MongoConnect{}
	mgoClient := mgoConnect.ConnectToMongo(&mgoConnInfo)
	if mgoClient.Err != nil {
		log.Fatal(mgoClient.Err)
		return
	}

	router := gin.Default()

	repo := repo.NewPingRepository(mgoClient.DB, mgoClient.Client)
	svc := service.NewPingService(repo)
	handler.NewPingHandler(router, svc)

	router.Run()
}
