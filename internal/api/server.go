package api

import (
	"LAB1/internal/app/handler"
	"LAB1/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("/Users/nachernev/Desktop/rip/templates/*")
	r.Static("/static", "/Users/nachernev/Desktop/rip/resources")

	r.GET("/indexes", handler.GetIndexes)
	r.GET("/index/:id", handler.GetIndex)
	r.GET("/request/:id", handler.RequestHandler)
	r.GET("/request/:id/up:num", handler.RequestHandler)
	r.GET("/request/:id/down:num", handler.RequestHandler)
	r.Run()
	log.Println("Server down")
}
