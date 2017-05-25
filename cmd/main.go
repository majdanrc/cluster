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
	par := flag.Int("par", 0, "[optional] number of workers - default 4")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		return
	}

	var maxClusters int
	if *max > 0 {
		maxClusters = *max
	} else {
		maxClusters = 24
	}

	var workers int
	if *max > 0 {
		workers = *par
	} else {
		workers = 4
	}

	chatLog, err := os.Open(*file)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	clusters := make(map[string][]cluster.ClusteredMessage)

	var wg sync.WaitGroup

	input := cluster.NewReader(chatLog).Read()
	output := make(chan cluster.ClusteredMessage)

	for ind := 0; ind <= workers; ind++ {
		wg.Add(1)
		go func() {
			cluster.Classify(input, output, maxClusters)
			wg.Done()
		}()
	}

	var count int
	go func() {
		for item := range output {
			clusters[item.ClusterNo] = append(clusters[item.ClusterNo], item)
			count++
			fmt.Println(count)
		}
	}()

	wg.Wait()

	// for k := range clusters {
	// 	fmt.Printf("cluster [%s]: count[%d]\n", k, len(clusters[k]))
	// }

	close(output)

	defer chatLog.Close()
}
