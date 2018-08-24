package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/fatih/structs"

	"github.com/koblas/askgo"

	"github.com/mitchellh/mapstructure"
)

// States for the Quiz to be in
type States int

const (
	// START ing state
	START States = iota
	// QUIZ in progress state
	QUIZ States = iota
)

// Attributes for the current user session
type Attributes struct {
	SessionID string `mapstructure:"-"`
	UserID    string `mapstructure:"-"`

	State         States
	Counter       int
	QuizScore     int
	QuizItemIndex int
	QuizProperty  string
	QuizAnswer    string
}

///
var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	welcomeMessage = `Welcome to the United States Quiz Game!
											You can ask me about any of the fifty states and their capitals, 
											or you can ask me to start a quiz.  What would you like to do?`
	startQuizMessage = `OK.  I will ask you 10 questions about the United States. `
	exitSkillMessage = `Thank you for playing the United States Quiz Game!  Let's play again soon!`
	repromptSpeech   = `Which other state or capital would you like to know about?`
	helpMessage      = `I know lots of things about the United States. 
											 You can ask me about a state or a capital, and I'll tell you what I know.
											 You can also test your knowledge by asking me to start a quiz.  
											 What would you like to do?`

	speechConsCorrect = []string{
		"Booya",
		"All righty",
		"Bam",
		"Bazinga",
		"Bingo",
		"Boom",
		"Bravo",
		"Cha Ching",
		"Cheers",
		"Dynomite",
		"Hip hip hooray",
		"Hurrah",
		"Hurray",
		"Huzzah",
		"Oh dear.  Just kidding.  Hurray",
		"Kaboom",
		"Kaching",
		"Oh snap",
		"Phew",
		"Righto",
		"Way to go",
		"Well done",
		"Whee",
		"Woo hoo",
		"Yay",
		"Wowza",
		"Yowsa",
	}
	speechConsWrong = []string{"Argh",
		"Aw man",
		"Blarg",
		"Blast",
		"Boo",
		"Bummer",
		"Darn",
		"D\"oh",
		"Dun dun dun",
		"Eek",
		"Honk",
		"Le sigh",
		"Mamma mia",
		"Oh boy",
		"Oh dear",
		"Oof",
		"Ouch",
		"Ruh roh",
		"Shucks",
		"Uh oh",
		"Wah wah",
		"Whoops a daisy",
		"Yikes",
	}
)

// QuizItem for the quiz
type QuizItem struct {
	StateName      string
	Abbreviation   string
	Capital        string
	StatehoodYear  int
	StatehoodOrder int
}

