/*
 * Copyright (c) 2019. Brickman.
 */

package alipay

import (
	"bytes"
)

func aliEscape(str string) string {
	data := []byte(str)
	var buf bytes.Buffer
	for _, c := range data {
		switch c {
		case '!':
			buf.Write([]byte(`%21`))
		case '*':
			buf.Write([]byte(`%2A`))
		case '\'':
			buf.Write([]byte(`%27`))
		case '(':
			buf.Write([]byte(`%28`))
		case ')':
			buf.Write([]byte(`%29`))
		case ';':
			buf.Write([]byte(`%3B`))
		case ':':
			buf.Write([]byte(`%3A`))
		case '@':
			buf.Write([]byte(`%40`))
		case '&':
			buf.Write([]byte(`%26`))
		case '=':
			buf.Write([]byte(`%3D`))
		case '+':
			buf.Write([]byte(`%2B`))
		case '$':
			buf.Write([]byte(`%24`))
		case ',':
			buf.Write([]byte(`%2C`))
		case '/':
			buf.Write([]byte(`%2F`))
		case '?':
			buf.Write([]byte(`%3F`))
		case '%':
			buf.Write([]byte(`%25`))
		case '#':
			buf.Write([]byte(`%23`))
		case '[':
			buf.Write([]byte(`%5B`))
		case ']':
			buf.Write([]byte(`%5D`))
		default:
			buf.WriteByte(c)
		}
	}
	return buf.String()
}
