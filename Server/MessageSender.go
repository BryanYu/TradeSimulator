package Server

import (
	"TradeSimulator/Models/Enum"
)

type MessageSender struct {
	socketServer ISocketServer
}

func NewTradeLogSender(server ISocketServer) *MessageSender {
	return &MessageSender{
		socketServer: server,
	}
}
func (sender *MessageSender) Send(channelName Enum.SocketChannel, event Enum.SocketEvent, argument interface{}) {
	sender.socketServer.Send(channelName, event, argument)
}
