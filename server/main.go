package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//Global Variable with in memory list of malwareDB
var MalwareInMemoryDB MalwareDB

func init() {
	//In-memory malwareData. Load the data from a file
	MalwareInMemoryDB = LoadMalwareDataInMemory("./malwaredata/blacklist.json")
}
func main() {

	//Create a http router to handle requests
	router := httprouter.New()

	//routes
	router.GET("/urlinfo/1/:hostname_with_port/*original_path", LookUpMalwareDB(MalwareInMemoryDB))

	//Start httpServer
	log.Fatal(http.ListenAndServe(":8080", router))
}
