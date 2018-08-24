package askgo

import (
	"errors"
	"log"
	"math"
	"strconv"
	"time"
)

var timestampTolerance = 150

// HandlerInput is the standard type for input
type HandlerInput interface {
	// GetRequestEnvelope get the full Alexa Request Envelope
	GetRequestEnvelope() RequestEnvelope
	// GetRequest is a shortcut to GetRequestEnvelope().Request
	GetRequest() *Request
	// Get the response structure
	GetResponse() *ResponseEnvelope
	/*
		Context? : any;
		AttributesManager : AttributesManager;
		ServiceClientFactory? : ServiceClientFactory;
	*/
}

// RequestHandler interface
type RequestHandler interface {
	CanHandle(input HandlerInput) bool
	Handle(input HandlerInput) (*ResponseEnvelope, error)
}

// ErrorHandler interface
type ErrorHandler interface {
	CanHandle(input HandlerInput, e error) bool
	Handle(input HandlerInput, e error) (*ResponseEnvelope, error)
}

// ProcessRequest Main entry point for request processing
func (skill *Skill) ProcessRequest(input HandlerInput) (interface{}, error) {
	envelope := input.GetRequestEnvelope()

	if skill.ApplicationID != "" {
		if err := skill.verifyApplicationID(envelope); err != nil {
			return nil, err
		}
	} else {
		log.Println("Ignoring application verification.")
	}
	if !skill.IgnoreTimestamp {
		if err := skill.verifyTimestamp(envelope); err != nil {
			return nil, err
		}
	} else {
		log.Println("Ignoring timestamp verification.")
	}

	for _, interceptor := range skill.RequestInterceptors {
		if err := interceptor.Process(input); err != nil {
			return skill.dispatchError(input, err)
		}
	}

	var response *ResponseEnvelope

	for _, handler := range skill.Handlers {
		if handler.CanHandle(input) {
			var err error
			response, err = handler.Handle(input)
			if err != nil {
				return skill.dispatchError(input, err)
			}
			break
		}
	}

	for _, interceptor := range skill.ResponseInterceptors {
		if err := interceptor.Process(input, response); err != nil {
			return skill.dispatchError(input, err)
		}
	}

	return response, nil
}

func (skill *Skill) dispatchError(input HandlerInput, err error) (interface{}, error) {
	for _, handler := range skill.ErrorHandlers {
		if handler.CanHandle(input, err) {
			return handler.Handle(input, err)
		}
	}

	return nil, err
}

// verifyApplicationId verifies that the ApplicationID sent in the request
// matches the one configured for this skill.
func (skill *Skill) verifyApplicationID(envelope RequestEnvelope) error {
	if appID := skill.ApplicationID; appID != "" {
		requestAppID := envelope.Session.Application.ApplicationID
		if requestAppID == "" {
			return errors.New("request Application ID was set to an empty string")
		}
		if appID != requestAppID {
			return errors.New("request Application ID does not match expected ApplicationId")
		}
	}

	return nil
}

// verifyTimestamp compares the request timestamp to the current timestamp
// and returns an error if they are too far apart.
func (skill *Skill) verifyTimestamp(envelope RequestEnvelope) error {
	timestamp, err := time.Parse(time.RFC3339, envelope.Request.Timestamp)
	if err != nil {
		return errors.New("Unable to parse request timestamp.  Err: " + err.Error())
	}

	now := time.Now()
	delta := now.Sub(timestamp)
	deltaSecsAbs := math.Abs(delta.Seconds())
	if deltaSecsAbs > float64(timestampTolerance) {
		return errors.New("Invalid Timestamp. The request timestap " + timestamp.String() + " was off the current time " + now.String() + " by more than " + strconv.FormatInt(int64(timestampTolerance), 10) + " seconds.")
	}

	return nil
}

// DefaultHandler for request processing
type DefaultHandler struct {
	Envelope *RequestEnvelope
	Response *ResponseEnvelope
}

// GetRequestEnvelope -- get the full envelope from the request
func (handler *DefaultHandler) GetRequestEnvelope() RequestEnvelope {
	return *handler.Envelope
}

// GetRequest -- quickly get to the request structure
func (handler *DefaultHandler) GetRequest() *Request {
	return handler.Envelope.Request
}

// GetResponse -- Get the response structure
func (handler *DefaultHandler) GetResponse() *ResponseEnvelope {
	if handler.Response == nil {
		handler.Response = &ResponseEnvelope{Version: "1.0"}
	}
	return handler.Response
}
