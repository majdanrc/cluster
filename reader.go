package cluster

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"
)

type Message struct {
	Timestamp int64
	Sender    string
	MessageID string
	Content   string
	Error     error
}

func newMessage(timestamp, sender, messageid, content string) Message {
	message := Message{
		Sender:    sender,
		MessageID: messageid,
		Content:   content,
	}

	if timestamp == "" || sender == "" || messageid == "" || content == "" {
		message.Error = errors.New("timestamp, sender, messageid and content are required")
		return message
	}

	var messageTimestamp int64
	messageTimestamp, message.Error = strconv.ParseInt(timestamp, 10, 64)
	message.Timestamp = messageTimestamp

	return message
}

type Reader struct {
	scanner *bufio.Scanner
}

func (r Reader) Read() <-chan Message {
	readChannel := make(chan Message)

	go func() {
		var wg sync.WaitGroup

		for r.scanner.Scan() {
			wg.Add(1)
			go func(b string) {
				defer wg.Done()
				r := strings.Split(b, `;`)
				readChannel <- newMessage(r[0], r[1], r[2], r[3])
			}(r.scanner.Text())
		}
		go func() {
			wg.Wait()
			close(readChannel)
		}()
	}()

	return readChannel
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(r),
	}
}
