package http

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

const maxBodyLen = 8 << 20 // 8MB

func ParseXml(r *http.Request, v any) error {
	reader := io.LimitReader(r.Body, maxBodyLen)
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)

	decoder := xml.NewDecoder(teeReader)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func Parse(r *http.Request, v any, fns ...func(*http.Request, any) error) error {
	if err := ParseXml(r, v); err != nil {
		return err
	}

	for _, fn := range fns {
		if err := fn(r, v); err != nil {
			return err
		}
	}

	return nil
}
