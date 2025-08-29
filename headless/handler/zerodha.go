package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"headless/utils"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

var ZerodhaRequestToken chan string = make(chan string)

func zAuth(w http.ResponseWriter, r *http.Request) {
	var body map[string]any

	close(ZerodhaRequestToken)
	ZerodhaRequestToken = make(chan string)

	json.NewDecoder(r.Body).Decode(&body)

	totp, err := utils.GetTOTP()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get TOTP: %v", err), http.StatusInternalServerError)
		return
	}

	body["totp"] = totp
	go getZerodhaTokenUsingBrowser(body)

	resp := make(map[string]any)
	resp["request_token"] = <-ZerodhaRequestToken
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

func zAuthenticated(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	requestToken := params.Get("request_token")
	ZerodhaRequestToken <- requestToken
	w.WriteHeader(http.StatusOK)
}

func getZerodhaTokenUsingBrowser(body map[string]any) {
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts, chromedp.DisableGPU)
	opts = append(opts, chromedp.Flag("blink-settings", "scriptEnabled=true"))
	opts = append(opts, chromedp.Flag("headless", false))
	opts = append(opts, chromedp.ExecPath("C:\\Users\\JenishJ\\Downloads\\chrome-win\\chrome.exe"))

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if err := chromedp.Run(ctx, zerodhaSendCredentials(body)); err != nil {
		log.Fatal(err)
	}
}

func zerodhaSendCredentials(body map[string]any) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(body["url"].(string)),
		chromedp.WaitVisible("div.login-form"),
		chromedp.SendKeys("input#userid", body["username"].(string)),
		chromedp.SendKeys("input#password", body["password"].(string)),
		chromedp.Click("#container > div > div > div.login-form > form > div.actions > button[type=submit]", chromedp.NodeVisible),
		chromedp.WaitVisible("input[label='External TOTP']", chromedp.ByQuery),
		chromedp.SendKeys("#userid", body["totp"].(string), chromedp.ByQuery),
		chromedp.Sleep(3 * time.Second),
	}
}
