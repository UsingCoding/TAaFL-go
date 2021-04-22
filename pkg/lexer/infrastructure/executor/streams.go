package executor

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

func newInStream() *inStream {
	return &inStream{
		channel: make(chan string),
	}
}

type inStream struct {
	isEOF       bool
	channel     chan string
	hasProducer bool
}

func (i *inStream) Write(data string) error {
	if !i.hasProducer {
		return errors.New("there is no producer for input")
	}
	i.channel <- data
	return nil
}

func (i *inStream) Read(p []byte) (n int, err error) {
	fmt.Println("READ REQUEST")

	if i.isEOF {
		return 0, io.EOF
	}

	i.hasProducer = true

	data := <-i.channel

	fmt.Println("RECEIVED FROM CHANNEL")

	i.hasProducer = false

	data += "\n"

	n = copy(p, data[0:])

	return n, nil
}

func newOutStream() *outStream {
	return &outStream{
		channel: make(chan string),
	}
}

type outStream struct {
	channel     chan string
	hasListener bool
}

func (o *outStream) Read() string {
	o.hasListener = true
	return <-o.channel
}

func (o *outStream) Write(p []byte) (n int, err error) {
	fmt.Println("From process:", string(p))

	if !o.hasListener {
		fmt.Println("No listener for output")
		return len(p), nil
	}

	o.channel <- string(p)

	return len(p), nil
}
