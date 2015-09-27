package datahandler.eventque

import ("fmt"
	"time")

/* Struct containing information about the media that will be sent to
 * the displays
 */
type DisplayContent struct {
	fileType string
	path string
	desc string
	duration int
}

//Not used at the moment
type EventSettings struct {
	imageDuration int
}

/* Struct containing a DisplayContent channel and an unique id. Used by
 * the queue to communicate with the displays
 */
type DisplayChannel struct {
	id int32
	channel chan DisplayContent
}

/* Struct containing the channels used iternaly by the eventQueue
 * to communicate with the outside world
 */
type EventQueChannels struct {
	inputChan           chan DisplayContent
	priorityChan        chan DisplayContent

	//queSettingsChan        chan EventSettings
	registerDisplayChan   chan DisplayChannel
	unregisterDisplayChan chan DisplayChannel
	shouldExitChan        chan bool
}

/* The eventQueue
 *
 */
type eventQueue struct {
	channels EventQueChannels

	content         []DisplayContent
	priorityContent []DisplayContent
	shownContent    []DisplayContent

	displayChannels []DisplayChannel
	displayChannelID int32
} 

func (e *eventQueue) initQueue() {
	e.channels.inputChan             = make(chan DisplayContent, 1024)
	e.channels.priorityChan          = make(chan DisplayContent, 1024)
	//e.channels.queSettingsChan        = make(chan EventSettings)
	e.channels.registerDisplayChan   = make(chan DisplayChannel)
	e.channels.unregisterDisplayChan = make(chan DisplayChannel)
	e.channels.shouldExitChan        = make(chan bool)

	e.displayChannelID = 0
}

func (e eventQueue) runEventQueue() {

	var newContent DisplayContent
	var currentContent DisplayContent
	//closedChanIds := []int{}
	for {
		select {
		case  shouldExit := <- e.channels.shouldExitChan:
			if shouldExit == true { 
				/* Before we close the queue we have to close all the
                                 * channels
                                 */
				for _, channel := range e.displayChannels {
					close(channel.channel)
				}

				return 
			}
		default:
		}

		//Add images from the normal channel
		select {
		case newContent = <- e.channels.inputChan:
			e.content = append(e.content, newContent)
		default:
		}

		//Add images from the priority channel
		select {
		case newContent = <- e.channels.priorityChan:
			e.priorityContent = append(e.priorityContent, newContent)
		default:
		}

		// Add new display channels (there should be a way too delete them to
		select {
		case newChannel := <- e.channels.registerDisplayChan:
			e.displayChannels = append(e.displayChannels, newChannel)
		default:
		}
		
		/* The content sent by the priority channel will always be prioritised
                 * then the new content and last the content already shown
                 */
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
				channel.channel <- currentContent
			}
		}
		time.Sleep(time.Second * 1)
	}
}

/* Add content to the queue */
func (e eventQueue) queContent(content DisplayContent) {
	e.channels.inputChan <- content
}

/* Force push content to the queue. Mainly intented for admins */
func (e eventQueue) quePriorityContent(content DisplayContent) {
	e.channels.priorityChan <- content
}

/* Register a channel to a new display */
func (e eventQueue) registerDisplay() DisplayChannel{
	c := make(chan DisplayContent)
	channel := DisplayChannel{e.displayChannelID, c}
	e.channels.registerDisplayChan <- channel

	e.displayChannelID++

	return channel
}

/* Exit the queue */
func (e eventQueue) exit() {
	e.channels.shouldExitChan <- true
}

func Display(name string, eq eventQueue) {
	fmt.Println("Registring channel")
	channel := eq.registerDisplay()
	fmt.Println(channel)
	for {
		var content DisplayContent
		content = <- channel.channel
		fmt.Println(name, ":", content.path)
	}
}

func main() {
	img1 := DisplayContent{"img", "/bin/img1", "img1", 5}
	img2 := DisplayContent{"img", "/bin/img2", "img2", 5}
	img3 := DisplayContent{"img", "/bin/img3", "img3", 5}
	img4 := DisplayContent{"img", "/bin/img4", "img4", 5}

	var myEventQueue eventQueue
	myEventQueue.initQueue()

	go myEventQueue.runEventQueue()
	
	go Display("display1", myEventQueue)
	go Display("display2", myEventQueue)

	myEventQueue.queContent(img1)
	myEventQueue.queContent(img2)
	myEventQueue.queContent(img3)

	myEventQueue.quePriorityContent(img4)

	time.Sleep(time.Second * 100)
}
