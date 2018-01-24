// +build debug

package basicblock

import "log"

// to use: go test -tags debug
func debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
