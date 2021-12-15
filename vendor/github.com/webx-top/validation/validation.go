package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/webx-top/tagfast"
)

func New(args ...func(*Validation)) *Validation {
	valid := &Validation{}
	if len(args) > 0 {
		for _, fn := range args {
			fn(valid)
		}
	}
	return valid
}

type ValidFormer interface {
	Valid(*Validation)
}

var NoError = &ValidationError{}

type ValidationError struct {
	Message    string      //错误信息
	Key        string      //验证键(比如：Title|Required)
	Name       string      //验证器名称
	Field      string      //字段名称
	Tmpl       string      //错误信息所使用的文本模板
	Value      interface{} //要验证的值
	LimitValue interface{}
	withField  bool
}

// Returns the Message.
func (e *ValidationError) String() string {
	if e == nil {
		return ""
	}
	if e.withField {
		return e.Field + `: ` + e.Message
	}
	return e.Message
}

func (e *ValidationError) WithField(args ...bool) *ValidationError {
	e.withField = true
	if len(args) > 0 {
		e.withField = args[0]
	}
	return e
}

func (e *ValidationError) Error() string {
	return e.String()
}

// A ValidationResult is returned from every validation method.
// It provides an indication of success, and a pointer to the Error (if any).
type ValidationResult struct {
	Error *ValidationError
	Ok    bool
	Valid *Validation
}

func (r *ValidationResult) Key(key string) *ValidationResult {
	if r.Error != nil {
		r.Error.Key = key
	}
	return r
}

func (r *ValidationResult) Message(message string, args ...interface{}) *ValidationResult {
	if r.Error != nil {
		if len(args) == 0 {
			r.Error.Message = message
		} else {
			r.Error.Message = fmt.Sprintf(message, args...)
		}
	}
	return r
}

// A Validation context manages data validation and error messages.
type Validation struct {
	Errors    []*ValidationError
	ErrorsMap map[string]*ValidationError
	SendError func(*ValidationError)
}

func (v *Validation) Clear() {
	v.Errors = []*ValidationError{}
	v.ErrorsMap = nil
}

func (v *Validation) HasError() bool {
	return len(v.Errors) > 0
}

func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 1
}

// ErrorMap Return the errors mapped by key.
// If there are multiple validation errors associated with a single key, the
// first one "wins".  (Typically the first validation will be the more basic).
func (v *Validation) ErrorMap() map[string]*ValidationError {
	return v.ErrorsMap
}

func (v *Validation) Error() *ValidationError {
	if v.HasError() {
		return v.Errors[0]
	}
	return NoError
}

// Required Test that the argument is non-nil and non-empty (if string or list)
func (v *Validation) Required(obj interface{}, key string) *ValidationResult {
	return v.apply(Required{key}, obj)
}

// Min Test that the obj is greater than min if obj's type is int
func (v *Validation) Min(obj interface{}, min float64, key string) *ValidationResult {
	return v.apply(Min{min, key}, obj)
}

// Max Test that the obj is less than max if obj's type is int
func (v *Validation) Max(obj interface{}, max float64, key string) *ValidationResult {
	return v.apply(Max{max, key}, obj)
}

// Range Test that the obj is between mni and max if obj's type is int
func (v *Validation) Range(obj interface{}, min, max float64, key string) *ValidationResult {
	return v.apply(Range{Min{Min: min}, Max{Max: max}, key}, obj)
}

func (v *Validation) MinSize(obj interface{}, min int, key string) *ValidationResult {
	return v.apply(MinSize{min, key}, obj)
}

func (v *Validation) MaxSize(obj interface{}, max int, key string) *ValidationResult {
	return v.apply(MaxSize{max, key}, obj)
}

func (v *Validation) Length(obj interface{}, n int, key string) *ValidationResult {
	return v.apply(Length{n, key}, obj)
}

func (v *Validation) Alpha(obj interface{}, key string) *ValidationResult {
	return v.apply(Alpha{key}, obj)
}

func (v *Validation) Numeric(obj interface{}, key string) *ValidationResult {
	return v.apply(Numeric{key}, obj)
}

func (v *Validation) AlphaNumeric(obj interface{}, key string) *ValidationResult {
	return v.apply(AlphaNumeric{key}, obj)
}

