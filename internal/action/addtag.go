package action

import (
	"bytes"
	"flag"
	"strings"
	"unicode"
)

var _ Action = (*tagsEdit)(nil)

type tagsEdit struct {
	fs     *flag.FlagSet
	tag    *string
	source *string
}

func NewTagsEdit() Action {
	tg := &tagsEdit{fs: flag.NewFlagSet("tags-edit", flag.ContinueOnError)}
	tg.tag = tg.fs.String("t", "", "tags")
	tg.source = tg.fs.String("s", "", "source")
	return tg
}

// Action implements Action.
func (t *tagsEdit) Action() string {
	return "add-tags"
}

// Run implements Action.
func (t *tagsEdit) Run(args []string) *Result {
	if len(args) < 1 {
		return NewFailResult(ErrIncorrectRequest)
	}

	if err := t.fs.Parse(args); err != nil {
		return NewFailResult(err)
	}

	return &Result{
		Code: SuccessCode,
		Data: addTags(*t.source, *t.tag),
	}
}

// Usege implements Action.
func (t *tagsEdit) Usege() string {
	panic("unimplemented")
}

func addTags(str, tag string) string {
	replace, s, ok := getChangeEdit(str, tag)
	if ok {
		return str[0:s] + replace + str[s:]
	}
	return str
}

func getChangeEdit(str string, tag string) (string, int, bool) {
	if len(str) == 0 {
		return "", 0, true
	}

	key := strings.Fields(strings.TrimSpace(str))
	result := ""
	switch tag {
	case "yaml":
		result = tag + ":\"" + toCamelCase(key[0]) + "\""
	case "form-req":
		result = "form:\"" + toUnderline(key[0]) + "\" binding:\"required\""
	case "bind":
		result = "binding:\"required\""
	default:
		result = tag + ":\"" + toUnderline(key[0]) + "\""
	}

	ci := strings.Index(str, "//")
	tStart := strings.Index(str, "`")
	if (ci != -1 && ci < tStart) || tStart == -1 {
		return " `" + result + "`", len(str), true
	}

	tEnd := strings.Index(str[tStart+1:], "`") + tStart + 1
	tags := strings.Fields(str[tStart+1 : tEnd])
	for _, t := range tags {
		if strings.Split(t, ":")[0] == tag {
			return "", 0, false
		}
	}
	return " " + result, tEnd, true
}

func toUnderline(str string) string {
	buff := bytes.Buffer{}
	for index, s := range str {
		if index == 0 {
			buff.WriteRune(unicode.ToLower(s))
			continue
		}

		if unicode.IsUpper(s) {
			if  index == 0 || !unicode.IsUpper(rune(str[index-1])) {
				buff.WriteByte('_')
			}
			buff.WriteRune(unicode.ToLower(s))
			continue
		}

		buff.WriteRune(s)
	}
	return buff.String()
}

func toCamelCase(str string) string {
	return string(unicode.ToLower(rune(str[0]))) + str[1:]
}
