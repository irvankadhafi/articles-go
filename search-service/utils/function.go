package utils

import "runtime"

// MyCaller will return the method caller. skip value defines how many steps to be skipped.
// skip=0 will always return the MyCaller
// skip=1 returns the caller of the MyCaller
// and so on...
func MyCaller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return details.Name()
	}
	return "failed to identify method caller"
}
