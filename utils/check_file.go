package utils

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

var httpClient = http.Client{
	Timeout: time.Duration(2) * time.Second,
}

func try(f func() error) (err error) {
	if err = f(); err == nil {
		return
	}

	for i := 1; i < 3; i++ {
		time.Sleep(time.Millisecond * time.Duration(10))

		if err = f(); err == nil {
			return
		}
	}

	return
}

func CheckFile(url string, fileType string, maxSize int) error {
	return try(func() error {
		resp, err := httpClient.Head(url)
		if err != nil {
			return err
		}

		if c := resp.StatusCode; !(c >= http.StatusOK && c < http.StatusMultipleChoices) {
			return errors.New("can't detect")
		}

		if fileType != "" {
			ct := resp.Header.Get("content-type")
			if !strings.Contains(strings.ToLower(ct), strings.ToLower(fileType)) {
				return errors.New("unknown file type")
			}
		}

		if maxSize > 0 {
			if resp.ContentLength == -1 {
				return errors.New("unknown file size")
			}

			if resp.ContentLength > int64(maxSize) {
				return errors.New("big file")
			}
		}

		return nil
	})
}
