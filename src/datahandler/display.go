package display

import (
	"net/http"

	eq "./../eventque"
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
	ContentChan chan *eq.Content
	ConnClosed  chan bool
	DisplayID   int
	RemoteAddr  string
}

//New returns a pointer to a newly created Display
func New(id int) *Display {
	return &Display{
		ContentChan: make(chan *eq.Content),
		ConnClosed:  make(chan bool),
		DisplayID:   id,
		RemoteAddr:  "",
	}

}

//Handle do
func (d *Display) Handle(w http.ResponseWriter, r *http.Request) error {

	content := make(chan *eq.Content, 10)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	//The go routine should close when the conection is closed and parent should be
	//notified
	go func() {
		for {
			select {
			case c := <-content:
				err = conn.WriteMessage(
					websocket.TextMessage,
					[]byte(c.Path))
				if err != nil {
					//handle error
				}
			}
		}
	}()
	return nil
}
