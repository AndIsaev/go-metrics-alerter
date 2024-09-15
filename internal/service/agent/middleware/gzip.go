package middleware

import (
	"github.com/go-resty/resty/v2"
)

// GzipRequestMiddleware - указывает заголовок данных в формате gzip.
func GzipRequestMiddleware(c *resty.Client, _ *resty.Request) error {
	c.Header.Set("Accept-Encoding", "gzip")
	return nil
}
