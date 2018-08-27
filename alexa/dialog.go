package alexa

// DisplayRenderTemplateDirective directive to render display text, images or items on an device with screen.
type dialogCommonDirective struct {
	Type          string  `json:"type"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
}

// DialogDelegateDirective sends Alexa a command to handle the next turn in the dialog with the user.
// You can use this directive if the skill has a dialog model and the current status of the dialog (dialogState) is either STARTED or IN_PROGRESS.
// You cannot return this directive if the dialogState is COMPLETED.
type DialogDelegateDirective struct {
	Type          string  `json:"type"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
}

// DialogElicitDirective sends Alexa a command to ask the user for the value of a specific slot.
// Specify the name of the slot to elicit in the slotToElicit property.
// Provide a prompt to ask the user for the slot value in an OutputSpeech object in the response.
type DialogElicitDirective struct {
	Type          string  `json:"type"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
	SlotToElicit  string  `json:"slotToElicit"`
}

// DialogConfirmSlotDirective sends Alexa a command to confirm the value of a specific slot before continuing with the dialog.
// Specify the name of the slot to confirm in the slotToConfirm property.
// Provide a prompt to ask the user for confirmation in an OutputSpeech object in the response.
// Be sure repeat back the value to confirm in the prompt.
type DialogConfirmSlotDirective struct {
	Type          string  `json:"type"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
	SlotToConfirm string  `json:"slotToConfirm"`
}

// DialogConfirmIntentDirective sends Alexa a command to confirm the all the information the user has provided for the intent before the skill takes action.
// Provide a prompt to ask the user for confirmation in an OutputSpeech object in the response.
// Be sure to repeat back all the values the user needs to confirm in the prompt.
type DialogConfirmIntentDirective struct {
	Type          string  `json:"type"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
}
