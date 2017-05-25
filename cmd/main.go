package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/majdanrc/cluster"
)

func main() {
	file := flag.String("file", "", "file of chat logs, each line is a chat log, include timestamp, message ID, sender and content.")
	max := flag.Int("max", 0, "[optional] max allowed cluster number within a day (24 hours) - default 24")
	workers := flag.Int("workers", 0, "[optional] number of workers - default 4")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		return
	}

	maxClusters := parseParam(*max, 24)
	workerCount := parseParam(*workers, 4)

	chatLog, err := os.Open(*file)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	clusters := make(map[string][]cluster.ClusteredMessage)

	var wg sync.WaitGroup

	input := cluster.NewReader(chatLog).Read()
	output := make([]<-chan cluster.ClusteredMessage, workerCount)

	for index := 0; index < workerCount; index++ {
		output[index] = cluster.Classify(input, maxClusters)
	}

	for item := range merge(output) {
		clusters[item.ClusterNo] = append(clusters[item.ClusterNo], item)
	}

	wg.Wait()

	for k := range clusters {
		fmt.Printf("cluster [%s]: count[%d]\n", k, len(clusters[k]))
	}

	defer chatLog.Close()
}

func merge(cs []<-chan cluster.ClusteredMessage) <-chan cluster.ClusteredMessage {
	var wg sync.WaitGroup
	out := make(chan cluster.ClusteredMessage)

	output := func(c <-chan cluster.ClusteredMessage) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func parseParam(value int, defVal int) int {
	var out int

	if value > 0 {
		out = value
	} else {
		out = defVal
	}

	return out
}
