package main

import (
	"fmt"
	"log"

	"github.com/steampoweredtaco/gotiktoklive"
)

func main() {
	t, err := gotiktoklive.NewTikTok()
	if err != nil {
		log.Fatalf("NewTikTok error: %v", err)
	}
	live, err := t.TrackUser("caothulienquan.official")
	if err != nil {
		log.Fatalf("TrackUser error: %v", err)
	}

	for event := range live.Events {
		switch e := event.(type) {
		case gotiktoklive.ChatEvent:
			// comment / chat
			fmt.Println("Chat:", e.User.Username, ":", e.Comment)
		// case gotiktoklive.GiftEvent:
		// 	fmt.Println("Gift:", e.User.Username, "gave", e.Gift.GetName())
		default:
			// fmt.Printf("Other event: %T %+v\n", e, e)
		}
	}
}

type LiveHub struct {
	RoomIds []string
	Rooms   []Room
}

type Room struct {
	*gotiktoklive.TikTok
	RoomId string
	Status string // connected / disconnected
}

type IRoom interface {
	Connect() error
	Disconnect() error
	Event() <-chan gotiktoklive.Event
	EventListen() error
}
