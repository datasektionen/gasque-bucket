/*
Author: Gustav N. Schneider

Queue used for holding DisplayContent items. The EventQueue is using two queues
with two different levels of priority. The items in the priorityQueue is always
chosen first.
*/
package datahandler

import "time"

//DisplayContent is used to hold information about the objects to be displayed.
type Content struct {
	ID       int
	FileType int
	Path     string
	Desc     string
	Duration time.Duration
}

type queueItem struct {
	data *Content
	next *queueItem
}

type EventQueue struct {
	head     *queueItem
	tail     *queueItem
	prioHead *queueItem
	prioTail *queueItem
	size     int
	prioSize int
}

//NewEventQueue is used to create a new EventQueue
func NewEventQueue() *EventQueue {
	return new(EventQueue)
}

//Size returns the total size of the queue.
func (e *EventQueue) Size() int {
	return e.size + e.prioSize
}

//Next returns the next item in the queue.
func (e *EventQueue) Next() *Content {
	if e.prioSize != 0 {
		temp := e.prioHead.data
		e.prioHead = e.prioHead.next
		e.prioSize--
		return temp
	}

	if e.size == 0 {
		return nil
	}
	temp := e.head.data
	e.head = e.head.next
	e.size--
	return temp
}

//Push is used to add new content to the queue. The content added with push will be
//less prioritised that content pushed with PushPriority.
func (e *EventQueue) Push(d *Content) {
	n := new(queueItem)
	n.data = d
	n.next = nil

	if e.size == 0 {
		e.head = n
	} else {
		e.tail.next = n
	}
	e.tail = n
	e.size++
}

//PushPriority is used for pushing content with higher priority.
func (e *EventQueue) PushPriority(d *Content) {
	n := new(queueItem)
	n.data = d
	n.next = nil

	if e.prioSize == 0 {
		e.prioHead = n
	} else {
		e.prioTail.next = n
	}
	e.prioTail = n
	e.prioSize++
}

//Returns the item in front of the queue.
func (e *EventQueue) Top() *Content {
	if e.prioSize != 0 {
		return e.prioHead.data
	}

	if e.size == 0 {
		return nil
	}
	return e.head.data
}

//DeleteContent removes content with a given id.
func (e *EventQueue) DeleteContent(id int) error {
	return nil
}

//GetContent returns a pointer to a slice containing all the content in the queue
//ordered in the same way as the queue.
func (e *EventQueue) GetContent() []*Content {
	content := make([]*Content, e.Size())
	index := 0

	temp := e.prioHead
	for temp != nil {
		content[index] = temp.data
		index++
		temp = temp.next
	}

	temp = e.head
	for temp != nil {
		content[index] = temp.data
		index++
		temp = temp.next
	}

	return content
}
