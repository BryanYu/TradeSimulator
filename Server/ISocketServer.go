package Server

import (
	"TradeSimulator/Models/Enum"
	socketio "github.com/googollee/go-socket.io"
)

type ISocketServer interface {
	InitialServer(orderBook IOrderBook)
	GetServer() *socketio.Server
	RegisterEvent(namespace string)
	Start()
	Send(channelName Enum.SocketChannel, eventName Enum.SocketEvent, argument interface{})
}
