/*
Package askgo is a pure go Alex Skills Kit developement package.

Int most cases clients will quickly build a skill with something like this


	package main

	import (
		"context"

		"github.com/aws/aws-lambda-go/lambda"
		"github.com/koblas/askgo"
	)

	func main() {
		skill := &askgo.Skill{
			ApplicationID: "xyzzy",

			// Our handler interface
			Handlers: []askgo.RequestHandler{
				&launchHandler{},
				&sessionEndHandler{},
				&exitHandler{},
				&skillInvokationHandler{},
			},
		}

		lambda.Start(func(ctx context.Context, envelope *askgo.RequestEnvelope) (interface{}, error) {
			return skill.ProcessRequest(&askgo.DefaultHandler{Envelope: envelope})
		})
	}

In general this follows the Alexa Skills Kit approach to building skill as used in the Java and JavaScript SDKs
as provided by amazon.  The key difference is that this is not using builder patterns, but focused on Go style
iterfaces.

You can find a complete, working example of a Skill at
https://github.com/koblas/askgo/tree/master/example/quiz.


*/
package askgo
