Marshal
===

<!-- .slide: data-background="#FFDFEF" -->
<!-- .slide: data-transition="zoom" -->

不給你的,自己來!!! :dizzy:

> [name=郭學聰 Hsueh-Tsung Kuo] [time=Sun, 06 Aug 2017] [color=red]

---

<!-- .slide: data-transition="convex" -->

## who am I?

![fieliapm](https://pbs.twimg.com/profile_images/591670980021387264/aZAYLRUe_400x400.png)

----

<!-- .slide: data-transition="convex" -->

* programmer from Rayark, a game company in Taiwan
* backend engineer
* usually develop something related to my work in Python, Ruby, Golang, C#
* built almost entire VOEZ game server by myself only

---

<!-- .slide: data-transition="convex" -->

## conclusion

#### Too Long; Didn't Read

----

<!-- .slide: data-transition="convex" -->

```go=
type UnixTime struct {
	//anonymous field
	time.Time `bson:",inline" json:",inline"`
}

func (unixTime UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(unixTime.Time.Unix())
}

func (unixTime *UnixTime) UnmarshalJSON(data []byte) error {
	var unixTimeInt int64

	err := json.Unmarshal(data, &unixTimeInt)
	if err != nil {
		return err
	}

	unixTime.Time = time.Unix(unixTimeInt, 0)
	return nil
}
```

----

<!-- .slide: data-transition="convex" -->

#### is this it?

----

<!-- .slide: data-transition="convex" -->

#### there seems to be much behind all this!

---

<!-- .slide: data-transition="convex" -->

## outline

----

<!-- .slide: data-transition="convex" -->

3. conclusion
4. outline
5. preface: type system
6. origin: json package
    1. json marshal/unmarshal package
    2. usage
    3. implementation & example
7. customization
	1. requirement
	2. try it
	3. result
	4. solution

----

<!-- .slide: data-transition="convex" -->

8. new coming problem
    1. case
	2. try it
	3. reason
	4. dirty solution
9. better way
10. summary
    1. slogan
    2. special thanks
11. appendix: code to test marshal/unmarshal
12. Q&A

---

<!-- .slide: data-transition="convex" -->

## preface: type system

----

<!-- .slide: data-transition="convex" -->

```go=
// struct
type Person struct {
	Name string
	Age  int
}

// interface
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

---

<!-- .slide: data-transition="convex" -->

## origin: json package

----

<!-- .slide: data-transition="convex" -->

### json marshal/unmarshal package

* https://golang.org/pkg/encoding/json/

----

<!-- .slide: data-transition="convex" -->

### usage

```go=
// encode
b, err := json.Marshal(&fromData)

enc := json.NewEncoder(writer)
err := enc.Encode(&fromData)

// decode
err := json.Unmarshal(b, &toData)

dec := json.NewDecoder(reader)
err := dec.Decode(&toData)
```

----

<!-- .slide: data-transition="convex" -->

### implementation

```go=
// encoder will call this method
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// decoder will call this method
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}
```

----

<!-- .slide: data-transition="convex" -->

### example

```go=
// can be value or pointer receiver
func (variable Type) MarshalJSON() ([]byte, error)

