package analytics_test

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nbalakrcloud/visitor-analytics/internal/analytics"
)

// Test the RecordVisitor method
func TestVisitorStore_RecordVisitor(t *testing.T) {
	store := analytics.NewVisitorStore()

	store.RecordVisitor("http://foo.com", "Alice")
	store.RecordVisitor("http://foo.com", "Bob")
	store.RecordVisitor("http://foo.com", "Alice") // Duplicate visit
	store.RecordVisitor("http://bar.com", "Alice")

	expectedFooVisitors := 2
	expectedBarVisitors := 1

	if len(store.Store["http://foo.com"]) != expectedFooVisitors {
		t.Errorf("Expected %d unique visitors for http://foo.com, got %d", expectedFooVisitors, len(store.Store["http://foo.com"]))
	}
	if len(store.Store["http://bar.com"]) != expectedBarVisitors {
		t.Errorf("Expected %d unique visitors for http://bar.com, got %d", expectedBarVisitors, len(store.Store["http://bar.com"]))
	}
}

// Test the GetUniqueVisitors method
func TestVisitorStore_GetUniqueVisitors(t *testing.T) {
	store := analytics.NewVisitorStore()

	store.RecordVisitor("http://foo.com", "Alice")
	store.RecordVisitor("http://foo.com", "Bob")
	store.RecordVisitor("http://bar.com", "Alice")

	visitorMap := store.GetUniqueVisitors()

	expectedFooVisitors := 2
	expectedBarVisitors := 1

	if visitorMap["http://foo.com"] != expectedFooVisitors {
		t.Errorf("Expected %d unique visitors for http://foo.com, got %d", expectedFooVisitors, visitorMap["http://foo.com"])
	}
	if visitorMap["http://bar.com"] != expectedBarVisitors {
		t.Errorf("Expected %d unique visitors for http://bar.com, got %d", expectedBarVisitors, visitorMap["http://bar.com"])
	}
}

func sortOutput(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// Test the getAnalytics method
func TestVisitorStore_getAnalytics(t *testing.T) {
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

	sortedExpectedBody := sortOutput(expectedBody)
	sortedResponseBody := sortOutput(rr.Body.String())

	if sortedResponseBody != sortedExpectedBody {
		t.Errorf("Handler returned unexpected body: got %v want %v", sortedResponseBody, sortedExpectedBody)
	}
}
