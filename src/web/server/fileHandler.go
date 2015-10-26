
/*
Author: Tomas Möre, 2015
Edited by:

*/
package server

import (
	//"net/http"
	"os"
	"fmt"
	md5 "crypto/md5"
	//"bytes"
	
)

var bufferSize int64 = 1024 * 8 // 8 kb buffer


/*
Standard file sender. Has a default buffer size. 
*/
func FileSender(file *os.File, webIO *WebIO, mimeType string) bool{
	// We need the file info to send a file.
	fileInfo, error := file.Stat()
	if error != nil{
		return false
	}

	webIO.AddHeader("Mime-Type", mimeType)
	
	fileSize := fileInfo.Size()

	webIO.AddHeader("Content-Length", fmt.Sprintf("%d", fileSize))


	webIO.WriteHeader(200)


	
	buffer   := make([]byte, bufferSize, bufferSize)

	var remainingBytes int64 = fileSize
	
	for remainingBytes != 0 {
		bytesRead, readError := file.Read(buffer)

		if readError != nil{
			return false
		}
		
		_,writeError := webIO.WriteBytes(buffer[:bytesRead])
	
		if (writeError != nil) {
			return false
		}
		remainingBytes -= int64(bytesRead)
	}
	// If all went ok we return true 
	return true


}

func FileSenderWithEtag(file *os.File, webIO *WebIO, mimeType string) bool {
	// We need the file info to send a file.
	fileInfo, error := file.Stat()
	if error != nil{
		return false
	}
	// If etag is valid we send the header and return 
	if validEtag(fileInfo, webIO){
		webIO.WriteHeader(304)
		return true
	}

	newEtag := makeEtag(fileInfo)


	// Currently hard coded to one day. This should be changed
        webIO.AddHeader("Cache-Control", fmt.Sprintf("public, %d", 86400))
	webIO.AddHeader("ETag", newEtag)

	// Finaly sends the file like normal (if we didn't chose to use the tag)
	return FileSender(file, webIO, mimeType) 

}





/* 
Checks for an ETag if it exist check if it still is valid. Else create a new one.
If the etag exists and is valid and exists return true else false.

Should make it such that if it return true it should NOT send the file (and use status code 304)
*/
func validEtag(fileInfo os.FileInfo, webIO *WebIO) bool {
	//eTag, ok := webIO.ReqHeader["If-None-Match"]
	return false
	
	
}



func makeEtag(fileInfo os.FileInfo) string{
	modTime := fileInfo.ModTime()
	hash := md5.Sum([]byte(modTime.String()))
	return fmt.Sprintf("%x", hash)
}
