// Package humanize provides methods for formatting and parsing values in human readable form.
package humanize

import (
	"fmt"
	"regexp"
	"strings"
)

// Humanizer is the main struct that provides the public methods.
type Humanizer struct {
	provider      LanguageProvider
	timeInputRe   *regexp.Regexp
	metricInputRe *regexp.Regexp
	language      string
}

// New creates a new humanizer for a given language.
func New(language string, defaults ...string) (*Humanizer, error) {
	language = strings.ToLower(language)
	if humanizer := humanizers.Get(language); humanizer != nil {
		return humanizer, nil
	}
	if provider, exists := languages[language]; exists {
		humanizer := &Humanizer{
			provider: provider,
			language: language,
		}
		humanizer.buildTimeInputRe()
		humanizer.buildMetricInputRe()
		humanizers.Set(language, humanizer)
		return humanizer, nil
	}
	if len(defaults) > 0 && len(defaults[0]) > 0 {
		return New(defaults[0])
	}
	if language != DefaultLanguage {
		return New(DefaultLanguage)
	}
	return nil, fmt.Errorf("Language not supported: %s", language)
}

func (humanizer *Humanizer) Language() string {
	return humanizer.language
}
