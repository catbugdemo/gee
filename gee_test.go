package gee

import (
	"net/http"
	"testing"
)

func TestGee(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		r := New()
		r.GET("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test"))
		})

		r.Run(":8000")
	})
}
