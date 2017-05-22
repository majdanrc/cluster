package cluster

import (
	"time"
)

type Topic struct {
	From     time.Time
	To       time.Time
	Messages []Message
}

func Classify(m Message, max int) {

}
