package mw

import (
	"compress/gzip"
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// GZIP - middleware that decompress request.body when header Content-Encoding = gzip
// and provides http.ResponseWriter with gzip Writer
// when header Accept-Encoding contains gzip and set w.header Content-Encoding = gzip
func GZIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// request part decode
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = gz
		}

		// response part
		if strings.Contains(r.Header.Get(resp.HeaderAcceptEncoding), "gzip") {
			// создаём gzip.Writer поверх текущего w
			gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
			if err != nil {
				resp.NoContent(w, http.StatusBadRequest)
				return
			}
			defer gz.Close()

			w.Header().Set(resp.HeaderContentEncoding, "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		}

		next.ServeHTTP(w, r)
	})
}
