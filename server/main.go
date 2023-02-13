package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	//Create a http router to handle requests
	router := httprouter.New()

	//routes
	//router.GET("/urlinfo/1/:hostname_with_port/*original_path", lookUpMalwareDB)
	router.GET("/urlinfo/1/:hostname_with_port/*original_path", lookupMalwareEtcD)

	//adminRouter creates http router to handle privileged operations (PUT, DELETE etc...)
	adminRouter := httprouter.New()
	//routes for adminRouter
	adminRouter.PUT("/urlinfo/1/:hostname_with_port/*original_path", putMalwareUrlToEtcD)
	adminRouter.DELETE("/urlinfo/1/:hostname_with_port/*original_path", deleteMalwareUrlInEtcD)
	//Start both http Servers
	go func() {
		log.Fatal(http.ListenAndServe(":8081", adminRouter))
	}()
	log.Fatal(http.ListenAndServe(":8080", router))
}
