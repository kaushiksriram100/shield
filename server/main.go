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
	//router.GET("/urlinfo/1/:hostname_with_port/*original_path", LookUpMalwareDB)
	router.GET("/urlinfo/1/:hostname_with_port/*original_path", LookupMalwareEtcD)

	//Start httpServer
	log.Fatal(http.ListenAndServe(":8080", router))
}
