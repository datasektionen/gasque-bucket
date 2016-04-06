package datahandler

import (
	"fmt"
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

type Display struct {
	ContentChan chan *Content
	ConnClosed  chan bool
	DisplayID   int
	RemoteAddr  string
}

//New returns a pointer to a newly created Display
func NewDisplay(id int) *Display {
	return &Display{
		ContentChan: make(chan *Content),
		ConnClosed:  make(chan bool),
		DisplayID:   id,
		RemoteAddr:  "",
	}

}

//Handle do
func (d *Display) Handle(w http.ResponseWriter, r *http.Request) error {

	d.RemoteAddr = r.RemoteAddr
	conn, err := upgrader.Upgrade(w, r, nil)
	_ = conn
	if err != nil {
		return err
	}
	//The go routine should close when the conection is closed and parent should be
	//notified
	fmt.Println("Sending data to client......")
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
			}
		}
	}()
	return nil
}
