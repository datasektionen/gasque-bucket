/*
Author: Tomas Möre 

Edited by: 


This module doesn't run anything at at all.
It loads all the pageHandlers such that they can register themselves into the system.

*/

package pages

import (
	_ "./public"
	_ "./restricted"
	"fmt"
)




func init(){
	fmt.Println("Pages loaded! :D")
}
