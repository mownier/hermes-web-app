package utils

import (
	"bytes"
)

func AppendString(b ...string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(b); i++ {
		var s string = b[i]
		buffer.WriteString(s)
	}
	return buffer.String()
}
