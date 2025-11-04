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
	r.GET("/indexes/:id", handler.GetIndex)
	r.GET("/queries/:id", handler.QueryHandler)
	r.GET("/queries/:id/up/:num", handler.QueryHandlerUp)
	r.GET("/queries/:id/down/:num", handler.QueryHandlerDown)
	r.Run()
	log.Println("Server down")
}
