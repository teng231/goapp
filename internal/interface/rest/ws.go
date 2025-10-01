package rest

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/steampoweredtaco/gotiktoklive"
)

func (a *API) WsBroadcast(roomId string, msg []byte) {
	for roomID, conn := range a.wsConn {
		if roomID != roomId {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	}
}

func (a *API) wsConnHandler(conn *websocket.Conn) {
	log.Print("connected")
	// Cấu hình ping/pong tự động (client phải reply pong)
	conn.SetPingHandler(func(appData string) error {
		log.Printf("Ping from %s", conn.RemoteAddr())
		return conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})
	conn.SetPongHandler(func(appData string) error {
		log.Printf("Pong from %s", conn.RemoteAddr())
		return nil
	})
	roomId := conn.Params("roomId")
	a.mt.Lock()
	a.wsConn[roomId] = conn
	a.mt.Unlock()
	// Lắng nghe message từ client
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		a.mt.Lock()
		delete(a.wsConn, roomId)
		a.mt.Unlock()
		log.Printf("Client %s disconnected", conn.RemoteAddr())
		a.liveHubApp.Close(context.Background(), roomId)
	}()
	cEvent := a.liveHubApp.Register(ctx, roomId)
	go func() {
		for {
			event := <-cEvent
			switch e := event.(type) {
			case gotiktoklive.ChatEvent:
				e.User.Badge = nil
				dat, _ := json.Marshal(e)
				err := conn.WriteMessage(websocket.TextMessage, dat)
				if err != nil {
					break
				}
			}
		}
	}()
	// Lắng nghe message từ client
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Print(err)
			// client ngắt hoặc lỗi đọc
			break
		}
		log.Printf("recv from %s: %s", conn.RemoteAddr(), msg)
		// echo lại
		if err = conn.WriteMessage(mt, msg); err != nil {
			break
		}
	}
}
