package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unicode/utf8"
)

var MessageTmpls = map[string]string{
	"Required":     "Can not be empty",
	"Min":          "Minimum is %d",
	"Max":          "Maximum is %d",
	"Range":        "Range is %d to %d",
	"MinSize":      "Minimum size is %d",
	"MaxSize":      "Maximum size is %d",
	"Length":       "Required length is %d",
	"Alpha":        "Must be valid alpha characters",
	"Numeric":      "Must be valid numeric characters",
	"AlphaNumeric": "Must be valid alpha or numeric characters",
	"Match":        "Must match %s",
	"NoMatch":      "Must not match %s",
	"AlphaDash":    "Must be valid alpha or numeric or dash(-_) characters",
	"Email":        "Must be a valid email address",
	"Ip":           "Must be a valid ip address",
	"Base64":       "Must be valid base64 characters",
	"Mobile":       "Must be valid mobile number",
	"Tel":          "Must be valid telephone number",
	"Phone":        "Must be valid telephone or mobile phone number",
	"ZipCode":      "Must be valid zipcode",
}

// SetDefaultMessage set default messages
// if not set, the default messages are
//  "Required":     "Can not be empty",
//  "Min":          "Minimum is %d",
//  "Max":          "Maximum is %d",
//  "Range":        "Range is %d to %d",
//  "MinSize":      "Minimum size is %d",
//  "MaxSize":      "Maximum size is %d",
//  "Length":       "Required length is %d",
//  "Alpha":        "Must be valid alpha characters",
//  "Numeric":      "Must be valid numeric characters",
//  "AlphaNumeric": "Must be valid alpha or numeric characters",
//  "Match":        "Must match %s",
//  "NoMatch":      "Must not match %s",
//  "AlphaDash":    "Must be valid alpha or numeric or dash(-_) characters",
//  "Email":        "Must be a valid email address",
//  "IP":           "Must be a valid ip address",
//  "Base64":       "Must be valid base64 characters",
//  "Mobile":       "Must be valid mobile number",
//  "Tel":          "Must be valid telephone number",
//  "Phone":        "Must be valid telephone or mobile phone number",
//  "ZipCode":      "Must be valid zipcode",
func SetDefaultMessage(msg map[string]string) {
	if len(msg) == 0 {
		return
	}

	for name := range msg {
		MessageTmpls[name] = msg[name]
	}
}

type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
	GetKey() string
	GetLimitValue() interface{}
}

type Required struct {
	Key string
}

func (r Required) IsSatisfied(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		return len(str) > 0
	}
	if _, ok := obj.(bool); ok {
		return true
	}
	if i, ok := obj.(int); ok {
		return i != 0
	}
	if i, ok := obj.(uint); ok {
		return i != 0
	}
	if i, ok := obj.(int8); ok {
		return i != 0
	}
	if i, ok := obj.(uint8); ok {
		return i != 0
	}
	if i, ok := obj.(int16); ok {
		return i != 0
	}
	if i, ok := obj.(uint16); ok {
		return i != 0
	}
	if i, ok := obj.(uint32); ok {
		return i != 0
	}
	if i, ok := obj.(int32); ok {
		return i != 0
	}
	if i, ok := obj.(int64); ok {
		return i != 0
	}
	if i, ok := obj.(uint64); ok {
		return i != 0
	}
	if t, ok := obj.(time.Time); ok {
		return !t.IsZero()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}

func (r Required) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Required"])
}

func (r Required) GetKey() string {
	return r.Key
}

func (r Required) GetLimitValue() interface{} {
	return nil
}

type Min struct {
	Min float64
	Key string
}

func (m Min) IsSatisfied(obj interface{}) bool {
	switch num := obj.(type) {
	case int:
		return float64(num) >= m.Min
	case int64:
		return float64(num) >= m.Min
	case float64:
		return num >= m.Min
	default:
		return false
	}
}

func (m Min) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Min"], m.Min)
}

func (m Min) GetKey() string {
	return m.Key
}

func (m Min) GetLimitValue() interface{} {
	return m.Min
}

type Max struct {
	Max float64
	Key string
}

func (m Max) IsSatisfied(obj interface{}) bool {
	switch num := obj.(type) {
	case int:
		return float64(num) <= m.Max
	case int64:
		return float64(num) <= m.Max
	case float64:
		return num <= m.Max
	default:
		return false
	}
}

func (m Max) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Max"], m.Max)
}

func (m Max) GetKey() string {
	return m.Key
}

func (m Max) GetLimitValue() interface{} {
	return m.Max
}

// Range Requires an integer to be within Min, Max inclusive.
type Range struct {
	Min
	Max
	Key string
}

func (r Range) IsSatisfied(obj interface{}) bool {
	return r.Min.IsSatisfied(obj) && r.Max.IsSatisfied(obj)
}

func (r Range) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Range"], r.Min.Min, r.Max.Max)
}

func (r Range) GetKey() string {
	return r.Key
}

func (r Range) GetLimitValue() interface{} {
	return []float64{r.Min.Min, r.Max.Max}
}

// MinSize Requires an array or string to be at least a given length.
type MinSize struct {
	Min int
	Key string
}

func (m MinSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) >= m.Min
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() >= m.Min
	}
	return false
}

func (m MinSize) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["MinSize"], m.Min)
}

func (m MinSize) GetKey() string {
	return m.Key
}

func (m MinSize) GetLimitValue() interface{} {
	return m.Min
}

// Requires an array or string to be at most a given length.
type MaxSize struct {
	Max int
	Key string
}

