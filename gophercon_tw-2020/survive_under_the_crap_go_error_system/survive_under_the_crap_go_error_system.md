---
title: Survive under the Crap Go Error System
tags: Slide, Go, error
description: View the slide with "Slide Mode".
slideOptions:
  spotlight:
    enabled: false
  allottedMinutes: 25
---
Survive under the Crap Go Error System
===

<!-- .slide: data-background-color="pink" -->
<!-- .slide: data-transition="zoom" -->

> [name=郭學聰 Hsueh-Tsung Kuo]
> [time=Sat, 14 Nov 2020] [color=red]

###### CC BY-SA 4.0

---

<!-- .slide: data-transition="convex" -->

## Who am I?

![fieliapm](https://www.gravatar.com/avatar/2aef78f04240a6ac9ccd473ba1cbd1e3?size=2048 =384x384)

<small>Someone (who?) said:
a game programmer should be able to draw cute anime character(?)</small>

----

<!-- .slide: data-transition="convex" -->

* A programmer from game company in Taiwan
* Backend (and temporary frontend) engineer
* Usually develop something related to my work in Python, Ruby, ECMAScript, Golang, C#
* ECMAScript hater since **Netscape** is dead
* Built CDN-aware game asset update system
* Professional small vehicle driver
* Draw cute anime character in spare time

---

<!-- .slide: data-transition="convex" -->

## Outline

----

<!-- .slide: data-transition="convex" -->

4. Error Usage
5. Error Design Flaw
    * Error 5 Ws
    * Go Error System
6. Package
    * errors
    * golang.org/x/xerrors
    * github.com/pkg/errors
    * Integration

----

<!-- .slide: data-transition="convex" -->

7. Advice
    * Wrapping
    * Inspection
    * Displaying
    * Sum Up
8. One More Thing
9. Conclusion
10. Reference
11. Q&A

---

<!-- .slide: data-transition="convex" -->

## Error Usage

----

<!-- .slide: data-transition="convex" -->

### Typical Example

```go=
func MyFunc(param1 int, param2 string) (string, error) {
	// do something...
	f, err := os.Open("filename.ext")
	if err != nil {
		return "", err
	}
	// do something...
	return result, nil
}
```

---

<!-- .slide: data-transition="convex" -->

## Error Design Flaw

----

<!-- .slide: data-transition="convex" -->

### Error 5 Ws

* What
  * Error title
* Why
  * Error description
* Which
  * Error type
* When
  * Line number
* Where
  * Call stack

----

<!-- .slide: data-transition="convex" -->

#### Who Need to Know What

* End-User
  * What
  * Why
* Package Importer
  * Which
* Package Developer & Maintainer
  * When
  * Where

----

<!-- .slide: data-transition="convex" -->

### :thumbsdown: Go Error System

* [x] Error title
* [x] Error description
* [ ] Error type
* [ ] Line number
* [ ] Call stack

==Why?== <!-- .element: class="fragment" data-fragment-index="1" -->

----

<!-- .slide: data-transition="convex" -->

#### Review Typical Example

```go=
func MyFunc(param1 int, param2 string) (string, error) {
	// do something...
	f, err := os.Open("filename.ext")
	if err != nil {
		return "", err
	}
	// do something...
	return result, nil
}
```

----

<!-- .slide: data-transition="convex" -->

#### Pitfalls

```go=
ErrSomething = errors.New("error description")
return ErrSomething

err := fmt.Errorf("error %s", detail)
return err
```

----

<!-- .slide: data-transition="convex" -->

#### Darkness

* How to classify built-in errors before Go 1.13
  * Unreliable error message
  * Useless error type (in most cases)
  * Nothing further :face_palm:

----

<!-- .slide: data-transition="convex" -->

#### Dawn?

* 3rd party error module
* Go 1.13 error system

---

<!-- .slide: data-transition="convex" -->

## Package

----

<!-- .slide: data-transition="convex" -->

### errors

* :no_entry_sign: Go 1.13 error system

----

<!-- .slide: data-transition="convex" -->

#### Wrap & Unwrap

```go=
// wrap error
err = fmt.Errorf("... %w ...", ..., cause, ...)

// unwrap and return NEXT error, or nil
cause = errors.Unwrap(err)
```

----

<!-- .slide: data-transition="convex" -->

#### Inspect Error Chain

```go=
errors.Is(err, ErrorInstance) // must be same address

var perr *ErrorType
errors.As(err, &perr)
```

----

<!-- .slide: data-transition="convex" -->

#### Example

```go=
package main

import (
	"errors"
	"fmt"
)

type MyError struct {
	msg string
	err error
}

func NewMyError(msg string, err error) *MyError {
	if err != nil {
		msg = fmt.Sprintf("%s: %s", msg, err.Error())
	}
	return &MyError{
		msg: msg,
		err: err,
	}
}

func (e *MyError) Error() string {
	return e.msg
}

func (e *MyError) Unwrap() error {
	return e.err
}

func main() {
	var err error
	err = errors.New("err1")
	err = fmt.Errorf("err2: %w", err)

	myErr := NewMyError("my err3", err)
	myErr2 := NewMyError("my err3", err)
	err = myErr

	err = fmt.Errorf("err4: %w", err)
	err = fmt.Errorf("err5: %w", err)

	if err != nil {
		fmt.Println("error is:", err)

		if errors.Is(err, myErr) {
			fmt.Println("it is myErr!")
		}
		if !errors.Is(err, myErr2) {
			fmt.Println("it is NOT myErr2!")
		}

		fmt.Println("")

		var myError *MyError
		if errors.As(err, &myError) {
			fmt.Println("failed:", err.Error())
			fmt.Println("MyError:", myError.Error())
		}
	}
}
```

----

<!-- .slide: data-transition="convex" -->

#### Result

```shell
error is: err5: err4: my err3: err2: err1
it is myErr!
it is NOT myErr2!

failed: err5: err4: my err3: err2: err1
MyError: my err3: err2: err1
```

----

<!-- .slide: data-transition="convex" -->

### golang.org/x/xerrors

* :no_entry_sign: golang.org/x/xerrors by Go team members
  * https://github.com/golang/xerrors
  * https://godoc.org/golang.org/x/xerrors

----

<!-- .slide: data-transition="convex" -->

#### Additional Features

* Frame
  * New & Wrap
    * Attach current function & line number
  * Inspect
    * Print current function & line number

----

<!-- .slide: data-transition="convex" -->

#### Usage

```go=
// new error
err = xerrors.New("...")
// wrap error
err = xerrors.Errorf("... %w ...", ..., cause, ...)

// print error chain and frames
fmt.Printf("%+v\n", err)
```

----

<!-- .slide: data-transition="convex" -->

#### Example

```go=
package main

import (
	"fmt"

	"golang.org/x/xerrors"
)

func MyFuncInner() error {
	return xerrors.New("inner error")
}

func MyFuncMiddle() error {
	err := MyFuncInner()
	return xerrors.Errorf("middle error: %w", err)
}

func MyFuncOuter() error {
	err := MyFuncMiddle()
	return xerrors.Errorf("outer error: %w", err)
}

func main() {
	err := MyFuncOuter()
	fmt.Printf("%v\n\n", err)
	fmt.Printf("%+v\n", err)
}
```

----

<!-- .slide: data-transition="convex" -->

#### Result

```shell
outer error: middle error: inner error

outer error:
    main.MyFuncOuter
        /.../go/src/survive_under_the_crap_go_error_system/xerrors_example.go:20
  - middle error:
    main.MyFuncMiddle
        /.../go/src/survive_under_the_crap_go_error_system/xerrors_example.go:15
  - inner error:
    main.MyFuncInner
        /.../go/src/survive_under_the_crap_go_error_system/xerrors_example.go:10
```

----

<!-- .slide: data-transition="convex" -->

### github.com/pkg/errors

* :thumbsup: github.com/pkg/errors by Dave Cheney
  * https://github.com/pkg/errors
  * https://godoc.org/github.com/pkg/errors

----

<!-- .slide: data-transition="convex" -->

#### Wrap & Unwrap

```go=
// new error
err = pkg_errors.New("...")
// wrap error
err = pkg_errors.Wrap(cause, "...")
err = pkg_errors.Wrapf(cause, "oh noes #%d", 2)

// unwrap the whole error chain and return the MOST INNER error, or return the error itself
cause = pkg_errors.Cause(err)
```

----

<!-- .slide: data-transition="convex" -->

#### Inspect Error Chain

```go=
fmt.Printf("%+v\n", err)
```

----

<!-- .slide: data-transition="convex" -->

#### Example

```go=
package main

import (
	"fmt"

	pkg_errors "github.com/pkg/errors"
)

func MyFuncInner() error {
	return pkg_errors.New("inner error")
}

func MyFuncMiddle() error {
	err := MyFuncInner()
	return pkg_errors.Wrap(err, "middle error")
}

func MyFuncOuter() error {
	err := MyFuncMiddle()
	return pkg_errors.Wrapf(err, "outer error - %d", 2)
}

func main() {
	err := MyFuncOuter()
	fmt.Printf("%v\n\n", err)
	fmt.Printf("%+v\n", err)
}
```

----

<!-- .slide: data-transition="convex" -->

#### Result

```shell
outer error - 2: middle error: inner error

inner error
main.MyFuncInner
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:10
main.MyFuncMiddle
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:14
main.MyFuncOuter
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:19
main.main
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:24
runtime.main
	/usr/local/go/src/runtime/proc.go:204
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1374
middle error
main.MyFuncMiddle
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:15
main.MyFuncOuter
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:19
main.main
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:24
runtime.main
	/usr/local/go/src/runtime/proc.go:204
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1374
outer error - 2
main.MyFuncOuter
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:20
main.main
	/.../go/src/survive_under_the_crap_go_error_system/pkg_errors_example.go:24
runtime.main
	/usr/local/go/src/runtime/proc.go:204
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1374
```

----

<!-- .slide: data-transition="convex" -->

### Integration

* for http server
  * use the gin, luke
    * <small>https://github.com/gin-gonic/gin#model-binding-and-validation</small>

```go=
// push error
c.Error(err)

// list error
c.Errors.ByType(errType)
```

---

<!-- .slide: data-transition="convex" -->

## Advice

----

<!-- .slide: data-transition="convex" -->

### Wrapping

Wrap as early as possible

```go=
import (
	"errors"
	"fmt"
	"golang.org/x/xerrors"
	pkg_errors "github.com/pkg/errors"
)

func MyFunc(param1 int, param2 string) (string, error) {
	// do something...
	r, err := foreign_package.ForeignFunc()
	if err != nil {
		err = fmt.Errorf("<error description>: %w", err)
		err = xerrors.Errorf("<error description>: %w", err)
		err = pkg_errors.Wrap(err, "<error description>")
		return "", err
	}
	// do something...
	return result, nil
}
```

----

<!-- .slide: data-transition="convex" -->

### Inspection

Collect and log error at entry function or middleware

```go=
// for Go 1.13 error system
func ExtractErrorMessageChain(err error) string {
	var b strings.Builder
	for e := err; e != nil; e = errors.Unwrap(e) {
		b.WriteString(e.Error())
		b.WriteString("\n")
	}
	return b.String()
}

func MyEntryFunc() {
	// do something...
	r, err := MyFunc(param1, param2)

	// for Go 1.13 error system
	fmt.Printf("%s", ExtractErrorMessageChain(err))

	// for golang.org/x/xerrors & github.com/pkg/errors
	fmt.Printf("%+v\n", err)
}
```

----

<!-- .slide: data-transition="convex" -->

### Displaying

Only show the most outer error message to end-user

```go=
// for Go 1.13 error system & golang.org/x/xerrors
func WrapError(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// for all error frameworks
func GetErrorMessageTitle(err error) string {
	return strings.SplitN(err.Error(), ":", 2)[0]
}
```

----

<!-- .slide: data-transition="convex" -->

### Sum Up

----

<!-- .slide: data-transition="convex" -->

#### Example

* Go 1.13 error system
* Go standard HTTP library
* github.com/julienschmidt/httprouter

```go=
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	//"reflect"

	"github.com/julienschmidt/httprouter"
)

// utility function

// for Go 1.13 error system
func ExtractErrorMessageChain(err error) string {
	var b strings.Builder
	for e := err; e != nil; e = errors.Unwrap(e) {
		b.WriteString(e.Error())
		b.WriteString("\n")
	}
	return b.String()
}

// for Go 1.13 error system & golang.org/x/xerrors
func WrapError(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// for all error frameworks
func GetErrorMessageTitle(err error) string {
	return strings.SplitN(err.Error(), ":", 2)[0]
}

// compare error using first part of message before colon, for example: "<first>: <rest>"
func IsErrorEqual(e1 error, e2 error) bool {
	if e1 == nil || e2 == nil {
		return e1 == e2
	} else {
		// CAUTION: error type is useless (in most cases)
		// return reflect.TypeOf(e1) == reflect.TypeOf(e2) && GetErrorMessageTitle(e1) == GetErrorMessageTitle(e2)
		return GetErrorMessageTitle(e1) == GetErrorMessageTitle(e2)
	}
}

// ApiError

type ApiError struct {
	statusCode int
	msg        string
	err        error
}

func NewApiError(statusCode int, msg string, err error) error {
	if err != nil {
		msg = fmt.Sprintf("%s: %s", msg, err.Error())
	}
	return &ApiError{statusCode: statusCode, msg: msg, err: err}
}

func (e *ApiError) Error() string {
	return e.msg
}

func (e *ApiError) Unwrap() error {
	return e.err
}

func (e *ApiError) StatusCode() int {
	return e.statusCode
}

func WrapApiError(template error, err error) error {
	apiErrorTemplate, ok := template.(*ApiError)
	if ok {
		err = NewApiError(apiErrorTemplate.StatusCode(), apiErrorTemplate.Error(), err)
	} else {
		err = WrapError(template.Error(), err)
	}
	return err
}

// handler

type HandleE func(http.ResponseWriter, *http.Request, httprouter.Params) error

func ErrAwareHandle(h HandleE) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := h(w, r, ps)
		if err != nil {
			// respond error with simplified message
			errorMessage := GetErrorMessageTitle(err)
			var statusCode int
			switch e := err.(type) {
			case *ApiError:
				statusCode = e.StatusCode()
			default:
				statusCode = http.StatusInternalServerError
			}
			http.Error(w, errorMessage, statusCode)

			// log detailed error
			messageChain := ExtractErrorMessageChain(err)
			fmt.Fprintf(os.Stderr, "error:\n%s====\n", messageChain)
		}
	}
}

// usage

var (
	ErrEmptyMessage  = errors.New("message is empty")
	ErrDirtyMessage  = errors.New("助兵衛")
	ErrHentaiMessage = errors.New("変態")
)

func detectNonDirtyMessage(msg string) error {
	if msg == "FQ" {
		return ErrDirtyMessage
	}
	return nil
}

func processMessage(msg string) (string, error) {
	if msg == "" {
		return "", ErrEmptyMessage
	}

	err := detectNonDirtyMessage(msg)
	if err != nil {
		return "", WrapApiError(ErrHentaiMessage, err)
	}

	return msg, nil
}

var (
	ErrMessageLost = NewApiError(http.StatusBadRequest, "message lost", nil)
	ErrYouDirty    = NewApiError(http.StatusForbidden, "you dirty", nil)
)

type Input struct {
	Msg string `json:msg`
}

func example(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	name := ps.ByName("name")
	if name != "say" {
		err := errors.New("must be /error_test/say")
		return NewApiError(http.StatusNotFound, "not found", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// this statement won't wrap error, it just concatenates error message
		return fmt.Errorf("body read failed: %v", err)
	}

	if len(body) < 2 {
		return errors.New("body is not enough")
	}

	var input Input
	if err = json.Unmarshal(body, &input); err != nil {
		return NewApiError(http.StatusBadRequest, "bad request", err)
	}

	resp, err := processMessage(input.Msg)
	if IsErrorEqual(err, ErrEmptyMessage) {
		return WrapApiError(ErrMessageLost, err)
	}
	if IsErrorEqual(err, ErrHentaiMessage) {
		return WrapApiError(ErrYouDirty, err)
	}
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write([]byte("you sent:\n"))
	w.Write([]byte(resp))
	w.Write([]byte("\n"))
	return nil
}

func main() {
	router := httprouter.New()
	router.POST("/error_test/:name", ErrAwareHandle(example))

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		panic(err)
	}
}

// client example
// curl -v -X POST --data '{"msg": "FQ2"}' http://localhost:8000/error_test/say
```

----

<!-- .slide: data-transition="convex" -->

#### Integrate with Error Type from Foreign Packages

* If we use internal error types from
  * golang.org/x/xerrors
  * github.com/pkg/errors
* We could store error message like below

```go=
"403|forbidden: <error description>: <error description>: ..."
```
<!-- .element: class="fragment" data-fragment-index="1" -->

---

<!-- .slide: data-transition="convex" -->

## One More Thing

----

<!-- .slide: data-transition="convex" -->

### Fail Fast

* When error occurs:
  * Stop immediately
  * Report the error as early as possible

----

<!-- .slide: data-transition="convex" -->

### Idiom Error Handling

```go=
func MyFunc(param1 int, param2 string) (string, error) {
	// do something...
	f, err := os.Open("filename.ext")
	if err != nil {
		return "", err
	}
	// do something...
	n, err := io.ReadFull(f, buf1)
	if err != nil {
		return "", err
	}
	// do something...
	n, err = f.Read(buf2)
	if err != nil {
		return "", err
	}
	// do something...
	return result, nil
}
```

----

<!-- .slide: data-transition="convex" -->

### Error Handling in Foreign Language

![rust](https://www.rust-lang.org/static/images/rust-logo-blk.svg)

----

<!-- .slide: data-transition="convex" -->

```rust=
panic!()

Result<T, E>
Option<T>

myobj.myfn()? // ? operator
```

<small>https://doc.rust-lang.org/book/ch09-00-error-handling.html</small>

----

<!-- .slide: data-transition="convex" -->

```rust=
e.source()
e.backtrace()
```

<small>https://doc.rust-lang.org/std/error/trait.Error.html</small>

----

<!-- .slide: data-transition="convex" -->

### Evolution?

```go=
func MyFunc(param1 int, param2 string) (string, error) {
	// do something...
	f := try(os.Open("filename.ext"))
	// do something...
	n := try(io.ReadFull(f, buf1))
	// do something...
	n = try(f.Read(buf2))
	// do something...
	return result, nil
}
```

----

<!-- .slide: data-transition="convex" -->

### Impossible!?

Because something mess up:

```go=
if err != nil {
	if err != io.EOF {
		return "", err
	} else {
		return result, nil
	}
}
```

----

<!-- .slide: data-transition="convex" -->

人生ｵﾜﾀ ＼(\^o\^)／ I am done!

----

<!-- .slide: data-transition="convex" -->

### Evolution Again

Add handle:

```go=
func MyFunc(param1 int, param2 string) (string, error) {
	handle err {
		if err != io.EOF {
			return "", err
		} else {
			return result, nil
		}
	}
	// do something...
	f := check os.Open("filename.ext")
	// do something...
	n := check io.ReadFull(f, buf1)
	// do something...
	n = check f.Read(buf2)
	// do something...
	return result, nil
}
```

<small>https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling.md</small>

----

<!-- .slide: data-transition="convex" -->

![讓我們看下去](https://stickershop.line-scdn.net/stickershop/v1/sticker/16846578/iPhone/sticker@2x.png =531x495)

<small>https://store.line.me/stickershop/product/8601/zh-Hant</small>

---

<!-- .slide: data-transition="convex" -->

## Conclusion

----

<!-- .slide: data-transition="convex" -->

### Advice

* Wrapping
  * Wrap as early as possible
* Inspection
  * Collect and log error at entry function or middleware
* Displaying
  * Only show the most outer error message to end-user

----

<!-- .slide: data-transition="convex" -->

### Promote

* :thumbsup: github.com/pkg/errors by Dave Cheney
  * https://github.com/pkg/errors
  * https://godoc.org/github.com/pkg/errors

----

<!-- .slide: data-transition="convex" -->

### Bless

:hash: {錯誤抓精光 碼農發大財|Errors get eliminated, coders get richer.}

> [name=郭學聰 Hsueh-Tsung Kuo] [time=2020_11_14] [color=red] :notebook:

---

<!-- .slide: data-transition="convex" -->

## Reference

----

<!-- .slide: data-transition="convex" -->

### Resources

* Go 2 Draft Designs
  * <small>https://go.googlesource.com/proposal/+/master/design/go2draft.md</small>
* Proposal: Go 2 Error Inspection
  * <small>https://go.googlesource.com/proposal/+/master/design/29934-error-values.md</small>

----

<!-- .slide: data-transition="convex" -->

### Resources

* Working with Errors in Go 1.13
  * https://blog.golang.org/go1.13-errors
* golang.org/x/xerrors by Go team members
  * https://github.com/golang/xerrors
  * https://godoc.org/golang.org/x/xerrors
* github.com/pkg/errors by Dave Cheney
  * https://github.com/pkg/errors
  * https://godoc.org/github.com/pkg/errors

---

<!-- .slide: data-transition="zoom" -->

## Q&A

---

<style>
.reveal {
    background: #FFDFEF;
    color: black;
}
.reveal h2,
.reveal h3,
.reveal h4,
.reveal h5,
.reveal h6 {
    color: black;
}
.reveal code {
    font-size: 18px !important;
    line-height: 1.2;
}

.progress div{
height:14px !important;
background: hotpink !important;
}

// original template

.rightpart{
    float:right;
    width:50%;
}

.leftpart{
    margin-right: 50% !important;
    height:50%;
}
.reveal section img { background:none; border:none; box-shadow:none; }
p.blo {
	font-size: 50px !important;
	background:#B6BDBB;
	border:1px solid silver;
	display:inline-block;
	padding:0.5em 0.75em;
	border-radius: 10px;
	box-shadow: 5px 5px 5px #666;
}

p.blo1 {
	background: #c7c2bb;
}
p.blo2 {
	background: #b8c0c8;
}
p.blo3 {
	background: #c7cedd;
}

p.bloT {
	font-size: 60px !important;
	background:#B6BDD3;
	border:1px solid silver;
	display:inline-block;
	padding:0.5em 0.75em;
	border-radius: 8px;
	box-shadow: 1px 2px 5px #333;
}
p.bloA {
	background: #B6BDE3;
}
p.bloB {
	background: #E3BDB3;
}

/*.slide-number{
	margin-bottom:10px !important;
	width:100%;
	text-align:center;
	font-size:25px !important;
	background-color:transparent !important;
}*/
iframe.myclass{
	width:100px;
	height:100px;
	bottom:0;
	left:0;
	position:fixed;
	border:none;
	z-index:99999;
}
h1.raw {
	color: #fff;
	background-image: linear-gradient(90deg,#f35626,#feab3a);
	-webkit-background-clip: text;
	-webkit-text-fill-color: transparent;
	animation: hue 5s infinite linear;
}
@keyframes hue {
	from {
	  filter: hue-rotate(0deg);
	}
	to {
	  filter: hue-rotate(360deg);
	}
}
.progress{
height:14px !important;
}

.progress span{
height:14px !important;
background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAAMCAIAAAAs6UAAAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyJpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMy1jMDExIDY2LjE0NTY2MSwgMjAxMi8wMi8wNi0xNDo1NjoyNyAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNiAoV2luZG93cykiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6QUNCQzIyREQ0QjdEMTFFMzlEMDM4Qzc3MEY0NzdGMDgiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6QUNCQzIyREU0QjdEMTFFMzlEMDM4Qzc3MEY0NzdGMDgiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDpBQ0JDMjJEQjRCN0QxMUUzOUQwMzhDNzcwRjQ3N0YwOCIgc3RSZWY6ZG9jdW1lbnRJRD0ieG1wLmRpZDpBQ0JDMjJEQzRCN0QxMUUzOUQwMzhDNzcwRjQ3N0YwOCIvPiA8L3JkZjpEZXNjcmlwdGlvbj4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gPD94cGFja2V0IGVuZD0iciI/PovDFgYAAAAmSURBVHjaYvjPwMAAxjMZmBhA9H8INv4P4TPM/A+m04zBNECAAQBCWQv9SUQpVgAAAABJRU5ErkJggg==") repeat-x !important;

}

.progress span:after,
.progress span.nyancat{
	content: "";
	background: url('data:image/gif;base64,R0lGODlhIgAVAKIHAL3/9/+Zmf8zmf/MmZmZmf+Z/wAAAAAAACH/C05FVFNDQVBFMi4wAwEAAAAh/wtYTVAgRGF0YVhNUDw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMy1jMDExIDY2LjE0NTY2MSwgMjAxMi8wMi8wNi0xNDo1NjoyNyAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wTU09Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9tbS8iIHhtbG5zOnN0UmVmPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvc1R5cGUvUmVzb3VyY2VSZWYjIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDpDMkJBNjY5RTU1NEJFMzExOUM4QUM2MDAwNDQzRERBQyIgeG1wTU06RG9jdW1lbnRJRD0ieG1wLmRpZDpCREIzOEIzMzRCN0IxMUUzODhEQjgwOTYzMTgyNTE0QiIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDpCREIzOEIzMjRCN0IxMUUzODhEQjgwOTYzMTgyNTE0QiIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgQ1M2IChXaW5kb3dzKSI+IDx4bXBNTTpEZXJpdmVkRnJvbSBzdFJlZjppbnN0YW5jZUlEPSJ4bXAuaWlkOkM1QkE2NjlFNTU0QkUzMTE5QzhBQzYwMDA0NDNEREFDIiBzdFJlZjpkb2N1bWVudElEPSJ4bXAuZGlkOkMyQkE2NjlFNTU0QkUzMTE5QzhBQzYwMDA0NDNEREFDIi8+IDwvcmRmOkRlc2NyaXB0aW9uPiA8L3JkZjpSREY+IDwveDp4bXBtZXRhPiA8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') !important;
   width: 34px !important;
   height: 21px !important;
   border: none !important;
   float:right;
   margin-top:-7px;
   margin-right:-10px;
}
</style>

