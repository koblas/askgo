package askgo

import (
	"fmt"
	"strings"

	"github.com/koblas/askgo/amazon"
)

type ResponseEnvelope struct {
	amazon.ResponseEnvelope
}

// ResponseBuilder that matches
type ResponseBuilder interface {
	Speak(speechOutput string) *ResponseEnvelope
	Reprompt(speechOutput string) *ResponseEnvelope
	WithSimpleCard(cardTitle, cardContent string) *ResponseEnvelope
	WithStandardCard(cardTitle, cardContent string, smallImageURL, largeImageURL *string) *ResponseEnvelope
	WithLinkAccountCard() *ResponseEnvelope
	WithAskForPermissionsConsentCard(permissions []string) *ResponseEnvelope
	AddDelegateDirective(updatedIntent *amazon.Intent) *ResponseEnvelope
	AddElicitSlotDirective(slotToElicit string, updatedIntent *amazon.Intent) *ResponseEnvelope
	AddConfirmSlotDirective(slotToConfirm string, updatedIntent *amazon.Intent) *ResponseEnvelope
	AddConfirmIntentDirective(updatedIntent *amazon.Intent) *ResponseEnvelope
	AddAudioPlayerPlayDirective(playBehavior, url, token string, offsetInMilliseconds int, expectedPreviousToken *string, audioItemMetadata *amazon.AudioItemMetadata) *ResponseEnvelope
	AddAudioPlayerStopDirective() *ResponseEnvelope
	AddAudioPlayerClearQueueDirective(clearBehavior string) *ResponseEnvelope
	AddRenderTemplateDirective(template string) *ResponseEnvelope
	AddHintDirective(text string) *ResponseEnvelope
	AddVideoAppLaunchDirective(source string, title, subtitle *string) *ResponseEnvelope
	WithShouldEndSession(val bool) *ResponseEnvelope
	AddDirective(directive interface{}) *ResponseEnvelope
	GetResponse() *ResponseEnvelope
}

// Verify that we're making the interface requirment
var _ ResponseBuilder = &ResponseEnvelope{}

func trimOutputSpeech(speechOutput string) string {
	speech := strings.TrimSpace(speechOutput)
	length := len(speech)

	if strings.HasPrefix(speech, "<speak>") && strings.HasSuffix(speech, "</speak>") {
		return speech[7 : length-8]
	}

	return speech
}

func (envelope *ResponseEnvelope) getResponse() *amazon.Response {
	if envelope.Response == nil {
		envelope.Response = &amazon.Response{}
	}
	return envelope.Response
}

// Speak - have Alexa say the provided speech to the user
func (envelope *ResponseEnvelope) Speak(speechOutput string) *ResponseEnvelope {
	response := envelope.getResponse()
	response.OutputSpeech = &amazon.OutputSpeech{
		Type: "SSML",
		SSML: fmt.Sprintf("<speak>%s</speak>", trimOutputSpeech(speechOutput)),
	}

	return envelope
}

// Reprompt - Has alexa listen for speech from the user. If the user doesn't respond
// within 8 seconds then has alexa reprompt with the provided reprompt speech
func (envelope *ResponseEnvelope) Reprompt(speechOutput string) *ResponseEnvelope {
	response := envelope.getResponse()
	response.Reprompt = &amazon.Reprompt{
		&amazon.OutputSpeech{
			Type: "SSML",
			SSML: fmt.Sprintf("<speak>%s</speak>", trimOutputSpeech(speechOutput)),
		},
	}

	return envelope
}

// WithSimpleCard renders a simple card with the following title and content
func (envelope *ResponseEnvelope) WithSimpleCard(cardTitle, cardContent string) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &amazon.Card{
		Type:    "Simple",
		Title:   cardTitle,
		Content: cardContent,
	}

	return envelope
}

// WithStandardCard - renders a standard card with the following title, content and image
func (envelope *ResponseEnvelope) WithStandardCard(cardTitle, cardContent string, smallImageUrl, largeImageUrl *string) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &amazon.Card{
		Type:  "Standard",
		Title: cardTitle,
		Text:  cardContent,
	}

	if smallImageUrl != nil || largeImageUrl != nil {
		response.Card.Image = &amazon.Image{}
		if smallImageUrl != nil {
			response.Card.Image.SmallImageURL = *smallImageUrl
		}
		if largeImageUrl != nil {
			response.Card.Image.LargeImageURL = *largeImageUrl
		}
	}

	return envelope
}

// WithLinkAccountCard - renders a link account card
func (envelope *ResponseEnvelope) WithLinkAccountCard() *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &amazon.Card{
		Type: "LinkAccount",
	}

	return envelope
}

