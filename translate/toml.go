package translate

import (
	"regexp"

	"github.com/naoina/toml"
)

// ToToml .() -> string
// Takes a Content struct and outputs it as TOML frontmatter followed by main-content.
func (s Content) ToToml() string {
	result := WriteTomlFrontmatter(s.Params)
	result += s.MainContent

	return result
}

// FromToml reads in a *.toml file and returns all mappings.
func FromToml(s string) (c map[string]interface{}, err error) {
	c = map[string]interface{}{}
	//regex := regexp.MustCompile(`(?<=^[+]{3})((?:\s*.+=.+\s*)+)(?=^[+]{3})`)
	regex := regexp.MustCompile(`((?:.+=.+\s*)+)`)
	regRes := regex.Find([]byte(s))

	err = toml.Unmarshal(regRes, &c)

	return
}

// WriteTomlFrontmatter (fm Map[]) -> string
// Converts a Map[] into a TOML string, pre and postfixing it with `+++` to designate frontmatter.
func WriteTomlFrontmatter(fm interface{}) string {
	result := "+++\n"
	output, err := toml.Marshal(fm)
	if err != nil {
		return "ERR"
	}

	result += string(output)
	result += "+++\n"

	return result
}
