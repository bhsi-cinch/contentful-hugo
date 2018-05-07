package translate

import (
	"fmt"
	"regexp"

	"github.com/bhsi-cinch/contentful-hugo/mapper"
)

type TranslationContext struct {
	Result      mapper.ItemResult
	TransConfig TransConfig
}

type Content struct {
	Params      map[string]interface{}
	MainContent string
	Slug        string
}

// UnionValuesAndOverride takes a defaults and an overrides map and assigns any missing
// values from the defaults to the overrides map.
func (tc *TranslationContext) UnionValuesAndOverride(itemDefault map[string]interface{}, itemOverride map[string]interface{}) (combinedItem map[string]interface{}) {
	for k, v := range itemDefault {
		if itemOverride[k] == nil {
			itemOverride[k] = v
		} else if s, ok := itemOverride[k].(string); ok && s == "" {
			itemOverride[k] = v
		} else if s, ok := itemOverride[k].([]string); ok && len(s) == 0 {
			itemOverride[k] = v
		} else if s, ok := itemOverride[k].([]interface{}); ok && len(s) == 0 {
			itemOverride[k] = v
		}
	}

	return itemOverride
}

// MapContentValuesToTypeNames takes the values map and the typefield map from contentful and merges the two.
func (tc *TranslationContext) MapContentValuesToTypeNames(Map map[string]interface{}, fields []mapper.TypeField) map[string]interface{} {
	fieldMap := map[string]interface{}{}
	for _, field := range fields {
		value := tc.translateField(Map[field.ID], field)
		if value != nil {
			fieldMap[field.ID] = value
		}
	}

	return fieldMap
}

func removeItem(Map map[string]interface{}, toDelete string) interface{} {
	value := Map[toDelete]
	if value == nil {
		return ""
	}
	delete(Map, toDelete)

	return value
}

// ConvertToContent takes a map of values and converts it to a Content struct
func (tc *TranslationContext) ConvertToContent(fieldMap map[string]interface{}) Content {
	mainContent := removeItem(fieldMap, "mainContent").(string)
	slug, _ := fieldMap["slug"].(string)

	return Content{
		fieldMap,
		mainContent,
		slug,
	}
}

// TranslateFromMarkdown takes a markdown file's contents and converts it to a map.
func (tc *TranslationContext) TranslateFromMarkdown(content string) (rawContent map[string]interface{}, err error) {
	switch tc.TransConfig.Encoding {
	case "yaml":
		return FromYaml(content)
	case "toml":
		return FromToml(content)
	default:
		return FromToml(content)
	}
}

// TranslateToMarkdown accepts a Content struct and converts it to markdown file contents.
func (tc *TranslationContext) TranslateToMarkdown(rawContent Content) (content string) {
	switch tc.TransConfig.Encoding {
	case "yaml":
		return rawContent.ToYaml()
	case "toml":
		return rawContent.ToToml()
	default:
		return rawContent.ToToml()
	}
}

func (tc *TranslationContext) translateArrayField(value interface{}) interface{} {
	if value == nil {
		return []interface{}{}
	}
	items := value.([]interface{})

	var array []string
	array = make([]string, len(items))

	for i, el := range items {
		s, isString := el.(string)
		if isString {
			array[i] = s
		} else {
			if s, ok := tc.translateLinkField(el).(string); ok {
				array[i] = s
			}
		}
	}

	return array
}

func (tc *TranslationContext) translateLinkField(value interface{}) interface{} {
	if value == nil {
		return value
	}
	item := value.(map[string]interface{})
	sys := item["sys"].(map[string]interface{})

	linkType := sys["linkType"]
	if linkType == "Entry" {
		return sys["id"].(string) + ".md"
	} else {
		assets := tc.Result.Includes["Asset"]
		for _, asset := range assets {
			if sys["id"].(string) == asset.Sys.ID {
				return asset.Fields
			}
		}
		// Look up asset - but from where???
	}

	return "ERR"
}

func (tc *TranslationContext) translateDateField(value interface{}) interface{} {
	re, err := regexp.Compile(`([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})(\+[0-9]{2}:[0-9]{2})?`) // want to know what is in front of 'at'
	if err != nil {
		fmt.Println(err)
	}

	res := re.FindAllStringSubmatch(value.(string), -1)
	if len(res) > 0 {
		value = fmt.Sprintf("%v:00%v", res[0][1], res[0][2])
	}

	return value
}

func (tc *TranslationContext) translateField(value interface{}, field mapper.TypeField) interface{} {
	if field.Type == "Array" {
		return tc.translateArrayField(value)

	} else if field.Type == "Link" {
		return tc.translateLinkField(value)

	} else if field.Type == "Date" {
		return tc.translateDateField(value)
	}

	return value
}