var data = []QuizItem{
	QuizItem{StateName: "Alabama", Abbreviation: "AL", Capital: "Montgomery", StatehoodYear: 1819, StatehoodOrder: 22},
	QuizItem{StateName: "Alaska", Abbreviation: "AK", Capital: "Juneau", StatehoodYear: 1959, StatehoodOrder: 49},
	QuizItem{StateName: "Arizona", Abbreviation: "AZ", Capital: "Phoenix", StatehoodYear: 1912, StatehoodOrder: 48},
	QuizItem{StateName: "Arkansas", Abbreviation: "AR", Capital: "Little Rock", StatehoodYear: 1836, StatehoodOrder: 25},
	QuizItem{StateName: "California", Abbreviation: "CA", Capital: "Sacramento", StatehoodYear: 1850, StatehoodOrder: 31},
	QuizItem{StateName: "Colorado", Abbreviation: "CO", Capital: "Denver", StatehoodYear: 1876, StatehoodOrder: 38},
	QuizItem{StateName: "Connecticut", Abbreviation: "CT", Capital: "Hartford", StatehoodYear: 1788, StatehoodOrder: 5},
	QuizItem{StateName: "Delaware", Abbreviation: "DE", Capital: "Dover", StatehoodYear: 1787, StatehoodOrder: 1},
	QuizItem{StateName: "Florida", Abbreviation: "FL", Capital: "Tallahassee", StatehoodYear: 1845, StatehoodOrder: 27},
	QuizItem{StateName: "Georgia", Abbreviation: "GA", Capital: "Atlanta", StatehoodYear: 1788, StatehoodOrder: 4},
	QuizItem{StateName: "Hawaii", Abbreviation: "HI", Capital: "Honolulu", StatehoodYear: 1959, StatehoodOrder: 50},
	QuizItem{StateName: "Idaho", Abbreviation: "ID", Capital: "Boise", StatehoodYear: 1890, StatehoodOrder: 43},
	QuizItem{StateName: "Illinois", Abbreviation: "IL", Capital: "Springfield", StatehoodYear: 1818, StatehoodOrder: 21},
	QuizItem{StateName: "Indiana", Abbreviation: "IN", Capital: "Indianapolis", StatehoodYear: 1816, StatehoodOrder: 19},
	QuizItem{StateName: "Iowa", Abbreviation: "IA", Capital: "Des Moines", StatehoodYear: 1846, StatehoodOrder: 29},
	QuizItem{StateName: "Kansas", Abbreviation: "KS", Capital: "Topeka", StatehoodYear: 1861, StatehoodOrder: 34},
	QuizItem{StateName: "Kentucky", Abbreviation: "KY", Capital: "Frankfort", StatehoodYear: 1792, StatehoodOrder: 15},
	QuizItem{StateName: "Louisiana", Abbreviation: "LA", Capital: "Baton Rouge", StatehoodYear: 1812, StatehoodOrder: 18},
	QuizItem{StateName: "Maine", Abbreviation: "ME", Capital: "Augusta", StatehoodYear: 1820, StatehoodOrder: 23},
	QuizItem{StateName: "Maryland", Abbreviation: "MD", Capital: "Annapolis", StatehoodYear: 1788, StatehoodOrder: 7},
	QuizItem{StateName: "Massachusetts", Abbreviation: "MA", Capital: "Boston", StatehoodYear: 1788, StatehoodOrder: 6},
	QuizItem{StateName: "Michigan", Abbreviation: "MI", Capital: "Lansing", StatehoodYear: 1837, StatehoodOrder: 26},
	QuizItem{StateName: "Minnesota", Abbreviation: "MN", Capital: "St. Paul", StatehoodYear: 1858, StatehoodOrder: 32},
	QuizItem{StateName: "Mississippi", Abbreviation: "MS", Capital: "Jackson", StatehoodYear: 1817, StatehoodOrder: 20},
	QuizItem{StateName: "Missouri", Abbreviation: "MO", Capital: "Jefferson City", StatehoodYear: 1821, StatehoodOrder: 24},
	QuizItem{StateName: "Montana", Abbreviation: "MT", Capital: "Helena", StatehoodYear: 1889, StatehoodOrder: 41},
	QuizItem{StateName: "Nebraska", Abbreviation: "NE", Capital: "Lincoln", StatehoodYear: 1867, StatehoodOrder: 37},
	QuizItem{StateName: "Nevada", Abbreviation: "NV", Capital: "Carson City", StatehoodYear: 1864, StatehoodOrder: 36},
	QuizItem{StateName: "New Hampshire", Abbreviation: "NH", Capital: "Concord", StatehoodYear: 1788, StatehoodOrder: 9},
	QuizItem{StateName: "New Jersey", Abbreviation: "NJ", Capital: "Trenton", StatehoodYear: 1787, StatehoodOrder: 3},
	QuizItem{StateName: "New Mexico", Abbreviation: "NM", Capital: "Santa Fe", StatehoodYear: 1912, StatehoodOrder: 47},
	QuizItem{StateName: "New York", Abbreviation: "NY", Capital: "Albany", StatehoodYear: 1788, StatehoodOrder: 11},
	QuizItem{StateName: "North Carolina", Abbreviation: "NC", Capital: "Raleigh", StatehoodYear: 1789, StatehoodOrder: 12},
	QuizItem{StateName: "North Dakota", Abbreviation: "ND", Capital: "Bismarck", StatehoodYear: 1889, StatehoodOrder: 39},
	QuizItem{StateName: "Ohio", Abbreviation: "OH", Capital: "Columbus", StatehoodYear: 1803, StatehoodOrder: 17},
	QuizItem{StateName: "Oklahoma", Abbreviation: "OK", Capital: "Oklahoma City", StatehoodYear: 1907, StatehoodOrder: 46},
	QuizItem{StateName: "Oregon", Abbreviation: "OR", Capital: "Salem", StatehoodYear: 1859, StatehoodOrder: 33},
	QuizItem{StateName: "Pennsylvania", Abbreviation: "PA", Capital: "Harrisburg", StatehoodYear: 1787, StatehoodOrder: 2},
	QuizItem{StateName: "Rhode Island", Abbreviation: "RI", Capital: "Providence", StatehoodYear: 1790, StatehoodOrder: 13},
	QuizItem{StateName: "South Carolina", Abbreviation: "SC", Capital: "Columbia", StatehoodYear: 1788, StatehoodOrder: 8},
	QuizItem{StateName: "South Dakota", Abbreviation: "SD", Capital: "Pierre", StatehoodYear: 1889, StatehoodOrder: 40},
	QuizItem{StateName: "Tennessee", Abbreviation: "TN", Capital: "Nashville", StatehoodYear: 1796, StatehoodOrder: 16},
	QuizItem{StateName: "Texas", Abbreviation: "TX", Capital: "Austin", StatehoodYear: 1845, StatehoodOrder: 28},
	QuizItem{StateName: "Utah", Abbreviation: "UT", Capital: "Salt Lake City", StatehoodYear: 1896, StatehoodOrder: 45},
	QuizItem{StateName: "Vermont", Abbreviation: "VT", Capital: "Montpelier", StatehoodYear: 1791, StatehoodOrder: 14},
	QuizItem{StateName: "Virginia", Abbreviation: "VA", Capital: "Richmond", StatehoodYear: 1788, StatehoodOrder: 10},
	QuizItem{StateName: "Washington", Abbreviation: "WA", Capital: "Olympia", StatehoodYear: 1889, StatehoodOrder: 42},
	QuizItem{StateName: "West Virginia", Abbreviation: "WV", Capital: "Charleston", StatehoodYear: 1863, StatehoodOrder: 35},
	QuizItem{StateName: "Wisconsin", Abbreviation: "WI", Capital: "Madison", StatehoodYear: 1848, StatehoodOrder: 30},
	QuizItem{StateName: "Wyoming", Abbreviation: "WY", Capital: "Cheyenne", StatehoodYear: 1890, StatehoodOrder: 44},
}

