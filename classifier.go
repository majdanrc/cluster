package cluster

import "time"

type Topic struct {
	From     time.Time
	To       time.Time
	Messages []Message
}

type ClusteredMessage struct {
	ClNo int
	Msg  Message
}

func Classify(mc <-chan Message, oc chan<- ClusteredMessage, max int) {
	for m := range mc {

		clno := m.CreatedAt.Day()

		oc <- ClusteredMessage{
			ClNo: clno,
			Msg:  m,
		}
	}
}
