package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func NewGzip() *Gzip {
	buf := &bytes.Buffer{}
	gz := gzip.NewWriter(buf)
	return &Gzip{
		buf: buf,
		gz:  gz,
	}
}

type Gzip struct {
	buf *bytes.Buffer
	gz  *gzip.Writer
}

func (g *Gzip) Write(p []byte) (n int, err error) {
	return g.gz.Write(p)
}

func (g *Gzip) Len() int {
	return g.buf.Len()
}

func (g *Gzip) Bytes() []byte {
	g.gz.Close()
	return g.buf.Bytes()
}

func (g *Gzip) Hash() string {
	return Hash(g.Bytes())
}

func UnzipString(b []byte) (string, error) {
	bt, err := UnzipByte(b)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}

func UnzipByte(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return []byte{}, err
	}
	bt, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	return bt, nil
}
