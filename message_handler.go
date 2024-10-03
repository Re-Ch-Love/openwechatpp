package openwechatpp

import (
	"fmt"

	ow "github.com/eatmoreapple/openwechat"
)

type Filter func(*ow.Message) (canHandle bool)

type Command struct {
	IsOnce  bool
	Name    string
	Usage   string
	Filter  Filter
	Handler func(*ow.Message) error
}

func (c Command) CheckAvailability() error {
	if c.Filter == nil {
		return fmt.Errorf("`Filter` field is nil")
	}
	if c.Handler == nil {
		return fmt.Errorf("`Handler` field is nil")
	}
	return nil
}
