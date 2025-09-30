package app

import (
	"context"
	"log"
	"teng231/goapp/internal/domain"
	"teng231/goapp/internal/infastructure/live"
	"teng231/goapp/internal/infastructure/sheetdb"

	"github.com/steampoweredtaco/gotiktoklive"
)

type LiveCommentApp struct {
	hub *live.LiveHub
	st  sheetdb.SheetDB
}

func NewLiveCommentApp(hub *live.LiveHub, st sheetdb.SheetDB) *LiveCommentApp {
	return &LiveCommentApp{hub: hub, st: st}
}

func (a *LiveCommentApp) Register(ctx context.Context, roomId string) chan gotiktoklive.Event {
	log.Print("addd ", roomId)
	room, err := a.hub.Add(roomId)
	if err != nil {
		log.Fatalf("GetLiveRoom error: %v", err)
	}
	ev := make(chan gotiktoklive.Event, 1000)
	go room.Event(ctx, ev)
	return ev
}

func (a *LiveCommentApp) Close(ctx context.Context, roomId string) error {
	err := a.hub.Remove(roomId)
	if err != nil {
		log.Fatalf("GetLiveRoom error: %v", err)
	}
	return err
}

func (a *LiveCommentApp) AppendSheet(ctx context.Context, comment *domain.Comment) error {
	err := a.st.AppendComment(comment)
	if err != nil {
		log.Fatalf("AppendSheet error: %v", err)
	}
	return err
}