func (v *Validation) Match(obj interface{}, regex *regexp.Regexp, key string) *ValidationResult {
	return v.apply(Match{regex, key}, obj)
}

func (v *Validation) NoMatch(obj interface{}, regex *regexp.Regexp, key string) *ValidationResult {
	return v.apply(NoMatch{Match{Regexp: regex}, key}, obj)
}

func (v *Validation) AlphaDash(obj interface{}, key string) *ValidationResult {
	return v.apply(AlphaDash{NoMatch{Match: Match{Regexp: alphaDashPattern}}, key}, obj)
}

func (v *Validation) Email(obj interface{}, key string) *ValidationResult {
	return v.apply(Email{Match{Regexp: emailPattern}, key}, obj)
}

func (v *Validation) Ip(obj interface{}, key string) *ValidationResult {
	return v.apply(Ip{Match{Regexp: ipPattern}, key}, obj)
}

func (v *Validation) Base64(obj interface{}, key string) *ValidationResult {
	return v.apply(Base64{Match{Regexp: base64Pattern}, key}, obj)
}

func (v *Validation) Mobile(obj interface{}, key string) *ValidationResult {
	return v.apply(Mobile{Match{Regexp: mobilePattern}, key}, obj)
}

func (v *Validation) Tel(obj interface{}, key string) *ValidationResult {
	return v.apply(Tel{Match{Regexp: telPattern}, key}, obj)
}

func (v *Validation) Phone(obj interface{}, key string) *ValidationResult {
	return v.apply(Phone{Mobile{Match: Match{Regexp: mobilePattern}},
		Tel{Match: Match{Regexp: telPattern}}, key}, obj)
}

func (v *Validation) ZipCode(obj interface{}, key string) *ValidationResult {
	return v.apply(ZipCode{Match{Regexp: zipCodePattern}, key}, obj)
}

func (v *Validation) apply(chk Validator, obj interface{}) *ValidationResult {
	if chk.IsSatisfied(obj) {
		return &ValidationResult{Ok: true, Valid: v}
	}

	// Add the error to the validation context.
	key := chk.GetKey()
	Field := key
	Name := ""

	parts := strings.Split(key, "|")
	if len(parts) == 2 {
		Field = parts[0]
		Name = parts[1]
	}

	err := &ValidationError{
		Message:    chk.DefaultMessage(),
		Key:        key,
		Name:       Name,
		Field:      Field,
		Value:      obj,
		Tmpl:       MessageTmpls[Name],
		LimitValue: chk.GetLimitValue(),
	}
	v.setError(err)

	// Also return it in the result.
	return &ValidationResult{
		Ok:    false,
		Error: err,
		Valid: v,
	}
}

func (v *Validation) setError(err *ValidationError) {
	v.Errors = append(v.Errors, err)
	if v.ErrorsMap == nil {
		v.ErrorsMap = make(map[string]*ValidationError)
	}
	if _, ok := v.ErrorsMap[err.Field]; !ok {
		v.ErrorsMap[err.Field] = err
	}
	if v.SendError != nil {
		v.SendError(err)
	}
}

func (v *Validation) SetError(fieldName string, errMsg string) *ValidationError {
	err := &ValidationError{Key: fieldName, Field: fieldName, Tmpl: errMsg, Message: errMsg}
	v.setError(err)
	return err
}

// Check Apply a group of validators to a field, in order, and return the
// ValidationResult from the first one that fails, or the last one that
// succeeds.
func (v *Validation) Check(obj interface{}, checks ...Validator) *ValidationResult {
	var result *ValidationResult
	for _, check := range checks {
		result = v.apply(check, obj)
		if !result.Ok {
			return result
		}
	}
	return result
}

