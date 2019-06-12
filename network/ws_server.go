package network

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)
type WsServer struct {
	Name string		   	//服务器的名称
	IPVersion string    //服务器绑定的ip版本
	Ip string	        //服务器监听的IP
	Port int  	        //服务器监听的端口
	OnOpen	  func(wsConn *WsConn)      //监听WebSocket连接打开事件
	OnMessage func(wsConn *WsConn, msg []byte)   //监听WebSocket消息事件
	OnClose   func(wsConn *WsConn)    //监听WebSocket关闭事件
}

//创建一个服务器句柄
func NewWsServer (ip string,Port int) *WsServer {
	s:= WsServer {
		Ip:ip,
		Port: Port ,
	}
	return &s
}

//启动服务器
func (wsServer *WsServer) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var   (
				conn *websocket.Conn
				wsConn *WsConn
				msg []byte
				err error
			)
			upgrader := websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			if conn, err = upgrader.Upgrade(w,r,nil); err != nil {
				return
			}
			if wsConn,err = NewWsConn(wsServer,conn);err != nil{
				goto ERR
			}
			wsServer.OnOpen(wsConn)
			for  {
				if msg ,err = wsConn.ReadMessage(); err != nil {
					goto ERR
				}
				wsServer.OnMessage(wsConn,msg)
			}
		ERR:
			conn.Close()
	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", wsServer.Ip, wsServer.Port), nil)
	if err != nil {
		fmt.Println("websocket is start error")
	}
}

//停止服务器
func (wsServer *WsServer) Stop() {

}

//监听WebSocket连接打开事件
func (wsServer *WsServer) SetOnOpen(handler func(wsConn *WsConn)) {
	wsServer.OnOpen = handler
}

//监听WebSocket消息事件
func (wsServer *WsServer) SetOnMessage(handler func(wsConn *WsConn, msg []byte)) {
	wsServer.OnMessage = handler
}

//监听WebSocket关闭事件
func (wsServer *WsServer) SetOnClose(handler func(wsConn *WsConn)) {
	wsServer.OnClose = handler
}