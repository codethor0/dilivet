// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	mldsa "github.com/codethor0/dilivet/code/clean"
	"github.com/codethor0/dilivet/code/clean/kats"
	"github.com/codethor0/dilivet/code/diag"
)

var version = "dev"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Load security configuration
	secCfg := loadSecurityConfig()

	// Build middleware chain
	var handler http.Handler = http.NewServeMux()
	mux := handler.(*http.ServeMux)

	// Apply security middleware (order matters: outermost first)
	handler = corsMiddleware(secCfg.allowedOrigins)(handler)
	handler = authMiddleware(secCfg.requireAuth, secCfg.authToken)(handler)
	handler = timeoutMiddleware(secCfg.requestTimeout)(handler)
	handler = maxBodySizeMiddleware(secCfg.maxBodySize)(handler)

	// Register routes
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/verify", handleVerify)
	mux.HandleFunc("/api/kat-verify", handleKATVerify)

	// Serve static files from web/ui/dist if they exist
	staticDir := "./web/ui/dist"
	if _, err := os.Stat(staticDir); err == nil {
		// Serve static assets
		mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir+"/assets"))))
		mux.Handle("/dilivet-logo.png", http.FileServer(http.Dir(staticDir)))

		// SPA fallback: serve index.html for all non-API routes
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Don't handle API routes
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			// Serve index.html for SPA routing
			http.ServeFile(w, r, staticDir+"/index.html")
		})
	}

	addr := ":" + port
	log.Printf("DiliVet Web Server starting on %s", addr)
	log.Printf("Version: %s", version)
	log.Printf("Security: auth=%v cors=%v maxBodySize=%d timeout=%v",
		secCfg.requireAuth, len(secCfg.allowedOrigins) > 0, secCfg.maxBodySize, secCfg.requestTimeout)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version,
	})
}

type verifyRequest struct {
	ParamSet     string `json:"paramSet"`
	PublicKeyHex string `json:"publicKeyHex"`
	SignatureHex string `json:"signatureHex"`
	MessageHex   string `json:"messageHex,omitempty"`
	Message      string `json:"message,omitempty"`
}

