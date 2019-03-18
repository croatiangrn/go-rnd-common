package go_rnd_common

import (
	"log"
	"runtime"
)

func LogError(err error) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] %s:%d %v", fn, line, err)
}
