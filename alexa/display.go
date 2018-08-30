package alexa

// DisplayRenderTemplateDirective directive to render display text, images or items on an device with screen.
type DisplayRenderTemplateDirective struct {
	Type string `json:"type,omitempty"`
	// Template is the body template to render
	Template DisplayTemplate `json:"template"`
}

// DisplayTemplate displays text and images. Types may either be BodyTemplate* or ListTemplate*.
// For a body template these images cannot be made selectable.
// List template displays a scrollable list of items, each with associated text and optional images.
// These images can be made selectable, as described in this reference.
type DisplayTemplate struct {
	Type  string `json:"type"`
	Token string `json:"token"`
	// BackButton state (e.g. 'VISIBLE' or 'HIDDEN')
	BackButton      string             `json:"backButton,omitempty"`
	BackgroundImage DisplayImageObject `json:"backgroundImage,omitempty"`
	Title           string             `json:"title,omitempty"`
	TextContent     *TextContent       `json:"textContent,omitempty"`
	// ListItems contains the text and images of the list items.
	ListItems []DisplayListItem `json:"listItems,omitempty"`
}

type TextContent struct {
	PrimaryText   DisplayTextContent  `json:"primaryText,omitempty"`
	SecondaryText *DisplayTextContent `json:"secondaryText,omitempty"`
	TertiaryText  *DisplayTextContent `json:"tertiaryText,omitempty"`
}

type DisplayListItem struct {
	Token       string      `json:"token"`
	TextContent TextContent `json:"textContent,omitempty"`
}

// DisplayTextContent contains text and a text type for displaying text with the Display interface.
type DisplayTextContent struct {
	//Type must be PlainText or RichtText
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

// DisplayImageObject references and describes the image. Multiple sources for the image can be provided.
type DisplayImageObject struct {
	ContentDescription string                `json:"contentDescription,omitempty"`
	Sources            []*DisplayImageSource `json:"sources"`
}

// DisplayImageSource describes the source url and size for a image.
type DisplayImageSource struct {
	URL          string `json:"url"`
	Size         string `json:"size,omitempty"`
	WidthPixels  int    `json:"widthPixels"`
	HeightPixels int    `json:"heightPixels"`
}

// AddImageSource adds source information for a image with the given size.
func (i *DisplayImageObject) AddImageSource(size, url string, heightPixels, widthPixels int) *DisplayImageSource {
	if i.Sources == nil {
		i.Sources = make([]*DisplayImageSource, 0)
	}
	displayImageSource := &DisplayImageSource{
		Size:         size,
		URL:          url,
		HeightPixels: heightPixels,
		WidthPixels:  widthPixels,
	}
	i.Sources = append(i.Sources, displayImageSource)
	return displayImageSource
}
