package main

import ("fmt"
	"time")

type DisplayContent struct {
	fileType string
	path string
	desc string
	duration int
}

type EventSettings struct {
	imageDuration int
}

type EventQueChannels struct {
	queInputChan           chan DisplayContent
	quePriorityChan        chan DisplayContent

	//queSettingsChan        chan EventSettings
	queReqisterDisplayChan chan (chan DisplayContent)
	queShouldExitChan      chan bool
}

type eventQueue struct {
	channels EventQueChannels

	content         []DisplayContent
	priorityContent []DisplayContent
	shownContent    []DisplayContent

	displayChannels []chan DisplayContent
} 

func (e *eventQueue) initQueue() {
	e.channels.queInputChan           = make(chan DisplayContent, 1024)
	e.channels.quePriorityChan        = make(chan DisplayContent, 1024)
	//e.channels.queSettingsChan        = make(chan EventSettings)
	e.channels.queReqisterDisplayChan = make(chan (chan DisplayContent))
	e.channels.queShouldExitChan      = make(chan bool)
}

func (e eventQueue) runEventQueue() {

	var newContent DisplayContent
	var currentContent DisplayContent
	for {
		select {
		case  shouldExit := <- e.channels.queShouldExitChan:
			if shouldExit == true { return }
		default:
		}

		//Add images from the normal channel
		select {
		case newContent = <- e.channels.queInputChan:
			e.content = append(e.content, newContent)
		default:
		}

		//Add images from the priority channel
		select {
		case newContent = <- e.channels.quePriorityChan:
			e.priorityContent = append(e.priorityContent, newContent)
		default:
		}

		// Add new display channels (there should be a way too delete them to
		select {
		case newChannel := <- e.channels.queReqisterDisplayChan:
			e.displayChannels = append(e.displayChannels, newChannel)
		default:
		}
		
		if len(e.priorityContent) > 0 {
			currentContent  = e.priorityContent[0]
			e.priorityContent = e.priorityContent[1:]
		} else if len(e.content) > 0 {
			currentContent =  e.content[0]
			e.content = e.content[1:]
		} else if len(e.shownContent) > 0 {
			currentContent = e.shownContent[0]
			e.shownContent = e.shownContent[1:]
		}

		if currentContent != (DisplayContent{}) {
			e.shownContent = append(e.shownContent, currentContent)

			for _, channel := range e.displayChannels {
				channel <- currentContent
			}
		}
		time.Sleep(time.Second * 1)
	}
}

/* Add content to the queue */
func (e eventQueue) queContent(content DisplayContent) {
	e.channels.queInputChan <- content
}

/* Force push content to the queue. Mainly intented for admins */
func (e eventQueue) quePriorityContent(content DisplayContent) {
	e.channels.quePriorityChan <- content
}

/* Register a channel to a new display */
func (e eventQueue) registerDisplay(c chan DisplayContent) {
	e.channels.queReqisterDisplayChan <- c
}

/* Exit the queue */
func (e eventQueue) exit() {
	e.channels.queShouldExitChan <- true
}

func Display(c chan DisplayContent, name string) {
	for {
		var content DisplayContent
		content = <- c
		fmt.Println(name, ":", content.path)
	}
}

func main() {
	img1 := DisplayContent{"img", "/bin/img1", "img1", 5}
	img2 := DisplayContent{"img", "/bin/img2", "img2", 5}
	img3 := DisplayContent{"img", "/bin/img3", "img3", 5}
	img4 := DisplayContent{"img", "/bin/img4", "img4", 5}

	myEventQueue := new(eventQueue)
	myEventQueue.initQueue()

	go myEventQueue.runEventQueue()
	
	c1 := make(chan DisplayContent)
	c2 := make(chan DisplayContent)
	go Display(c1, "display1")
	go Display(c2, "display2")
	myEventQueue.registerDisplay(c1)
	myEventQueue.registerDisplay(c2)

	myEventQueue.queContent(img1)
	myEventQueue.queContent(img2)
	myEventQueue.queContent(img3)

	myEventQueue.quePriorityContent(img4)

	time.Sleep(time.Second * 100)
}
