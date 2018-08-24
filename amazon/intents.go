package amazon

var (
	// AMAZON Skill Intents
	AMAZON = struct {
		StartOverIntent string
		CancelIntent    string
		PauseIntent     string
		HelpIntent      string
		StopIntent      string
		RepeatIntent    string
	}{
		StartOverIntent: "AMAZON.StartOverIntent",
		CancelIntent:    "AMAZON.CancelIntent",
		PauseIntent:     "AMAZON.PauseIntent",
		HelpIntent:      "AMAZON.HelpIntent",
		StopIntent:      "AMAZON.StopIntent",
		RepeatIntent:    "AMAZON.RepeatIntent",
	}
)
