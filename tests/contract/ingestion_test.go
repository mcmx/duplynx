package contract_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mcmx/duplynx/internal/ingestion"
)

func TestIngestionRejectsInvalidSignature(t *testing.T) {
	h := ingestion.Handler{TenantSecrets: map[string]string{"tenant-a": "secret"}}
	req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader("{}"))
	req.Header.Set("X-Duplynx-Tenant", "tenant-a")
	req.Header.Set("X-Duplynx-Signature", hex.EncodeToString([]byte("bogus")))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestIngestionAcceptsValidSignature(t *testing.T) {
	payload := "{}"
	secret := "secret"
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	h := ingestion.Handler{TenantSecrets: map[string]string{"tenant-a": secret}}
	req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader(payload))
	req.Header.Set("X-Duplynx-Tenant", "tenant-a")
	req.Header.Set("X-Duplynx-Signature", signature)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}
}
