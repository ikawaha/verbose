⚠️ This library is archived and move to [github.com/ikawaha/errors](https://github.com/ikawaha/errors). 

verbose
===

The verbose is a library that adds context to errors.

Setting a value to an error:
```
err := errors.New("some error")
err = verbose.WithValue(err, "hello", "goodbye")
```

Getting the value associated with this error for key:
```
val, ok := verbose.Value(err, "hello")
fmt.Println(val, ok)

OUTPUT:
goodbye, true
```

Logging stack trace to an error:
```
err = verbose.WithStackTrace(err)
```

Retrieving the record from error:
```
trace, ok := verbose.StackTrace(err)
fmt.Println(trace)

OUTPUT:
error caused by fn1:
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:12 github.com/ikawaha/verbose_test.fn1
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:16 github.com/ikawaha/verbose_test.fn2
    /Users/ikawaha/go/src/github.com/ikawaha/verbose/error_test.go:20 github.com/ikawaha/verbose_test.TestWithStackTrace
    /usr/local/opt/go/libexec/src/testing/testing.go:1439 testing.tRunner
    /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1571 runtime.goexit
```

---
MIT
