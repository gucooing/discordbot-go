package pgk

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var conn *websocket.Conn

func Wsls() {
	addr := GetConfig().Addr
	port := GetConfig().Port
	path := GetConfig().Path
	url := "ws://" + addr + ":" + port + "/" + path

	http.HandleFunc("/"+path, handleWebSocket)

	fmt.Println("WebSocket 服务器启动在：", url)
	err := http.ListenAndServe(addr+":"+port, nil)
	if err != nil {
		fmt.Printf("启动失败: ", err)
		return
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	newConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("无法升级到 WebSocket :", err)
		return
	}
	conn = newConn
	defer conn.Close()

	// 在这里处理 WebSocket 连接
	for {
		// 读取客户端发送的消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("读取消息时出错:", err)
			break
		}

		// 处理消息
		fmt.Printf("Received message: %s\n", msg)
	}
}

func SendMessage(message interface{}) string {
	// 检查是否已经存在连接
	if conn == nil {
		return "ws客户端未连接"
	}
	err := conn.WriteJSON(message)
	if err != nil {
		fmt.Printf("发送消息时出错:", err)
		return "发送命令时出错"
	}
	return "发送成功！"
}
