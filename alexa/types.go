// Package alexa is the JSON structure for Alexa Request and response types
package alexa

// ResponseEnvelope contains the Response and additional attributes.
type ResponseEnvelope struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          *Response              `json:"response"`
}

// Response contains the body of the response.
type Response struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
	ShouldSessionEnd bool          `json:"shouldEndSession"`
}

// OutputSpeech contains the data the defines what Alexa should say to the user.
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

// Card contains the data displayed to the user by the Alexa app.
type Card struct {
	Type        string   `json:"type"`
	Title       string   `json:"title,omitempty"`
	Content     string   `json:"content,omitempty"`
	Text        string   `json:"text,omitempty"`
	Image       *Image   `json:"image,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// Image provides URL(s) to the image to display in resposne to the request.
type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// Reprompt contains data about whether Alexa should prompt the user for more data.
type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech,omitempty"`
}

// PlainTextHint -
type PlainTextHint struct {
	Type string `json:"type"`
	Text string `json:"template,omitempty"`
}

// HintDirective -
type HintDirective struct {
	Type string        `json:"type"`
	Hint PlainTextHint `json:"hint,omitempty"`
}

// VideoItemMetadata -
type VideoItemMetadata struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

// VideoItem -
type VideoItem struct {
	Source   string             `json:"source"`
	Metadata *VideoItemMetadata `json:"metadata,omitempty"`
}

// LaunchDirective -
type LaunchDirective struct {
	Type      string    `json:"type"`
	VideoItem VideoItem `json:"videoItem,omitempty"`
}
