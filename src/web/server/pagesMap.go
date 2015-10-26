/*
Author: Tomas Möre 
Edited by: 

This file is meant to contain all first-level pages.

To add a page to the system the page file should run in it's initator one of the "Add" funtions

AdPublic
AdResticted

*/

package server

import (
	//"net/http"
	"fmt"
)

type PageHandler func(*WebIO)

type pageHandlerMap map[string]PageHandler

var publicPages pageHandlerMap = make(pageHandlerMap)
var restrictedPages pageHandlerMap = make(pageHandlerMap)



func GetPage(pathFragment string) PageHandler{
	pageHandler, ok := publicPages[pathFragment]
	if ok {
		return pageHandler;
	} else {
		panic("page doesn't exist")
	}
}



/*
function for adding a cosure to the page map
*/
func addToPageMap(pageMap pageHandlerMap, path string, hFunc PageHandler){
	ifExistsPanic(pageMap, path)

	pageMap[path] = hFunc
	
}

/*
This is only meant to be run at startup.
Adds a pageHandler to the map.
*/
func AddPublic(path string, hFunc PageHandler){
	addToPageMap(publicPages, path, hFunc)
	
}

func AddRestricted(path string, hFunc PageHandler){
	addToPageMap(restrictedPages, path, hFunc)
}
 
func ifExistsPanic(pageMap pageHandlerMap, path string){
	_, exists := pageMap[path]

	// If the pageHandler allready has been set. panic!
	if exists {
		panic(fmt.Sprintf("Page handler for \"%v\" allreadyExists", path))
	}
}


