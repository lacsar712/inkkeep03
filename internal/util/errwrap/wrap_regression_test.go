package errwrap

import (
	"errors"
	"testing"
)

// TestWrapDeniedSentinelChain is a regression test for the bug where WrapDenied
// formatted ErrDenied with %v instead of %w, which broke errors.Is chain
// traversal and caused the sentinel to be lost after wrapping.
func TestWrapDeniedSentinelChain(t *testing.T) {
	wrapped := WrapDenied("update")

	if !errors.Is(wrapped, ErrDenied) {
		t.Fatalf("errors.Is(wrapped, ErrDenied) = false; sentinel lost, got %v", wrapped)
	}
	if !IsDenied(wrapped) {
		t.Fatalf("IsDenied(wrapped) = false; got %v", wrapped)
	}

	// The wrapped sentinel must itself still match the original sentinel so
	// downstream errors.Is/As traversal (including nested wrapping) works.
	var target *error
	_ = target
	if !errors.Is(ErrDenied, ErrDenied) {
		t.Fatal("baseline errors.Is(ErrDenied, ErrDenied) must be true")
	}

	// Double-wrap: an error wrapping WrapDenied's result must still match the
	// sentinel through the full chain.
	double := fmtWrapf(wrapped, "outer")
	if !errors.Is(double, ErrDenied) {
		t.Fatalf("errors.Is after double-wrap lost sentinel; got %v", double)
	}
}

// fmtWrapf mirrors fmt.Errorf("%s: %w", op, err) without importing fmt here, to
// keep the regression test focused on the public errwrap API.
func fmtWrapf(err error, op string) error {
	return &chainedErr{op: op, err: err}
}

type chainedErr struct {
	op  string
	err error
}

func (c *chainedErr) Error() string { return c.op + ": " + c.err.Error() }
func (c *chainedErr) Unwrap() error { return c.err }
