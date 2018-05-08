package translate

import (
	"strings"

	"gopkg.in/yaml.v2"
)

const _YAMLdelimiter string = "---"

// ToYaml .() -> string
// Takes a Content struct and outputs it as YAML frontmatter followed by main-content.
func (s Content) ToYaml() string {
	result := WriteYamlFrontmatter(s.Params)
	result += s.MainContent

	return result
}

// FromYaml reads in a *.yaml file and returns all mappings.
func FromYaml(s string) (c map[string]interface{}, err error) {
	frontmatter := []byte(strings.TrimRight(strings.SplitAfter(s, _YAMLdelimiter)[1], _YAMLdelimiter))

	c = map[string]interface{}{}
	err = yaml.Unmarshal(frontmatter, &c)

	return
}

// WriteYamlFrontmatter (fm Map[]) -> string
// Converts a Map[] into a YAML string, pre and postfixing it with `---` to designate frontmatter.
func WriteYamlFrontmatter(fm interface{}) string {
	result := _YAMLdelimiter + "\n"
	output, err := yaml.Marshal(fm)
	if err != nil {
		return "ERR"
	}

	result += string(output)
	result += _YAMLdelimiter + "\n"
	return result
}
