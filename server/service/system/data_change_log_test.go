package system

import "testing"

func TestScrubDataJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"password string", `{"password":"secret123","user":"alice"}`, `{"password":"[REDACTED]","user":"alice"}`},
		{"token field", `{"token":"eyJ.abc"}`, `{"token":"[REDACTED]"}`},
		{"accessToken camel", `{"accessToken":"ey..."}`, `{"accessToken":"[REDACTED]"}`},
		{"signingKey", `{"signingKey":"abc"}`, `{"signingKey":"[REDACTED]"}`},
		{"nested password", `{"user":{"password":"hunter2"}}`, `{"user":{"password":"[REDACTED]"}}`},
		{"no sensitive", `{"name":"alice","role":"admin"}`, `{"name":"alice","role":"admin"}`},
		{"empty", ``, ``},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := scrubDataJSON(c.in)
			if got != c.want {
				t.Fatalf("scrubDataJSON(%q)\n got  %q\n want %q", c.in, got, c.want)
			}
		})
	}
}
