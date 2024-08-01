package analytics_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nbalakrcloud/visitor-analytics/internal/analytics"
)

// Test the /track endpoint
func TestAnalyticsService_TrackVisitor(t *testing.T) {
	store := analytics.NewVisitorStore()
	service := &analytics.AnalyticsService{VisitorStore: *store}
	router := mux.NewRouter()
	service.WireUpAnalytics(router)

	req, err := http.NewRequest("GET", "/track?url=http://foo.com&visitorID=Alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if len(store.Store["http://foo.com"]) != 1 {
		t.Errorf("Expected 1 unique visitor for http://foo.com, got %d", len(store.Store["http://foo.com"]))
	}
}

// Test the /analytics endpoint
func TestAnalyticsService_GetAnalytics(t *testing.T) {
	store := analytics.NewVisitorStore()
	service := &analytics.AnalyticsService{VisitorStore: *store}
	router := mux.NewRouter()
	service.WireUpAnalytics(router)

	store.RecordVisitor("http://foo.com", "Alice")
	store.RecordVisitor("http://foo.com", "Bob")
	store.RecordVisitor("http://bar.com", "Alice")

	req, err := http.NewRequest("GET", "/analytics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	expectedBody := "http://foo.com:2\nhttp://bar.com:1\n"

	// if rr.Body.String() != expectedBody {
	// 	t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	// }
	sortedExpectedBody := sortOutput(expectedBody)
	sortedResponseBody := sortOutput(rr.Body.String())

	if sortedResponseBody != sortedExpectedBody {
		t.Errorf("Handler returned unexpected body: got %v want %v", sortedResponseBody, sortedExpectedBody)
	}
}
