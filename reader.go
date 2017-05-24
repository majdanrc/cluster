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
	m := Message{
		Sender:    sender,
		MessageID: messageid,
		Content:   content,
	}

	if timestamp == "" || sender == "" || messageid == "" || content == "" {
		m.Error = errors.New("timestamp, sender, messageid and content are required")
		return m
	}

	var d int64
	d, m.Error = strconv.ParseInt(timestamp, 10, 64)
	m.Timestamp = d

	return m
}

type Reader struct {
	scanner *bufio.Scanner
}

func (r Reader) Read() <-chan Message {
	c := make(chan Message)

	go func() {
		var wg sync.WaitGroup

		for r.scanner.Scan() {
			wg.Add(1)
			go func(b string) {
				defer wg.Done()
				r := strings.Split(b, `;`)
				c <- newMessage(r[0], r[1], r[2], r[3])
			}(r.scanner.Text())
		}
		go func() {
			wg.Wait()
			close(c)
		}()
	}()

	return c
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(r),
	}
}
