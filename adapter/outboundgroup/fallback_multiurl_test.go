package outboundgroup

import (
	"reflect"
	"testing"

	C "github.com/metacubex/mihomo/constant"
)

type fallbackHealthProxy struct {
	C.Proxy
	alive map[string]bool
}

func (p *fallbackHealthProxy) AliveForTestUrl(url string) bool {
	return p.alive[url]
}

func TestFallbackAliveForAllTestUrls(t *testing.T) {
	f := &Fallback{testUrls: []string{
		"https://example-a.test/",
		"https://example-b.test/",
		"https://example-c.test/",
	}}
	p := &fallbackHealthProxy{alive: map[string]bool{
		"https://example-a.test/": true,
		"https://example-b.test/": true,
		"https://example-c.test/": true,
	}}

	if !f.aliveForAllTestUrls(p) {
		t.Fatal("expected proxy to be alive when all test URLs pass")
	}

	p.alive["https://example-b.test/"] = false
	if f.aliveForAllTestUrls(p) {
		t.Fatal("expected proxy to be unhealthy when any test URL fails")
	}
}

func TestFallbackSingleURLCompatibility(t *testing.T) {
	f := &Fallback{testUrls: []string{"https://example.test/"}}
	p := &fallbackHealthProxy{alive: map[string]bool{"https://example.test/": true}}

	if !f.aliveForAllTestUrls(p) {
		t.Fatal("expected legacy single-URL health state to remain valid")
	}

	p.alive["https://example.test/"] = false
	if f.aliveForAllTestUrls(p) {
		t.Fatal("expected legacy single-URL failure to remain unhealthy")
	}
}

func TestNormalizeTestURLs(t *testing.T) {
	got, err := normalizeTestURLs([]string{
		" https://example-a.test/ ",
		"https://example-b.test/",
		"https://example-a.test/",
	})
	if err != nil {
		t.Fatalf("normalizeTestURLs returned error: %v", err)
	}

	want := []string{
		"https://example-a.test/",
		"https://example-b.test/",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeTestURLs() = %#v, want %#v", got, want)
	}
}

func TestNormalizeTestURLsRejectsEmptyEntry(t *testing.T) {
	if _, err := normalizeTestURLs([]string{"https://example.test/", "  "}); err == nil {
		t.Fatal("expected empty URL entry to be rejected")
	}
}

func TestExtraTestURLsLegacySafe(t *testing.T) {
	if got := extraTestURLs(nil); got != nil {
		t.Fatalf("extraTestURLs(nil) = %#v, want nil", got)
	}
	if got := extraTestURLs([]string{"https://example.test/"}); got != nil {
		t.Fatalf("extraTestURLs(single) = %#v, want nil", got)
	}
}
