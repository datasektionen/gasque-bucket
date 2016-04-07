package datahandler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type display struct {
	ContentChan      chan *Content
	ConnClosed       chan struct{}
	CloseDisplayChan chan struct{}
	DisplayID        int
	RemoteAddr       string
}

//NewDisplay returns a pointer to a newly created Display
func newDisplay(id int, w http.ResponseWriter, r *http.Request) (*display, error) {

	d := display{
		ContentChan:      make(chan *Content),
		ConnClosed:       make(chan struct{}),
		CloseDisplayChan: make(chan struct{}),
		DisplayID:        id,
		RemoteAddr:       "",
	}

	d.RemoteAddr = r.RemoteAddr
	conn, err := upgrader.Upgrade(w, r, nil)
	_ = conn
	if err != nil {
		return nil, err
	}
	//The go routine should close when the conection is closed and parent should be
	//notified
	go func() {
		for {
			select {
			case c := <-d.ContentChan:
				if conn != nil && c != nil {
					err = conn.WriteMessage(
						websocket.TextMessage,
						[]byte(c.Path))
				}
				if err != nil {
					//handle error
				}
			case <-d.CloseDisplayChan:
				return
			}
		}
	}()

	return &d, nil
}