func (m MaxSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) <= m.Max
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() <= m.Max
	}
	return false
}

func (m MaxSize) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["MaxSize"], m.Max)
}

func (m MaxSize) GetKey() string {
	return m.Key
}

func (m MaxSize) GetLimitValue() interface{} {
	return m.Max
}

// Requires an array or string to be exactly a given length.
type Length struct {
	N   int
	Key string
}

func (l Length) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) == l.N
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() == l.N
	}
	return false
}

func (l Length) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Length"], l.N)
}

func (l Length) GetKey() string {
	return l.Key
}

func (l Length) GetLimitValue() interface{} {
	return l.N
}

type Alpha struct {
	Key string
}

func (a Alpha) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
				return false
			}
		}
		return true
	}
	return false
}

func (a Alpha) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Alpha"])
}

func (a Alpha) GetKey() string {
	return a.Key
}

func (a Alpha) GetLimitValue() interface{} {
	return nil
}

type Numeric struct {
	Key string
}

func (n Numeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if '9' < v || v < '0' {
				return false
			}
		}
		return true
	}
	return false
}

func (n Numeric) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Numeric"])
}

func (n Numeric) GetKey() string {
	return n.Key
}

func (n Numeric) GetLimitValue() interface{} {
	return nil
}

type AlphaNumeric struct {
	Key string
}

func (a AlphaNumeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
				return false
			}
		}
		return true
	}
	return false
}

func (a AlphaNumeric) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["AlphaNumeric"])
}

func (a AlphaNumeric) GetKey() string {
	return a.Key
}

func (a AlphaNumeric) GetLimitValue() interface{} {
	return nil
}

// Requires a string to match a given regex.
type Match struct {
	Regexp *regexp.Regexp
	Key    string
}

func (m Match) IsSatisfied(obj interface{}) bool {
	return m.Regexp.MatchString(fmt.Sprintf("%v", obj))
}

func (m Match) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Match"], m.Regexp.String())
}

func (m Match) GetKey() string {
	return m.Key
}

func (m Match) GetLimitValue() interface{} {
	return m.Regexp.String()
}

// Requires a string to not match a given regex.
type NoMatch struct {
	Match
	Key string
}

func (n NoMatch) IsSatisfied(obj interface{}) bool {
	return !n.Match.IsSatisfied(obj)
}

func (n NoMatch) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["NoMatch"], n.Regexp.String())
}

func (n NoMatch) GetKey() string {
	return n.Key
}

func (n NoMatch) GetLimitValue() interface{} {
	return n.Regexp.String()
}

var (
	alphaDashPattern *regexp.Regexp
	emailPattern     *regexp.Regexp
	zipCodePattern   *regexp.Regexp
	telPattern       *regexp.Regexp
	mobilePattern    *regexp.Regexp
	ipPattern        *regexp.Regexp
	base64Pattern    *regexp.Regexp
)

func init() {
	alphaDashPattern = regexp.MustCompile(DefaultRule.AlphaDash)
	emailPattern = regexp.MustCompile(DefaultRule.Email)
	zipCodePattern = regexp.MustCompile(DefaultRule.ZipCode)
	telPattern = regexp.MustCompile(DefaultRule.Telephone)
	mobilePattern = regexp.MustCompile(DefaultRule.Mobile)
	ipPattern = regexp.MustCompile(DefaultRule.IPv4)
	base64Pattern = regexp.MustCompile(DefaultRule.Base64)
}

type AlphaDash struct {
	NoMatch
	Key string
}

func (a AlphaDash) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["AlphaDash"])
}

func (a AlphaDash) GetKey() string {
	return a.Key
}

func (a AlphaDash) GetLimitValue() interface{} {
	return nil
}

type Email struct {
	Match
	Key string
}

func (e Email) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Email"])
}

func (e Email) GetKey() string {
	return e.Key
}

func (e Email) GetLimitValue() interface{} {
	return nil
}

type Ip struct {
	Match
	Key string
}

func (i Ip) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Ip"])
}

func (i Ip) GetKey() string {
	return i.Key
}

func (i Ip) GetLimitValue() interface{} {
	return nil
}

type Base64 struct {
	Match
	Key string
}

func (b Base64) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Base64"])
}

func (b Base64) GetKey() string {
	return b.Key
}

func (b Base64) GetLimitValue() interface{} {
	return nil
}

// just for chinese mobile phone number

type Mobile struct {
	Match
	Key string
}

func (m Mobile) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Mobile"])
}

func (m Mobile) GetKey() string {
	return m.Key
}

func (m Mobile) GetLimitValue() interface{} {
	return nil
}

// just for chinese telephone number

type Tel struct {
	Match
	Key string
}

func (t Tel) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Tel"])
}

func (t Tel) GetKey() string {
	return t.Key
}

func (t Tel) GetLimitValue() interface{} {
	return nil
}

// just for chinese telephone or mobile phone number
type Phone struct {
	Mobile
	Tel
	Key string
}

func (p Phone) IsSatisfied(obj interface{}) bool {
	return p.Mobile.IsSatisfied(obj) || p.Tel.IsSatisfied(obj)
}

func (p Phone) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Phone"])
}

func (p Phone) GetKey() string {
	return p.Key
}

func (p Phone) GetLimitValue() interface{} {
	return nil
}

// just for chinese zipcode

type ZipCode struct {
	Match
	Key string
}

func (z ZipCode) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["ZipCode"])
}

func (z ZipCode) GetKey() string {
	return z.Key
}

func (z ZipCode) GetLimitValue() interface{} {
	return nil
}
