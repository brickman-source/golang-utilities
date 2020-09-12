package http

import "bytes"

type UploadFile struct {
	Name string
	Path string
	Data []byte
}

type ProgressReader struct {
	bytes.Buffer
	Reporter func(r int64)
	sent     int64
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Buffer.Read(p)
	pr.sent += int64(n)
	if pr.Reporter != nil {
		pr.Reporter(pr.sent)
	}
	return
}

func (pr *ProgressReader) Write(p []byte) (n int, err error) {
	n, err = pr.Buffer.Write(p)
	return
}
