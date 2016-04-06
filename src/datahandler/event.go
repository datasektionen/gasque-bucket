package display

import "time"

type Event struct {
	Displays      []*Display
	EventID       int
	AddDisplay    chan *Display
	AddContent    chan *Content
	AddPriority   chan *Content
	RemoveContent chan int

	queue EventQueue
}

func New(id int) *Event {
	return &Event{
		Displays:   new([]*Display),
		EventID:    id,
		AddDisplay: make(chan *Display),
	}
}

func (d *event) Run() error {
	go func() {
		for {
			wait := time.Second

			select {
			case d := <-d.AddDisplay:
				append(Displays, d)
			case c := <-d.AddContent:
				d.queue.Push(c)
			case c := <-d.AddPriority:
				d.queue.PushPriority(c)
			case <-time.After(wait):
				next := queue.Next()
				queue.Push(next)
				if next != nil {
					sendToDisplays(next)
					wait = next.Duration
				}
			}
		}
	}()
}
func (d *event) sendToDisplays(content *Content) {
	for id, disp := range d.Displays {
		disp.ContentChan <- content
	}
}
