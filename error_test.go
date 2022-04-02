package verbose_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ikawaha/verbose"
)

func fn1() error {
	return verbose.WithStackTrace(errors.New("error caused by fn1"))
}

func fn2() error {
	return fmt.Errorf("fn2 has error: %w", fn1())
}

func TestWithStackTrace(t *testing.T) {
	err := fn2()
	if err == nil {
		t.Fatal("expected error, but nil")
	}
	t.Run("wrapped error", func(t *testing.T) {
		want := "fn2 has error: error caused by fn1"
		if got := err.Error(); want != err.Error() {
			t.Errorf("want %v, got %v", want, got)
		}
	})
	t.Run("stack trace", func(t *testing.T) {
		got, ok := verbose.StackTrace(err)
		if !ok {
			t.Errorf("expected true, but false")
		}
		want := `error caused by fn1:
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:12 github.com/ikawaha/verbose_test.fn1
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:16 github.com/ikawaha/verbose_test.fn2
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:20 github.com/ikawaha/verbose_test.TestWithStackTrace
    /usr/local/opt/go/libexec/src/testing/testing.go:1439 testing.tRunner
    /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1571 runtime.goexit
`
		if got != want {
			t.Errorf("want: %q, got: %q", got, want)
		}
	})
}

func TestWithValue(t *testing.T) {
	err := errors.New("")
	err = verbose.WithValue(err, "key1", "value1")
	err = verbose.WithValue(err, "key2", "value2")
	err = verbose.WithValue(err, "key3", "value3")
	err = verbose.WithValue(err, "key1", "overwritten")

	testdata := []struct {
		name    string
		key     string
		wantOk  bool
		wantVal string
	}{
		{name: "ok: key2", key: "key2", wantOk: true, wantVal: "value2"},
		{name: "ok: key3", key: "key3", wantOk: true, wantVal: "value3"},
		{name: "ok: key1 (overwritten)", key: "key1", wantOk: true, wantVal: "overwritten"},
		{name: "ng (empty key)", key: "", wantOk: false, wantVal: ""},
		{name: "ng", key: "hello", wantOk: false, wantVal: ""},
	}
	for _, d := range testdata {
		t.Run(d.name, func(t *testing.T) {
			got, ok := verbose.Value(err, d.key)
			if ok != d.wantOk {
				t.Errorf("want: %t, got: %t", ok, d.wantOk)
			}
			if got != d.wantVal {
				t.Errorf("want: %v, got: %v", d.wantVal, got)
			}
		})
	}
}
