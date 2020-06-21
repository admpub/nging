package filter

import (
	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

func init() {
	caddy.RegisterPlugin("filter", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(controller *caddy.Controller) error {
	handler, err := parseConfiguration(controller)
	if err != nil {
		return err
	}

	config := httpserver.GetConfig(controller)
	config.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		handler.next = next
		return handler
	})

	return nil
}

func parseConfiguration(controller *caddy.Controller) (*filterHandler, error) {
	handler := new(filterHandler)
	handler.rules = []*rule{}
	handler.maximumBufferSize = defaultMaxBufferSize

	for controller.Next() {
		err := evalFilterBlock(controller, handler)
		if err != nil {
			return nil, err
		}
	}

	if len(handler.rules) <= 0 {
		return nil, controller.Err("No rule block provided.")
	}
	return handler, nil
}

func evalFilterBlock(controller *caddy.Controller, target *filterHandler) error {
	args := controller.RemainingArgs()
	if len(args) == 0 {
		return evalDefaultFilterBlock(controller, target)
	}
	return evalNamedBlock(controller, args, target)
}

func evalDefaultFilterBlock(controller *caddy.Controller, target *filterHandler) error {
	for controller.NextBlock() {
		args := []string{controller.Val()}
		args = append(args, controller.RemainingArgs()...)

		err := evalNamedBlock(controller, args, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func evalNamedBlock(controller *caddy.Controller, args []string, target *filterHandler) error {
	switch args[0] {
	case "rule":
		return evalRule(controller, args[1:], target)
	case "max_buffer_size":
		return evalMaximumBufferSize(controller, args[1:], target)
	}
	return controller.Errf("Unknown directive: %v", args[0])
}

func evalRule(controller *caddy.Controller, args []string, target *filterHandler) (err error) {
	if len(args) > 0 {
		return controller.Errf("No more arguments for filter block 'rule' supported.")
	}
	targetRule := new(rule)
	targetRule.pathAndContentTypeCombination = pathAndContentTypeAndCombination
	for controller.NextBlock() {
		optionName := controller.Val()
		switch optionName {
		case "path":
			err = evalPath(controller, targetRule)
		case "content_type":
			err = evalContentType(controller, targetRule)
		case "path_content_type_combination":
			err = evalPathAndContentTypeCombination(controller, targetRule)
		case "search_pattern":
			err = evalSearchPattern(controller, targetRule)
		case "replacement":
			err = evalReplacement(controller, targetRule)
		default:
			err = controller.Errf("Unknown option: %v", optionName)
		}
		if err != nil {
			return err
		}
	}
	if targetRule.path == nil && targetRule.contentType == nil {
		return controller.Errf("Neither 'path' nor 'content_type' definition was provided for filter rule block.")
	}
	if targetRule.searchPattern == nil {
		return controller.Errf("No 'search_pattern' definition was provided for filter rule block.")
	}
	target.rules = append(target.rules, targetRule)
	return nil
}

func evalPath(controller *caddy.Controller, target *rule) error {
	return evalRegexpOption(controller, func(value *regexp.Regexp) error {
		target.path = value
		return nil
	})
}

func evalContentType(controller *caddy.Controller, target *rule) error {
	return evalRegexpOption(controller, func(value *regexp.Regexp) error {
		target.contentType = value
		return nil
	})
}

func evalPathAndContentTypeCombination(controller *caddy.Controller, target *rule) error {
	return evalSimpleOption(controller, func(plainValue string) error {
		for _, candidate := range possiblePathAndContentTypeCombination {
			if string(candidate) == plainValue {
				target.pathAndContentTypeCombination = candidate
				return nil
			}
		}
		return controller.Errf("Illegal value for 'path_content_type_combination': %v", plainValue)
	})
}

func evalSearchPattern(controller *caddy.Controller, target *rule) error {
	return evalRegexpOption(controller, func(value *regexp.Regexp) error {
		target.searchPattern = value
		return nil
	})
}

func evalReplacement(controller *caddy.Controller, target *rule) error {
	return evalSimpleOption(controller, func(value string) error {
		target.replacement = []byte(value)
		if len(target.replacement) > 1 && target.replacement[0] == '@' {
			targetFilename := string(target.replacement[1:])
			content, err := ioutil.ReadFile(targetFilename)
			if err != nil {
				if !os.IsNotExist(err) {
					return controller.Errf("Could not read file provided in 'replacement' definition. Got: %v", err)
				}
			} else {
				target.replacement = content
			}
		}
		return nil
	})
}

func evalSimpleOption(controller *caddy.Controller, setter func(string) error) error {
	args := controller.RemainingArgs()
	if len(args) != 1 {
		return controller.ArgErr()
	}
	return setter(args[0])
}

func evalRegexpOption(controller *caddy.Controller, setter func(*regexp.Regexp) error) error {
	return evalSimpleOption(controller, func(plainValue string) error {
		value, err := regexp.Compile(plainValue)
		if err != nil {
			return err
		}
		return setter(value)
	})
}

func evalMaximumBufferSize(controller *caddy.Controller, args []string, target *filterHandler) (err error) {
	if len(args) != 1 {
		return controller.Errf("There are exact one argument for filter directive 'max_buffer_size' expected.")
	}
	value, err := strconv.Atoi(args[0])
	if err != nil {
		return controller.Errf("There is no valid value for filter directive 'max_buffer_size' provided. Got: %v", err)
	}
	target.maximumBufferSize = value
	return nil
}
