package chiprometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Test_Logger(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	m := NewMiddleware("test")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/users/JoeBob", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequest("GET", "http://localhost:3000/users/Misty", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)
	body := recorder.Body.String()
	if !strings.Contains(body, reqsName) {
		t.Errorf("body does not contain request total entry '%s'", reqsName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", latencyName)
	}

	req1Count := `chi_request_duration_seconds_count{code="200",method="GET",path="/ok",service="test"} 1`
	req2Count := `chi_request_duration_seconds_count{code="200",method="GET",path="/users/JoeBob",service="test"} 1`
	req3Count := `chi_request_duration_seconds_count{code="200",method="GET",path="/users/Misty",service="test"} 1`

	if !strings.Contains(body, req1Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req1Count)
	}
	if !strings.Contains(body, req2Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req2Count)
	}
	if !strings.Contains(body, req3Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req3Count)
	}
}

func Test_PatternLogger(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	m := NewPatternMiddleware("patternOnlyTest")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/users/JoeBob", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequest("GET", "http://localhost:3000/users/Misty", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)

	body := recorder.Body.String()

	if !strings.Contains(body, patternReqsName) {
		t.Errorf("body does not contain request total entry '%s'", patternReqsName)
	}
	if !strings.Contains(body, patternLatencyName) {
		t.Errorf("body does not contain request duration entry '%s'", patternLatencyName)
	}

	req1Count := `chi_pattern_request_duration_seconds_count{code="200",method="GET",path="/ok",service="patternOnlyTest"} 1`
	joeBobCount := `chi_pattern_request_duration_seconds_count{code="200",method="GET",path="/users/JoeBob",service="patternOnlyTest"} 1`
	mistyCount := `chi_pattern_request_duration_seconds_count{code="200",method="GET",path="/users/Misty",service="patternOnlyTest"} 1`
	firstNamePatternCount := `chi_pattern_request_duration_seconds_count{code="200",method="GET",path="/users/{firstName}",service="patternOnlyTest"} 2`

	if !strings.Contains(body, req1Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req1Count)
	}
	if strings.Contains(body, joeBobCount) {
		t.Errorf("body should not contain Joe Bob count summary '%s'", joeBobCount)
	}
	if strings.Contains(body, mistyCount) {
		t.Errorf("body should not contain Misty count summary '%s'", mistyCount)
	}
	if !strings.Contains(body, firstNamePatternCount) {
		t.Errorf("body does not contain first name pattern count summary '%s'", firstNamePatternCount)
	}
}

func Test_MultipleLoggers(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	mid := NewMiddleware("pathTest")
	m := NewPatternMiddleware("patternTest")

	n.Use(mid)
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/users/JoeBob", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequest("GET", "http://localhost:3000/users/Misty", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)

	body := recorder.Body.String()

	if !strings.Contains(body, patternReqsName) {
		t.Errorf("body does not contain request total entry '%s'", patternReqsName)
	}
	if !strings.Contains(body, patternLatencyName) {
		t.Errorf("body does not contain request duration entry '%s'", patternLatencyName)
	}

	req1Count := `chi_pattern_request_duration_seconds_count{code="200",method="GET",path="/ok",service="patternTest"} 1`
	joeBobCount := `chi_request_duration_seconds_count{code="200",method="GET",path="/users/JoeBob",service="pathTest"} 1`
	mistyCount := `chi_request_duration_seconds_count{code="200",method="GET",path="/users/Misty",service="pathTest"} 1`
	firstNamePatternCount := `chi_pattern_request_duration_seconds_count{code="200",method="GET",path="/users/{firstName}",service="patternTest"} 2`

	if !strings.Contains(body, req1Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req1Count)
	}
	if !strings.Contains(body, joeBobCount) {
		t.Errorf("body does not contain Joe Bob count summary '%s'", joeBobCount)
	}
	if !strings.Contains(body, mistyCount) {
		t.Errorf("body does not contain Misty count summary '%s'", mistyCount)
	}
	if !strings.Contains(body, firstNamePatternCount) {
		t.Errorf("body does not contain first name pattern count summary '%s'", firstNamePatternCount)
	}
}

