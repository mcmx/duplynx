package perf_test

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

func BenchmarkIngestionAcknowledgement(b *testing.B) {
	payload := "{}"
	secret := "secret"
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	h := ingestion.Handler{TenantSecrets: map[string]string{"tenant-a": secret}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader(payload))
		req.Header.Set("X-Duplynx-Tenant", "tenant-a")
		req.Header.Set("X-Duplynx-Signature", signature)

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			b.Fatalf("expected 202, got %d", rec.Code)
		}
	}
}
