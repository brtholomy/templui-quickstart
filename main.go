package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/axzilla/templui-quickstart/assets"
	"github.com/axzilla/templui-quickstart/ui/pages"
	"github.com/joho/godotenv"
)

func main() {
	InitDotEnv()
	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)
	mux.Handle("GET /", templ.Handler(pages.Landing()))

	// https://templ.guide/server-side-rendering/creating-an-http-server-with-templ#displaying-dynamic-data
	output := MakeRequest()
	mux.Handle("GET /qbo", templ.Handler(pages.Qbo(output)))

	mux.Handle("GET /foo", templ.Handler(pages.Foo()))
	mux.Handle("GET /file", templ.Handler(pages.File()))
	mux.Handle("GET /time", templ.Handler(pages.Time(time.Now())))
	mux.HandleFunc("GET /task/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling task with id=%v\n", id)
	})
	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}

func InitDotEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Error loading .env file")
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
