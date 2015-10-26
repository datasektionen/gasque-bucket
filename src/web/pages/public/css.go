/*
Author: Tomas Möre 
Edited by:

This is the css server
*/
package public

import(
	server "../../server"
)


// Hardcoded should change
var cssPath string = "data/css"



func handler(wIO *server.WebIO){
	wIO.Write("Hello world!")
}









func init(){
	server.AddPublic("css", handler)
}
