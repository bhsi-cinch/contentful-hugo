package translate

import (
	"strings"

	"github.com/naoina/toml"
)

const _TOMLdelimiter string = "+++"

// ToToml .() -> string
// Takes a Content struct and outputs it as TOML frontmatter followed by main-content.
func (s Content) ToToml() string {
	result := WriteTomlFrontmatter(s.Params)
	result += s.MainContent

	return result
}

// FromToml reads in a *.toml file and returns all mappings.
func FromToml(s string) (c map[string]interface{}, err error) {
	frontmatter := []byte(strings.TrimRight(strings.SplitAfter(s, _TOMLdelimiter)[1], _TOMLdelimiter))

	c = map[string]interface{}{}
	err = toml.Unmarshal(frontmatter, &c)

	return
}

// WriteTomlFrontmatter (fm Map[]) -> string
// Converts a Map[] into a TOML string, pre and postfixing it with `+++` to designate frontmatter.
func WriteTomlFrontmatter(fm interface{}) string {
	result := _TOMLdelimiter + "\n"
	output, err := toml.Marshal(fm)
	if err != nil {
		return "ERR"
	}

	result += string(output)
	result += _TOMLdelimiter + "\n"

	return result
}