// WithAskForPermissionsConcentCard - renders an askForPermissionsConsent card
func (envelope *ResponseEnvelope) WithAskForPermissionsConsentCard(permissions []string) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &amazon.Card{
		Type:        "AskForPermissionsConsent",
		Permissions: permissions,
	}

	return envelope
}

func (envelope *ResponseEnvelope) AddDelegateDirective(updatedIntent *amazon.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&amazon.DialogDirective{
		Type:          "Dialog.Delegate",
		UpdatedIntent: updatedIntent,
	})
}

func (envelope *ResponseEnvelope) AddElicitSlotDirective(slotToElicit string, updatedIntent *amazon.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&amazon.DialogDirective{
		Type:          "Dialog.ElicitSlot",
		SlotToElicit:  slotToElicit,
		UpdatedIntent: updatedIntent,
	})
}

func (envelope *ResponseEnvelope) AddConfirmSlotDirective(slotToConfirm string, updatedIntent *amazon.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&amazon.DialogDirective{
		Type:          "Dialog.ConfirmSlot",
		SlotToConfirm: slotToConfirm,
		UpdatedIntent: updatedIntent,
	})
}

func (envelope *ResponseEnvelope) AddConfirmIntentDirective(updatedIntent *amazon.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&amazon.DialogDirective{
		Type:          "Dialog.ConfirmIntent",
		UpdatedIntent: updatedIntent,
	})
}

func (envelope *ResponseEnvelope) AddAudioPlayerPlayDirective(
	playBehavior string,
	url string,
	token string,
	offsetInMilliseconds int,
	expectedPreviousToken *string,
	audioItemMetadata *amazon.AudioItemMetadata) *ResponseEnvelope {

	stream := amazon.Stream{
		Token:                token,
		URL:                  url,
		OffsetInMilliseconds: offsetInMilliseconds,
	}

	if expectedPreviousToken != nil {
		stream.ExpectedPreviousToken = *expectedPreviousToken
	}

	audioItem := &amazon.AudioItem{
		Stream:            stream,
		AudioItemMetadata: audioItemMetadata,
	}

	return envelope.AddDirective(&amazon.AudioPlayerDirective{
		Type:         "AudioPlayer.Play",
		PlayBehavior: playBehavior,
		AudioItem:    audioItem,
	})
}

func (envelope *ResponseEnvelope) AddAudioPlayerStopDirective() *ResponseEnvelope {
	return envelope.AddDirective(&amazon.AudioPlayerDirective{
		Type: "AudioPlayer.Stop",
	})
}

func (envelope *ResponseEnvelope) AddAudioPlayerClearQueueDirective(clearBehavior string) *ResponseEnvelope {
	return envelope.AddDirective(&amazon.AudioPlayerDirective{
		Type:          "AudioPlayer.ClearQueue",
		ClearBehavior: clearBehavior,
	})
}

func (envelope *ResponseEnvelope) AddRenderTemplateDirective(template string) *ResponseEnvelope {
	return envelope.AddDirective(&amazon.RenderTemplateDirective{
		Type:     "Display.RenderTemplate",
		Template: template,
	})
}

func (envelope *ResponseEnvelope) AddHintDirective(text string) *ResponseEnvelope {
	return envelope.AddDirective(&amazon.HintDirective{
		Type: "Hint",
		Hint: amazon.PlainTextHint{
			Type: "PlainText",
			Text: text,
		},
	})
}

func (envelope *ResponseEnvelope) AddVideoAppLaunchDirective(source string, title, subtitle *string) *ResponseEnvelope {
	videoItem := amazon.VideoItem{
		Source: source,
	}

	if title != nil || subtitle != nil {
		videoItem.Metadata = &amazon.VideoItemMetadata{}
		if title != nil {
			videoItem.Metadata.Title = *title
		}
		if subtitle != nil {
			videoItem.Metadata.Subtitle = *subtitle
		}
	}

	envelope.Response.ShouldSessionEnd = false

	return envelope.AddDirective(&amazon.LaunchDirective{
		Type:      "VideoApp.Launch",
		VideoItem: videoItem,
	})
}

// AddDirective - helper method for adding directives to responses
func (envelope *ResponseEnvelope) AddDirective(directive interface{}) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Directives = append(response.Directives, directive)

	return envelope
}

// WithShouldEndSession set the session end flag
func (envelope *ResponseEnvelope) WithShouldEndSession(val bool) *ResponseEnvelope {
	response := envelope.getResponse()

	// If we're launch a video session cannot end
	for _, d := range response.Directives {
		if launch, ok := d.(amazon.LaunchDirective); ok {
			if launch.Type == "VideoApp.Launch" {
				return envelope
			}
		}
	}

	envelope.getResponse().ShouldSessionEnd = val

	return envelope
}

// GetResponse - just return ourself
func (envelope *ResponseEnvelope) GetResponse() *ResponseEnvelope {
	return envelope
}
