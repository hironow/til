package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_indexHandler(t *testing.T) {
	tests := []struct {
		name string
		path string
		wantStatus int
		wantBody string
	}{
		{"ok", "/", http.StatusOK, "Hello, World!"},
		{"404", "/404", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(indexHandler)
			handler.ServeHTTP(rr, req)

			if got, want := rr.Code, tt.wantStatus; got != want {
				t.Errorf("unexpected status: got: %+v, want: %+v", got, want)
			}

			if tt.wantBody != "" {
				if got, want := rr.Body.String(), tt.wantBody; got != want {
					t.Errorf("unexpected body: got: %+v, want: %+v", got, want)
				}
			}
		})
	}
}
