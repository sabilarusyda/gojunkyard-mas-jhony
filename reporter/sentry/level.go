package sentry

const (
	DEBUG uint8 = iota + 1
	INFO
	WARNING
	ERROR
	FATAL
)

func isDebug(level uint8) bool {
	return level <= DEBUG
}

func isInfo(level uint8) bool {
	return level <= INFO
}
func isWarning(level uint8) bool {
	return level <= WARNING
}
func isError(level uint8) bool {
	return level <= ERROR
}
func isFatal(level uint8) bool {
	return level <= FATAL
}
