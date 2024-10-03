package openwechatpp

import (
	"strings"

	ow "github.com/eatmoreapple/openwechat"
)

func AcceptEverything(*ow.Message) bool {
	return true
}

func AcceptText(msg *ow.Message) bool {
	return msg.IsText()
}

func AcceptImage(msg *ow.Message) bool {
	return msg.IsPicture()
}

func AcceptSamePrefix(prefix string) Filter {
	return func(msg *ow.Message) bool {
		if !msg.IsText() {
			return false
		}
		return strings.HasPrefix(msg.Content, prefix)
	}
}

func AcceptSameContent(content string) Filter {
	return func(msg *ow.Message) bool {
		return msg.Content == content
	}
}

func AcceptRegexMatching(regex string) Filter {
	panic("TODO")
}

func extractSenderInfo(msg *ow.Message) (groupId string, userId string, err error) {
	sender, err := msg.Sender()
	if err != nil {
		return
	}
	if sender.IsGroup() {
		groupId = sender.AvatarID()
		var user *ow.User
		user, err = msg.SenderInGroup()
		if err != nil {
			return
		}
		userId = user.AvatarID()
	} else {
		userId = sender.AvatarID()
	}
	return
}

// 接受同源消息
// 即对于好友消息，是同一好友发送的
// 对于群聊消息，是同一群聊中的同一个人发送的
func ConstructSameOriginFilter(originMsg *ow.Message) (Filter, error) {
	originGroupId, originUserId, originErr := extractSenderInfo(originMsg)
	if originErr != nil {
		return nil, originErr
	}
	return func(msg *ow.Message) bool {
		groupID, userID, err := extractSenderInfo(msg)
		// 用if写出判断逻辑
		if err != nil {
			return false
		}
		return originGroupId == groupID && originUserId == userID
	}, nil
}
