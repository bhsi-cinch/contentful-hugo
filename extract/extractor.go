package extract

import (
	"github.com/bhsi-cinch/contentful-hugo/mapper"
	"github.com/bhsi-cinch/contentful-hugo/read"
	"github.com/bhsi-cinch/contentful-hugo/translate"
	"github.com/bhsi-cinch/contentful-hugo/write"
)

// Extractor enables the automated tests to replace key functionalities
// with fakes, mocks and stubs by parameterizing the Reader Configuration,
// the HTTP Getter and the File Store.
type Extractor struct {
	ReadConfig  read.ReadConfig
	Getter      read.Getter
	RStore      read.Store
	TransConfig translate.TransConfig
	WStore      write.Store
	stats       *Stats
}

type Stats struct {
	HTTP       []string
	Types      int
	IndexFiles []string
	Entries    []string
}

// ProcessAll goes through all stages: Read, Map, Translate and Write.
// Underwater, it uses private function processItems to allow reading
// through multiple pages of items being returned from Contentful.
func (e *Extractor) ProcessAll() (Stats, error) {

	e.stats = new(Stats)
	e.stats.HTTP = make([]string, 0)
	e.stats.IndexFiles = make([]string, 0)
	e.stats.Entries = make([]string, 0)

	cf := read.Contentful{
		Getter:     e.Getter,
		ReadConfig: e.ReadConfig,
	}
	types, err := cf.Types()
	if err != nil {
		return *e.stats, err
	}

	typeResult, err := mapper.MapTypes(types)
	if err != nil {
		return *e.stats, err
	}

	e.stats.Types = len(typeResult.Items)

	writer := write.Writer{
		Store: e.WStore,
		Files: make([]string, 0),
	}
	for _, t := range typeResult.Items {
		fileName, content := translate.EstablishDirLevelConf(t, e.TransConfig)
		if fileName != "" && content != "" {
			err = writer.SaveToFile(fileName, content)
			if err != nil {
				return *e.stats, err
			}
		}
	}
	e.stats.IndexFiles = writer.Files
	skip := 0

	err = e.processItems(cf, typeResult, skip)
	if err != nil {
		return *e.stats, err
	}
	e.stats.HTTP = cf.Getter.Stats()
	return *e.stats, nil
}

// processItems is a recursive function which goes through all pages
// returned by Contentful and creates a markdownfile for each.
func (e *Extractor) processItems(cf read.Contentful, typeResult mapper.TypeResult, skip int) error {
	items, err := cf.Items(skip)
	if err != nil {
		return err
	}

	itemResult, err := mapper.MapItems(items)
	if err != nil {
		return err
	}

	archetypeDataMap := make(map[string]map[string]interface{})
	reader := read.Reader{Store: e.RStore}
	writer := write.Writer{
		Store: e.WStore,
		Files: make([]string, 0),
	}
	tc := translate.TranslationContext{Result: itemResult, TransConfig: e.TransConfig}
	for _, item := range itemResult.Items {
		contentType := item.ContentType()
		itemType, err := typeResult.GetType(contentType)
		if err != nil {
			return err
		}

		if archetypeDataMap[contentType] == nil {
			result, err := reader.ViewFromFile(translate.GetArchetypeFilename(contentType))
			if err == nil {
				archeMap, err := tc.TranslateFromMarkdown(result)
				if err != nil {
					return err
				}

				archetypeDataMap[contentType] = archeMap
			} else {

				archetypeDataMap[contentType] = make(map[string]interface{})
			}
		}

		contentMap := tc.MapContentValuesToTypeNames(item.Fields, itemType.Fields)
		overriddenContentmap := tc.MergeMaps(archetypeDataMap[contentType], contentMap)
		contentMarkdown := tc.TranslateToMarkdown(tc.ConvertToContent(overriddenContentmap))
		fileName := translate.Filename(item)
		err = writer.SaveToFile(fileName, contentMarkdown)
		if err != nil {
			return err
		}
	}

	e.stats.Entries = append(e.stats.Entries, writer.Files...)
	nextPage := itemResult.Skip + itemResult.Limit
	if nextPage < itemResult.Total {
		return e.processItems(cf, typeResult, nextPage)
	}
	return nil
}
