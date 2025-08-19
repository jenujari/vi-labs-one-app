package router

import (
	"net/http"

	"server/tpl"
)

func SetRoutes(r *http.ServeMux) {
	r.HandleFunc("/", Index)
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := tpl.GetTemplateExecutor()
	err := t.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		return
	}
}
