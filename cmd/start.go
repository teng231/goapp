package cmd

import (
	"log"
	"os"
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
	api.Router().Listen(os.Getenv("PORT"))
}
