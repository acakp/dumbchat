package ws

import (
	"net"
	"net/http"
	"strings"

	"github.com/acakp/dumbchat/pkg/render"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

func HandleWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			render.Error(w, err, http.StatusUpgradeRequired, "Error upgrading to websockets")
			return
		}

		clientIp := clientIP(r)
		err = hub.trackConnection(clientIp)
		if err != nil {
			render.Error(w, err, http.StatusTooManyRequests, "Too many connections")
			return
		}
		client := &Client{
			ip:   clientIp,
			hub:  hub,
			conn: conn,
			send: make(chan []byte),
			rate: rate.NewLimiter(1, 5),
		}
		hub.Register <- client

		go client.writePump(hub)
		go client.readPump(hub)
	}
}

func clientIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	if ip == "127.0.0.1" {
		if real := r.Header.Get("X-Real-IP"); real != "" {
			return real
		}
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			return strings.Split(xff, ",")[0]
		}
	}

	return ip
}
