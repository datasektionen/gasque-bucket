package eventque
import ("fmt")

type Handler struct {
	queues map[int]*EventQueue
	queueCount int
}

func (*Handler) getQueueByID(queueID int) (*EventQueue, error) {
	if q, ok := Handler.queues[queueID]; !ok {
		return nil, fmt.Errorf("No queue with id %d", queueID)
	}
	return q, nil

}

func (*Handler) changeIdByID(oldID int, newID int) error {
	if _, ok := Handler.queues[newID]; ok {
		return fmt.Errorf("Queue with id %d already exists!", newID)
	}
	temp := Handler.queues[oldID]
	delete(Handler.queues, oldID)
	Handler.queues[newID] := temp
	return nil
}
func (*Handler) addDisplayContent(queueID int, content *DisplayContent) error {
	if q, ok := Handler.queue[queueID]; !ok {
		return fmt.Errorf("No qeueue with id %d", queueID)
	}
	q.Push(content)
	return nil
}
func (*Handler) getContent(queueID int) ([]*DisplayContent, error) {
	if q, ok := Handler.queue[queueID]; !ok {
		return fmt.Errorf("No qeueue with id %d", queueID)
	}
	return q.GetContent()
}

//getQueuesIDs -> array of ids
//saveQueuesToDisk -> error
//deleteQueue: id -> error
//deleteContent: queID, contentID -> error

