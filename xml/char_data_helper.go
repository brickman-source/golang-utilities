package xml

import "strings"

type CharData struct {
	Text []byte `xml:",innerxml"`
}

func NewCharData(s string) CharData {
	//s = strings.Replace(s, "[", "&#91;", -1)
	//s = strings.Replace(s, "]", "&#93;", -1)
	return CharData{[]byte("<![CDATA[" + strings.Replace(s, "[", "&#91;", -1) + "]]>")}
}

func WrapCharData(s string) string {
	return "<![CDATA[" + s + "]]>"
}

type CData struct {
	Value string `xml:",cdata"`
}

func NewCData(s string) CData {
	return CData{s}
}
