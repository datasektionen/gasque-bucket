
/*
Author: Tomas Möre, 2015
Edited by:

*/
package server

import (
	"net/http"
	"time"
	"log"
	
)


type Handler struct{
	// DB connection
	
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request){
	webIO := MakeWebIO(writer, req)
	GetPage("")(webIO)
}






func Start(){

	server := &http.Server{
		Addr:           ":8080",   
		Handler:        Handler{}, // custom handler.
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	log.Fatal(server.ListenAndServe())   
	

}



