package securecookie

import "errors"

func IsValueTooLong(err error) bool {
	return errors.Is(err, errEncodedValueTooLong) || errors.Is(err, errValueToDecodeTooLong)
}

func IsEncodedValueTooLong(err error) bool {
	return errors.Is(err, errEncodedValueTooLong)
}

func IsDecodeValueTooLong(err error) bool {
	return errors.Is(err, errValueToDecodeTooLong)
}

func IsNoCodecs(err error) bool {
	return errors.Is(err, errNoCodecs)
}

func IsHashKeyNotSet(err error) bool {
	return errors.Is(err, errHashKeyNotSet)
}

func IsBlockKeyNotSet(err error) bool {
	return errors.Is(err, errBlockKeyNotSet)
}

func IsTimestampInvalid(err error) bool {
	return errors.Is(err, errTimestampInvalid)
}

func IsTimestampTooNew(err error) bool {
	return errors.Is(err, errTimestampTooNew)
}

func IsTimestampExpired(err error) bool {
	return errors.Is(err, errTimestampExpired)
}

func IsDecryptionFailed(err error) bool {
	return errors.Is(err, errDecryptionFailed)
}

func IsValueNotByte(err error) bool {
	return errors.Is(err, errValueNotByte)
}

func IsValueNotBytePtr(err error) bool {
	return errors.Is(err, errValueNotBytePtr)
}

func IsMacInvalid(err error) bool {
	return errors.Is(err, ErrMacInvalid)
}
