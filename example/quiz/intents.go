package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/structs"
	"github.com/koblas/askgo"
)

//  -----------------------
var attributeContext struct{}

type unpackAttributes struct{}

func (h *unpackAttributes) Process(input askgo.HandlerInput) error {
	attributes := getAttributes(input)

	log.Printf("Got Attributes")

	input.SetContext(context.WithValue(input.GetContext(), &attributeContext, attributes))

	return nil
}

//  -----------------------
type saveAttributes struct{}

func (h *saveAttributes) Process(input askgo.HandlerInput, envelope *askgo.ResponseEnvelope) error {
	if !envelope.Response.ShouldSessionEnd {
		attributes, ok := input.GetContext().Value(&attributeContext).(*Attributes)

		if !ok {
			log.Printf("Error: Attributes not correct type")
		} else {
			envelope.SessionAttributes = structs.Map(attributes)
		}
	}

	return nil
}

//  -----------------------
type errorHandler struct{}

func (h *errorHandler) CanHandle(input askgo.HandlerInput) bool {
	return true
}
func (h *errorHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	builder := input.GetResponse().WithShouldEndSession(false)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("ErrorHandler requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	return builder.Speak(helpMessage).Reprompt(helpMessage), nil
}

//  -----------------------
type helpHandler struct{}

func (h *helpHandler) CanHandle(input askgo.HandlerInput) bool {
	request := input.GetRequest()
	return request.Intent.Name == askgo.AMAZON.HelpIntent
}
func (h *helpHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	builder := input.GetResponse().WithShouldEndSession(false)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("HelpHandler requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	return builder.Speak(helpMessage).Reprompt(helpMessage), nil
}

//  -----------------------
type exitHandler struct{}

func (h *exitHandler) CanHandle(input askgo.HandlerInput) bool {
	request := input.GetRequest()
	return request.Intent.Name == askgo.AMAZON.StopIntent ||
		request.Intent.Name == askgo.AMAZON.PauseIntent ||
		request.Intent.Name == askgo.AMAZON.CancelIntent
}
func (h *exitHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	builder := input.GetResponse().WithShouldEndSession(true)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("ExitHandler requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	return builder.Speak(exitSkillMessage), nil
}

//  -----------------------
type sessionEndHandler struct{}

func (h *sessionEndHandler) CanHandle(input askgo.HandlerInput) bool {
	return input.GetRequest().Type == "SessionEndedRequest"
}
func (h *sessionEndHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("SessionEnd requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	return input.GetResponse().WithShouldEndSession(true), nil
}

//  -----------------------
type launchHandler struct{}

func (h *launchHandler) CanHandle(input askgo.HandlerInput) bool {
	return input.GetRequest().Type == "LaunchRequest"
}
func (h *launchHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	response := input.GetResponse().WithShouldEndSession(false)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("LaunchRequest requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	return response.Speak(welcomeMessage).Reprompt(helpMessage), nil
}

//  -----------------------
type repeatHandler struct{}

func (h *repeatHandler) CanHandle(input askgo.HandlerInput) bool {
	request := input.GetRequest()
	return request.Intent.Name == askgo.AMAZON.RepeatIntent
}
func (h *repeatHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	builder := input.GetResponse().WithShouldEndSession(false)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("RepeatHandler requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	question := getQuestion(attributes)

	return builder.Speak(question).Reprompt(question), nil
}

//  -----------------------
type quizHandler struct{}

func (h *quizHandler) CanHandle(input askgo.HandlerInput) bool {
	request := input.GetRequest()

	return request.Intent.Name == "QuizIntent" || request.Intent.Name == askgo.AMAZON.StartOverIntent
}
func (h *quizHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	builder := input.GetResponse().WithShouldEndSession(false)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("QuizHandler requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	attributes.State = QUIZ
	attributes.Counter = 0
	askQuestion(request, attributes)
	question := getQuestion(attributes)

	return builder.Speak(fmt.Sprintf("%s %s", startQuizMessage, question)).Reprompt(question), nil

	/*
		if supportsDisplay(acontext) {
			// response.SetSimpleCard(fmt.Sprintf("Question #$d", session.Attributes.String["counter"]), "")
			// * TODO * more interesting display
		}
	*/
}

//  -----------------------
type definitionHandler struct{}

func (h *definitionHandler) CanHandle(input askgo.HandlerInput) bool {
	request := input.GetRequest()
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	return attributes.State != QUIZ && request.Intent.Name == "AnswerIntent"
}

func (h *definitionHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	response := input.GetResponse().WithShouldEndSession(false)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("DefinitionHandler requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	overlap := make(map[string]int)

	var slotItem string
	for k, v := range request.Intent.Slots {
		if v.Value != "" {
			overlap[k] = 1
			slotItem = k
		}
	}

	s := structs.New(&data[0])
	for _, n := range s.Names() {
		if _, found := overlap[n]; found {
			overlap[n]++
		} else {
			overlap[n] = 1
		}
	}

	keys := make([]string, 0)
	for k, v := range overlap {
		if v == 2 {
			keys = append(keys, k)
		}
	}

	var match *QuizItem

	if len(keys) != 0 {
		key := keys[0]

		if item, ok := request.Intent.Slots[key]; ok {
			for _, entry := range data {
				s := structs.New(entry)
				v := s.Field(key).Value()
				if strings.EqualFold(fmt.Sprintf("%v", v), fmt.Sprintf("%v", item.Value)) {
					match = &entry
					break
				}
			}
		}
	}

	if match != nil {
		msg := getSpeechDescription(*match)

		// @TODO -- msg is <speak>tag...
		response.Speak(msg).Reprompt(msg)
	} else {
		msg := fmt.Sprintf("I'm sorry. %s is not something I know very much about in this skill. %s", formatCasing(slotItem), helpMessage)

		response.Speak(msg).Reprompt(msg)
	}

	return response, nil
}

//  -----------------------
type quizAnswerHandler struct{}

func (h *quizAnswerHandler) CanHandle(input askgo.HandlerInput) bool {
	request := input.GetRequest()
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	return attributes.State == QUIZ && request.Intent.Name == "AnswerIntent"
}
func (h *quizAnswerHandler) Handle(input askgo.HandlerInput) (*askgo.ResponseEnvelope, error) {
	request := input.GetRequest()
	response := input.GetResponse().WithShouldEndSession(false)
	attributes := input.GetContext().Value(&attributeContext).(*Attributes)

	log.Printf("QuizAnswerHandler requestId=%s, sessionId=%s", request.RequestID, attributes.SessionID)

	var isCorrect bool

	if prop, ok := request.Intent.Slots[attributes.QuizProperty]; ok {
		isCorrect = strings.EqualFold(prop.Value, attributes.QuizAnswer)
	}

	var cons string

	if isCorrect {
		attributes.QuizScore++

		cons = speechConsCorrect[random.Intn(len(speechConsCorrect))]
	} else {
		cons = speechConsWrong[random.Intn(len(speechConsWrong))]
	}

	output := []string{fmt.Sprintf("<say-as interpret-as='interjection'>%s</say-as><break strength='strong'/>", cons)}

	if attributes.Counter < 10 {
		askQuestion(request, attributes)
		question := getQuestion(attributes)

		output = append(output, question)
		response.Reprompt(question)
	} else {
		output = append(output, getFinalScore(attributes))
		output = append(output, exitSkillMessage)

		attributes.State = START
	}

	response.Speak(strings.Join(output, " "))

	return response, nil
}
