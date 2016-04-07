package datahandler

import (
	"net/http"
	"time"
)

type Event struct {
	EventID int

	displays          []*display
	addDisplayChan    chan *display
	addContentChan    chan *Content
	addPriorityChan   chan *Content
	removeContentChan chan int
	closeEventChan    chan struct{}

	queue *EventQueue
}

//NewEvent creates and returns a pointer to an Event
func NewEvent(id int) *Event {

	e := Event{
		EventID:           id,
		addDisplayChan:    make(chan *display),
		addContentChan:    make(chan *Content),
		addPriorityChan:   make(chan *Content),
		removeContentChan: make(chan int),
		closeEventChan:    make(chan struct{}),
		queue:             NewEventQueue(),
	}

	go func() {
		wait := time.Second
		var next *Content
		for {
			select {
			case d := <-e.addDisplayChan:
				e.displays = append(e.displays, d)
				//send the current image to the display connected
				if next != nil {
					d.ContentChan <- next
				}
			case c := <-e.addContentChan:
				e.queue.Push(c)
			case c := <-e.addPriorityChan:
				e.queue.PushPriority(c)
			case <-e.closeEventChan:
				for _, d := range e.displays {
					d.CloseDisplayChan <- struct{}{}
				}
				return
			case <-time.After(wait):
				next = e.queue.Next()
				e.queue.Push(next)
				if next != nil {
					e.sendToDisplays(next)
					wait = next.Duration
				}
			}
		}
	}()

	return &e
}
func (e *Event) sendToDisplays(content *Content) {
	for _, disp := range e.displays {
		disp.ContentChan <- content
	}
}

//AddDisplay will create a new display with the id specified.
func (e *Event) AddDisplay(id int, w http.ResponseWriter, r *http.Request) error {
	d, err := newDisplay(id, w, r)
	if err != nil {
		return err
	}
	e.addDisplayChan <- d
	return nil
}

//AddContent adds content to the end of the queue if priority is true
//the content will be prioritised.
func (e *Event) AddContent(c *Content, priority bool) {
	if priority {
		e.addPriorityChan <- c
	} else {
		e.addContentChan <- c
	}
}

//RemoveContent removes contet with specified id.
func (e *Event) RemoveContent(id int) {

}

func (e *Event) Close() {
	e.closeEventChan <- struct{}{}
}
