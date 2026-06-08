package errors

import (
	stderrors "errors"
	"testing"
)

func TestNewAndError(t *testing.T) {
	e := New("boom")
	if e.Error() != "boom" {
		t.Fatalf("Error() = %q, want %q", e.Error(), "boom")
	}
}

func TestWithMethodsAreImmutable(t *testing.T) {
	base := New("base")
	derived := base.
		WithUpstreamName("alipay").
		WithUpstreamCode("CODE").
		WithUpstreamMessage("msg").
		WithFrameworkName("fiber")

	// The base sentinel must stay clean.
	if base.UpstreamName() != "" || base.UpstreamCode() != "" ||
		base.UpstreamMessage() != "" || base.FrameworkName() != "" {
		t.Errorf("base error was mutated: %+v", base)
	}

	// The derived error carries the metadata.
	if derived.UpstreamName() != "alipay" {
		t.Errorf("UpstreamName() = %q, want alipay", derived.UpstreamName())
	}
	if derived.UpstreamCode() != "CODE" {
		t.Errorf("UpstreamCode() = %q, want CODE", derived.UpstreamCode())
	}
	if derived.UpstreamMessage() != "msg" {
		t.Errorf("UpstreamMessage() = %q, want msg", derived.UpstreamMessage())
	}
	if derived.FrameworkName() != "fiber" {
		t.Errorf("FrameworkName() = %q, want fiber", derived.FrameworkName())
	}
}

func TestWithUpstreamResponse(t *testing.T) {
	type resp struct{ Code int }
	r := resp{Code: 200}
	e := New("x").WithUpstreamResponse(r)
	got, ok := e.UpstreamResponse().(resp)
	if !ok || got.Code != 200 {
		t.Errorf("UpstreamResponse() = %v, want %v", e.UpstreamResponse(), r)
	}
}

func TestErrorsIsMatchesSentinelAfterChaining(t *testing.T) {
	// A derived error must still match the original sentinel via errors.Is,
	// so callers can use errors.Is(err, ErrTradeNotExist) even when a driver
	// attached upstream metadata.
	derived := ErrTradeNotExist.WithUpstreamName("alipay").WithUpstreamCode("X")
	if !stderrors.Is(derived, ErrTradeNotExist) {
		t.Errorf("errors.Is(derived, ErrTradeNotExist) = false, want true")
	}
	if stderrors.Is(derived, ErrUnsupportedCurrency) {
		t.Errorf("errors.Is(derived, ErrUnsupportedCurrency) = true, want false")
	}
}
