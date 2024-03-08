package triggers

import (
	"sync"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

type outgoingQueue struct {
	responses []*core2.ScheduledMessage
	lock      sync.Mutex
}

func (dr *outgoingQueue) Push(rsp *core2.ScheduledMessage) {
	dr.lock.Lock()
	dr.responses = append(dr.responses, rsp)
	dr.lock.Unlock()
}

func (dr *outgoingQueue) PopAll() []*core2.ScheduledMessage {
	dr.lock.Lock()
	output := dr.responses
	dr.responses = make([]*core2.ScheduledMessage, 0)
	dr.lock.Unlock()

	return output
}

var baseQueue *outgoingQueue

func usingOutgoingQueue() *outgoingQueue {
	if baseQueue == nil {
		baseQueue = &outgoingQueue{
			responses: make([]*core2.ScheduledMessage, 0),
			lock:      sync.Mutex{},
		}
	}

	return baseQueue
}

func outputQueueScheduler(svc *core2.Service) []*core2.ScheduledMessage {
	rsp := usingOutgoingQueue().PopAll()

	if rsp == nil || len(rsp) == 0 {
		return nil
	}

	return rsp
}
