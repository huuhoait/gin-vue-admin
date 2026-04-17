package system

import (
	"testing"

	"github.com/huuhoait/gin-vue-admin/server/model/system"
)

func TestHashChainEntryDeterministic(t *testing.T) {
	row := &system.SysPolicyChangeLog{
		Actor: 1, ActorUserID: 2,
		Action: "update", AuthorityID: "888",
		Before: `{"a":1}`, After: `{"a":2}`,
		IP: "1.2.3.4", RequestID: "req-123", Note: "test",
	}
	h1 := hashChainEntry(row, "prev000")
	h2 := hashChainEntry(row, "prev000")
	if h1 != h2 {
		t.Fatalf("hash not deterministic: %q vs %q", h1, h2)
	}
	if len(h1) != 64 {
		t.Fatalf("expected 64-char hex SHA-256, got len=%d", len(h1))
	}
}

func TestHashChainEntryChangesWithPrev(t *testing.T) {
	row := &system.SysPolicyChangeLog{Action: "update"}
	h1 := hashChainEntry(row, "aaa")
	h2 := hashChainEntry(row, "bbb")
	if h1 == h2 {
		t.Fatal("different prevHash must produce different digest")
	}
}

func TestHashChainEntryChangesWithField(t *testing.T) {
	row := &system.SysPolicyChangeLog{Action: "update", Note: "original"}
	h1 := hashChainEntry(row, "")
	row.Note = "tampered"
	h2 := hashChainEntry(row, "")
	if h1 == h2 {
		t.Fatal("mutating a field must change the digest")
	}
}
