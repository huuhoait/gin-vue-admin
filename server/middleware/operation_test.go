package middleware

import "testing"

func TestScrubSensitive(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"password field", `{"password":"hunter2","user":"a"}`, `{"password":"[REDACTED]","user":"a"}`},
		{"token field", `{"token":"abc.def.ghi"}`, `{"token":"[REDACTED]"}`},
		{"accessToken camel", `{"accessToken":"ey..."}`, `{"accessToken":"[REDACTED]"}`},
		{"x-token dashed", `{"x-token":"ey..."}`, `{"x-token":"[REDACTED]"}`},
		{"nested password", `{"user":{"password":"hunter2"}}`, `{"user":{"password":"[REDACTED]"}}`},
		{"no sensitive", `{"name":"alice"}`, `{"name":"alice"}`},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := scrubSensitive(c.in)
			if got != c.want {
				t.Fatalf("scrubSensitive(%q) = %q want %q", c.in, got, c.want)
			}
		})
	}
}

func TestSummarizeResponse(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"envelope", `{"code":0,"msg":"ok","data":{"secret":1}}`, `{"code":0,"msg":"ok"}`},
		{"non json short", `OK`, `OK`},
		{"empty", ``, ``},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := summarizeResponse(c.in)
			if got != c.want {
				t.Fatalf("summarizeResponse(%q) = %q want %q", c.in, got, c.want)
			}
		})
	}
}
