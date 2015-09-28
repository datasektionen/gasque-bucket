package eventque

import ("time")

/* Struct containing information about the media that will be sent to
 * the displays
 */
type DisplayContent struct {
	FileType string
	Path string
	Desc string
	Duration int //How long is the content intendet to be shown
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

/* Close the channel contained in a DisplayChannel */
func (dc DisplayChannel) closeChannel() {
	close(dc.channel)
}

/* Struct containing the channels used iternaly by the EventQueue
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

/* The EventQueue
 *
 */
type EventQueue struct {
	//Struct holding all the channels used to communicate with other threads
	channels EventQueChannels
	//Content not already shown
	content         []DisplayContent
	//Force pushed content intendet to be used by admin
	priorityContent []DisplayContent 
	//Content already shown which is stored to be shown again
	shownContent    []DisplayContent

	// Registred displays
	displayChannels []DisplayChannel
	/* Every display will get an unique id incremented by 1 from the last one
	created */
	displayChannelID int32
}

func (e *EventQueue) init() {
	//Channel used to add new content to the queue
	e.channels.inputChan = make(chan DisplayContent, 1024)
	//Channel used to force push new content to the queue
	e.channels.priorityChan = make(chan DisplayContent, 1024)
	//e.channels.queSettingsChan        = make(chan EventSettings)
	//Channel used to register new displays which will recieve content
	e.channels.registerDisplayChan = make(chan DisplayChannel)
	//Channel used to unregister a display
	e.channels.unregisterDisplayChan = make(chan DisplayChannel)
	//Channel used to send an exit signal to the queue
	e.channels.shouldExitChan = make(chan bool)

	e.displayChannelID = 0 //Used to give every DisplayChannel an unique id
}

func (e EventQueue) Start() {
	e.init()
	go e.startQueue()
}

func (e EventQueue) startQueue() {

	var newContent DisplayContent
	var currentContent DisplayContent
	//closedChanIds := []int{}
	for {
		select {
		case  shouldExit := <- e.channels.shouldExitChan:
			if shouldExit == true {
				/* Before we close the queue we have to close all
                                 * the channels
                                 */
				for _, channel := range e.displayChannels {
					channel.closeChannel()
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

		/* Add new display channels (there should be a way too delete
		them to */
		select {
		case newChannel := <- e.channels.registerDisplayChan:
			e.displayChannels = append(e.displayChannels, newChannel)
		default:
		}

		/* The content sent by the priority channel will always be 
                 * prioritised then the new content and last the content already
                 * shown
                 */
		currentContent = e.getNext()

		// If currentContent is empty and appended to the queue kittens
		// will die
		e.sendToDisplays(currentContent)

		time.Sleep(time.Second * 1)
	}
}

/* Get the next item in the queue */
func (e EventQueue) getNext() DisplayContent {
	
	var next DisplayContent
	if len(e.priorityContent) > 0 {
		next  = e.priorityContent[0]
		e.priorityContent = e.priorityContent[1:]
	} else if len(e.content) > 0 {
		next =  e.content[0]
		e.content = e.content[1:]
	} else if len(e.shownContent) > 0 {
		next  = e.shownContent[0]
		e.shownContent = e.shownContent[1:]
	}

	return next

}

/* Append next to shownContent and then send the struct to the registred displays */
func (e EventQueue) sendToDisplays(next DisplayContent) {
	/* Check if the struct is empty. We dont want an empty struct being appended
         * to the showContent slice */
	if next != (DisplayContent{}) {
		e.shownContent = append(e.shownContent, next)
		for _, channel := range e.displayChannels {
			channel.channel <- next
		}
	}
}

/* Add content to the queue */
func (e EventQueue) QueueContent(content DisplayContent) {
	e.channels.inputChan <- content
}

/* Force push content to the queue. Mainly intented for admins */
func (e EventQueue) QueuePriorityContent(content DisplayContent) {
	e.channels.priorityChan <- content
}

/* Called from an display to register its existanse to the EventQueue
 * a DisplayChannel struct will be returned containing the channel and
 * an unique ID
 */
func (e EventQueue) RegisterDisplay() DisplayChannel{
	c := make(chan DisplayContent)
	channel := DisplayChannel{e.displayChannelID, c}
	e.channels.registerDisplayChan <- channel

	e.displayChannelID++

	return channel
}

func (e EventQueue) UnregisterDisplay() {
}

/* Exit the queue */
func (e EventQueue) Exit() {
	e.channels.shouldExitChan <- true
}
