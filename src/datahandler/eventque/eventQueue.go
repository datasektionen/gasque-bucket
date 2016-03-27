package eventque
/*
Author: Gustav N. Schneider

Queue used for holding DisplayContent items. The EventQueue is using two queues
with two different levels of priority. The items in the priorityQueue is always
chosen first.
*/
type DisplayContent struct {
	FileType string
	Path string
	Desc string
	Duration int
}

type queueItem struct {
	data *DisplayContent
	next *queueItem
}

type EventQueue struct {
	head *queueItem
	tail *queueItem
	prioHead *queueItem
	prioTail *queueItem
	size int
	prioSize int
}

func (e *EventQueue) Size() int {
	return e.size + e.prioSize
}

func (e *EventQueue) Next() *DisplayContent {
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

func (e *EventQueue) Push(d *DisplayContent) {
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

func (e *EventQueue) PushPriority(d *DisplayContent) {
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

func (e *EventQueue) Top() *DisplayContent {
	if e.prioSize != 0 {
		return e.prioHead.data
	}

	if e.size == 0 {
		return nil
	}
	return e.head.data
}

func (e *EventQueue) GetContent() []*DisplayContent {
	content := make([]*DisplayContent, e.Size())
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
