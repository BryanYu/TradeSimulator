package Server

import (
	"TradeSimulator/Models/Enum"
	"github.com/gofrs/uuid"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"slices"
	"sync"
)

var socketServerInstance *SocketServer
var once sync.Once

type SocketServer struct {
	server     *socketio.Server
	contextIds []string
	ID         string
}

func getSocketServerInstance() *SocketServer {
	id, _ := uuid.NewV4()
	once.Do(func() {
		socketServerInstance = &SocketServer{
			ID: id.String(),
		}
	})
	return socketServerInstance
}

func NewSocketServer() ISocketServer {
	instance := getSocketServerInstance()
	return instance
}

func (socket SocketServer) InitialServer() {
	instance := getSocketServerInstance()
	instance.server = socketio.NewServer(nil)
	instance.contextIds = make([]string, 0)

}
func (socket SocketServer) GetServer() *socketio.Server {
	instance := getSocketServerInstance()
	return instance.server
}

func (socket SocketServer) Start() {
	instance := getSocketServerInstance()
	go instance.server.Serve()
	defer instance.server.Close()
}

func (socket SocketServer) AddContextId(contextId string) {
	instance := getSocketServerInstance()
	instance.contextIds = append(socket.contextIds, contextId)
}

func (socket SocketServer) RemoveContextId(contextId string) {
	instance := getSocketServerInstance()
	index := slices.Index(socket.contextIds, contextId)
	instance.contextIds = append(socket.contextIds[:index], socket.contextIds[index+1:]...)
}
func (socket SocketServer) RegisterEvent(namespace string) {
	instance := getSocketServerInstance()
	instance.server.OnConnect(namespace, onConnect)
	instance.server.OnEvent(namespace, "send", onEvent)
	instance.server.OnError(namespace, onError)
	instance.server.OnDisconnect(namespace, onDisConnect)

}
func onConnect(s socketio.Conn) error {
	id := s.ID()
	instance := getSocketServerInstance()
	instance.AddContextId(id)
	s.Join("Stock1_LatestPrice")
	s.Join("Stock1_TradeLog")
	return nil
}

func onEvent(s socketio.Conn, socketEvent string) {
	switch socketEvent {
	case string(Enum.Join):
		break
	case string(Enum.Order):
		break
	default:
		break
	}
}

func onError(s socketio.Conn, err error) {
	log.Println("error:", err)
}

func onDisConnect(s socketio.Conn, reason string) {
	log.Println("disconnected:", s.ID(), reason)
	instance := getSocketServerInstance()
	instance.RemoveContextId(s.ID())
	s.LeaveAll()
}
