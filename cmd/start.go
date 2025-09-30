package cmd

import (
	"log"
	"teng231/goapp/internal/app"
	"teng231/goapp/internal/infastructure/live"
	"teng231/goapp/internal/interface/rest"

	"github.com/spf13/cobra"
)

func StartCLI(cmd *cobra.Command, args []string) {
	hub, err := live.NewLiveHub()
	if err != nil {
		log.Fatalf("NewLiveHub error: %v", err)
	}
	app := app.NewLiveCommentApp(hub)
	api := rest.New(app)
	api.Router().Listen(":3005")
	// room, err := hub.Add("ZSHWxnUK61JFg-6zFu7")
	// if err != nil {
	// 	log.Fatalf("NewLiveHub error: %v", err)
	// }
	// ev := make(chan gotiktoklive.Event, 1000)
	// go room.Event(context.Background(), ev)
	// for e := range ev {
	// 	log.Printf("Event: %v", e)
	// }

}
