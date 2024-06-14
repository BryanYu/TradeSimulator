package Enum

type SocketEvent string

const (
	Join  SocketEvent = "join"  //加入報價頻道
	Order SocketEvent = "order" //下單
)
