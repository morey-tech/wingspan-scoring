package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware_LogsRequest(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr) // Reset to default after test

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap with logging middleware
	handler := loggingMiddleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test-path", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	// Execute request
	handler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test response", w.Body.String())

	// Verify log output
	logOutput := buf.String()
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test-path")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "127.0.0.1:12345")
}

func TestLoggingMiddleware_LogsStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"Not Found", http.StatusNotFound},
		{"Internal Server Error", http.StatusInternalServerError},
		{"Created", http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			// Create test handler that returns specific status code
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			// Wrap with logging middleware
			handler := loggingMiddleware(testHandler)

			// Create and execute test request
			req := httptest.NewRequest("POST", "/api/test", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Verify status code in response
			assert.Equal(t, tt.statusCode, w.Code)

			// Verify status code in log
			logOutput := buf.String()
			assert.Contains(t, logOutput, "POST")
			assert.Contains(t, logOutput, "/api/test")
			// Convert status code to string and check it's in the log
			statusStr := string('0' + byte(tt.statusCode/100))
			assert.Contains(t, logOutput, statusStr)
		})
	}
}

func TestLoggingMiddleware_LogsDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with logging middleware
			handler := loggingMiddleware(testHandler)

			// Create and execute test request
			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Verify method in log
			logOutput := buf.String()
			assert.Contains(t, logOutput, method)
		})
	}
}

func TestLoggingMiddleware_DoesNotInterfereWithHandler(t *testing.T) {
	// Create a handler that writes specific data
	expectedBody := "Hello, World!"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(expectedBody))
	})

	// Wrap with logging middleware
	handler := loggingMiddleware(testHandler)

	// Create and execute test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify the handler behavior is preserved
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

func TestResponseWriter_CapturesStatusCode(t *testing.T) {
	// Create a response recorder
	baseWriter := httptest.NewRecorder()

	// Wrap it
	wrapped := &responseWriter{
		ResponseWriter: baseWriter,
		statusCode:     http.StatusOK,
	}

	// Write a status code
	wrapped.WriteHeader(http.StatusNotFound)

	// Verify it was captured
	assert.Equal(t, http.StatusNotFound, wrapped.statusCode)
	assert.Equal(t, http.StatusNotFound, baseWriter.Code)
}

func TestResponseWriter_DefaultStatusCode(t *testing.T) {
	// Create a wrapped response writer
	baseWriter := httptest.NewRecorder()
	wrapped := &responseWriter{
		ResponseWriter: baseWriter,
		statusCode:     http.StatusOK,
	}

	// Write data without explicitly setting status code
	wrapped.Write([]byte("test"))

	// Verify default status code is preserved
	assert.Equal(t, http.StatusOK, wrapped.statusCode)
}

func TestLoggingMiddleware_LogsRequestDuration(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with logging middleware
	handler := loggingMiddleware(testHandler)

	// Create and execute test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify duration is logged (should contain time units like µs, ms, or s)
	logOutput := buf.String()
	hasDuration := strings.Contains(logOutput, "µs") ||
		strings.Contains(logOutput, "ms") ||
		strings.Contains(logOutput, "s")
	assert.True(t, hasDuration, "Log should contain request duration")
}
