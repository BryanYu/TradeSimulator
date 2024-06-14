package Server

import (
	socketio "github.com/googollee/go-socket.io"
)

type ISocketServer interface {
	InitialServer(orderBook IOrderBook)
	GetServer() *socketio.Server
	RegisterEvent(namespace string)
	Start()
	Send(channelName string, argument interface{})
}
