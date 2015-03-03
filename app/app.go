package app

import (
	. "fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"

	. "kofalt.com/unce/def"
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
		},{
			Name:   "whelp",
			Usage:  "Log for eternity",
			Action: Whelp,
		},
	}
}

func Run(c *cli.Context) {
	// Resolve to configured and db-open state
	config := LoadorCreate()
	seenDB, logDB := Bees()

	// Sup dawgs
	ghP := github.New(config.Producers.Github)
	ns := notifysend.New()

	// Later, add a loop over non-nil producers
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
func Whelp(c *cli.Context) {
	seenDB, logDB := Bees()

	var db *bolt.DB

	if len(c.Args()) != 2 {
		Println("Requires two arguments: database, and bucket.")
		os.Exit(1)
	}

	switch c.Args()[0] {
		case "seen":
			db = seenDB
		case "log":
			db = logDB
		default:
			Println("First argument must be 'seen' or 'log'.")
			os.Exit(1)
	}

	PrintKeys(db, c.Args()[1])
}
