package filter

import (
	"net/http"
	"regexp"
)

type rule struct {
	path                          *regexp.Regexp
	contentType                   *regexp.Regexp
	pathAndContentTypeCombination pathAndContentTypeCombination
	searchPattern                 *regexp.Regexp
	replacement                   []byte
}

type pathAndContentTypeCombination string

const (
	pathAndContentTypeAndCombination = pathAndContentTypeCombination("and")
	pathAndContentTypeOrCombination  = pathAndContentTypeCombination("or")
)

var possiblePathAndContentTypeCombination = []pathAndContentTypeCombination{
	pathAndContentTypeAndCombination,
	pathAndContentTypeOrCombination,
}

func (instance rule) evaluatePathAndContentTypeResult(pathMatch bool, contentTypeMatch bool) bool {
	combination := instance.pathAndContentTypeCombination
	if combination == pathAndContentTypeCombination("") {
		combination = pathAndContentTypeAndCombination
	}
	if combination == pathAndContentTypeAndCombination {
		return pathMatch && contentTypeMatch
	}
	if combination == pathAndContentTypeOrCombination {
		return pathMatch || contentTypeMatch
	}
	return false
}

func (instance *rule) matches(request *http.Request, responseHeader *http.Header) bool {
	var pathMatch, contentTypeMatch bool

	if instance.path != nil {
		pathMatch = request != nil && instance.path.MatchString(request.URL.Path)
	} else {
		pathMatch = true
	}

	if instance.contentType != nil {
		contentTypeMatch = responseHeader != nil && instance.contentType.MatchString(responseHeader.Get("Content-Type"))
	} else {
		contentTypeMatch = true
	}

	return instance.evaluatePathAndContentTypeResult(pathMatch, contentTypeMatch)
}

func (instance *rule) execute(request *http.Request, responseHeader *http.Header, input []byte) []byte {
	pattern := instance.searchPattern
	if pattern == nil {
		return input
	}
	action := &ruleReplaceAction{
		request:        request,
		responseHeader: responseHeader,
		searchPattern:  instance.searchPattern,
		replacement:    instance.replacement,
	}
	output := pattern.ReplaceAllFunc(input, action.replacer)
	return output
}
