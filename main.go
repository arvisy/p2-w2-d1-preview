package main

import (
	"log"
	"net/http"
	"preview/config"
	"preview/handler"

	"github.com/julienschmidt/httprouter"
)

func main() {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	router := httprouter.New()

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	router.GET("/branches", handler.GetBranch)
	router.GET("/branches/:id", handler.GetBranchByID)
	router.POST("/branches", handler.CreateBranch)
	router.DELETE("/branches/:id", handler.DeleteBranchByID)
	router.PUT("/branches/:id", handler.UpdateBranchByID)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
