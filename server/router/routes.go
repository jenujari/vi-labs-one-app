package router

import (
	"net/http"

	"server/helpers"
	"server/tpl"
)

func SetRoutes(r *http.ServeMux) {
	r.HandleFunc("/", Index)


	r.HandleFunc("GET /api/zerodha/token_refresh", refreshZerodhaToken)
	r.HandleFunc("GET /api/zerodha/token", checkZerodhaTokenValidity)
	r.HandleFunc("GET /api/zerodha/rsi_obv", runRsiObvFilterForDay)

	r.HandleFunc("/test", test)
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := tpl.GetTemplateExecutor()
	err := t.ExecuteTemplate(w, "test", nil)
	if err != nil {
		return
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	helpers.DumpHTML(w, "test")
}