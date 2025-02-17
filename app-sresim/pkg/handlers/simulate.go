package handlers

import (
	"fmt"
	"net/http"
)

// SimulateHandler responds with a basic success message.
func SimulateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Request processed successfully!")
}
