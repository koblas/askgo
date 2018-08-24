# askgo

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://en.wikipedia.org/wiki/MIT_License)

Alexa Skill Kit for GoLang (askgo) is a lightweight impelementation of the Alexa Skill Kit patterns that 
is found in both the JavaScript and 
[Java SDKs](https://alexa-skills-kit-sdk-for-java.readthedocs.io/en/latest/).  While keeping in the 
spirit of Go.

> accept interfaces, return structs

## Usage

This explanation assumes familiarity with with AWS Documentation.  Please review 
[Developing an Alexa Skill as a Lambda Function](https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-lambda-function) before proceeding. This SDK addresses some of the steps documented here for you, but you should be familiar with the entire process.

The examples directory provides example usage.

The Alexa struct is the initial interface point with the SDK.  Alexa must be
 initialized first.  The struct is defined as:

```Go
type Skill struct {
    ApplicationID       string
    IgnoreTimestamp     bool

    RequestInterceptors  []RequestInterceptor
    Handlers             []RequestHandler
    ResponseInterceptors []ResponseInterceptor
    ErrorHandlers        []ErrorHandler
}
```

The ApplicationID must match the ApplicationID defined in the Alexa Skills, if it is the empty string it is ignored.

IgnoreTimestamp should be used during debugging to test with hard-coded requests.

Requests from Alexa should be passed into the ```ProcessRequest``` method.  The ```askgo.DefaultHandler``` is a standard wrapper for generating an interface that is compatible with HandleInput.

*Sample code from a lambda main function*
```Go
lambda.Start(func(ctx context.Context, envelope *askgo.RequestEnvelope) (interface{}, error) {
    return skill.ProcessRequest(&askgo.DefaultHandler{Envelope: envelope})
})
```

To be consistent with the AWS skills kit, request handling is broken up into some clear steps.
* Preprocessing -- RequestInterceptor
* Handling -- RequstHandler
* Preprocessing -- ResponseInterceptor
* Errors -- if any of the Pre/Handle/Post processors return an error, this is passed to the Error Handler

There is no magic support for SessionEnd or OnLaunch, please make sure you're handling those events.

```Go
// RequestHandler interface
type RequestHandler interface {
    CanHandle(input HandlerInput) bool
    Handle(input HandlerInput) (*ResponseEnvelope, error)
}
```

```Go
// ErrorHandler interface
type ErrorHandler interface {
    CanHandle(input HandlerInput, e error) bool
    Handle(input HandlerInput, e error) (*ResponseEnvelope, error)
}
```

```Go
// RequestInterceptor interface
type RequestInterceptor interface {
    Process(input HandlerInput) error
}
```

```Go
// ResponseInterceptor interface
type ResponseInterceptor interface {
    Process(input HandlerInput, response *ResponseEnvelope) error
}
```

Response generation, to generate a response it's similar to the Builder model that is found in the Java/JavaScript SDKs however, we're subscribing to the *pass interfaces return structs* model of Go the response object is fully constructed as
methods are called

For example a response can be as simple as this

```Go
return input.GetResponse().WithShouldEndSession(false).Speak("Shall we play a game?"), nil
```

## samples

[Quiz Game](https://github.com/koblas/askgo/tree/master/example)

## Limitations

This version does not support use as a standalone web server as it does not implement
any of the HTTPS validation.  It was developed to be used as an AWS Lambda function
using AWS Labda Go support.
