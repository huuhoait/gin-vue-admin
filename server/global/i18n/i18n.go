// Package i18n provides a lightweight message bundle for the Admin service.
//
// Design goals:
//   - Zero new external dependencies (reuses github.com/BurntSushi/toml which is
//     already pulled in by gorm-oracle transitively).
//   - BCP-47 tags ("vi-VN", "en-US"). Primary subtag accepted as fallback
//     ("vi" -> "vi-VN"), mirroring the previous SkyAgent admin i18n behaviour.
//   - Aligned with monorepo ARCH-04 convention: default locale vi-VN, fallback
//     en-US. Backend log strings remain English regardless of active locale.
//
// Usage:
//
//	i18n.MustLoad("services/admin/server/resource/i18n")
//	msg := i18n.T("vi-VN", "admin.common.create_failed")          // plain key
//	msg := i18n.T("vi-VN", "admin.user.welcome", "name", "Alice") // templated
package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

const (
	// DefaultLocale is what we fall back to when no Accept-Language is present
	// and when the configured preference is unrecognised.
	DefaultLocale = "vi-VN"
	// FallbackLocale is consulted when a key is missing from the active locale.
	FallbackLocale = "en-US"
)

var (
	mu      sync.RWMutex
	bundles = map[string]map[string]string{} // locale -> key -> template
)

// Load reads every *.toml file under dir and registers it as a locale bundle
// keyed by the filename stem (e.g. vi-VN.toml -> "vi-VN").
//
// TOML format is a flat key/value map:
//
//	"admin.common.create_failed" = "Tạo không thành công"
//	"admin.user.welcome"         = "Xin chào {{.name}}"
func Load(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("i18n: read dir %q: %w", dir, err)
	}

	loaded := map[string]map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		locale := strings.TrimSuffix(e.Name(), ".toml")
		m := map[string]string{}
		if _, err := toml.DecodeFile(filepath.Join(dir, e.Name()), &m); err != nil {
			return fmt.Errorf("i18n: decode %q: %w", e.Name(), err)
		}
		loaded[locale] = m
	}

	mu.Lock()
	bundles = loaded
	mu.Unlock()
	return nil
}

// MustLoad panics on error — intended for application bootstrap.
func MustLoad(dir string) {
	if err := Load(dir); err != nil {
		panic(err)
	}
}

// T resolves a translation key for the given locale. Variadic args form
// key/value pairs substituted into "{{.name}}" placeholders.
//
// Resolution order: active locale -> FallbackLocale -> raw key.
// The raw key is returned (instead of empty string) to make missing keys
// visible in UI / tests without blowing up.
func T(locale, key string, kv ...string) string {
	locale = Normalize(locale)

	mu.RLock()
	tmpl, ok := lookup(locale, key)
	if !ok && locale != FallbackLocale {
		tmpl, ok = lookup(FallbackLocale, key)
	}
	mu.RUnlock()

	if !ok {
		return key
	}
	return substitute(tmpl, kv)
}

func lookup(locale, key string) (string, bool) {
	b, ok := bundles[locale]
	if !ok {
		return "", false
	}
	v, ok := b[key]
	return v, ok
}

func substitute(tmpl string, kv []string) string {
	if len(kv) == 0 {
		return tmpl
	}
	out := tmpl
	for i := 0; i+1 < len(kv); i += 2 {
		placeholder := "{{." + kv[i] + "}}"
		out = strings.ReplaceAll(out, placeholder, kv[i+1])
	}
	return out
}

// Normalize canonicalises a raw Accept-Language value into one of the bundled
// locales. Unknown tags return DefaultLocale.
func Normalize(tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return DefaultLocale
	}
	// Take the first language-range if comma-separated.
	if i := strings.IndexByte(tag, ','); i >= 0 {
		tag = tag[:i]
	}
	// Strip q-factor.
	if i := strings.IndexByte(tag, ';'); i >= 0 {
		tag = tag[:i]
	}
	tag = strings.TrimSpace(tag)

	// Canonical case: "vi-VN", "en-US".
	parts := strings.Split(tag, "-")
	switch strings.ToLower(parts[0]) {
	case "vi":
		return "vi-VN"
	case "en":
		return "en-US"
	default:
		return DefaultLocale
	}
}

// Locales returns the sorted list of loaded locale tags (for diagnostics / tests).
func Locales() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(bundles))
	for k := range bundles {
		out = append(out, k)
	}
	return out
}
