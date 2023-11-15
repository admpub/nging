package echo

import (
	"strings"

	"github.com/webx-top/echo/param"
)

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(interface{}, Context, ...FormDataFilter) error
		BindAndValidate(interface{}, Context, ...FormDataFilter) error
		MustBind(interface{}, Context, ...FormDataFilter) error
		MustBindAndValidate(interface{}, Context, ...FormDataFilter) error

		BindWithDecoder(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error
		BindAndValidateWithDecoder(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error
		MustBindWithDecoder(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error
		MustBindAndValidateWithDecoder(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error
	}
	binder struct {
		decoders map[string]func(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error
	}

	BeforeBind interface {
		BeforeBind(Context) error
	}

	// for tag

	BinderValueDecoder func(field string, values []string, params string) (interface{}, error)
	BinderValueEncoder func(field string, value interface{}, params string) []string

	// for function argument

	BinderValueCustomDecoder  func(values []string) (interface{}, error)
	BinderValueCustomDecoders map[string]BinderValueCustomDecoder // 这里的 key 为结构体字段或map的key层级路径
	BinderValueCustomEncoder  func(interface{}) []string
	BinderValueCustomEncoders map[string]BinderValueCustomEncoder // 这里的 key 为表单字段层级路径

	ValueDecodersGetter interface {
		ValueDecoders(Context) BinderValueCustomDecoders
	}

	ValueEncodersGetter interface {
		ValueEncoders(Context) BinderValueCustomEncoders
	}

	ValueStringersGetter interface {
		ValueStringers(Context) param.StringerMap
	}

	FormNameFormatterGetter interface {
		FormNameFormatter(Context) FieldNameFormatter
	}

	BinderFormTopNamer interface {
		BinderFormTopName() string
	}

	BinderKeyNormalizer interface {
		BinderKeyNormalizer(string) string
	}

	// FromConversion a struct implements this interface can be convert from request param to a struct
	FromConversion interface {
		FromString(content string) error
	}

	// ToConversion a struct implements this interface can be convert from struct to template variable
	// Not Implemented
	ToConversion interface {
		ToString() string
	}
)

func NewBinder(e *Echo) Binder {
	return &binder{
		decoders: DefaultBinderDecoders,
	}
}

func (b *binder) MustBind(i interface{}, c Context, filter ...FormDataFilter) error {
	return b.MustBindWithDecoder(i, c, nil, filter...)
}

func (b *binder) MustBindAndValidate(i interface{}, c Context, filter ...FormDataFilter) error {
	return b.MustBindAndValidateWithDecoder(i, c, nil, filter...)
}

func (b *binder) Bind(i interface{}, c Context, filter ...FormDataFilter) (err error) {
	return b.BindWithDecoder(i, c, nil, filter...)
}

func (b *binder) BindAndValidate(i interface{}, c Context, filter ...FormDataFilter) error {
	return b.BindAndValidateWithDecoder(i, c, nil, filter...)
}

func (b *binder) MustBindWithDecoder(i interface{}, c Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
	if f, y := i.(BeforeBind); y {
		if err := f.BeforeBind(c); err != nil {
			return err
		}
	}
	contentType := c.Request().Header().Get(HeaderContentType)
	contentType = strings.ToLower(strings.TrimSpace(strings.SplitN(contentType, `;`, 2)[0]))
	if decoder, ok := b.decoders[contentType]; ok {
		return decoder(i, c, valueDecoders, filter...)
	}
	if decoder, ok := b.decoders[`*`]; ok {
		return decoder(i, c, valueDecoders, filter...)
	}
	return ErrUnsupportedMediaType
}

func (b *binder) MustBindAndValidateWithDecoder(i interface{}, c Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
	if err := b.MustBindWithDecoder(i, c, valueDecoders, filter...); err != nil {
		return err
	}
	return ValidateStruct(c, i)
}

func (b *binder) BindWithDecoder(i interface{}, c Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) (err error) {
	err = b.MustBindWithDecoder(i, c, valueDecoders, filter...)
	if err == ErrUnsupportedMediaType {
		err = nil
	}
	return
}

func (b *binder) BindAndValidateWithDecoder(i interface{}, c Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
	if err := b.MustBindWithDecoder(i, c, valueDecoders, filter...); err != nil {
		if err != ErrUnsupportedMediaType {
			return err
		}
	}
	return ValidateStruct(c, i)
}

func (b *binder) SetDecoders(decoders map[string]func(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error) {
	b.decoders = decoders
}

func (b *binder) AddDecoder(mime string, decoder func(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error) {
	b.decoders[mime] = decoder
}
