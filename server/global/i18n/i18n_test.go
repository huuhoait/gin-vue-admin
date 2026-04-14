package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"":               "vi-VN",
		"vi":             "vi-VN",
		"vi-VN":          "vi-VN",
		"en":             "en-US",
		"en-US":          "en-US",
		"en-GB":          "en-US",
		"fr":             "vi-VN", // unknown -> default
		"en-US,vi;q=0.9": "en-US",
		"  vi-VN  ":      "vi-VN",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLoadAndT(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("vi-VN.toml", `"admin.common.create_failed" = "Tạo không thành công"
"admin.user.welcome" = "Xin chào {{.name}}"`)
	write("en-US.toml", `"admin.common.create_failed" = "Create failed"
"admin.common.only_en" = "english-only"
"admin.user.welcome" = "Hello {{.name}}"`)

	if err := Load(dir); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got := T("vi-VN", "admin.common.create_failed"); got != "Tạo không thành công" {
		t.Errorf("vi-VN hit: got %q", got)
	}
	if got := T("vi-VN", "admin.common.only_en"); got != "english-only" {
		t.Errorf("vi-VN fallback to en-US: got %q", got)
	}
	if got := T("vi-VN", "admin.user.welcome", "name", "Nam"); got != "Xin chào Nam" {
		t.Errorf("template substitution: got %q", got)
	}
	if got := T("en-US", "admin.does.not.exist"); got != "admin.does.not.exist" {
		t.Errorf("missing key should echo: got %q", got)
	}
}
