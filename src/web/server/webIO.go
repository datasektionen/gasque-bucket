/*
Author: Tomas Möre 
Edited by: 

This request stuct is made such that we can make the reuqest data into what we want it to be.
All native data is stored in Req
*/

package server

import ("net/http"
	"net/url"
	"fmt"
//	str "strings"
)

/*
The main struct that holds all web related data. 

Handles both input and output.

Please use "MakeWebIO" if you want to make one of these
*/
type WebIO struct{
	ReqURL *url.URL
	//ReqCookies  map[string]http.cookie
	ReqHeader   http.Header

	
	Request    *http.Request
	Response    http.ResponseWriter
	
	RespHeader  http.Header

	headerWritten bool 
}




func (r *WebIO)parseCookies(){
	//cookieStr, ok := r.ReqHeader["Cookie"]
	//if !ok{
	//	return 
	//}
	return
	
	
	
}


func MakeWebIO(writer http.ResponseWriter, req *http.Request) *WebIO{
	webIO := &WebIO{
		Request          : req,
		Response         : writer,
		ReqHeader        : req.Header,
		ReqURL           : req.URL,
		RespHeader       : writer.Header(),
		headerWritten    : false,
	}
	webIO.parseCookies()
	return webIO
	

}


func (r *WebIO) WriteBytes(slc []byte) (int,error){
	return r.Response.Write(slc)
}
func (r *WebIO) Write(str string) {
	fmt.Fprint(r.Response, str)
}

func (r *WebIO) SetCookie(cookie *http.Cookie){
	http.SetCookie(r.Response, cookie)
}

func (r *WebIO) AddHeader(key, value string){
	r.RespHeader.Add(key, value)
}


func (r *WebIO) WriteHeader(statusCode int){
	defer func(){r.headerWritten = true}()
	r.Response.WriteHeader(statusCode)
	
	
}


func (r *WebIO) RunNative(handler func(http.ResponseWriter, *http.Request)){
	handler(r.Response,r.Request)
}
