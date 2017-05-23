package main

import (
	"flag"
	"os"

	"github.com/majdanrc/cluster"
	"github.com/majdanrc/cluster/log"
)

func main() {
	file := flag.String("file", "", "file of chat logs, each line is a chat log, include timestamp, message ID, sender and content.")
	max := flag.Int("max", 0, "[optional] max allowed cluster number within a day (24 hours). This parameter allow to control the generated number of cluster within a 24 hours.")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		return
	}

	log.Info("cluster.parameters.file", "%s", *file)
	log.Info("cluster.parameters.max", "%d", *max)

	chat, err := os.Open(*file)
	if err != nil {
		log.Error("error.file", err.Error())
		return
	}

	output := make(chan cluster.ClusteredMessage, 10)

	inc := cluster.NewReader(chat).Read()

	for ind := 0; ind < 4; ind++ {
		go cluster.Classify(inc, output, 6)
	}

	for item := range output {
		log.Info("nic nie musze", "%v", item)
	}

	defer chat.Close()
}
