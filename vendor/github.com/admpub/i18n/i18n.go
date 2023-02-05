package i18n

import (
	// standard library
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	// third party

	"github.com/admpub/confl"
)

// TranslatorFactory is a struct which contains the info necessary for creating
// Translator "instances". It also "caches" previously created Translators.
// Because of this caching, you can request a Translator for a specific locale
// multiple times and always get a pointer to the same Translator instance.
type TranslatorFactory struct {
	messagesPaths []string
	rulesPaths    []string
	lock          sync.RWMutex
	translators   map[string]*Translator
	fallback      *Translator
	getFileSystem func(file string) http.FileSystem
	project       string
}

var pathSeparator string

func init() {
	p := path.Join("a", "b")
	pathSeparator = p[1 : len(p)-1]
}

// NewTranslatorFactory returns a TranslatorFactory instance with the specified
// paths and fallback locale.  If a fallback locale is specified, it
// automatically creates the fallback Translator instance. Several errors can
// occur during this process, and those are all returned in the errors slice.
// Even if errors are returned, this function should still return a working
// Translator if the fallback works.
//
// If multiple rulesPaths or messagesPaths are provided, they loaded in the
// order they appear in the slice, with values added later overriding any rules
// or messages loaded earlier.
//
// One lat thing about the messagesPaths. You can organize your locale messages
// files in this messagesPaths directory in 2 different ways.
//
//  1. Place *.yaml files in that directory directly, named after locale codes -
//
//     messages/
//     en.yaml
//     fr.yaml
//
//  2. Place subdirectores in that directory, named after locale codes and
//     containing *.yaml files
//
//     messages/
//     en/
//     front-end.yaml
//     email.yaml
//     fr/
//     front-end.yaml
//     email.yaml
//
//     Using the second way allows you to organize your messages into multiple
//     files.
func NewTranslatorFactory(rulesPaths []string, messagesPaths []string, fallbackLocale string, fs ...func(file string) http.FileSystem) (f *TranslatorFactory, errors []error) {
	return NewTranslatorFactoryWith(``, rulesPaths, messagesPaths, fallbackLocale, fs...)
}

func MakeOpeners(getFileSystem func(string) http.FileSystem) []func(string) http.FileSystem {
	var openers []func(string) http.FileSystem
	if getFileSystem != nil {
		openers = append(openers, getFileSystem)
	}
	openers = append(openers, func(p string) http.FileSystem {
		return http.Dir(p)
	})
	return openers
}

func NewTranslatorFactoryWith(project string, rulesPaths []string, messagesPaths []string, fallbackLocale string, fs ...func(file string) http.FileSystem) (f *TranslatorFactory, errors []error) {
	f = new(TranslatorFactory)
	f.project = project
	if len(fs) > 0 {
		f.getFileSystem = fs[0]
	}

	if len(rulesPaths) == 0 {
		errors = append(errors, translatorError{message: "rules paths empty"})
	}

	if len(messagesPaths) == 0 {
		errors = append(errors, translatorError{message: "messages paths empty"})
	}

	foundRules := fallbackLocale == ""
	foundMessages := fallbackLocale == ""
	openers := MakeOpeners(f.getFileSystem)

	for _, p := range rulesPaths {
		p = strings.TrimRight(p, pathSeparator)
		var errs []string
		for _, opener := range openers {
			fs := opener(p)
			file, err := fs.Open(".")
			if err != nil {
				message := "can't read rules path <" + p + ">: " + err.Error()
				errs = append(errs, message)
				continue
			}
			_, err = file.Stat()
			if err != nil {
				message := "can't read rules path <" + p + ">: " + err.Error()
				errs = append(errs, message)
			}
			file.Close()

			if foundRules {
				continue
			}
			file, err = fs.Open(fallbackLocale + ".yaml")
			if err != nil {
				continue
			}
			if fi, err := file.Stat(); err == nil && !fi.IsDir() {
				foundRules = true
			}
			file.Close()
		}
		if len(errs) == len(openers) {
			errors = append(errors, translatorError{message: errs[0]})
		}
	}

	for _, p := range messagesPaths {
		p = strings.TrimRight(p, pathSeparator)
		var errs []string
		for _, opener := range openers {
			fs := opener(p)
			file, err := fs.Open(".")
			if err != nil {
				message := "can't read messages path <" + p + ">: " + err.Error()
				errs = append(errs, message)
				continue
			}
			_, err = file.Stat()
			if err != nil {
				message := "can't read messages path " + p + ": " + err.Error()
				errs = append(errs, message)
			}
			file.Close()

			if foundMessages {
				continue
			}
			file, err = fs.Open(fallbackLocale + ".yaml")
			if err == nil {
				if fi, err := file.Stat(); err == nil && !fi.IsDir() {
					foundMessages = true
				}
				file.Close()
				continue
			}
			file, err = fs.Open(fallbackLocale)
			if err != nil {
				continue
			}
			if fi, err := file.Stat(); err == nil && fi.IsDir() {
				dirList, err := file.Readdir(-1)
				if err == nil {
					for _, fileInfo := range dirList {
						if fileInfo.IsDir() {
							continue
						}
						if !strings.HasSuffix(fileInfo.Name(), ".yaml") {
							continue
						}
						foundMessages = true
						break
					}
				}
			}
			file.Close()
		}
		if len(errs) == len(openers) {
			errors = append(errors, translatorError{message: errs[0]})
		}
	}

	if !foundRules {
		errors = append(errors, translatorError{message: "found no rules for fallback locale"})
	}

	if !foundMessages {
		errors = append(errors, translatorError{message: "found no messages for fallback locale"})
	}

	f.rulesPaths = rulesPaths
	f.messagesPaths = messagesPaths
	f.translators = map[string]*Translator{}

	// load and check the fallback locale
	if len(fallbackLocale) > 0 {
		var errs []error
		f.fallback, errs = f.GetTranslator(fallbackLocale)
		errors = append(errors, errs...)
	}

	return
}

