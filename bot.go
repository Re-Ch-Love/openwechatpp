package openwechatpp

import (
	"fmt"
	"time"

	ow "github.com/eatmoreapple/openwechat"
)

type Dispatcher struct {
	Commands []*Command
}

func (e *Dispatcher) AsMessageHandler() ow.MessageHandler {
	return func(msg *ow.Message) {
		e.HandleMessage(msg)
	}
}

func (e *Dispatcher) AddCommand(cmd Command) {
	if err := cmd.CheckAvailability(); err != nil {
		panic(fmt.Sprintf("command is unavailable because %s", err.Error()))
	}
	e.Commands = append(e.Commands, &cmd)
}

func (e *Dispatcher) HandleMessage(msg *ow.Message) {
	for i, cmd := range e.Commands {
		if !cmd.Filter(msg) {
			continue
		}
		go cmd.Handler(msg)
		if cmd.IsOnce {
			e.Commands = append(e.Commands[:i], e.Commands[i+1:]...)
		}
		return
	}
}

func (e *Dispatcher) AwaitMatchingMessage(filter Filter, maxWaitingTime time.Duration) (*ow.Message, error) {
	ch := make(chan *ow.Message, 1)
	timer := make(chan interface{}, 1)

	go func() {
		time.Sleep(maxWaitingTime)
		timer <- struct{}{}
	}()

	newCmd := &Command{
		Filter: filter,
		Handler: func(msg *ow.Message) error {
			ch <- msg
			close(ch)
			return nil
		},
	}

	e.Commands = append([]*Command{newCmd}, e.Commands...)

	select {
	case msg := <-ch:
		return msg, nil
	case <-timer:
		return &ow.Message{}, fmt.Errorf("timeout")
	}
}

// 等待给定msg的发送者的下一条消息直到超出maxInputTime
func (e *Dispatcher) WaitForNext(msg *ow.Message, maxInputTime time.Duration) (*ow.Message, error) {
	sameOriginFilter, err := ConstructSameOriginFilter(msg)
	if err != nil {
		return nil, err
	}
	return e.AwaitMatchingMessage(sameOriginFilter, maxInputTime)
}
