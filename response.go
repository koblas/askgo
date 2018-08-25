package askgo

import (
	"fmt"
	"strings"

	"github.com/koblas/askgo/alexa"
)

// ResponseEnvelope wrapper around askgo.alexa type
type ResponseEnvelope struct {
	alexa.ResponseEnvelope
}

// ResponseBuilder interface for building requests
type ResponseBuilder interface {
	Speak(speechOutput string) *ResponseEnvelope
	Reprompt(speechOutput string) *ResponseEnvelope
	WithSimpleCard(cardTitle, cardContent string) *ResponseEnvelope
	WithStandardCard(cardTitle, cardContent string, smallImageURL, largeImageURL *string) *ResponseEnvelope
	WithLinkAccountCard() *ResponseEnvelope
	WithAskForPermissionsConsentCard(permissions []string) *ResponseEnvelope
	AddDelegateDirective(updatedIntent *alexa.Intent) *ResponseEnvelope
	AddElicitSlotDirective(slotToElicit string, updatedIntent *alexa.Intent) *ResponseEnvelope
	AddConfirmSlotDirective(slotToConfirm string, updatedIntent *alexa.Intent) *ResponseEnvelope
	AddConfirmIntentDirective(updatedIntent *alexa.Intent) *ResponseEnvelope
	AddAudioPlayerPlayDirective(playBehavior, url, token string, offsetInMilliseconds int, expectedPreviousToken *string, audioItemMetadata *alexa.AudioItemMetadata) *ResponseEnvelope
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

func (envelope *ResponseEnvelope) getResponse() *alexa.Response {
	if envelope.Response == nil {
		envelope.Response = &alexa.Response{}
	}
	return envelope.Response
}

// Speak - have Alexa say the provided speech to the user
func (envelope *ResponseEnvelope) Speak(speechOutput string) *ResponseEnvelope {
	response := envelope.getResponse()
	response.OutputSpeech = &alexa.OutputSpeech{
		Type: "SSML",
		SSML: fmt.Sprintf("<speak>%s</speak>", trimOutputSpeech(speechOutput)),
	}

	return envelope
}

// Reprompt - Has alexa listen for speech from the user. If the user doesn't respond
// within 8 seconds then has alexa reprompt with the provided reprompt speech
func (envelope *ResponseEnvelope) Reprompt(speechOutput string) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Reprompt = &alexa.Reprompt{
		OutputSpeech: &alexa.OutputSpeech{
			Type: "SSML",
			SSML: fmt.Sprintf("<speak>%s</speak>", trimOutputSpeech(speechOutput)),
		},
	}

	return envelope
}

// WithSimpleCard renders a simple card with the following title and content
func (envelope *ResponseEnvelope) WithSimpleCard(cardTitle, cardContent string) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &alexa.Card{
		Type:    "Simple",
		Title:   cardTitle,
		Content: cardContent,
	}

	return envelope
}

// WithStandardCard - renders a standard card with the following title, content and image
func (envelope *ResponseEnvelope) WithStandardCard(cardTitle, cardContent string, smallImageURL, largeImageURL *string) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &alexa.Card{
		Type:  "Standard",
		Title: cardTitle,
		Text:  cardContent,
	}

	if smallImageURL != nil || largeImageURL != nil {
		response.Card.Image = &alexa.Image{}
		if smallImageURL != nil {
			response.Card.Image.SmallImageURL = *smallImageURL
		}
		if largeImageURL != nil {
			response.Card.Image.LargeImageURL = *largeImageURL
		}
	}

	return envelope
}

// WithLinkAccountCard - renders a link account card
func (envelope *ResponseEnvelope) WithLinkAccountCard() *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &alexa.Card{
		Type: "LinkAccount",
	}

	return envelope
}

// WithAskForPermissionsConsentCard - renders an askForPermissionsConsent card
func (envelope *ResponseEnvelope) WithAskForPermissionsConsentCard(permissions []string) *ResponseEnvelope {
	response := envelope.getResponse()

	response.Card = &alexa.Card{
		Type:        "AskForPermissionsConsent",
		Permissions: permissions,
	}

	return envelope
}

