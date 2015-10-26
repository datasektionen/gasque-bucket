/*
Author: Tomas Möre 
Edited by:

This is the home page
*/
package public

import(
	server "../../server"
) 


func homeHandler(wIO *server.WebIO){
	wIO.Write("Hello world!")
}


func init(){
	server.AddPublic("", homeHandler)
}
