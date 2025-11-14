package catalog

import (
	"os"
	"time"

	"github.com/martencassel/opsmsg/message"
	"gopkg.in/yaml.v3"
)

type CatalogEntry struct {
	ID       string   `yaml:"id"`
	Severity string   `yaml:"severity"`
	Text     string   `yaml:"text"`
	Help     string   `yaml:"help"`
	Replies  []string `yaml:"replies"`
}

type Catalog map[string]CatalogEntry

func Load(path string) (Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []CatalogEntry
	if err := yaml.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	catalog := make(Catalog)
	for _, e := range entries {
		catalog[e.ID] = e
	}
	return catalog, nil
}

func (c Catalog) New(id string, ctx map[string]string) message.Message {
	e := c[id]
	return message.Message{
		ID:        e.ID,
		Severity:  message.Severity(e.Severity),
		Text:      e.Text,
		Context:   ctx,
		Timestamp: time.Now(),
		Help:      e.Help,
		Replies:   e.Replies,
	}
}

func Merge(catalogs ...Catalog) Catalog {
	merged := make(Catalog)
	for _, catalog := range catalogs {
		for id, entry := range catalog {
			merged[id] = entry
		}
	}
	return merged
}
