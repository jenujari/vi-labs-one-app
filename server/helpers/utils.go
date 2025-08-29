package helpers

import (
	"net/http"

	"github.com/goforj/godump"
)

func DumpHTML(w http.ResponseWriter, data any) {
	html := godump.DumpHTML(data)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