type verifyResponse struct {
	OK     bool   `json:"ok"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req verifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logSecurityEvent("verify_error", "/api/verify", "invalid_json")
		respondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate paramSet
	paramSet := strings.ToUpper(req.ParamSet)
	if paramSet != "ML-DSA-44" && paramSet != "ML-DSA-65" && paramSet != "ML-DSA-87" {
		respondError(w, http.StatusBadRequest, "paramSet must be ML-DSA-44, ML-DSA-65, or ML-DSA-87")
		return
	}

	// Decode public key
	pub, err := decodeHex(req.PublicKeyHex, "public key")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Decode signature
	sig, err := decodeHex(req.SignatureHex, "signature")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Decode message (either hex or UTF-8)
	var msg []byte
	if req.MessageHex != "" {
		msg, err = decodeHex(req.MessageHex, "message")
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	} else if req.Message != "" {
		msg = []byte(req.Message)
	} else {
		respondError(w, http.StatusBadRequest, "either messageHex or message must be provided")
		return
	}

	// Verify signature
	valid, verr := mldsa.Verify(pub, msg, sig)
	if verr != nil {
		// Log verification failure (sanitized)
		logSecurityEvent("verify_failure", "/api/verify", sanitizeError(verr))
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Verification error: %v", verr))
		return
	}

	// Log successful verification (metadata only)
	logSecurityEvent("verify_success", "/api/verify", fmt.Sprintf("paramSet=%s valid=%v", paramSet, valid))

	w.Header().Set("Content-Type", "application/json")
	if valid {
		json.NewEncoder(w).Encode(verifyResponse{
			OK:     true,
			Result: "valid",
		})
	} else {
		json.NewEncoder(w).Encode(verifyResponse{
			OK:     true,
			Result: "invalid",
		})
	}
}

type katVerifyRequest struct {
	VectorsPath string `json:"vectorsPath,omitempty"`
}

type katVerifyResponse struct {
	OK             bool              `json:"ok"`
	TotalVectors   int               `json:"totalVectors,omitempty"`
	Passed         int               `json:"passed,omitempty"`
	Failed         int               `json:"failed,omitempty"`
	DecodeFailures int               `json:"decodeFailures,omitempty"`
	Error          string            `json:"error,omitempty"`
	Details        []katVerifyDetail `json:"details,omitempty"`
}

type katVerifyDetail struct {
	CaseID       int    `json:"caseId"`
	Passed       bool   `json:"passed"`
	ParameterSet string `json:"parameterSet,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

func handleKATVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req katVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Empty body is OK, use default path
		req.VectorsPath = ""
	}

	vectorsPath := req.VectorsPath
	if vectorsPath == "" {
		vectorsPath = "code/clean/testdata/kats/ml-dsa/ML-DSA-sigVer-FIPS204-internalProjection.json"
	}

	vectors, err := kats.LoadSigVerVectors(vectorsPath)
	if err != nil {
		logSecurityEvent("kat_error", "/api/kat-verify", sanitizeError(err))
		respondError(w, http.StatusInternalServerError, "Failed to load test vectors")
		return
	}

	logSecurityEvent("kat_start", "/api/kat-verify", fmt.Sprintf("vectorsPath=%s", vectorsPath))

	report := &diag.Report{}
	var details []katVerifyDetail

	for _, tg := range vectors.TestGroups {
		for _, tc := range tg.Tests {
			report.TotalTests++

			pk, err := hex.DecodeString(tc.Public)
			if err != nil {
				report.DecodeFailures++
				details = append(details, katVerifyDetail{
					CaseID:       tc.CaseID,
					Passed:       false,
					ParameterSet: tg.ParameterSet,
					Reason:       fmt.Sprintf("Failed to decode public key: %v", err),
				})
				continue
			}

			sig, err := hex.DecodeString(tc.Signature)
			if err != nil {
				report.DecodeFailures++
				details = append(details, katVerifyDetail{
					CaseID:       tc.CaseID,
					Passed:       false,
					ParameterSet: tg.ParameterSet,
					Reason:       fmt.Sprintf("Failed to decode signature: %v", err),
				})
				continue
			}

			msg, err := hex.DecodeString(tc.Message)
			if err != nil {
				report.DecodeFailures++
				details = append(details, katVerifyDetail{
					CaseID:       tc.CaseID,
					Passed:       false,
					ParameterSet: tg.ParameterSet,
					Reason:       fmt.Sprintf("Failed to decode message: %v", err),
				})
				continue
			}

			ok, verr := mldsa.Verify(pk, msg, sig)
			if verr == nil && ok {
				report.StrictPasses++
				details = append(details, katVerifyDetail{
					CaseID:       tc.CaseID,
					Passed:       true,
					ParameterSet: tg.ParameterSet,
				})
			} else {
				report.StructuralFailures++
				reason := "Signature verification failed"
				if verr != nil {
					reason = fmt.Sprintf("Verification error: %v", verr)
				}
				details = append(details, katVerifyDetail{
					CaseID:       tc.CaseID,
					Passed:       false,
					ParameterSet: tg.ParameterSet,
					Reason:       reason,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(katVerifyResponse{
		OK:             true,
		TotalVectors:   report.TotalTests,
		Passed:         report.StrictPasses,
		Failed:         report.StructuralFailures + report.DecodeFailures,
		DecodeFailures: report.DecodeFailures,
		Details:        details,
	})

	// Log KAT completion (metadata only)
	logSecurityEvent("kat_complete", "/api/kat-verify",
		fmt.Sprintf("total=%d passed=%d failed=%d", report.TotalTests, report.StrictPasses, report.StructuralFailures+report.DecodeFailures))
}

func decodeHex(s, fieldName string) ([]byte, error) {
	clean := stripWhitespace(s)
	if clean == "" {
		return nil, fmt.Errorf("%s cannot be empty", fieldName)
	}
	// Validate hex length (reasonable bounds)
	if len(clean) > 100000 { // ~50KB hex = ~25KB binary
		return nil, fmt.Errorf("%s too long (max 100000 hex chars)", fieldName)
	}
	buf, err := hex.DecodeString(clean)
	if err != nil {
		return nil, fmt.Errorf("invalid hex in %s: %w", fieldName, err)
	}
	return buf, nil
}

func stripWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if !isWhitespace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(verifyResponse{
		OK:    false,
		Error: message,
	})
}
