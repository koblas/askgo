package askgo

import (
	"context"
	"errors"
	"log"
	"math"
	"strconv"
	"time"
)

var timestampTolerance = 150

// Skill Alexa defines the primary interface to use to create an Alexa request handler.
type Skill struct {
	// ApplicationID must match the ApplicationID defined in the Alexa Skills,
	// if it is the empty string it is ignored.
	ApplicationID string
	// IgnoreTimestamp should be used during debugging to test with hard-coded requests
	IgnoreTimestamp bool

	// Request interceptors are invoked immediately prior to execution of the request handler
	// for an incoming request. Request attributes provide a way for request interceptors to
	// pass data and entities on to request handlers.
	RequestInterceptors []RequestInterceptor

	// Request handlers are responsible for handling one or more types of incoming requests.
	Handlers []RequestHandler

	// Response interceptors are invoked immediately after execution of the request handler.
	// Because response interceptors have access to the output generated from execution of the
	// request handler, they are ideal for tasks such as response sanitization and validation.
	ResponseInterceptors []ResponseInterceptor

	// ErrorHandlers are similar to request handlers, but
	// are instead responsible for handling one or more types of errors.
	// They are invoked by the SDK when an error is returned during the
	// course of request processing.
	ErrorHandlers []ErrorHandler
}

// HandlerInput is the standard type for input for request handlers,
// request and response interceptors, and exception handlers are all passed
// a HandlerInput instance when invoked. This class exposes various entities useful in
// request processing
type HandlerInput interface {
	// GetRequestEnvelope get the full Alexa Request Envelope
	GetRequestEnvelope() RequestEnvelope

	// GetRequest is a shortcut to GetRequestEnvelope().Request
	GetRequest() *Request

	// Get the response structure
	GetResponse() *ResponseEnvelope

	// Provides the context object passed in by the host container. For example, for skills
	// running on AWS Lambda, this is the context object for the AWS Lambda function.
	GetContent() context.Context
}

// RequestInterceptor are invoked immediately prior to execution of the request handler
// for an incoming request. Request attributes provide a way for request interceptors to
// pass data and entities on to request handlers.
type RequestInterceptor interface {
	Process(input HandlerInput) error
}

// ResponseInterceptor are called after the main request handler has been triggered
// these can make any further updates or inspections of the response (e.g. logging)
type ResponseInterceptor interface {
	Process(input HandlerInput, response *ResponseEnvelope) error
}

// RequestHandler are responsible for handling one or more types of incoming requests.
type RequestHandler interface {
	// CanHandle, which is called by the SDK to determine if the given handler is capable of
	// processing the incoming request. This method returns **true** if the handler can handle the
	// request, or **false** if not. You have the flexibility to choose the conditions on which to
	// base this determination, including the type or parameters of the incoming request, or
	// skill attributes.
	CanHandle(input HandlerInput) bool

	// Handle, which is called by the SDK when invoking the request handler. This method contains
	// the handler’s request processing logic, and returns an optional Response.
	Handle(input HandlerInput) (*ResponseEnvelope, error)
}

// ErrorHandler handlers are similar to request handlers, but
// are instead responsible for handling one or more types of errors.
// They are invoked by the SDK when an error is returned during the
// course of request processing.
type ErrorHandler interface {
	// CanHandle, which is called by the SDK to determine if the given handler is capable of
	// handling the error. This method returns **true** if the handler can handle the exception,
	// or **false** if not. A catch-all handler can be easily introduced by simply returning **true**
	// in all cases.
	CanHandle(input HandlerInput, e error) bool
	// Handle, which is called by the SDK when invoking the error handler. This
	// method contains all exception handling logic, and returns an output which
	// optionally may contain a Response.
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
	envelope *RequestEnvelope
	response *ResponseEnvelope
	context  context.Context
}

// NewDefaultHandler builds a structure that supports the default HandlerInput methods
func NewDefaultHandler(ctx context.Context, envelope *RequestEnvelope) *DefaultHandler {
	return &DefaultHandler{envelope: envelope, context: ctx}
}

// GetRequestEnvelope -- get the full envelope from the request
func (handler *DefaultHandler) GetRequestEnvelope() RequestEnvelope {
	return *handler.envelope
}

// GetRequest -- quickly get to the request structure
func (handler *DefaultHandler) GetRequest() *Request {
	return handler.envelope.Request
}

// GetResponse -- Get the response structure
func (handler *DefaultHandler) GetResponse() *ResponseEnvelope {
	if handler.response == nil {
		handler.response = &ResponseEnvelope{Version: "1.0"}
	}
	return handler.response
}

// GetContext returns the default context from construction
func (handler *DefaultHandler) GetContext() context.Context {
	return handler.context
}
