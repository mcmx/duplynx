package app

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateAssetDirectory ensures the supplied directory exists and contains the expected Tailwind bundle.
func ValidateAssetDirectory(root string) error {
	if root == "" {
		return fmt.Errorf("assets directory must be provided")
	}

	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("unable to stat assets directory %q: %w", root, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("assets path %q is not a directory", root)
	}

	bundle := filepath.Join(root, "tailwind.css")
	if stat, err := os.Stat(bundle); err != nil || stat.IsDir() {
		return fmt.Errorf("tailwind bundle missing at %q; run `npm run build:tailwind` before serving", bundle)
	}

	return nil
}
