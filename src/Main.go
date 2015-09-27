package main

import (
	"./pages"
	"net/http"
	"time"
	"log"
)

type gasqueHandler struct{
	something int
}

func (f gasqueHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request){
	pages.MakePage(writer,req)
}
 

func main(){
	//var myHandler http.Handler = WebHandler{}
	
	server := &http.Server{
		Addr:           ":8080",
		Handler:        gasqueHandler{5}, //myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	log.Fatal(server.ListenAndServe())   
	
}



