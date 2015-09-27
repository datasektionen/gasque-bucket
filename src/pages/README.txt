package pages

import(
	"net/http"
	"fmt"
)




func MakePage(writer http.ResponseWriter, request *http.Request){
	fmt.Fprintf(writer, "Hello world")
}




