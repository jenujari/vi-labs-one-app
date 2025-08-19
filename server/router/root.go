package router

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"server/config"
	"server/helpers"
	"server/tpl"
)

var (
	server *http.Server
	router *http.ServeMux
)

func init() {

	cfg := config.GetConfig()

	server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		MaxHeaderBytes:    0,
	}

	router = http.NewServeMux()

	FileServer(router, "/assets/")
	SetRoutes(router)

	server.Handler = router
	config.GetLogger().Println("server initialization complete.")
}

func RunServer(ctx *helpers.ProcessContext) {
	defer ctx.CompleteOneWorker()

	go func(cmdx *helpers.ProcessContext) {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			cmdx.FatalErrorChan <- fmt.Errorf("ListenAndServe(): %v", err)
		}
	}(ctx)

	// helper.HitBrowser("http://localhost:5456", j)

	<-ctx.CTX.Done()
	config.GetLogger().Println("shutting down server...")
	if err := server.Shutdown(ctx.CTX); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	config.GetLogger().Println("server shutdown complete...")
}

func GetServer() *http.Server {
	return server
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r *http.ServeMux, path string) {
	sub, err := fs.Sub(tpl.GetAssetsFs(), "assets")
	if err != nil {
		panic(err)
	}

	fsx := http.FileServer(http.FS(sub))
	r.Handle("GET "+path, http.StripPrefix(path, fsx))
}
