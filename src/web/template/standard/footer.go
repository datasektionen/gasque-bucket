
/*
Author: Tomas Möre, 2015
Edited by:

*/

package main

import (
	server "./web/server"
	_ "./web/pages" // Loads all the pages 
	
)


func main(){
	
	server.Start()
	
}



