package network

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type WsConn struct {
	wsServer *WsServer
	conn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan  byte
	isClosed bool
	mutex sync.Mutex
}

func NewWsConn(wsServer *WsServer,conn *websocket.Conn) (wsConn *WsConn,err error) {
	wsConn = &WsConn{
		wsServer:wsServer,
		conn:conn,
		inChan:make(chan []byte,1000),
		outChan:make(chan []byte,1000),
	}
	go wsConn.readLoop()
	go wsConn.writeLoop()
	return wsConn,nil
}

//读消息
func (wsConn *WsConn)ReadMessage()(data []byte,err error)  {
	select {
	case data = <-wsConn.inChan:
	case <- wsConn.closeChan:
		err  = errors.New("connection is closed")
	}
	return
}

//写消息
func (wsConn *WsConn)WirteMessage(data []byte)(err error)  {
	select {
	case wsConn.outChan <- data:
	case <- wsConn.closeChan:
		err  = errors.New("connection is closed")
	}
	return
}

//关闭连接
func (wsConn *WsConn)Close(){
	wsConn.conn.Close()
	wsConn.mutex.Lock()
	wsConn.wsServer.OnClose(wsConn)
	wsConn.mutex.Unlock()
}

//循环读消息
func (wsConn *WsConn)readLoop() {
	var(
		data []byte
		err error
	)
	for {
		if _,data,err = wsConn.conn.ReadMessage();err != nil {
			goto ERR
		}
		select {
		case wsConn.inChan <- data:
		case <- wsConn.closeChan:
			goto ERR
		}
	}
ERR:
	wsConn.Close()
}

//循环写消息
func (wsConn *WsConn) writeLoop()  {
	var(
		data []byte
		err error
	)
	for {
		data = <-wsConn.outChan
		if err = wsConn.conn.WriteMessage(websocket.TextMessage,data); err != nil{
			goto ERR
		}
	}
	select {
	case data = <-wsConn.outChan:
	case <- wsConn.closeChan:
		goto ERR
	}
ERR:
	wsConn.Close()
}