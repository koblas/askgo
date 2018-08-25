// Alexa Request and response types
package alexa

// RequestEnvelope contains the data passed from Alexa to the request handler.
type RequestEnvelope struct {
	Version string   `json:"version"`
	Session *Session `json:"session"`
	Request *Request `json:"request"`
	Context *Context `json:"context"`
}

// Session containes the session data from the Alexa request.
type Session struct {
	New        bool                   `json:"new"`
	SessionID  string                 `json:"sessionId"`
	Attributes map[string]interface{} `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken"`
	} `json:"user"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
}

// Context for the alexa request
type Context struct {
	AudioPlayer struct {
		PlayerActivity string `json:"playerActivity"`
	} `json:"AudioPlayer"`
	Display struct {
		Token string `json:"token"`
	} `json:"Display"`
	System struct {
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
		User struct {
			UserID string `json:"userId"`
		} `json:"user"`
		Device struct {
			DeviceID            string `json:"deviceId"`
			SupportedInterfaces struct {
				AudioPlayer struct {
				} `json:"AudioPlayer"`
				Display struct {
					TemplateVersion string `json:"templateVersion"`
					MarkupVersion   string `json:"markupVersion"`
				} `json:"Display"`
			} `json:"supportedInterfaces"`
		} `json:"device"`
		APIEndpoint    string `json:"apiEndpoint"`
		APIAccessToken string `json:"apiAccessToken"`
	} `json:"System"`
}

// Request contines the data in the request within the main request.
type Request struct {
	Locale      string `json:"locale"`
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	DialogState string `json:"dialogState"`
	Intent      Intent `json:"intent"`
	Name        string `json:"name"`
}

// Intent contains the data about the Alexa Intent requested.
type Intent struct {
	Name               string                `json:"name"`
	ConfirmationStatus string                `json:"confirmationStatus,omitempty"`
	Slots              map[string]IntentSlot `json:"slots"`
}

// IntentSlot contains the data for one Slot
type IntentSlot struct {
	Name               string `json:"name"`
	ConfirmationStatus string `json:"confirmationStatus,omitempty"`
	Value              string `json:"value"`
	ID                 string `json:"id,omitempty"`
}

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

// AudioPlayerDirective contains device level instructions on how to handle the response.
type AudioPlayerDirective struct {
	Type          string     `json:"type"`
	PlayBehavior  string     `json:"playBehavior,omitempty"`
	AudioItem     *AudioItem `json:"audioItem,omitempty"`
	ClearBehavior string     `json:"clearBehavior,omitempty"`
}

// AudioItem contains an audio Stream definition for playback.
type AudioItem struct {
	Stream            Stream             `json:"stream,omitempty"`
	AudioItemMetadata *AudioItemMetadata `json:"audioItemMetadata,omitempty"`
}

type AudioItemMetadata struct {
	Title    string `json:"title,omitempty`
	Subtitle string `json:"subtitle,omitempty`
	// Art             string `json:"art,omitempty`
	// BackgroundImage string `json:"backgroundImage,omitempty`
}

// Stream contains instructions on playing an audio stream.
type Stream struct {
	Token                 string `json:"token"`
	URL                   string `json:"url"`
	OffsetInMilliseconds  int    `json:"offsetInMilliseconds"`
	ExpectedPreviousToken string `json:"expectedPreviousToken,omitempty"`
}

// DialogDirective contains directives for use in Dialog prompts.
type DialogDirective struct {
	Type          string  `json:"type"`
	SlotToElicit  string  `json:"slotToElicit,omitempty"`
	SlotToConfirm string  `json:"slotToConfirm,omitempty"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
}

type RenderTemplateDirective struct {
	Type     string `json:"type"`
	Template string `json:"template,omitempty"`
}

type PlainTextHint struct {
	Type string `json:"type"`
	Text string `json:"template,omitempty"`
}

type HintDirective struct {
	Type string        `json:"type"`
	Hint PlainTextHint `json:"hint,omitempty"`
}

type VideoItemMetadata struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

type VideoItem struct {
	Source   string             `json:"source"`
	Metadata *VideoItemMetadata `json:"metadata,omitempty"`
}

type LaunchDirective struct {
	Type      string    `json:"type"`
	VideoItem VideoItem `json:"videoItem,omitempty"`
}
