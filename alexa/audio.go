package alexa

// AudioPlayerPlayDirective sends Alexa a command to stream the audio file identified by the specified audioItem.
// Use the playBehavior parameter to determine whether the stream begins playing immediately, or is added to the queue.
// shouldEndSession should be set to false otherwise playback will pause immediately
type AudioPlayerPlayDirective struct {
	Type         string    `json:"type"`
	PlayBehavior string    `json:"playBehavior"`
	AudioItem    AudioItem `json:"audioItem"`
}

// AudioItem described the stream to be played
type AudioItem struct {
	Stream   AudioStream        `json:"stream"`
	Metadata *AudioItemMetadata `json:"metadata,omitempty"`
}

// AudioItemMetadata described the additional attributes of the playable stream
type AudioStream struct {
	URL                   string `json:"url"`
	Token                 string `json:"token"`
	ExpectedPreviousToken string `json:"expectedPreviousToken,omitempty"`
	OffsetInMilliseconds  int    `json:"offsetInMilliseconds"`
}

// AudioItemMetadata contains an object providing metadata about the audio to be displayed on the Echo Show and Echo Spot.
type AudioItemMetadata struct {
	Title           string              `json:"title,omitempty"`
	Subtitle        string              `json:"subtitle,omitempty"`
	Art             *DisplayImageObject `json:"art,omitempty"`
	BackgroundImage *DisplayImageObject `json:"backgroundImage,omitempty"`
}

// AudioPlayerStopDirective stopts the current audio playback
type AudioPlayerStopDirective struct {
	Type string `json:"type"`
}

// AudioPlayerClearQueueDirective clears the audio playback queue. You can set this directive to clear the queue without
// stopping the currently playing stream, or clear the queue and stop any currently playing stream.
type AudioPlayerClearQueueDirective struct {
	Type          string `json:"type"`
	ClearBehavior string `json:"clearBehavior"`
}
