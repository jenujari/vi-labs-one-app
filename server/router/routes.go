package router

import (
	"log"
	"net/http"

	"server/config"
	"server/helpers"
	"server/tpl"

	"github.com/goforj/godump"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func SetRoutes(r *http.ServeMux) {
	r.HandleFunc("/", Index)
	r.HandleFunc("/dd", DD)
	r.HandleFunc("/authenticated", Authenticated)
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := tpl.GetTemplateExecutor()
	err := t.ExecuteTemplate(w, "test", nil)
	if err != nil {
		return
	}
}

func DD(w http.ResponseWriter, r *http.Request) {
	kc := helpers.GetKiteClient()

	dump := struct {
		KC    *kiteconnect.Client
		Url      string
	}{
		KC:   kc,
		Url:  kc.GetLoginURL(),
	}

	html := godump.DumpHTML(dump)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func Authenticated(w http.ResponseWriter, r *http.Request) {
	kc := helpers.GetKiteClient()
	apiSecret := config.GetConfig().Secret.ApiSecret
	params := r.URL.Query()

	// status := params.Get("status")
	requestToken := params.Get("request_token")


	// // Get user details and access token
	data, err := kc.GenerateSession(requestToken, apiSecret)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// // Set access token
	kc.SetAccessToken(data.AccessToken)
	log.Println("data.AccessToken", data.AccessToken)

	// Get margins
	margins, err := kc.GetUserMargins()
	if err != nil {
		http.Error(w, "Failed to get user margins", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	html := godump.DumpHTML(margins)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}