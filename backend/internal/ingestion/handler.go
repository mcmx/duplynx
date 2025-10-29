package ingestion

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
)

// Handler validates signed scan manifests before queuing processing work.
type Handler struct {
	TenantSecrets map[string]string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Duplynx-Tenant")
	if tenant == "" {
		http.Error(w, "missing tenant header", http.StatusBadRequest)
		return
	}

	secret, ok := h.TenantSecrets[tenant]
	if !ok || secret == "" {
		http.Error(w, "tenant not allowed", http.StatusForbidden)
		return
	}

	signature := r.Header.Get("X-Duplynx-Signature")
	if signature == "" {
		http.Error(w, "missing signature", http.StatusBadRequest)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read payload", http.StatusInternalServerError)
		return
	}

	if !validateSignature(secret, payload, signature) {
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	// TODO: enqueue payload for processing in Phase 3+.
	log.Printf("ingestion accepted tenant=%s bytes=%d", tenant, len(payload))
	w.WriteHeader(http.StatusAccepted)
}

func validateSignature(secret string, payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := mac.Sum(nil)
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return hmac.Equal(expected, sigBytes)
}
