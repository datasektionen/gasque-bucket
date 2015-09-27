package main

import (
	"./pages"
	"net/http"
)



func main(){
	http.HandleFunc("/Hello", pages.MakePage)     
	http.ListenAndServe(":8080", nil)
}




