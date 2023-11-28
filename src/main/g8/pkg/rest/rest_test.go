package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gympass/gcore/v3/glog"
	uuid "github.com/gofrs/uuid"
)

/*
Profiler endpoints:
http://localhost:7070/debug/pprof

curl -sK -v http://localhost:7070/debug/pprof/heap > /tmp/heap.out

go test -v ./internal/rest/
go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out -v ./internal/rest/
go tool pprof cpuprofile
go tool pprof memprofile
*/

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a default Rest configuration with logrus:info and configs/dev.yaml
	rm := New(Config{Service: "test-service", Logger: glog.Noop()})

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rm.Health)

	// Handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body format is what we expect
	var hc = struct {
		ID      uuid.UUID `json:"id"`
		Service string    `json:"service"`
	}{}
	err = json.NewDecoder(rr.Body).Decode(&hc)
	if err != nil {
		t.Errorf("handler returned wrong json schema: got %v",
			rr.Body.String())
	}
}