///

func askQuestion(request *askgo.Request, attributes *Attributes) {
	itemIndex := random.Intn(len(data))

	item := data[itemIndex]

	s := structs.New(&item)

	fields := make([]*structs.Field, 0)

	for _, f := range s.Fields() {
		if f.Name() != "StateName" {
			fields = append(fields, f)
		}
	}

	propIndex := rand.Intn(len(fields))

	name := fields[propIndex].Name()
	value := fields[propIndex].Value()

	attributes.Counter++
	attributes.QuizItemIndex = itemIndex
	attributes.QuizProperty = name
	attributes.QuizAnswer = fmt.Sprintf("%v", value)
}

func formatCasing(name string) string {
	re := regexp.MustCompile("([A-Z])")

	return re.ReplaceAllString(name, " $1")
}

func getQuestion(attributes *Attributes) string {
	title := formatCasing(attributes.QuizProperty)
	item := data[attributes.QuizItemIndex]

	return fmt.Sprintf("Here is your %dth question.  What is the %s of %s?", attributes.Counter, title, item.StateName)
}

func getAnswer(attributes *Attributes) string {
	title := formatCasing(attributes.QuizProperty)
	item := data[attributes.QuizItemIndex]

	s := structs.New(&item)

	switch attributes.QuizProperty {
	case "Abbreviation":
		return fmt.Sprintf("The %s of %s is <say-as interpret-as='spell-out'>%s</say-as>.", title, item.StateName, item.Abbreviation)
	default:
		return fmt.Sprintf("The %s of %s is %v.", title, item.StateName, s.Field(attributes.QuizProperty).Value())
	}
}

func getFinalScore(attributes *Attributes) string {
	return fmt.Sprintf("Your final score is %d out of %d.", attributes.QuizScore, attributes.Counter)
}

func getSpeechDescription(item QuizItem) string {
	return fmt.Sprintf(`
			%s is the %vth state, admitted to the Union in %v.  
			The capital of %s is %s, and the abbreviation for %s is 
			<break strength='strong'/><say-as interpret-as='spell-out'>%s</say-as>.  
			I've added %s to your Alexa app.  
			Which other state or capital would you like to know about?
			`,
		item.StateName, item.StatehoodOrder, item.StatehoodYear,
		item.StateName, item.Capital, item.StateName, item.Abbreviation, item.StateName)
}

///
func supportsDisplay(acontext *askgo.Context) bool {
	return acontext.System.Device.SupportedInterfaces.Display.TemplateVersion != ""
}

func getAttributes(input askgo.HandlerInput) *Attributes {
	session := input.GetRequestEnvelope().Session
	attributes := &Attributes{State: START}

	if session.Attributes != nil {
		mapstructure.Decode(session.Attributes, attributes)
		log.Printf("Attributes = %+v", attributes)
	} else {
		log.Printf("Attributes = DEFAULT")
	}

	attributes.SessionID = session.SessionID
	// attributes.UserID = session.UserID

	return attributes
}

func setAttributes(session *askgo.Session, attributes *Attributes) {
	log.Printf("SET Attributes = %+v", attributes)
	session.Attributes = structs.Map(attributes)
}
