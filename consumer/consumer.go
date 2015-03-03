package consumer

type Consumer interface {

	// Give the user help setting up.
	SetupHelp() string

	// Create & consume a test event
	Test()

	// Consume an event
	Consume(*Event)

}
