package datahandler

import "time"

type Event struct {
	Displays      []*Display
	EventID       int
	AddDisplay    chan *Display
	AddContent    chan *Content
	AddPriority   chan *Content
	RemoveContent chan int

	queue *EventQueue
}

func NewEvent(id int) *Event {
	return &Event{
		EventID:       id,
		AddDisplay:    make(chan *Display),
		AddContent:    make(chan *Content),
		AddPriority:   make(chan *Content),
		RemoveContent: make(chan int),
		queue:         NewEventQueue(),
	}
}

func (e *Event) Run() {
	go func() {
		wait := time.Second
		var next *Content
		for {
			select {
			case d := <-e.AddDisplay:
				e.Displays = append(e.Displays, d)
				if next != nil {
					d.ContentChan <- next
				}
			case c := <-e.AddContent:
				e.queue.Push(c)
			case c := <-e.AddPriority:
				e.queue.PushPriority(c)
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
}
func (e *Event) sendToDisplays(content *Content) {
	for _, disp := range e.Displays {
		disp.ContentChan <- content
	}
}