// GetTranslator returns an Translator instance for the requested locale. If you
// request the same locale multiple times, a pointed to the same Translator will
// be returned each time.
func (f *TranslatorFactory) GetTranslator(localeCode string) (t *Translator, errors []error) {

	f.lock.RLock()
	t, ok := f.translators[localeCode]
	f.lock.RUnlock()
	if ok {
		return t, nil
	}

	return f.load(localeCode)
}

func (f *TranslatorFactory) load(localeCode string) (t *Translator, errors []error) {

	exists, errs := f.LocaleExists(localeCode)
	if !exists {
		errors = append(errors, translatorError{message: "could not find rules and messages for locale " + localeCode})
	}
	errors = append(errors, errs...)

	rules := new(TranslatorRules)
	files := []string{}

	// TODO: the rules loading logic is fairly complex, and there are some
	// specific cases we are not testing for yet. We need to test that the
	// fallback locale rules do not influence the rules loaded, and that the
	// base rules do.

	// the load the base (default) rule values
	// the step above
	for _, p := range f.rulesPaths {
		p = strings.TrimRight(p, pathSeparator)
		files = append(files, p+pathSeparator+"root.yaml")
	}

	// load less specific fallback locale rules
	parts := strings.Split(localeCode, "-")
	if len(parts) > 1 {
		for i := range parts {
			fb := strings.Join(parts[0:i+1], "-")
			for _, p := range f.rulesPaths {
				p = strings.TrimRight(p, pathSeparator)
				files = append(files, p+pathSeparator+fb+".yaml")
				if len(f.project) > 0 {
					files = append(files, p+pathSeparator+f.project+`.`+fb+".yaml")
				}
			}
		}
	}

	// finally load files for this specific locale
	for _, p := range f.rulesPaths {
		p = strings.TrimRight(p, pathSeparator)
		files = append(files, p+pathSeparator+localeCode+".yaml")
		if len(f.project) > 0 {
			files = append(files, p+pathSeparator+f.project+`.`+localeCode+".yaml")
		}
	}

	errs = rules.load(files, f.getFileSystem)
	errors = append(errors, errs...)

	messages, errs := loadMessages(f.project, localeCode, f.messagesPaths, f.getFileSystem)
	errors = append(errors, errs...)

	t = new(Translator)
	t.locale = localeCode
	t.messages = messages
	t.fallback = f.getFallback(localeCode)
	t.rules = rules

	f.lock.Lock()
	f.translators[localeCode] = t
	f.lock.Unlock()

	return
}

// getFallback returns the best fallback for this locale. It first checks for
// less specific versions of the locale before falling back to the global
// fallback if it exists.
func (f *TranslatorFactory) getFallback(localeCode string) *Translator {

	if f.fallback != nil && localeCode == f.fallback.locale {
		return nil
	}

	separator := "-"

	// for a "multipart" locale code, find the most appropriate fallback
	// start by taking off the last "part"
	// if you run out of parts, use the factory's fallback

	fallback := f.fallback
	parts := strings.Split(localeCode, separator)
	for len(parts) > 1 {
		parts = parts[0 : len(parts)-1]
		fb := strings.Join(parts, separator)

		if exists, _ := f.LocaleExists(fb); exists {
			fallback, _ = f.GetTranslator(fb)
			break
		}
	}

	return fallback
}

