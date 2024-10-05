package openwechatpp

import (
	"fmt"
	"strings"
	"time"

	ow "github.com/eatmoreapple/openwechat"
)

type Dispatcher struct {
	Commands []*Command
}

func (d *Dispatcher) AsMessageHandler() ow.MessageHandler {
	return func(msg *ow.Message) {
		d.HandleMessage(msg)
	}
}

func (d *Dispatcher) AddCommand(cmd Command) {
	if err := cmd.CheckAvailability(); err != nil {
		panic(fmt.Sprintf("command is unavailable because %s", err.Error()))
	}
	d.Commands = append(d.Commands, &cmd)
}

func (d *Dispatcher) HandleMessage(msg *ow.Message) {
	for i, cmd := range d.Commands {
		if !cmd.Filter(msg) {
			continue
		}
		go cmd.Handler(msg)
		if cmd.IsOnce {
			d.Commands = append(d.Commands[:i], d.Commands[i+1:]...)
		}
		return
	}
}

func (d *Dispatcher) AwaitMatchingMessage(filter Filter, maxWaitingTime time.Duration) (*ow.Message, error) {
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

	d.Commands = append([]*Command{newCmd}, d.Commands...)

	select {
	case msg := <-ch:
		return msg, nil
	case <-timer:
		return &ow.Message{}, fmt.Errorf("timeout")
	}
}

// 等待给定msg的发送者的下一条消息直到超出maxInputTime
func (d *Dispatcher) WaitForNext(msg *ow.Message, maxInputTime time.Duration) (*ow.Message, error) {
	sameOriginFilter, err := ConstructSameOriginFilter(msg)
	if err != nil {
		return nil, err
	}
	return d.AwaitMatchingMessage(sameOriginFilter, maxInputTime)
}

// 输出帮助信息
// 格式如下：
// 【命令名】说明
// 【命令名】说明
func (d *Dispatcher) HelpText() string {
	builder := strings.Builder{}
	for _, cmd := range d.Commands {
		builder.WriteString(fmt.Sprintf("【%s】%s\n", cmd.Name, cmd.Usage))
	}
	return builder.String()
}
