package app

import (
	"log"
	"os"
	"strings"
)

// Config represents runtime configuration for DupLynx services.
type Config struct {
	DatabasePath  string
	Addr          string
	EmbedStatic   bool
	Mode          string
	TenantSecrets map[string]string
}

// LoadConfig builds a Config from environment variables with sensible defaults.
func LoadConfig() Config {
	cfg := Config{
		DatabasePath:  getEnv("DUPLYNX_DB_FILE", "var/duplynx.db"),
		Addr:          getEnv("DUPLYNX_ADDR", ":8080"),
		EmbedStatic:   strings.ToLower(getEnv("DUPLYNX_EMBED_STATIC", "true")) != "false",
		Mode:          getEnv("DUPLYNX_MODE", "server"),
		TenantSecrets: parseTenantSecrets(os.Getenv("DUPLYNX_TENANT_SECRETS")),
	}

	if len(cfg.TenantSecrets) == 0 {
		log.Println("warning: no tenant HMAC secrets configured; ingestion endpoints will reject unsigned payloads")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func (c Config) SQLiteDSN() string {
	params := []string{"_foreign_keys=1", "_busy_timeout=5000"}
	if strings.EqualFold(c.Mode, "gui") {
		params = append(params, "mode=ro")
	}
	return "file:" + c.DatabasePath + "?" + strings.Join(params, "&")
}

func parseTenantSecrets(raw string) map[string]string {
	secrets := make(map[string]string)
	pairs := strings.Split(raw, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			log.Printf("invalid tenant secret pair: %q", pair)
			continue
		}
		tenant := strings.TrimSpace(parts[0])
		secret := strings.TrimSpace(parts[1])
		if tenant == "" || secret == "" {
			log.Printf("invalid tenant secret pair: %q", pair)
			continue
		}
		secrets[tenant] = secret
	}
	return secrets
}