// LocaleExists checks to see if any messages files exist for the requested
// locale string.
func (f *TranslatorFactory) LocaleExists(localeCode string) (exists bool, errs []error) {
	localeCodes := []string{localeCode}
	if len(f.project) > 0 {
		localeCodes = append(localeCodes, f.project+`.`+localeCode)
	}
	openers := MakeOpeners(f.getFileSystem)
	for _, p := range f.messagesPaths {
		p = strings.TrimRight(p, pathSeparator)
		for _, opener := range openers {
			fs := opener(p)
			var err error
			for _, lc := range localeCodes {
				var file http.File
				file, err = fs.Open(lc + ".yaml")
				if err == nil {
					fi, e := file.Stat()
					if e != nil {
						err = e
					} else if fi.IsDir() {
						err = errors.New(filepath.Join(p, lc+".yaml") + " is dir")
					}
					file.Close()
				}
				if err == nil {
					exists = true
					return
				}
			}
			if !os.IsNotExist(err) {
				errs = append(errs, translatorError{message: "error getting file info: " + err.Error()})
			}

			for _, lc := range localeCodes {
				var file http.File
				file, err = fs.Open(lc)
				if err == nil {
					info, err := file.Stat()
					if err == nil && info.IsDir() {
						dirList, err := file.Readdir(-1)
						if err == nil {
							for _, fileInfo := range dirList {
								if fileInfo.IsDir() {
									continue
								}
								if !strings.HasSuffix(fileInfo.Name(), ".yaml") {
									continue
								}
								exists = true
								file.Close()
								return
							}
						} else if !os.IsNotExist(err) {
							errs = append(errs, translatorError{message: "error getting file info: " + err.Error()})
						}
					}
					file.Close()
				}
			}
		}
	}

	return
}

// loadMessages loads all messages from the properly named locale message yaml
// files in the requested messagesPaths.  if multiple paths are provided, paths
// further down the list take precedence over earlier paths.
func loadMessages(project string, locale string, messagesPaths []string, fsFunc func(string) http.FileSystem) (messages map[string]string, errorList []error) {
	messages = make(map[string]string)
	localeCodes := []string{locale}
	if len(project) > 0 {
		localeCodes = append(localeCodes, project+`.`+locale)
	}
	found := false
	openers := MakeOpeners(fsFunc)
	for _, p := range messagesPaths {
		p = strings.TrimRight(p, pathSeparator)
		for _, opener := range openers {
			fs := opener(p)
			var err error
			for _, lc := range localeCodes {
				var file http.File
				file, err = fs.Open(lc + ".yaml")
				if err != nil {
					continue
				}
				fi, e := file.Stat()
				if e != nil {
					err = e
				} else if fi.IsDir() {
					err = errors.New(filepath.Join(p, lc+".yaml") + " is dir")
				}
				var contents []byte
				newmap := map[string]string{}
				if err == nil {
					contents, err = ioutil.ReadAll(file)
				}
				if err == nil {
					err = confl.Unmarshal(contents, &newmap)
				}
				if err == nil {
					found = true
					for key, value := range newmap {
						messages[key] = value
					}
				} else {
					errorList = append(errorList, translatorError{message: "can't load messages YAML: " + err.Error()})
				}
				file.Close()
			}

			// now look for a directory named after this locale and get an *.yaml children
			for _, lc := range localeCodes {
				var file http.File
				file, err = fs.Open(lc)
				var info os.FileInfo
				if err == nil {
					info, err = file.Stat()
				}
				if err == nil && info.IsDir() {
					// found the directory - now look for *.yaml files
					dirList, err := file.Readdir(-1)
					if err == nil {
						for _, fileInfo := range dirList {
							if fileInfo.IsDir() {
								continue
							}
							if !strings.HasSuffix(fileInfo.Name(), ".yaml") {
								continue
							}
							var fp http.File
							fp, err = fs.Open(path.Join(lc, fileInfo.Name()))
							if err != nil {
								errorList = append(errorList, translatorError{message: "can't open messages file: " + err.Error()})
								fp.Close()
								continue
							}
							var contents []byte
							contents, err = ioutil.ReadAll(fp)
							if err != nil {
								errorList = append(errorList, translatorError{message: "can't open messages file: " + err.Error()})
								fp.Close()
								continue
							}
							newmap := map[string]string{}
							yamlErr := confl.Unmarshal(contents, &newmap)
							if yamlErr != nil {
								errorList = append(errorList, translatorError{message: "can't load messages YAML: " + yamlErr.Error()})
							} else {
								found = true
								for key, value := range newmap {
									messages[key] = value
								}
							}
							fp.Close()
						}
					} else if !os.IsNotExist(err) {
						errorList = append(errorList, translatorError{message: err.Error()})
					}
					file.Close()
				}
			}
		}
	}

	if !found {
		errorList = append(errorList, translatorError{message: "no messages files found: " + locale})
	}

	return
}
