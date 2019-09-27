package serializer

import (
	"bytes"
	"compress/zlib"
	"io"
)

func Compress(src []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(src); err != nil {
		return nil, err
	}
	w.Close() // nolint: errcheck

	return b.Bytes(), nil
}

func Decompress(src []byte) ([]byte, error) {
	b := bytes.NewReader(src)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}

	var dst bytes.Buffer
	if _, err = io.Copy(&dst, r); err != nil {
		return nil, err
	}
	r.Close()

	return dst.Bytes(), nil
}