// AddDelegateDirective -
func (envelope *ResponseEnvelope) AddDelegateDirective(updatedIntent *alexa.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&alexa.DialogDirective{
		Type:          "Dialog.Delegate",
		UpdatedIntent: updatedIntent,
	})
}

// AddElicitSlotDirective -
func (envelope *ResponseEnvelope) AddElicitSlotDirective(slotToElicit string, updatedIntent *alexa.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&alexa.DialogDirective{
		Type:          "Dialog.ElicitSlot",
		SlotToElicit:  slotToElicit,
		UpdatedIntent: updatedIntent,
	})
}

// AddConfirmSlotDirective -
func (envelope *ResponseEnvelope) AddConfirmSlotDirective(slotToConfirm string, updatedIntent *alexa.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&alexa.DialogDirective{
		Type:          "Dialog.ConfirmSlot",
		SlotToConfirm: slotToConfirm,
		UpdatedIntent: updatedIntent,
	})
}

// AddConfirmIntentDirective -
func (envelope *ResponseEnvelope) AddConfirmIntentDirective(updatedIntent *alexa.Intent) *ResponseEnvelope {
	return envelope.AddDirective(&alexa.DialogDirective{
		Type:          "Dialog.ConfirmIntent",
		UpdatedIntent: updatedIntent,
	})
}

// AddAudioPlayerPlayDirective -
func (envelope *ResponseEnvelope) AddAudioPlayerPlayDirective(
	playBehavior string,
	url string,
	token string,
	offsetInMilliseconds int,
	expectedPreviousToken *string,
	audioItemMetadata *alexa.AudioItemMetadata) *ResponseEnvelope {

	stream := alexa.Stream{
		Token:                token,
		URL:                  url,
		OffsetInMilliseconds: offsetInMilliseconds,
	}

	if expectedPreviousToken != nil {
		stream.ExpectedPreviousToken = *expectedPreviousToken
	}

	audioItem := &alexa.AudioItem{
		Stream:            stream,
		AudioItemMetadata: audioItemMetadata,
	}

	return envelope.AddDirective(&alexa.AudioPlayerDirective{
		Type:         "AudioPlayer.Play",
		PlayBehavior: playBehavior,
		AudioItem:    audioItem,
	})
}

// AddAudioPlayerStopDirective -
func (envelope *ResponseEnvelope) AddAudioPlayerStopDirective() *ResponseEnvelope {
	return envelope.AddDirective(&alexa.AudioPlayerDirective{
		Type: "AudioPlayer.Stop",
	})
}

// AddAudioPlayerClearQueueDirective -
func (envelope *ResponseEnvelope) AddAudioPlayerClearQueueDirective(clearBehavior string) *ResponseEnvelope {
	return envelope.AddDirective(&alexa.AudioPlayerDirective{
		Type:          "AudioPlayer.ClearQueue",
		ClearBehavior: clearBehavior,
	})
}

// AddRenderTemplateDirective -
func (envelope *ResponseEnvelope) AddRenderTemplateDirective(template string) *ResponseEnvelope {
	return envelope.AddDirective(&alexa.RenderTemplateDirective{
		Type:     "Display.RenderTemplate",
		Template: template,
	})
}

// AddHintDirective -
func (envelope *ResponseEnvelope) AddHintDirective(text string) *ResponseEnvelope {
	return envelope.AddDirective(&alexa.HintDirective{
		Type: "Hint",
		Hint: alexa.PlainTextHint{
			Type: "PlainText",
			Text: text,
		},
	})
}

// AddVideoAppLaunchDirective -
func (envelope *ResponseEnvelope) AddVideoAppLaunchDirective(source string, title, subtitle *string) *ResponseEnvelope {
	videoItem := alexa.VideoItem{
		Source: source,
	}

	if title != nil || subtitle != nil {
		videoItem.Metadata = &alexa.VideoItemMetadata{}
		if title != nil {
			videoItem.Metadata.Title = *title
		}
		if subtitle != nil {
			videoItem.Metadata.Subtitle = *subtitle
		}
	}

	envelope.Response.ShouldSessionEnd = false

	return envelope.AddDirective(&alexa.LaunchDirective{
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
		if launch, ok := d.(alexa.LaunchDirective); ok {
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
