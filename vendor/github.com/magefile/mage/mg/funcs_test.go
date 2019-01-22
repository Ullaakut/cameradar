package mg

import (
	"context"
	"testing"
)

func TestFuncCheck(t *testing.T) {
	// we can ignore errors here, since the error is always the same
	// and the FuncType will be InvalidType if there's an error.
	f, _ := funcCheck(func() {})
	if f != voidType {
		t.Errorf("expected func() to be a valid VoidType, but was %v", f)
	}
	f, _ = funcCheck(func() error { return nil })
	if f != errorType {
		t.Errorf("expected func() error to be a valid ErrorType, but was %v", f)
	}
	f, _ = funcCheck(func(context.Context) {})
	if f != contextVoidType {
		t.Errorf("expected func(context.Context) to be a valid ContextVoidType, but was %v", f)
	}
	f, _ = funcCheck(func(context.Context) error { return nil })
	if f != contextErrorType {
		t.Errorf("expected func(context.Context) error to be a valid ContextErrorType but was %v", f)
	}

	f, _ = funcCheck(Foo.Bare)
	if f != namespaceVoidType {
		t.Errorf("expected Foo.Bare to be a valid NamespaceVoidType but was %v", f)
	}

	f, _ = funcCheck(Foo.Error)
	if f != namespaceErrorType {
		t.Errorf("expected Foo.Error to be a valid NamespaceErrorType but was %v", f)
	}
	f, _ = funcCheck(Foo.BareCtx)
	if f != namespaceContextVoidType {
		t.Errorf("expected Foo.BareCtx to be a valid NamespaceContextVoidType but was %v", f)
	}
	f, _ = funcCheck(Foo.CtxError)
	if f != namespaceContextErrorType {
		t.Errorf("expected Foo.CtxError to be a valid NamespaceContextErrorType but was %v", f)
	}

	// Test the Invalid case
	f, err := funcCheck(func(int) error { return nil })
	if f != invalidType {
		t.Errorf("expected func(int) error to be InvalidType but was %v", f)
	}
	if err == nil {
		t.Error("expected func(int) error to not be a valid FuncType, but got nil error.")
	}
}

type Foo Namespace

func (Foo) Bare() {}

func (Foo) Error() error { return nil }

func (Foo) BareCtx(context.Context) {}

func (Foo) CtxError(context.Context) error { return nil }

func TestFuncTypeWrap(t *testing.T) {
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("Expected a panic, but didn't get one")
			}
		}()
		if funcTypeWrap(voidType, func(i int) {}) != nil {
			t.Errorf("expected func(int) to return nil")
		}
	}()

	if funcTypeWrap(voidType, func() {}) == nil {
		t.Errorf("expected func() to return a function")
	}

	if funcTypeWrap(errorType, func() error { return nil }) == nil {
		t.Errorf("expected func() error to return a function")
	}

	if funcTypeWrap(contextVoidType, func(context.Context) {}) == nil {
		t.Errorf("expected func(context.Context) to return a function")
	}

	if funcTypeWrap(contextErrorType, func(context.Context) error { return nil }) == nil {
		t.Errorf("expected func(context.Context) error to return a function")
	}
}
