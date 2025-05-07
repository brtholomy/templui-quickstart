package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/brtholomy/templui-quickstart/assets"
	"github.com/brtholomy/templui-quickstart/ui/pages"
	"github.com/joho/godotenv"
)

// convenience allocation function.
func NewNowHandler(now func() time.Time) NowHandler {
	return NowHandler{Now: now}
}

// the type which will implement the interface.
type NowHandler struct {
	Now func() time.Time
}

// implement the HTTP handler interface
// https://pkg.go.dev/net/http#Handler
func (nh NowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages.Time(nh.Now()).Render(r.Context(), w)
}

func main() {
	InitDotEnv()
	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)
	mux.Handle("GET /", templ.Handler(pages.Landing()))

	// dynamic:
	// https://templ.guide/server-side-rendering/creating-an-http-server-with-templ#displaying-dynamic-data
	mux.Handle("GET /qbo", NewQboHandler(QboRequest))

	// just for clarity's sake, the New* factory method thing is not necessary:
	nh := NowHandler{Now: time.Now}
	mux.Handle("GET /now", nh)

	mux.Handle("GET /time", templ.Handler(pages.Time(time.Now())))

	mux.Handle("GET /foo", templ.Handler(pages.Foo()))
	mux.Handle("GET /file", templ.Handler(pages.File()))

	mux.HandleFunc("GET /task/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling task with id=%v\n", id)
	})
	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}

func InitDotEnv() {
	err := godotenv.Load("/etc/secrets/.env")
	if err != nil {
		fmt.Println("Error loading prod .env file")
	}
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading local .env file")
	}
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDevelopment = os.Getenv("GO_ENV") != "production"

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))
}
