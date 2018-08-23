package askgo

import "fmt"

// ResponseBuilder that matches
type ResponseBuilder interface {
	Speak(speechOutput string) *ResponseEnvelope
	Reprompt(speechOutput string) *ResponseEnvelope
	// WithSimpleCard(cardTitle string, cardContent string) ResponseInterface
	// WithStandardCard(cardTitle string, cardContent string, smallImageUrl *string, largeImageUrl *string) ResponseInterface
	// WithLinkAccountCard() ResponseInterface
	// WithAskForPermissionsConsentCard(permissionArray []string) ResponseInterface
	// AddDelegateDirective(updatedIntent *Intent) ResponseInterface;
	// AddElicitSlotDirective(slotToElicit string, updatedIntent *Intent) ResponseInterface;
	// AddConfirmSlotDirective(slotToConfirm string, updatedIntent *Intent) ResponseInterface;
	// AddConfirmIntentDirective(updatedIntent *Intent) ResponseInterface;
	// AddAudioPlayerPlayDirective(playBehavior interfaces.audioplayer.PlayBehavior, url string, token string, offsetInMilliseconds number, expectedPreviousToken *string, audioItemMetadata *AudioItemMetadata) ResponseInterface;
	// AddAudioPlayerStopDirective() ResponseInterface
	// AddAudioPlayerClearQueueDirective(clearBehavior interfaces.audioplayer.ClearBehavior) ResponseInterface;
	// AddRenderTemplateDirective(template interfaces.display.Template) ResponseInterface;
	// AddHintDirective(text string) ResponseInterface;
	// AddVideoAppLaunchDirective(source string, title *string, subtitle *string) ResponseInterface;
	WithShouldEndSession(val bool) *ResponseEnvelope
	// AddDirective(directive Directive) ResponseInterface
	GetResponse() *ResponseEnvelope
}

func (envelope *ResponseEnvelope) getResponse() *Response {
	if envelope.Response == nil {
		envelope.Response = &Response{}
	}
	return envelope.Response
}

// Speak the output in plaintext
func (envelope *ResponseEnvelope) Speak(speechOutput string) *ResponseEnvelope {
	response := envelope.getResponse()
	response.OutputSpeech = &OutputSpeech{
		Type: "SSML",
		SSML: fmt.Sprintf("<speak>%s</speak>", speechOutput),
	}

	return envelope
}

// Reprompt the output in plaintext
func (envelope *ResponseEnvelope) Reprompt(speechOutput string) *ResponseEnvelope {
	response := envelope.getResponse()
	if response.Reprompt == nil {
		response.Reprompt = &Reprompt{}
	}
	response.Reprompt.OutputSpeech = &OutputSpeech{
		Type: "SSML",
		SSML: fmt.Sprintf("<speak>%s</speak>", speechOutput),
	}

	return envelope
}

// WithShouldEndSession set the session end flag
func (envelope *ResponseEnvelope) WithShouldEndSession(val bool) *ResponseEnvelope {
	envelope.getResponse().ShouldSessionEnd = val

	return envelope
}
