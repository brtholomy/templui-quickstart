package main

import (
	"net/http"

	"github.com/brtholomy/templui-quickstart/ui/pages"
)

type GlobalState struct {
	Count int
}

var global GlobalState

func CounterGetHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.Counter(global.Count, 0)
	component.Render(r.Context(), w)
}

func CounterPostHandler(w http.ResponseWriter, r *http.Request) {
	// Update state.
	r.ParseForm()

	// Check to see if the global button was pressed.
	if r.Form.Has("global") {
		global.Count++
	}
	//TODO: Update session.

	// Display the form.
	CounterGetHandler(w, r)
}
