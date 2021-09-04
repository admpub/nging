package dnsdomain

var DelimLeft = `#{`
var DelimRight = `}`

func Tag(tag string) string {
	return DelimLeft + tag + DelimRight
}
