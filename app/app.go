package app

import (
	. "fmt"

	"github.com/codegangsta/cli"

	"kofalt.com/unce/producer/github"
	"kofalt.com/unce/consumer/notify-send"
)

var App *cli.App

func init() {
	App = cli.NewApp()

	App.Name = "unce"
	App.Usage = "Unce Unce Unce Unce"
	App.Version = "0.0.1"

	App.Action = Run

	App.Commands = []cli.Command{
		{
			Name:   "unce",
			Usage:  "Become the Unce",
			Action: Run,
		},{
			Name:   "test",
			Usage:  "Test a consumer or producer",
			Action: Test,
		},{
			Name:   "setup",
			Usage:  "Get help setting up a consumer",
			Action: Setup,
		},
	}
}

func Run(c *cli.Context) {
	// Resolve to configured and db-open state
	config := LoadorCreate()
	seenDB, logDB := Bees()

	ghP := github.New(config.Producers.Github)

	ns := notifysend.New()

	if ghP != nil {
		events := ghP.Poll(seenDB, logDB)

		for _, event := range events {
			Println(event)
			ns.Consume(event)
		}
	}
}

func Test(c *cli.Context) {
}
func Setup(c *cli.Context) {
}
