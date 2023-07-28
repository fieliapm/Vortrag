---
title: Go to Generics = Go to hell?
tags: Slide, Go, Generics Programming
description: View the slide with "Slide Mode".
slideOptions:
  spotlight:
    enabled: false
  allottedMinutes: 25
---
<small>Go to Generics = Go to hell?</small>
===

<!-- .slide: data-background-color="pink" -->
<!-- .slide: data-transition="zoom" -->

What I experienced
after Generics Introduced into Go

> <small>skilled :angry:</small>

> [name=郭學聰 Hsueh-Tsung Kuo]
> [time=Sun, 30 Jul 2023] [color=red]

###### CC BY-SA 4.0

---

<!-- .slide: data-transition="convex" -->

## Who am I?

![fieliapm](https://www.gravatar.com/avatar/2aef78f04240a6ac9ccd473ba1cbd1e3?size=2048 =384x384)

<small>Someone (who?) said:
a game programmer should be able to draw cute anime character(?)</small>

----

<!-- .slide: data-transition="convex" -->

* A ~~programmer~~ coding peasant from game company in Taiwan.
* Backend (and temporary frontend) engineer.
* Usually develop something related to my work in Python, Ruby, ECMAScript, Golang, C#.
* Built CDN-aware game asset update system.
* Business large passenger vehicle driver. :bus:
* Ride bike to help traffic jam. :racing_motorcycle:
* Care about chaotic traffic in Taiwan.
* Draw cute anime character in spare time.

---

<!-- .slide: data-transition="convex" -->

## Outline

----

<!-- .slide: data-transition="convex" -->

4. Introduction to Generics
    * Languages
    * Boxing &amp; Monomorphization
5. Generics in Go
    * Usage
    * Type Constraint
    * Mixed Constraint

----

<!-- .slide: data-transition="convex" -->

6. Pitfalls
    * Usage
      * Generic Methods
      * Underlying Type
      * Structure Members
    * Performance
      * Monomorphic vs Interface vs Generics
        * GC Shape Stenciling
      * Callbacks

----

<!-- .slide: data-transition="convex" -->

7. Possible Use Cases
    * Collections
    * Callback
    * I/O of Data Structures
      * Serializers
      * Datastores
    * Plugins
    * Mocks
8. Conclusion
9. Resource
10. Q&A

---

<!-- .slide: data-transition="convex" -->

## Introduction to Generics

----

<!-- .slide: data-transition="convex" -->

### Generics

"algorithms are written in terms of data types to-be-specified-later that are then instantiated when needed for specific types provided as parameters."

----

<!-- .slide: data-transition="convex" -->

### Languages

C++

```cpp=
template <typename T>
T max (T a, T b) {
    return a>b?a:b;
}

max(1, 2);

template<typename T>
class list {
public:
    void append(T data) {
        ...
    }
}

list<string> ls;
```

----

<!-- .slide: data-transition="convex" -->

### Languages

Java

```java=
public class List<T> {
    public void append(T data) {
        ...
    }
}

List<String> list = new List<String>();
```

----

<!-- .slide: data-transition="convex" -->

### Languages

Rust

```rust=
let mut language_codes: HashMap<&str, &str> = HashMap::new();
```

----

<!-- .slide: data-transition="convex" -->

### Boxing &amp; Monomorphization

* Boxing
  * run-time
  * placing a primitive type within an object.
* Monomorphization
  * compile-time
  * polymorphic functions :arrow_right: monomorphic functions for each unique instantiation.

----

<!-- .slide: data-transition="convex" -->

#### Boxing

C#

```csharp=
int i = 300;
object o = i;   // boxing
int j = (int)o; // unboxing
```

----

<!-- .slide: data-transition="convex" -->

#### Monomorphization

Rust

```rust=
fn id<T>(x: T) -> T {
    return x;
}

fn main() {
    let int = id(10);
    let string = id("some text");
    println!("{int}, {string}");
}
```

```rust=
fn id_i32(x: i32) -> i32 {
    return x;
}

fn id_str(x: &str) -> &str {
    return x;
}

fn main() {
    let int = id_i32(10);
    let string = id_str("some text");
    println!("{int}, {string}");
}
```

---

<!-- .slide: data-transition="convex" -->

## Generics in Go

----

<!-- .slide: data-transition="convex" -->

### Usage

```go=
package main

import "fmt"

func MapKeys[K comparable, V any](m map[K]V) []K {
    r := make([]K, 0, len(m))
    for k := range m {
        r = append(r, k)
    }
    return r
}

type List[T any] struct {
    head, tail *element[T]
}

type element[T any] struct {
    next *element[T]
    val  T
}

func (lst *List[T]) Push(v T) {
    if lst.tail == nil {
        lst.head = &element[T]{val: v}
        lst.tail = lst.head
    } else {
        lst.tail.next = &element[T]{val: v}
        lst.tail = lst.tail.next
    }
}

func (lst *List[T]) GetAll() []T {
    var elems []T
    for e := lst.head; e != nil; e = e.next {
        elems = append(elems, e.val)
    }
    return elems
}

func main() {
    var m = map[int]string{1: "2", 2: "4", 4: "8"}

    fmt.Println("keys:", MapKeys(m))

    _ = MapKeys[int, string](m)

    lst := List[int]{}
    lst.Push(10)
    lst.Push(13)
    lst.Push(23)
    fmt.Println("list:", lst.GetAll())
}
```

<small>https://gobyexample.com/generics</small>

----

<!-- .slide: data-transition="convex" -->

### Type Constraint

```go=
import "golang.org/x/exp/constraints"

func GMin[T constraints.Ordered](x, y T) T {
    if x < y {
        return x
    }
    return y
}

x := GMin[int](2, 3)
y := GMin(2, 3) // type inference
```

----

<!-- .slide: data-transition="convex" -->

### Type Constraint

```go=
type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
    Signed | Unsigned
}

type Float interface {
    ~float32 | ~float64
}

type Ordered interface {
    Integer | Float | ~string
}
```

----

<!-- .slide: data-transition="convex" -->

### Mixed Constraint

```go=
package main

import (
    "fmt"
)

type MyType[T any] interface {
    MyType1 | MyType2
    // and
    Concat(T) T
}

type MyType1 struct {
    Str string
}

func (mt MyType1) Concat(s string) string {
    fmt.Println("Concat MyType1")
    return mt.Str + s
}

type MyType2 struct {
    Int int
}

func (mt MyType2) Concat(i int) int {
    fmt.Println("Concat MyType2")
    return mt.Int + i
}

func ConcatMyType[T ~string | ~int, MT MyType[T]](mt MT, v T) T {
    return mt.Concat(v)
}

func main() {
    mt1 := MyType1{Str: "Sak"}
    mt2 := MyType2{Int: 30}

    fmt.Println(ConcatMyType[string](mt1, "ura"))
    fmt.Println(ConcatMyType[int](mt2, 5))

    // type inference
    fmt.Println(ConcatMyType(mt1, "ura"))
    fmt.Println(ConcatMyType(mt2, 5))
}
```

<small>好油喔peko :rabbit:</small> <!-- .element: class="fragment" data-fragment-index="1" -->

---

<!-- .slide: data-transition="convex" -->

## Pitfalls

----

<!-- .slide: data-transition="convex" -->

### Usage

----

<!-- .slide: data-transition="convex" -->

#### Generic Methods

```go=
package main

import (
    "fmt"

    "golang.org/x/exp/constraints"
)

type List[T constraints.Ordered] struct {
    list []T
}

func (ls *List[T]) Append(v T) {
    ls.list = append(ls.list, v)
}

func (ls *List[T]) All() []T {
    return ls.list
}

func (ls *List[T]) String() string {
    return fmt.Sprintf("%v", ls.list)
}

// syntax error: method must have no type parameters
func (ls *List[T]) String[R ~string | ~[]byte]() R {
    return fmt.Sprintf("%v", ls.list)
}

func ToString[T constraints.Ordered, R ~string | ~[]byte](ls List[T]) R {
    return R(fmt.Sprintf("%v", ls.list))
}

func main() {
    var ls List[int]
    ls.Append(1)
    ls.Append(2)
    ls.Append(3)

    fmt.Println(ls.All())
    fmt.Println(ls.String())
    fmt.Println(ToString[int, []byte](ls))
}
```

----

<!-- .slide: data-transition="convex" -->

#### Underlying Type

```go=
package main

type Type struct{}

type MyType Type

func F[T ~struct{}](t T) {}

// invalid use of ~ (underlying type of Type is struct{})
func F[T ~Type](t T) {}

func main() {
    x := Type{}
    y := MyType{}
    F(x)
    F(y)
}
```

----

<!-- .slide: data-transition="convex" -->

#### Structure Members

```go=
package main

import "fmt"

type Type struct {
    Str string
    Int int
}

type MyType Type

func F[T ~struct {
    Str string
    Int int
}](t T) string {
    // t.Str undefined (type T has no field or method Str)
    // t.Int undefined (type T has no field or method Int)
    return fmt.Sprintf("%v %v", t.Str, t.Int)
}

func main() {
    x := Type{}
    y := MyType{}
    fmt.Println(F(x))
    fmt.Println(F(y))
}
```

but permit to access interface methods <!-- .element: class="fragment" data-fragment-index="1" -->

----

<!-- .slide: data-transition="convex" -->

### Performance

go 1.18~1.20

----

<!-- .slide: data-transition="convex" -->

#### Monomorphic vs Interface vs Generics

```go=
package main

import (
    "bytes"
    "fmt"
    "io"
    "time"
)

type Data struct {
    Name  string
    Value string
}

func SerializeM(w *bytes.Buffer, d Data) {
    b := []byte(d.Name + " " + d.Value)
    w.Write(b)
}

func SerializeI(w io.Writer, d Data) {
    b := []byte(d.Name + " " + d.Value)
    w.Write(b)
}

func SerializeG[W io.Writer](w W, d Data) {
    b := []byte(d.Name + " " + d.Value)
    w.Write(b)
}

func main() {
    var buf bytes.Buffer
    var bw io.ReadWriter
    bw = &buf
    data := Data{Name: "Sakura", Value: "35"}
    startTime := time.Now()
    for i := 0; i < 100000; i++ {
        //SerializeM(&buf, data) // 3.614264ms
        //SerializeI(&buf, data) // 3.667853ms
        //SerializeG(&buf, data) // 6.934544ms
        SerializeG(bw, data)   // 7.680705ms
    }
    endTime := time.Now()
    fmt.Println(endTime.Sub(startTime))
}
```

----

<!-- .slide: data-transition="convex" -->

#### Monomorphic vs Interface vs Generics

Monomorphic

```go=
"".SerializeM STEXT size=197 args=0x28 locals=0x88 funcid=0x0 align=0x0

0x0065 00101 (main.go:17)       MOVQ    CX, DI
0x0068 00104 (main.go:17)       MOVQ    BX, CX
0x006b 00107 (main.go:17)       MOVQ    AX, BX
0x006e 00110 (main.go:17)       MOVQ    ""..autotmp_9+120(SP), AX
0x0073 00115 (main.go:17)       PCDATA  $1, $2
0x0073 00115 (main.go:17)       CALL    bytes.(*Buffer).Write(SB)
```

----

<!-- .slide: data-transition="convex" -->

#### Monomorphic vs Interface vs Generics

Interface

```go=
"".SerializeI STEXT size=200 args=0x30 locals=0x60 funcid=0x0 align=0x0

0x0021 00033 (main.go:22)       MOVQ    AX, "".w+104(SP)
0x0026 00038 (main.go:22)       MOVQ    BX, "".w+112(SP)
0x002b 00043 (main.go:22)       PCDATA  $3, $2

0x005a 00090 (main.go:22)       MOVQ    "".w+104(SP), DX
0x005f 00095 (main.go:22)       MOVQ    24(DX), DX
0x0063 00099 (main.go:22)       MOVQ    CX, DI
0x0066 00102 (main.go:22)       MOVQ    BX, CX
0x0069 00105 (main.go:22)       MOVQ    AX, BX
0x006c 00108 (main.go:22)       MOVQ    "".w+112(SP), AX
0x0071 00113 (main.go:22)       PCDATA  $1, $2
0x0071 00113 (main.go:22)       CALL    DX
```

----

<!-- .slide: data-transition="convex" -->

#### Monomorphic vs Interface vs Generics

Generics (pointer arg)

```go=
"".SerializeG[go.shape.*uint8_0] STEXT dupok size=200 args=0x30 locals=0x60 funcid=0x0 align=0x0

0x005a 00090 (main.go:27)       PCDATA  $0, $-2
0x005a 00090 (main.go:27)       MOVQ    ""..dict+104(SP), DX
0x005f 00095 (main.go:27)       PCDATA  $0, $-1
0x005f 00095 (main.go:27)       MOVQ    16(DX), DX
0x0063 00099 (main.go:27)       MOVQ    24(DX), DX
0x0067 00103 (main.go:27)       MOVQ    CX, DI
0x006a 00106 (main.go:27)       MOVQ    BX, CX
0x006d 00109 (main.go:27)       MOVQ    AX, BX
0x0070 00112 (main.go:27)       MOVQ    "".w+112(SP), AX
0x0075 00117 (main.go:27)       PCDATA  $1, $2
0x0075 00117 (main.go:27)       CALL    DX
```

----

<!-- .slide: data-transition="convex" -->

#### Monomorphic vs Interface vs Generics

Generics (interface arg)

```go=
"".main STEXT size=395 args=0x0 locals=0xb0 funcid=0x0 align=0x0

0x00b3 00179 (main.go:27)       LEAQ    go.itab.*bytes.Buffer,io.ReadWriter(SB), BX

"".SerializeG[go.shape.interface { Read([]uint8) (int, error); Write([]uint8) (int, error) }_0] STEXT dupok size=253 args=0x38 locals=0x80 funcid=0x0 align=0x0

0x0030 00048 (main.go:27)       MOVQ    CX, "".w+152(SP)
0x0038 00056 (main.go:27)       PCDATA  $3, $2
0x0038 00056 (main.go:27)       MOVQ    BX, ""..autotmp_10+112(SP)

0x0075 00117 (main.go:27)       LEAQ    type.io.Writer(SB), AX
0x007c 00124 (main.go:27)       MOVQ    ""..autotmp_10+112(SP), BX
0x0081 00129 (main.go:27)       PCDATA  $1, $2
0x0081 00129 (main.go:27)       CALL    runtime.assertI2I(SB)
0x0086 00134 (main.go:27)       MOVQ    24(AX), DX
0x008a 00138 (main.go:27)       MOVQ    "".w+152(SP), AX
0x0092 00146 (main.go:27)       MOVQ    "".b.ptr+104(SP), BX
0x0097 00151 (main.go:27)       MOVQ    "".b.len+56(SP), CX
0x009c 00156 (main.go:27)       MOVQ    "".b.cap+64(SP), DI
0x00a1 00161 (main.go:27)       PCDATA  $1, $3
0x00a1 00161 (main.go:27)       CALL    DX
```

----

<!-- .slide: data-transition="convex" -->

##### GC Shape Stenciling

* Stencil the code for each different GC shape of the instantiated types.
* Using a dictionary to handle differing behaviors of types that have the same shape.
  * `*AnyType` :arrow_right: `SerializeG[go.shape.*uint8_0]`
  * `Interface` :arrow_right: `SerializeG[go.shape.interface { ... }_0]`

----

<!-- .slide: data-transition="convex" -->

##### Interface Table

```go=
type iface struct {
    tab *itab
    data unsafe.Pointer
}

type itab struct {
    inter *interfacetype // offset 0
    _type *_type // offset 8
    hash  uint32 // offset 16
    _     [4]byte
    fun   [1]uintptr // offset 24...
}
```

----

<!-- .slide: data-transition="convex" -->

##### :shit: :shit: :shit:

pointer arg :arrow_right: `dict` + `iface.itab.fun` :arrow_right: :shit:
interface arg :arrow_right: `runtime.assertI2I()` + `iface.itab.fun` :arrow_right: :shit:

----

<!-- .slide: data-transition="convex" -->

##### Exception

`string` :arrow_right: `SequenceToString[go.shape.string_0]`
`[]byte` :arrow_right: `SequenceToString[go.shape.[]uint8_0]`

```go=
type Byteseq interface {
    ~string | ~[]byte
}

func SequenceToString[T Byteseq](x T) string {
    ix := any(x)
    switch ix.(type) {
    case string:
        return ix.(string)
    case []byte:
        p := ix.([]byte)
        return *(*string)(unsafe.Pointer(&p))
    default:
        return ""
    }
}
```

<small>https://github.com/koykov/byteseq</small>

----

<!-- .slide: data-transition="convex" -->

#### Callbacks

----

<!-- .slide: data-transition="convex" -->

##### Callbacks

plain

```go=
package main

import "fmt"

func MapInt(a []int, callback func(int) int) []int {
    for n, elem := range a {
        a[n] = callback(elem)
    }
    return a
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := MapInt(input, func(i int) int {
        return i + 7
    })
    fmt.Println(output)
}
```

----

<!-- .slide: data-transition="convex" -->

##### Callbacks

plain (inlined)

```go=
0x0074 00116 (main.go:14)       XORL    CX, CX
0x0076 00118 (main.go:6)        JMP     128
0x0078 00120 (main.go:7)        ADDQ    $7, (AX)(CX*8)
0x007d 00125 (main.go:6)        INCQ    CX
0x0080 00128 (main.go:6)        CMPQ    CX, $10
0x0084 00132 (main.go:6)        JLT     120

"".MapInt STEXT size=170 args=0x20 locals=0x20 funcid=0x0 align=0x0
```

----

<!-- .slide: data-transition="convex" -->

#### Callbacks

parametrize callback parameter

```go=
package main

import "fmt"

func MapAny[I any](a []I, callback func(I) I) []I {
    for n, elem := range a {
        a[n] = callback(elem)
    }
    return a
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := MapAny(input, func(i int) int {
        return i + 7
    })
    fmt.Println(output)
}
```

----

<!-- .slide: data-transition="convex" -->

#### Callbacks

parametrize callback parameter (not inlined)

```go=
0x0079 00121 (main.go:14)       XORL    CX, CX
0x007b 00123 (main.go:6)        JMP     173
0x007d 00125 (main.go:6)        MOVQ    CX, "".n+40(SP)
0x0082 00130 (main.go:6)        MOVQ    (AX)(CX*8), BX
0x0086 00134 (main.go:7)        MOVQ    "".main.func1·f(SB), SI
0x008d 00141 (main.go:7)        LEAQ    "".main.func1·f(SB), DX
0x0094 00148 (main.go:7)        MOVQ    BX, AX
0x0097 00151 (main.go:7)        PCDATA  $1, $1
0x0097 00151 (main.go:7)        CALL    SI
0x0099 00153 (main.go:7)        MOVQ    "".n+40(SP), CX
0x009e 00158 (main.go:7)        MOVQ    ""..autotmp_38+48(SP), BX
0x00a3 00163 (main.go:7)        MOVQ    AX, (BX)(CX*8)
0x00a7 00167 (main.go:6)        INCQ    CX
0x00aa 00170 (main.go:6)        MOVQ    BX, AX
0x00ad 00173 (main.go:6)        CMPQ    CX, $10
0x00b1 00177 (main.go:6)        JLT     125

"".MapAny[go.shape.int_0] STEXT dupok size=186 args=0x28 locals=0x20 funcid=0x0 align=0x0
```

----

<!-- .slide: data-transition="convex" -->

#### Callbacks

parametrize callback

```go=
package main

import "fmt"

func MapInt[F func(int) int](a []int, callback F) []int {
    for n, elem := range a {
        a[n] = callback(elem)
    }
    return a
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := MapInt(input, func(i int) int {
        return i + 7
    })
    fmt.Println(output)
}
```

----

<!-- .slide: data-transition="convex" -->

#### Callbacks

parametrize callback (inlined)

```go=
0x0074 00116 (main.go:14)       XORL    CX, CX
0x0076 00118 (main.go:6)        JMP     128
0x0078 00120 (main.go:7)        ADDQ    $7, (AX)(CX*8)
0x007d 00125 (main.go:6)        INCQ    CX
0x0080 00128 (main.go:6)        CMPQ    CX, $10
0x0084 00132 (main.go:6)        JLT     120

"".MapInt[go.shape.func(int) int_0] STEXT dupok size=186 args=0x28 locals=0x20 funcid=0x0 align=0x0
```

----

<!-- .slide: data-transition="convex" -->

#### Callbacks

parametrize callback and its parameter

```go=
package main

import "fmt"

func MapAny[I any, F func(I) I](a []I, callback F) []I {
    for n, elem := range a {
        a[n] = callback(elem)
    }
    return a
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := MapAny(input, func(i int) int {
        return i + 7
    })
    fmt.Println(output)
}
```

----

<!-- .slide: data-transition="convex" -->

#### Callbacks

parametrize callback and its parameter (inlined)

```go=
0x0074 00116 (main.go:14)       XORL    CX, CX
0x0076 00118 (main.go:6)        JMP     135
0x0078 00120 (main.go:6)        MOVQ    (AX)(CX*8), DX
0x007c 00124 (<unknown line number>)    NOP
0x007c 00124 (main.go:15)       ADDQ    $7, DX
0x0080 00128 (main.go:7)        MOVQ    DX, (AX)(CX*8)
0x0084 00132 (main.go:6)        INCQ    CX
0x0087 00135 (main.go:6)        CMPQ    CX, $10
0x008b 00139 (main.go:6)        JLT     120

"".MapAny[go.shape.int_0,go.shape.func(int) int_1] STEXT dupok size=186 args=0x28 locals=0x20 funcid=0x0 align=0x0
```

---

<!-- .slide: data-transition="convex" -->

## Possible Use Cases

----

<!-- .slide: data-transition="convex" -->

### Collections

```go=
type List[T constraints.Ordered] struct {
    list []T
}

func (ls *List[T]) Append(v T) {
    ls.list = append(ls.list, v)
}

func (ls *List[T]) All() []T {
    return ls.list
}
```

good practice :o:
- T is pointer or interface :x:

----

<!-- .slide: data-transition="convex" -->

### Callback

```go=
package main

import "fmt"

func MapAny[I any, F func(I) I](a []I, callback F) []I {
    for n, elem := range a {
        a[n] = callback(elem)
    }
    return a
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := MapAny(input, func(i int) int {
        return i + 7
    })
    fmt.Println(output)
}
```

good practice :o:
- do parametrize callback type :warning:

----

<!-- .slide: data-transition="convex" -->

### I/O of Data Structures

----

<!-- .slide: data-transition="convex" -->

#### Serializers

```go=
func SerializeG[W io.Writer](w W, d Data) {
    b := []byte(d.Name + " " + d.Value)
    w.Write(b)
}
```

bad practice :x:
- performance down

----

<!-- .slide: data-transition="convex" -->

#### Datastores

```go=
func SerializeG[W io.Writer](w W, d Data) {
    b := []byte(d.Name + " " + d.Value)
    w.Write(b)
}

func Datastore(d Data) {
    SerializeG[TCPWriter](tcpW, d)
}
```

bad practice :x:
- performance down
- I/O latency >>>>>> code performance

----

<!-- .slide: data-transition="convex" -->

### Plugins

```go=
type Plugin[T any] interface {
    Method(T) T
}
```

acceptable? :thinking_face:
- T is pointer or interface :x:
- however...

----

<!-- .slide: data-transition="convex" -->

### Mocks

```go=
type Foo interface {
    Bar(x int) int
}

func SUT(f Foo) {
    // ...
}

func TestFoo(t *testing.T) {
    ctrl := gomock.NewController(t)

    // Assert that Bar() is invoked.
    defer ctrl.Finish()

    m := NewMockFoo(ctrl)

    // Asserts that the first and only call to Bar() is passed 99.
    // Anything else will fail.
    m.
        EXPECT().
        Bar(gomock.Eq(99)).
        Return(101)

    SUT(m)
}
```

practically impossible :x:
- not working in many cases :-1:

---

<!-- .slide: data-transition="convex" -->

## Conclusion

----

<!-- .slide: data-transition="convex" -->

### Prefer

* Use generics in data structures.
  * Collections
  * Vectors &amp; Matrices &amp; Tensors
* Pass value types to generic functions.
* Parametrize callback types.
* Parametrize constraint:
  * `type Byteseq interface{ ~string | ~[]byte }`

----

<!-- .slide: data-transition="convex" -->

### Forbidden

* :x: Parametrize any method parameter.
* :x: Attempt to pass any pointer arg :arrow_right: expect generic function monomorphized &amp; inlined.
  * Twice dereference when calling interface method.
* :no_entry: Pass interface arg to generic function.
  * Convert interface then dereference when calling interface method.

----

<!-- .slide: data-transition="convex" -->

### Slogan

:hash: {善用泛型,不濫用泛型|<big>Make good use of generics, don't abuse generics.</big>}

> [name=郭學聰 Hsueh-Tsung Kuo] [time=2023_07_30] [color=red] :notebook:

---

<!-- .slide: data-transition="convex" -->

## Resource 

----

<!-- .slide: data-transition="convex" -->

### Reference

* Generics can make your Go code slower[^ref1]
* An Introduction To Generics[^ref2]
* When To Use Generics[^ref3]
* No parameterized methods[^ref4]

[^ref1]:<small>https://planetscale.com/blog/generics-can-make-your-go-code-slower</small>
[^ref2]:<small>https://go.dev/blog/intro-generics</small>
[^ref3]:<small>https://go.dev/blog/when-generics</small>
[^ref4]:<small>https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods</small>

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

