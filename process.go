package askgo

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
