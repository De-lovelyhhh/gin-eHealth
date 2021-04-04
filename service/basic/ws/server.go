package ws

import (
	"e_healthy/pkg/setting"
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type wsBind struct {
	conn           *websocket.Conn
	userAccount    map[string]string
	patientAccount map[string]string
}

var addr = flag.String("addr", setting.IP_ADDR+":7003", "http server")

func server(w http.ResponseWriter, r *http.Request) {
	u := url.URL{
		Scheme: "ws",
		Host:   *addr,
		Path:   "/ws",
	}

	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		_, message, errRead := conn.ReadMessage()
		if errRead != nil {
			log.Println(errRead)
			return
		}

		// todo 把message传出去
		log.Println(message)
	}
}
