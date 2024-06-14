package Server

import (
	"TradeSimulator/Models"
	"TradeSimulator/Models/Enum"
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"slices"
	"sync"
	"time"
)

var socketServerInstance *SocketServer
var once sync.Once

type SocketServer struct {
	server     *socketio.Server
	contextIds []string
	orderBook  IOrderBook
}

func getSocketServerInstance() *SocketServer {
	once.Do(func() {
		socketServerInstance = &SocketServer{}
	})
	return socketServerInstance
}

func NewSocketServer() ISocketServer {
	instance := getSocketServerInstance()
	return instance
}

func (socket SocketServer) InitialServer(orderBook IOrderBook) {
	instance := getSocketServerInstance()
	instance.server = socketio.NewServer(nil)
	instance.contextIds = make([]string, 0)
	instance.orderBook = orderBook

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
	instance.server.OnEvent(namespace, string(Enum.Order), onOrderEvent)
	instance.server.OnError(namespace, onError)
	instance.server.OnDisconnect(namespace, onDisConnect)

}

func (socket SocketServer) Send(channelName string, argument interface{}) {
	instance := getSocketServerInstance()
	instance.server.BroadcastToRoom("/", channelName, "broadcast", argument)
}
func onConnect(s socketio.Conn) error {
	id := s.ID()
	instance := getSocketServerInstance()
	instance.AddContextId(id)
	s.Join("Stock1_LatestPrice")
	s.Join("Stock1_TradeLog")
	return nil
}

func onOrderEvent(s socketio.Conn, orderJson string) {
	instance := getSocketServerInstance()
	log.Println(orderJson)
	var order Models.Order
	_ = json.Unmarshal([]byte(orderJson), &order)
	order.Timestamp = time.Now()
	instance.orderBook.AddOrder(&order)
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