// Valid the obj parameter must be a struct or a struct pointer
func (v *Validation) Valid(obj interface{}, args ...string) (ok bool, err error) {
	err = v.validExec(obj, "", args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !v.HasError() {
		if form, ok := obj.(ValidFormer); ok {
			form.Valid(v)
		}
	}
	ok = v.HasError() == false
	return
}

// ValidResult the obj parameter must be a struct or a struct pointer
func (v *Validation) ValidResult(obj interface{}, args ...string) (ok bool, errs map[string]string) {
	ok, _ = v.Valid(obj, args...)
	if !ok {
		errs = v.ErrMap()
	}
	return
}

func (v *Validation) ErrMap() map[string]string {
	errs := make(map[string]string)
	for _, err := range v.Errors {
		errs[err.Field] = err.Message
	}
	return errs
}

func (v *Validation) validExec(obj interface{}, baseName string, args ...string) (err error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		err = fmt.Errorf("%v must be a struct or a struct pointer", obj)
		return
	}
	chkFields := make(map[string][]string)
	pNum := len(args)
	//fmt.Println(objT.Name(), ":[Struct NumIn]", pNum)
	if pNum > 0 {
		//aa.b.c,ab.b.c
		for _, v := range args {
			arr := strings.SplitN(v, ".", 2)
			if _, ok := chkFields[arr[0]]; !ok {
				chkFields[arr[0]] = make([]string, 0)
			}
			if len(arr) > 1 {
				chkFields[arr[0]] = append(chkFields[arr[0]], arr[1])
			}
		}
	}
	args = make([]string, 0)
	if len(chkFields) > 0 { //检测指定字段
		for field, args := range chkFields {
			f, ok := objT.FieldByName(field)
			if !ok {
				err = fmt.Errorf("No name for the '%s' field", field)
				return
			}
			tag := tagfast.Value(objT, f, VALIDTAG)
			if tag == "-" {
				continue
			}
			var vfs []ValidFunc

			var fName string
			if baseName == "" {
				fName = f.Name
			} else {
				fName = strings.Join([]string{baseName, f.Name}, ".")
			}
			fv := objV.FieldByName(field)
			if isStruct(f.Type) || isStructPtr(f.Type) {
				if fv.CanInterface() {
					err = v.validExec(fv.Interface(), fName, args...)
				}
				continue
			}
			if vfs, err = getValidFuncs(f, objT, fName); err != nil {
				return
			}
			for _, vf := range vfs {
				if _, err = funcs.Call(vf.Name,
					mergeParam(v, fv.Interface(), vf.Params)...); err != nil {
					return
				}
			}
		}
	} else { //检测全部字段
		for i := 0; i < objT.NumField(); i++ {
			tag := tagfast.Value(objT, objT.Field(i), VALIDTAG)
			if tag == "-" {
				continue
			}
			var vfs []ValidFunc

			var fName string
			if baseName == "" {
				fName = objT.Field(i).Name
			} else {
				fName = strings.Join([]string{baseName, objT.Field(i).Name}, ".")
			}
			//fmt.Println(fName, ":[Type]:", objT.Field(i).Type.Kind())
			if isStruct(objT.Field(i).Type) || isStructPtr(objT.Field(i).Type) {
				if objV.Field(i).CanInterface() {
					err = v.validExec(objV.Field(i).Interface(), fName)
				}
				continue
			}
			if vfs, err = getValidFuncs(objT.Field(i), objT, fName); err != nil {
				return
			}
			for _, vf := range vfs {
				if _, err = funcs.Call(vf.Name,
					mergeParam(v, objV.Field(i).Interface(), vf.Params)...); err != nil {
					return
				}
			}
		}
	}
	return
}

func (v *Validation) ValidSimple(name string, val interface{}, rule string) (b bool, err error) {
	err = v.validSimpleExec(val, rule, name)
	if err != nil {
		fmt.Println(err)
		return
	}
	b = !v.HasError()
	return
}

func (v *Validation) ValidField(name string, val string, rule string) (b bool) {
	b, _ = v.ValidSimple(name, val, rule)
	return
}

func (v *Validation) ValidOk(obj interface{}, args ...string) (b bool) {
	b, _ = v.Valid(obj, args...)
	return
}

func (v *Validation) validSimpleExec(val interface{}, rule string, fName string) (err error) {
	var vfs []ValidFunc
	if vfs, rule, err = getRegFuncs(rule, fName); err != nil {
		return
	}
	fs := strings.Split(rule, ";")
	for _, vfunc := range fs {
		var vf ValidFunc
		if len(vfunc) == 0 {
			continue
		}
		vf, err = parseFunc(vfunc, fName)
		if err != nil {
			return
		}
		vfs = append(vfs, vf)
	}
	for _, vf := range vfs {
		if _, err = funcs.Call(vf.Name,
			mergeParam(v, val, vf.Params)...); err != nil {
			return
		}
	}
	return
}
