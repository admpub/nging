package param

func IsNotEmptyString(s *string) bool {
	if s == nil {
		return false
	}
	return len(*s) > 0
}

func IsNotEmptyStringElement(_ int, s string) bool {
	return len(s) > 0
}

func IsTrueBoolElement(_ int, s bool) bool {
	return s
}

func IsFalseBoolElement(_ int, s bool) bool {
	return !s
}