// must be pointer receiver (need to overwrite variable)
func (variable *Type) UnmarshalJSON([]byte) error
```

---

<!-- .slide: data-transition="convex" -->

## customization

----

<!-- .slide: data-transition="convex" -->

### requirement

* time.Time.MarshalJSON() -> 2017-08-06T13:40:00+08:00
* prefer to use UNIX epoch time -> 1501998000

----

<!-- .slide: data-transition="convex" -->

### try it

----

<!-- .slide: data-transition="convex" -->

```go=
type AccountCoreData struct {
	Name      string
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
// AccountCoreData reflection redefinition
type _AuxAccountCoreData struct {
	*AccountCoreData `bson:",inline" json:",inline"`
	CreatedAt        int64
	UpdatedAt        int64
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (accountCoreData AccountCoreData) MarshalJSON() ([]byte, error) {
	return json.Marshal(_AuxAccountCoreData{
		AccountCoreData: &accountCoreData,
		CreatedAt:       accountCoreData.CreatedAt.Unix(),
		UpdatedAt:       accountCoreData.UpdatedAt.Unix(),
	})
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (accountCoreData *AccountCoreData) UnmarshalJSON(data []byte) error {
	aux := &_AuxAccountCoreData{
		AccountCoreData: accountCoreData,
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	accountCoreData.CreatedAt = time.Unix(aux.CreatedAt, 0)
	accountCoreData.UpdatedAt = time.Unix(aux.UpdatedAt, 0)
	return nil
}
```

----

<!-- .slide: data-transition="convex" -->

### result

#### 然後它就死掉了

* infinite call -> stack overflow!
  * MarshalJSON() -> json.Marshal() -> MarshalJSON() -> json.Marshal() -> ...
  * UnmarshalJSON() -> json.Unmarshal() -> UnmarshalJSON() -> json.Unmarshal() -> ...

----

<!-- .slide: data-transition="convex" -->

### solution

#### avoid inheritance

----

<!-- .slide: data-transition="convex" -->

```go=
// avoid inheritance!!!
type _AccountCoreDataAlias AccountCoreData

// AccountCoreData reflection redefinition
type _AuxAccountCoreData struct {
	*_AccountCoreDataAlias `bson:",inline" json:",inline"`
	CreatedAt              int64
	UpdatedAt              int64
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (accountCoreData AccountCoreData) MarshalJSON() ([]byte, error) {
	return json.Marshal(_AuxAccountCoreData{
		_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountCoreData),
		CreatedAt:             accountCoreData.CreatedAt.Unix(),
		UpdatedAt:             accountCoreData.UpdatedAt.Unix(),
	})
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (accountCoreData *AccountCoreData) UnmarshalJSON(data []byte) error {
	aux := &_AuxAccountCoreData{
		_AccountCoreDataAlias: (*_AccountCoreDataAlias)(accountCoreData),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	accountCoreData.CreatedAt = time.Unix(aux.CreatedAt, 0)
	accountCoreData.UpdatedAt = time.Unix(aux.UpdatedAt, 0)
	return nil
}
```

---

<!-- .slide: data-transition="convex" -->

## new coming problem

----

### case

<!-- .slide: data-transition="convex" -->

```go=
type AccountInfo struct {
	AccountCoreData `bson:",inline" json:",inline"`
	BirthYear       int64
	Description     string
}
```

----

<!-- .slide: data-transition="convex" -->

### try it

call json.Marshal() or json.Unmarshal() ...

* BirthYear and Description are gone

----

<!-- .slide: data-transition="convex" -->

### reason

#### why?

----

<!-- .slide: data-transition="convex" -->

#### receiver fallback

* AccountInfo.MarshalJSON() = AccountCoreData.MarshalJSON()
* AccountInfo.UnmarshalJSON() = AccountCoreData.UnmarshalJSON()

----

<!-- .slide: data-transition="convex" -->

### dirty solution

----

<!-- .slide: data-transition="convex" -->

```go=
// AccountInfo reflection redefinition
type _AuxAccountInfo struct {
	*_AuxAccountCoreData `bson:",inline" json:",inline"`
	BirthYear            *int64
	Description          *string
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (accountInfo AccountInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(_AuxAccountInfo{
		_AuxAccountCoreData: &_AuxAccountCoreData{
			_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountInfo.AccountCoreData),
			CreatedAt:             accountInfo.AccountCoreData.CreatedAt.Unix(),
			UpdatedAt:             accountInfo.AccountCoreData.UpdatedAt.Unix(),
		},
		BirthYear:   &accountInfo.BirthYear,
		Description: &accountInfo.Description,
	})
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (accountInfo *AccountInfo) UnmarshalJSON(data []byte) error {
	aux := &_AuxAccountInfo{
		_AuxAccountCoreData: &_AuxAccountCoreData{
			_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountInfo.AccountCoreData),
		},
		BirthYear:   &accountInfo.BirthYear,
		Description: &accountInfo.Description,
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	accountInfo.AccountCoreData.CreatedAt = time.Unix(aux._AuxAccountCoreData.CreatedAt, 0)
	accountInfo.AccountCoreData.UpdatedAt = time.Unix(aux._AuxAccountCoreData.UpdatedAt, 0)
	return nil
}
```

----

<!-- .slide: data-transition="convex" -->

#### code is like a sh*t!!!

---

<!-- .slide: data-transition="convex" -->

## better way

----

<!-- .slide: data-transition="convex" -->

```go=
type UnixTime time.Time
// MarshalJSON() and UnmarshalJSON() is gone
```

----

<!-- .slide: data-transition="convex" -->

```go=
type UnixTime struct {
	//anonymous field
	time.Time `bson:",inline" json:",inline"`
}
// 1. MarshalJSON() and UnmarshalJSON() is inherited
// 2. keep bson or xml Marshal/Unmarshal
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (unixTime UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(unixTime.Time.Unix())
}
```

----

<!-- .slide: data-transition="convex" -->

```go=
func (unixTime *UnixTime) UnmarshalJSON(data []byte) error {
	var unixTimeInt int64

	err := json.Unmarshal(data, &unixTimeInt)
	if err != nil {
		return err
	}

	unixTime.Time = time.Unix(unixTimeInt, 0)
	return nil
}
```

----

<!-- .slide: data-transition="convex" -->

#### Panacea?

* :x: UnixTime as anonymous field of struct

---

<!-- .slide: data-transition="convex" -->

## summary

----

<!-- .slide: data-transition="convex" -->

### slogan

----

<!-- .slide: data-transition="convex" -->

![而是你要享受這個過程](http://i.imgur.com/SWwsRq3.jpg)
http://i.imgur.com/SWwsRq3.jpg

----

<!-- .slide: data-transition="convex" -->

![不要太看重得與失, 而是你要享受這個過程](http://i.imgur.com/XDq9TWd.gif)
http://i.imgur.com/XDq9TWd.gif

----

<!-- .slide: data-transition="convex" -->

> :hash: "不要太看重得與失, 而是你要享受這個過程"
> [name=田中謙介] [color=red]

there are big stories in type, struct field, interface, and receiver

----

<!-- .slide: data-transition="convex" -->

### special thanks

* Rayark Inc.
  * CTO & CIO
  * backend team
  * other teams

---

<!-- .slide: data-transition="convex" -->

## appendix

----

<!-- .slide: data-transition="convex" -->

### code to test marshal/unmarshal

----

<!-- .slide: data-transition="convex" -->

```go=
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

type UnixTime struct {
	//anonymous field
	time.Time `bson:",inline" json:",inline"`
}

func (unixTime UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(unixTime.Time.Unix())
}

func (unixTime *UnixTime) UnmarshalJSON(data []byte) error {
	var unixTimeInt int64

	err := json.Unmarshal(data, &unixTimeInt)
	if err != nil {
		return err
	}

	unixTime.Time = time.Unix(unixTimeInt, 0)
	return nil
}

// AccountCoreData

type AccountCoreData struct {
	Name      string
	UUID      string
	CreatedAt UnixTime
	UpdatedAt UnixTime
	//CreatedAt time.Time
	//UpdatedAt time.Time
}

/*
// avoid inheritance!!!
type _AccountCoreDataAlias AccountCoreData

// AccountCoreData reflection redefinition
type _AuxAccountCoreData struct {
	//*AccountCoreData       `bson:",inline" json:",inline"`
	*_AccountCoreDataAlias `bson:",inline" json:",inline"`
	CreatedAt              int64
	UpdatedAt              int64
}

func (accountCoreData AccountCoreData) MarshalJSON() ([]byte, error) {
	return json.Marshal(_AuxAccountCoreData{
		//AccountCoreData:       &accountCoreData,
		_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountCoreData),
		CreatedAt:             accountCoreData.CreatedAt.Unix(),
		UpdatedAt:             accountCoreData.UpdatedAt.Unix(),
	})
}

func (accountCoreData *AccountCoreData) UnmarshalJSON(data []byte) error {
	aux := &_AuxAccountCoreData{
		//AccountCoreData: accountCoreData,
		_AccountCoreDataAlias: (*_AccountCoreDataAlias)(accountCoreData),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	accountCoreData.CreatedAt = time.Unix(aux.CreatedAt, 0)
	accountCoreData.UpdatedAt = time.Unix(aux.UpdatedAt, 0)
	return nil
}
*/

// AccountInfo

type AccountInfo struct {
	AccountCoreData `bson:",inline" json:",inline"`
	BirthYear       int64
	Description     string
}

/*
// avoid inheritance!!!
//type _AccountInfoAlias AccountInfo

// AccountInfo reflection redefinition
type _AuxAccountInfo struct {
	*_AuxAccountCoreData `bson:",inline" json:",inline"`
	BirthYear            *int64
	Description          *string
}

func (accountInfo AccountInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(_AuxAccountInfo{
		_AuxAccountCoreData: &_AuxAccountCoreData{
			//AccountCoreData:       &accountInfo.AccountCoreData,
			_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountInfo.AccountCoreData),
			CreatedAt:             accountInfo.AccountCoreData.CreatedAt.Unix(),
			UpdatedAt:             accountInfo.AccountCoreData.UpdatedAt.Unix(),
		},
		BirthYear:   &accountInfo.BirthYear,
		Description: &accountInfo.Description,
	})
}

func (accountInfo *AccountInfo) UnmarshalJSON(data []byte) error {
	aux := &_AuxAccountInfo{
		_AuxAccountCoreData: &_AuxAccountCoreData{
			//AccountCoreData: &accountInfo.AccountCoreData,
			_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountInfo.AccountCoreData),
		},
		BirthYear:   &accountInfo.BirthYear,
		Description: &accountInfo.Description,
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	accountInfo.AccountCoreData.CreatedAt = time.Unix(aux._AuxAccountCoreData.CreatedAt, 0)
	accountInfo.AccountCoreData.UpdatedAt = time.Unix(aux._AuxAccountCoreData.UpdatedAt, 0)
	return nil
}
*/

func main() {
	accountCoreData := AccountCoreData{
		Name:      "Rasmus Faber",
		UUID:      "123e4567-e89b-12d3-a456-426655440000",
		CreatedAt: UnixTime{Time: time.Now()},
		UpdatedAt: UnixTime{Time: time.Unix(1501998000, 0)},
		//CreatedAt: time.Now(),
		//UpdatedAt: time.Unix(1501998000, 0),
	}

	accountInfo := AccountInfo{
		AccountCoreData: accountCoreData,
		BirthYear:       1979,
		Description:     "Swedish pianist, DJ, remixer, composer, record producer, sound engineer, and founder of the record label Farplane Records.",
	}

	fmt.Println("account core data:")
	fmt.Println("original:")
	fmt.Println(accountCoreData)
	accountCoreDataJson, err := json.Marshal(accountCoreData)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON:")
	fmt.Println(string(accountCoreDataJson))
	var accountCoreDataParsed AccountCoreData
	json.Unmarshal(accountCoreDataJson, &accountCoreDataParsed)
	fmt.Println("parsed:")
	fmt.Println(accountCoreDataParsed)

	accountCoreDataXML, err := xml.Marshal(accountCoreData)
	if err != nil {
		panic(err)
	}
	fmt.Println("XML:")
	fmt.Println(string(accountCoreDataXML))

	fmt.Println("account info:")
	fmt.Println("original:")
	fmt.Println(accountInfo)
	accountInfoJson, err := json.Marshal(accountInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON:")
	fmt.Println(string(accountInfoJson))
	var accountInfoParsed AccountInfo
	json.Unmarshal(accountInfoJson, &accountInfoParsed)
	fmt.Println("parsed:")
	fmt.Println(accountInfoParsed)

	accountInfoXML, err := xml.Marshal(accountInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("XML:")
	fmt.Println(string(accountInfoXML))
}
```

---

<!-- .slide: data-transition="zoom" -->

## Q&A

---

<style>

.reveal code {
    font-size: 12px !important;
    line-height: 1.2;
}

body {
    background-color: Indigo;
}

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

.slide-number{
	margin-bottom:10px !important;
	width:100%;
	text-align:center;
	font-size:25px !important;
	background-color:transparent !important;
}
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

