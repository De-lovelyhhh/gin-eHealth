package wshandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type PushMessage struct {
	From   string
	Target string
}

type tokenMessage struct {
	Token  string
	Target string
}

type Client struct {
	Id       string
	Identity int
	Socket   *websocket.Conn
	Send     chan []byte
}

type ClientManager struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	message    []PushMessage
	token      []tokenMessage
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

var upgrader = websocket.Upgrader{}

func (wh *ClientManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			log.Println(err)
			return
		}

		conn.WriteMessage(websocket.TextMessage, []byte("wxm.alming"))
	}
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			if conn.Identity == 1 {
				for _, msg := range manager.message {
					if msg.Target == conn.Id {
						conn.Send <- []byte(msg.From)
					}
				}
			} else {
				for _, msg := range manager.token {
					if msg.Target == conn.Id {
						conn.Send <- []byte(msg.Token)
					}
				}
			}
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				conn.Socket.Close()
				close(conn.Send)
				delete(manager.clients, conn)
			}
		case message := <-manager.Broadcast:
			match := false
			rm := make(map[string]interface{})
			if err := json.Unmarshal(message, &rm); err != nil {
				log.Println("receive message error:", err)
				break
			}
			for conn := range manager.clients {
				if conn.Id == rm["target"].(string) {
					if rm["token"] != nil {
						conn.Send <- []byte(rm["token"].(string))
					} else {
						conn.Send <- []byte(rm["account"].(string))
					}
					match = true
				}
			}
			if !match {
				if rm["token"] != nil {
					if len(manager.token) == 0 {
						manager.token = make([]tokenMessage, 0)
					}
					tm := tokenMessage{
						Token:  rm["token"].(string),
						Target: rm["target"].(string),
					}
					manager.token = append(manager.token, tm)
				} else {
					if len(manager.message) == 0 {
						manager.message = make([]PushMessage, 0)
					}
					pm := PushMessage{
						From:   rm["account"].(string),
						Target: rm["target"].(string),
					}
					manager.message = append(manager.message, pm)
				}
			}
		}
	}
}

func Register(client *Client) {
	Manager.register <- client
}

func (c *Client) Write() {
	defer func() {
		Manager.unregister <- c
		c.Socket.Close()
		fmt.Println("写关闭了")
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("发送关闭")
				return
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				Manager.unregister <- c
				c.Socket.Close()
				fmt.Println("写不成功数据就关闭了")
				break
			}
			fmt.Println("写数据")
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.unregister <- c
		c.Socket.Close()
		fmt.Println("读关闭")
	}()
	for {
		_, message, errRead := c.Socket.ReadMessage()
		if errRead != nil {
			log.Println(errRead)
			return
		}

		rm := make(map[string]interface{})
		if err := json.Unmarshal(message, &rm); err != nil {
			log.Println("receive message error:", err)
			break
		}

		log.Println("ws receive data:", rm)
		c.Id = rm["account"].(string)
		c.Identity = int(rm["identity"].(float64))
		Register(c)
		if rm["target"] != "" {
			Manager.Broadcast <- []byte(message)
		}

	}
}
