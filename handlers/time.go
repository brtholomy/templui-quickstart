package handlers

import (
	"net/http"
	"time"

	"github.com/brtholomy/templui-quickstart/ui/pages"
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
