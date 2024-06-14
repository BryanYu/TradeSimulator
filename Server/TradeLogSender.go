package Server

import (
	"TradeSimulator/Models"
	"TradeSimulator/Models/Enum"
	"sync"
)

type TradeLogSender struct {
	latestRecords []Models.TradeLog
	socketServer  ISocketServer
	lock          sync.Mutex
}

func NewTradeLogSender(server ISocketServer) *TradeLogSender {
	return &TradeLogSender{
		latestRecords: make([]Models.TradeLog, 5),
		socketServer:  server,
	}
}

func (sender *TradeLogSender) Send(channelName Enum.SocketChannel, tradeLog Models.TradeLog) {
	sender.socketServer.Send(string(channelName), tradeLog)
}

func (sender *TradeLogSender) enqueue(tradeLog Models.TradeLog) {
	if len(sender.latestRecords) == 5 {
		sender.dequeue()
	}
	sender.lock.Lock()
	defer sender.lock.Unlock()
	sender.latestRecords = append(sender.latestRecords, tradeLog)
}

func (sender *TradeLogSender) dequeue() Models.TradeLog {
	sender.lock.Lock()
	defer sender.lock.Unlock()
	if len(sender.latestRecords) == 0 {
		return Models.TradeLog{}
	}
	lastRecord := sender.latestRecords[0]
	sender.latestRecords = sender.latestRecords[1:]
	return lastRecord
}
