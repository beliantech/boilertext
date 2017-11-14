package logger

import (
	"fmt"
)

const DEBUG = false

func Println(a ...interface{}) {
	if DEBUG {
		fmt.Println(a...)
	}
}
