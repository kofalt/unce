package producer

type Producer interface {

	// Give the user help setting up.
	SetupHelp() string

	// Check for new UNCE.
	Poll(seen, log *bolt.DB) []*Event
}
