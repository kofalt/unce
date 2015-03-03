package github

import (
	. "fmt"
	// "encoding/json"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	. "kofalt.com/unce/def"
)

// create struct for the token source
type tokenSource struct {
	token *oauth2.Token
}

// add Token() method to satisfy oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error){
	return t.token, nil
}

type GithubProducer struct {
	client *github.Client
}

func New(gc *GithubConfig) *GithubProducer {

	if gc.AccessToken == "example" { return nil }

	ts := &tokenSource{
		&oauth2.Token{AccessToken: gc.AccessToken},
		// &oauth2.Token{AccessToken: "f60ee6eaf926304d8009d91c92af90f3da690f3b"},
	}
	tc := oauth2.NewClient(oauth2.NoContext, ts)


	return &GithubProducer{
		client: github.NewClient(tc),
	}
}

// URLS ARE AN ENTIRELY REASONABLE WAY TO SEND IDS
func GetLastSplit(s string) string {
	bits := strings.Split(s, "/")
	return bits[len(bits) - 1]
}

func (g *GithubProducer) Poll(seen, log *bolt.DB) []*Event {

	Println("Polling Github")
	notifications, response, err := g.client.Activity.ListNotifications(& github.NotificationListOptions{})

	if err != nil {
		Println("Github producer reports get notifications error:", err)
		return nil
	}

	// TODO: bother to paginate?
	_ = response

	var events []*Event
	for _, n := range notifications {
		event := &Event{}

		// TODO: save to bolt, don't send this ID forever
		Println("Processing github", *n.ID)

		if IsSeen(seen, "github", *n.ID) {
			Println("Already seen"); continue
		}

		// Log base JSON
		StoreJSON(log, "github", *n.ID + "-base", n)

		event.Summary = *n.Repository.FullName
		event.Message = "Activity"

		// Why would a notification have its content in the API call?
		// That would be convenient and efficient!
		// Let's totes dump 300 URLs in there. NBD.

		switch *n.Subject.Type {
			// blind faith that issues works too
			case "PullRequest", "Issue":
				// In soviet russia, URL parses you
				commentIdString := GetLastSplit(*n.Subject.LatestCommentURL)
				IdString := GetLastSplit(*n.Subject.URL)

				// For message formatting
				shortType := "PR #"
				if *n.Subject.Type == "Issue" { shortType = "issue #" }

				id, err := strconv.Atoi(commentIdString)
				if err != nil { Println("Github parsing PR comment ID fail:", err); continue }

				// Get comment content with comment ID
				prComment, _, err := g.client.Issues.GetComment(*n.Repository.Owner.Login, *n.Repository.Name, id)

				if err != nil && strings.Contains(err.Error(), "404 Not Found") {
					Println("Comment 404 on type", *n.Subject.Type, "not sure if actually error")


				} else if err != nil {
					Println("Github producer reports PR comment error:", err); continue
				} else {

					// Log comment JSON
					StoreJSON(log, "github", *n.ID + "-comment", prComment)

					event.Summary = *prComment.User.Login
					event.Message = *n.Repository.FullName + " " + shortType + IdString + ": " + *prComment.Body
				}

				MarkSeen(seen, "github", *n.ID)

			default:
				Println("Github, unknown subject type:", *n.Subject.Type)

				event.Message = "whelp! " + *n.Subject.Type
		}

		events = append(events, event)
	}

	return events
}
