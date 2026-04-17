package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/huuhoait/gin-vue-admin/server/global/i18n"
	"github.com/gin-gonic/gin"
)

// loadTestBundles seeds minimal vi-VN + en-US bundles covering the keys
// exercised by the helpers under test. Kept inline so the test does not
// depend on files outside this package.
func loadTestBundles(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	vi := `"admin.common.ok"            = "Thành công"
"admin.common.fail"          = "Thất bại"
"admin.common.create_failed" = "Tạo không thành công"
"admin.user.welcome"         = "Xin chào {{.name}}"`
	en := `"admin.common.ok"            = "Success"
"admin.common.fail"          = "Failed"
"admin.common.create_failed" = "Create failed"
"admin.user.welcome"         = "Hello {{.name}}"`
	if err := os.WriteFile(filepath.Join(dir, "vi-VN.toml"), []byte(vi), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "en-US.toml"), []byte(en), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := i18n.Load(dir); err != nil {
		t.Fatal(err)
	}
}

// withLocale builds a gin test context preloaded with an Accept-Language
// and the locale key that the middleware would have set.
func withLocale(locale string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set(localeKey, locale)
	return c, w
}

func parseMsg(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var r Response
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
		t.Fatalf("decode response: %v (body=%q)", err, w.Body.String())
	}
	return r.Msg
}

func TestOkAndFailDefaultsRespectLocale(t *testing.T) {
	loadTestBundles(t)

	cVi, wVi := withLocale("vi-VN")
	Ok(cVi)
	if got := parseMsg(t, wVi); got != "Thành công" {
		t.Errorf("vi-VN Ok() msg = %q, want Thành công", got)
	}

	cEn, wEn := withLocale("en-US")
	Fail(cEn)
	if got := parseMsg(t, wEn); got != "Failed" {
		t.Errorf("en-US Fail() msg = %q, want Failed", got)
	}
}

func TestFailWithCodeTemplating(t *testing.T) {
	loadTestBundles(t)

	c, w := withLocale("vi-VN")
	OkWithCode(c, "admin.user.welcome", "name", "Nam")
	if got := parseMsg(t, w); got != "Xin chào Nam" {
		t.Errorf("template substitution failed: got %q", got)
	}
}

func TestMissingLocaleFallsBackToDefault(t *testing.T) {
	loadTestBundles(t)

	// No localeKey set — should still resolve via DefaultLocale (vi-VN).
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	FailWithCode(c, "admin.common.create_failed")
	if got := parseMsg(t, w); got != "Tạo không thành công" {
		t.Errorf("default locale resolution: got %q", got)
	}
}
