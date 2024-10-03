# OpenWeChatPP

OpenWeChatPP是对[openwechat](https://github.com/eatmoreapple/openwechat)的扩展，让编写机器人变得更方便。

## 示例

一个简单的天气查询机器人

![效果图](assets/image.png)

```go
bot := openwechat.DefaultBot(openwechat.Desktop)

// 你只需要构造一个调度器
dispatcher := openwechatpp.Dispatcher{}

// 把你想要的指令以声明式的语法添加进去
dispatcher.AddCommand(openwechatpp.Command{
    Name:   "天气查询",
    Usage:  "发送“/天气”查询天气",
    Filter: openwechatpp.AcceptSameContent("/天气"),
    Handler: func(msg *openwechat.Message) error {
        msg.ReplyText("请发送要查询的城市。")
        // 等待msg发送方的下一条消息
        replyMsg, err := dispatcher.WaitForNext(msg, time.Second*30)
        if err != nil {
            return err
        }
        _, err = replyMsg.ReplyText(fmt.Sprintf("%s今天的天气是……", replyMsg.Content))
        return err
    },
})

// 最后将调度器设置为bot的MessageHandler即可
bot.MessageHandler = dispatcher.AsMessageHandler()
```