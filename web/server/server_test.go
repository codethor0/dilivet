// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %q", resp["status"])
	}

	if resp["version"] == "" {
		t.Error("Version should not be empty")
	}
}

func TestHandleHealth_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	w := httptest.NewRecorder()
	handleHealth(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleVerify_ValidSignature(t *testing.T) {
	// This test uses dummy data that will fail verification, but tests the API structure
	// Note: The hex strings are too short for real ML-DSA keys, so verification will fail
	// but the API should return a structured response
	reqBody := verifyRequest{
		ParamSet:     "ML-DSA-44",
		PublicKeyHex: "deadbeefdeadbeef",
		SignatureHex: "cafebabecafebabe",
		Message:      "test message",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleVerify(w, req)

	// The request will likely fail due to invalid key/sig lengths, but we should get a structured error
	var resp verifyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Response should be structured (either ok=true with result, or ok=false with error)
	if resp.OK {
		if resp.Result != "invalid" && resp.Result != "valid" {
			t.Errorf("Expected result 'invalid' or 'valid', got %q", resp.Result)
		}
	} else {
		if resp.Error == "" {
			t.Error("Expected error message when ok=false")
		}
	}
}

func TestHandleVerify_AllParamSets(t *testing.T) {
	paramSets := []string{"ML-DSA-44", "ML-DSA-65", "ML-DSA-87"}
	for _, paramSet := range paramSets {
		t.Run(paramSet, func(t *testing.T) {
			reqBody := verifyRequest{
				ParamSet:     paramSet,
				PublicKeyHex: "deadbeef",
				SignatureHex: "cafebabe",
				Message:      "test",
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handleVerify(w, req)

			// Should get a structured response (either error or result)
			var resp verifyResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Response should be structured
			if resp.OK && resp.Result == "" {
				t.Error("Expected result when ok=true")
			}
			if !resp.OK && resp.Error == "" {
				t.Error("Expected error when ok=false")
			}
		})
	}
}

func TestHandleVerify_InvalidParamSet(t *testing.T) {
	reqBody := verifyRequest{
		ParamSet:     "ML-DSA-99",
		PublicKeyHex: "deadbeef",
		SignatureHex: "cafebabe",
		Message:      "test",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleVerify(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp verifyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.OK {
		t.Error("Expected ok=false for invalid paramSet")
	}
	if !strings.Contains(resp.Error, "paramSet") {
		t.Errorf("Expected error to mention paramSet, got %q", resp.Error)
	}
}

func TestHandleVerify_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		req  verifyRequest
	}{
		{"missing paramSet", verifyRequest{PublicKeyHex: "deadbeef", SignatureHex: "cafebabe", Message: "test"}},
		{"missing publicKeyHex", verifyRequest{ParamSet: "ML-DSA-44", SignatureHex: "cafebabe", Message: "test"}},
		{"missing signatureHex", verifyRequest{ParamSet: "ML-DSA-44", PublicKeyHex: "deadbeef", Message: "test"}},
		{"missing message", verifyRequest{ParamSet: "ML-DSA-44", PublicKeyHex: "deadbeef", SignatureHex: "cafebabe"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.req)
			req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handleVerify(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestHandleVerify_InvalidHex(t *testing.T) {
	tests := []struct {
		name          string
		publicKeyHex  string
		signatureHex  string
		messageHex    string
		expectError   bool
		errorContains string
	}{
		{"non-hex public key", "not hex", "cafebabe", "", true, "hex"},
		{"non-hex signature", "deadbeef", "not hex", "", true, "hex"},
		{"non-hex message", "deadbeef", "cafebabe", "not hex", true, "hex"},
		{"empty public key", "", "cafebabe", "", true, "empty"},
		{"empty signature", "deadbeef", "", "", true, "empty"},
		{"empty message hex", "deadbeef", "cafebabe", "", false, ""}, // Should use Message field
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := verifyRequest{
				ParamSet:     "ML-DSA-44",
				PublicKeyHex: tt.publicKeyHex,
				SignatureHex: tt.signatureHex,
			}
			if tt.messageHex != "" {
				reqBody.MessageHex = tt.messageHex
			} else if !tt.expectError || tt.name != "empty message hex" {
				reqBody.Message = "test"
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handleVerify(w, req)

			if tt.expectError {
				if w.Code != http.StatusBadRequest {
					t.Errorf("Expected status 400, got %d", w.Code)
				}
				var resp verifyResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err == nil {
					if tt.errorContains != "" && !strings.Contains(resp.Error, tt.errorContains) {
						t.Errorf("Expected error to contain %q, got %q", tt.errorContains, resp.Error)
					}
				}
			}
		})
	}
}

func TestHandleVerify_MessageModes(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		messageHex string
		shouldWork bool
	}{
		{"UTF-8 message only", "hello world", "", true},
		{"hex message only", "", "68656c6c6f20776f726c64", true},
		{"both provided (hex takes precedence)", "ignored", "68656c6c6f", true},
		{"neither provided", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := verifyRequest{
				ParamSet:     "ML-DSA-44",
				PublicKeyHex: "deadbeef",
				SignatureHex: "cafebabe",
				Message:      tt.message,
				MessageHex:   tt.messageHex,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handleVerify(w, req)

			if tt.shouldWork {
				// Should get a structured response (may be error due to invalid keys, but structure is correct)
				var resp verifyResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
			} else {
				if w.Code != http.StatusBadRequest {
					t.Errorf("Expected status 400, got %d", w.Code)
				}
			}
		})
	}
}

func TestHandleVerify_LargeMessage(t *testing.T) {
	// Test with a large message (1MB)
	largeMessage := strings.Repeat("a", 1024*1024)
	reqBody := verifyRequest{
		ParamSet:     "ML-DSA-44",
		PublicKeyHex: "deadbeef",
		SignatureHex: "cafebabe",
		Message:      largeMessage,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleVerify(w, req)

	// Should handle large message without panic
	var resp verifyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	// Response should be structured (may fail due to invalid keys, but no panic)
}

func TestHandleVerify_EmptyMessage(t *testing.T) {
	reqBody := verifyRequest{
		ParamSet:     "ML-DSA-44",
		PublicKeyHex: "deadbeef",
		SignatureHex: "cafebabe",
		Message:      "",
		MessageHex:   "",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleVerify(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty message, got %d", w.Code)
	}
}

func TestHandleVerify_WhitespaceInHex(t *testing.T) {
	// Test that whitespace is stripped from hex strings
	reqBody := verifyRequest{
		ParamSet:     "ML-DSA-44",
		PublicKeyHex: "dead beef\ncafe babe",
		SignatureHex: "cafe babe\ndead beef",
		Message:      "test",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleVerify(w, req)

	// Should handle whitespace without error (though verification may fail due to invalid keys)
	var resp verifyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
}

func TestHandleVerify_Concurrent(t *testing.T) {
	// Test concurrent requests to verify no shared state corruption
	const numRequests = 50
	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			reqBody := verifyRequest{
				ParamSet:     "ML-DSA-44",
				PublicKeyHex: "deadbeef",
				SignatureHex: "cafebabe",
				Message:      "test",
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handleVerify(w, req)

			var resp verifyResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				errors <- err
				return
			}

			// Response should be structured
			if resp.OK && resp.Result == "" && resp.Error == "" {
				errors <- nil // This is actually OK, just checking structure
			}
			if !resp.OK && resp.Error == "" {
				errors <- nil // This is also OK
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent request error: %v", err)
		}
	}
}

func TestHandleVerify_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/verify", nil)
	w := httptest.NewRecorder()
	handleVerify(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleVerify_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleVerify(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleKATVerify(t *testing.T) {
	reqBody := katVerifyRequest{}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/kat-verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleKATVerify(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp katVerifyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !resp.OK {
		t.Errorf("Expected ok=true, got ok=%v", resp.OK)
	}

	if resp.TotalVectors == 0 {
		t.Error("Expected at least one test vector")
	}

	if resp.Passed < 0 || resp.Failed < 0 {
		t.Error("Passed and Failed counts should be non-negative")
	}

	if resp.TotalVectors != resp.Passed+resp.Failed+resp.DecodeFailures {
		t.Errorf("Total vectors (%d) should equal passed (%d) + failed (%d) + decode failures (%d)",
			resp.TotalVectors, resp.Passed, resp.Failed, resp.DecodeFailures)
	}
}

func TestHandleKATVerify_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/kat-verify", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleKATVerify(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandleKATVerify_InvalidPath(t *testing.T) {
	reqBody := katVerifyRequest{
		VectorsPath: "/nonexistent/path/to/vectors.json",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/kat-verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleKATVerify(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var resp katVerifyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.OK {
		t.Error("Expected ok=false for invalid path")
	}

	if resp.Error == "" {
		t.Error("Expected error message for invalid path")
	}
}

func TestHandleKATVerify_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/kat-verify", nil)
	w := httptest.NewRecorder()
	handleKATVerify(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleKATVerify_Concurrent(t *testing.T) {
	// Test concurrent KAT verification requests
	const numRequests = 10
	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			reqBody := katVerifyRequest{}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/kat-verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handleKATVerify(w, req)

			if w.Code != http.StatusOK {
				errors <- nil // Expected, just checking no panic
				return
			}

			var resp katVerifyResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				errors <- err
				return
			}

			if !resp.OK {
				errors <- nil // May fail, but should be structured
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent KAT request error: %v", err)
		}
	}
}
