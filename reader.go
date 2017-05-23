package cluster

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Message struct {
	ID        string
	User      string
	Text      string
	CreatedAt time.Time
	Error     error
}

func newMessage(id, user, text, timestamp string) Message {
	m := Message{
		ID:   id,
		User: user,
		Text: text,
	}

	if id == "" || user == "" || text == "" || timestamp == "" {
		m.Error = errors.New("id, user, text and time is required")
		return m
	}

	var d int64
	d, m.Error = strconv.ParseInt(timestamp, 10, 64)
	m.CreatedAt = time.Unix(d, 0).UTC()

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
				r := strings.Split(b, `,`)
				c <- newMessage(r[1], r[2], r[3], r[0])
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
