package momolog

type TextFormatter string

const (
	LayerTextFormatter TextFormatter = "%s %s : %s"
	HttpTextFormatter  TextFormatter = "%s %s %s://%s%s - %d in %v"
)

func (textFormatter TextFormatter) String() string {
	return string(textFormatter)
}
