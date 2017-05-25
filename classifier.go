package cluster

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type ClusteredMessage struct {
	ClusterNo string
	Msg       Message
}

var secondsInDay = 24 * 60 * 60

func round(f float64) int {
	return int(math.Floor(f + .5))
}

func Classify(input <-chan Message, max int) <-chan ClusteredMessage {
	out := make(chan ClusteredMessage)

	go func() {
		for message := range input {

			clusterSeconds := round(float64(secondsInDay) / float64(max))

			itemDate := time.Unix(message.Timestamp, 0).UTC()

			clusterPrefix := strconv.Itoa(itemDate.Year()) + strconv.Itoa(int(itemDate.Month())) + strconv.Itoa(itemDate.Day())
			secondsSinceMidnight := (itemDate.Hour() * 60 * 60) + (itemDate.Minute() * 60) + itemDate.Second()
			clusterSlot := secondsSinceMidnight / clusterSeconds

			cluster := clusterPrefix + "_" + strconv.Itoa(clusterSlot)

			fmt.Println(cluster)

			out <- ClusteredMessage{
				ClusterNo: cluster,
				Msg:       message,
			}
		}
		close(out)
	}()

	return out
}
