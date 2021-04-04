package api

import (
	"e_healthy/handler/wshandler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WsAuthApi(c *gin.Context) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		return
	}

	client := &wshandler.Client{
		Id:     "",
		Socket: conn,
		Send:   make(chan []byte),
	}

	go wshandler.Manager.Start()
	go client.Read()
	go client.Write()
}
