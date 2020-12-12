package resttools

import "fmt"

const (
	CharsField   = `-_.9a-zA-Z`
	CharsLiteral = `-_.~0-9a-zA-Z%`
)

var (
	CharsFieldPath                           string
	RegexField, RegexFieldPath, RegexLiteral string
)

func init() {
	RegexField = fmt.Sprintf(`[%s]+`, CharsField)

	CharsFieldPath = CharsField + `.`
	RegexFieldPath = fmt.Sprintf(`^[%s]+`, CharsFieldPath)

	RegexLiteral = fmt.Sprintf(`^[%s]+`, CharsLiteral)
}