func Test_NumericStatusCodes(t *testing.T) {
	// Reset prometheus registry for clean test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.DefaultRegisterer.(prometheus.Gatherer)

	n := chi.NewRouter()
	m := NewMiddleware("statusTest")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	
	// Handler for 200 OK
	n.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	
	// Handler for 201 Created
	n.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "created")
	})
	
	// Handler for 404 Not Found
	n.Get("/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found")
	})
	
	// Handler for 400 Bad Request
	n.Get("/bad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad request")
	})
	
	// Handler for 500 Internal Server Error
	n.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal server error")
	})
	
	// Handler for 401 Unauthorized
	n.Get("/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "unauthorized")
	})

	// Make requests to each endpoint
	testCases := []struct {
		path       string
		method     string
		statusCode int
	}{
		{"/ok", "GET", 200},
		{"/create", "POST", 201},
		{"/notfound", "GET", 404},
		{"/bad", "GET", 400},
		{"/error", "GET", 500},
		{"/unauthorized", "GET", 401},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, "http://localhost:3000"+tc.path, nil)
		if err != nil {
			t.Error(err)
		}
		recorder := httptest.NewRecorder()
		n.ServeHTTP(recorder, req)
	}

	// Get metrics
	req, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	n.ServeHTTP(recorder, req)
	body := recorder.Body.String()

	// Verify each status code appears as a number in metrics
	statusCodeChecks := []struct {
		statusCode string
		path       string
		method     string
	}{
		{"200", "/ok", "GET"},
		{"201", "/create", "POST"},
		{"404", "/notfound", "GET"},
		{"400", "/bad", "GET"},
		{"500", "/error", "GET"},
		{"401", "/unauthorized", "GET"},
	}

	for _, check := range statusCodeChecks {
		// Check in request total metric
		expectedReqMetric := fmt.Sprintf(`chi_requests_total{code="%s",method="%s",path="%s",service="statusTest"} 1`,
			check.statusCode, check.method, check.path)
		if !strings.Contains(body, expectedReqMetric) {
			t.Errorf("body does not contain expected request metric: %s", expectedReqMetric)
		}
		
		// Check in duration metric
		expectedDurationMetric := fmt.Sprintf(`chi_request_duration_seconds_count{code="%s",method="%s",path="%s",service="statusTest"} 1`,
			check.statusCode, check.method, check.path)
		if !strings.Contains(body, expectedDurationMetric) {
			t.Errorf("body does not contain expected duration metric: %s", expectedDurationMetric)
		}
		
		// Ensure old text status values are NOT present
		oldTextStatuses := map[string]string{
			"200": "OK",
			"201": "Created",
			"404": "Not Found",
			"400": "Bad Request",
			"500": "Internal Server Error",
			"401": "Unauthorized",
		}
		
		if oldText, exists := oldTextStatuses[check.statusCode]; exists {
			badMetric := fmt.Sprintf(`code="%s"`, oldText)
			if strings.Contains(body, badMetric) {
				t.Errorf("body should not contain text status code: %s", badMetric)
			}
		}
	}
}

func Test_PatternMiddlewareNumericStatusCodes(t *testing.T) {
	// Reset prometheus registry for clean test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.DefaultRegisterer.(prometheus.Gatherer)

	n := chi.NewRouter()
	m := NewPatternMiddleware("patternStatusTest")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	
	// Handler for 200 OK with pattern
	n.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "user")
	})
	
	// Handler for 404 Not Found with pattern
	n.Get("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found")
	})
	
	// Handler for 500 error with pattern
	n.Get("/error/{type}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error")
	})

	// Make requests
	testCases := []struct {
		path       string
		statusCode int
	}{
		{"/users/123", 200},
		{"/users/456", 200},
		{"/items/789", 404},
		{"/items/999", 404},
		{"/error/db", 500},
		{"/error/network", 500},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", "http://localhost:3000"+tc.path, nil)
		if err != nil {
			t.Error(err)
		}
		recorder := httptest.NewRecorder()
		n.ServeHTTP(recorder, req)
	}

	// Get metrics
	req, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	n.ServeHTTP(recorder, req)
	body := recorder.Body.String()

	// Verify patterns are grouped correctly with numeric status codes
	expectedMetrics := []string{
		`chi_pattern_requests_total{code="200",method="GET",path="/users/{id}",service="patternStatusTest"} 2`,
		`chi_pattern_requests_total{code="404",method="GET",path="/items/{id}",service="patternStatusTest"} 2`,
		`chi_pattern_requests_total{code="500",method="GET",path="/error/{type}",service="patternStatusTest"} 2`,
		`chi_pattern_request_duration_seconds_count{code="200",method="GET",path="/users/{id}",service="patternStatusTest"} 2`,
		`chi_pattern_request_duration_seconds_count{code="404",method="GET",path="/items/{id}",service="patternStatusTest"} 2`,
		`chi_pattern_request_duration_seconds_count{code="500",method="GET",path="/error/{type}",service="patternStatusTest"} 2`,
	}

	for _, expectedMetric := range expectedMetrics {
		if !strings.Contains(body, expectedMetric) {
			t.Errorf("body does not contain expected metric: %s", expectedMetric)
		}
	}
	
	// Ensure text status codes are not present
	badTexts := []string{`code="OK"`, `code="Not Found"`, `code="Internal Server Error"`}
	for _, badText := range badTexts {
		if strings.Contains(body, badText) {
			t.Errorf("body should not contain text status: %s", badText)
		}
	}
}