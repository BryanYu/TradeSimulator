package Server

import (
	"TradeSimulator/Models"
	"TradeSimulator/Models/Enum"
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var socketServerInstance *SocketServer
var once sync.Once

type SocketServer struct {
	server    *socketio.Server
	orderBook IOrderBook
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
	var allowOriginFunc = func(r *http.Request) bool {
		return true
	}
	instance := getSocketServerInstance()
	instance.server = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		}})
	instance.orderBook = orderBook
}

func (socket SocketServer) GetServer() *socketio.Server {
	instance := getSocketServerInstance()
	return instance.server
}

func (socket SocketServer) Start() {
	instance := getSocketServerInstance()

	go func() {
		if err := instance.server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer instance.server.Close()
}

func (socket SocketServer) RegisterEvent(namespace string) {
	instance := getSocketServerInstance()
	instance.server.OnConnect(namespace, onConnect)
	instance.server.OnEvent(namespace, string(Enum.Order), onOrderEvent)
	instance.server.OnError(namespace, onError)
	instance.server.OnDisconnect(namespace, onDisConnect)
}

func (socket SocketServer) Send(channelName Enum.SocketChannel, eventName Enum.SocketEvent, argument interface{}) {
	instance := getSocketServerInstance()
	instance.server.BroadcastToRoom("/", string(channelName), string(eventName), &argument)
}
func onConnect(s socketio.Conn) error {
	instance := getSocketServerInstance()

	s.Join(string(Enum.Price))
	s.Join(string(Enum.TradeLog))
	s.Join(string(Enum.BetterFivePrice))
	s.Join(s.ID())
	time.AfterFunc(100*time.Millisecond, func() {
		latestPrice := instance.orderBook.GetLatestPrice()
		if latestPrice != nil {
			instance.server.BroadcastToRoom("/", s.ID(), string(Enum.GetLatestPrice), latestPrice)
		}
		tradeLogs := instance.orderBook.GetTradeLogs()
		if len(tradeLogs) > 0 {
			instance.server.BroadcastToRoom("/", s.ID(), string(Enum.GetTradeLogs), tradeLogs)
		}

		betterFivePrice := instance.orderBook.GetBetterFivePrice()
		if len(betterFivePrice.Buy) > 0 || len(betterFivePrice.Sell) > 0 {
			instance.server.BroadcastToRoom("/", s.ID(), string(Enum.GetBetterFivePrice), betterFivePrice)
		}

	})
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
	s.LeaveAll()
}
