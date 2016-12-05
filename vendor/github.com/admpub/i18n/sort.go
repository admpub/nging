package i18n

import (
	"bytes"
	"sort"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// i18nSorted is a stuct which satisfies the sort interface for sorting
// strings alphabetically according to locale
type i18nSorter struct {
	toBeSorted             []interface{}
	getComparisonValueFunc func(interface{}) string
	lessThanFunc           func(i, j int) bool
	collator               *collate.Collator
}

// Len satisfied the sort interface
func (s *i18nSorter) Len() int {
	return len(s.toBeSorted)
}

// Swap satisfied the sort interface
func (s *i18nSorter) Swap(i, j int) {
	s.toBeSorted[i], s.toBeSorted[j] = s.toBeSorted[j], s.toBeSorted[i]
}

// Less satisfied the sort interface.  It uses a collator if available to do a
// string comparison.  Otherwise is uses unicode normalization.
func (s *i18nSorter) Less(i, j int) bool {

	iValue := strings.ToLower(s.getComparisonValueFunc(s.toBeSorted[i]))
	jValue := strings.ToLower(s.getComparisonValueFunc(s.toBeSorted[j]))

	// if it's a local sort, use the collator
	if s.collator != nil {
		return s.collator.CompareString(iValue, jValue) == -1
	}

	// for universal sorts, normalize the unicode to sort
	normalizer := norm.NFKD
	iValue = normalizer.String(iValue)
	jValue = normalizer.String(jValue)

	return bytes.Compare([]byte(iValue), []byte(jValue)) == -1
}

// SortUniversal sorts a generic slice alphabetically in such a way that it
// should be mostly correct for most locales. It should be used in the following
// 2 cases:
//     - As a fallback for SortLocale, when a collator for a specific locale
//       cannot be found
//     - When a locale is not available, or a sorting needs to be done in a
//       locale-agnostic way
// It uses unicode normalization.
// The func argument tells this function what string value to do the comparisons
// on.
func SortUniversal(toBeSorted []interface{}, getComparisonValueFunction func(interface{}) string) {
	sorter := &i18nSorter{
		toBeSorted:             toBeSorted,
		getComparisonValueFunc: getComparisonValueFunction,
	}

	sort.Sort(sorter)
}

// SortLocal sorts a generic slice alphabetically for a specific locale. It
// uses collation information if available for the specific locale requested.
// It falls back to SortUniversal otherwise. The func argument tells this
// function what string value to do the comparisons on.
func SortLocal(locale string, toBeSorted []interface{}, getComparisonValueFunction func(interface{}) string) {

	if locale == "" {
		SortUniversal(toBeSorted, getComparisonValueFunction)
		return
	}

	collator := getCollator(locale)

	if collator == nil {
		SortUniversal(toBeSorted, getComparisonValueFunction)
		return
	}

	sorter := &i18nSorter{
		toBeSorted:             toBeSorted,
		getComparisonValueFunc: getComparisonValueFunction,
		collator:               collator,
	}

	sort.Sort(sorter)
}

// getCollator returns a collate package Collator pointer. This can result in a
// panic, so this function must recover from that if it happens.
func getCollator(locale string) *collate.Collator {

	defer func() {
		recover()
	}()

	tag := language.Make(locale)

	if tag == language.Und {
		return nil
	}
	return collate.New(tag)
}

// Sort sorts a generic slice alphabetically for this translator's locale. It
// uses collation information if available for the specific locale requested.
// It falls back to SortUniversal otherwise. The func argument tells this
// function what string value to do the comparisons on.
func (t *Translator) Sort(toBeSorted []interface{}, getComparisonValueFunction func(interface{}) string) {
	SortLocal(t.locale, toBeSorted, getComparisonValueFunction)
}
