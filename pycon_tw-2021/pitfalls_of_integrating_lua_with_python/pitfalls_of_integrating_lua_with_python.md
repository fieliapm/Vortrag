---
title: The Pitfalls of Integrating Lua with Python
tags: Slide, Lua, Python
description: View the slide with "Slide Mode".
slideOptions:
  spotlight:
    enabled: false
  allottedMinutes: 25
---
<small>The Pitfalls of Integrating Lua
with Python</small>
===

<!-- .slide: data-background-color="pink" -->
<!-- .slide: data-transition="zoom" -->

> [name=郭學聰 Hsueh-Tsung Kuo]
> [time=Sat, 2 Oct 2021] [color=red]

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

4. Introducing Lua
    * Features
    * Applications
5. lupa
    * What is lupa?
    * Overwhelming Performance
    * Unbelievable Integration
      * Lua table <-> Python dict
6. Pitfalls
    * Variable Exposure of Python Runtime
    * Proxy Object of Lua table and Python dict
    * Utility functions

----

<!-- .slide: data-transition="convex" -->

7. Workaround
    * Isolated Lua Runtime
    * Lua table <-> Python dict
    * Utility Functions
8. Live demo
9. Conclusion
10. Resource
11. Q&A

---

<!-- .slide: data-transition="convex" -->

## Introducing Lua

----

<!-- .slide: data-transition="convex" -->

### Introducing Lua

```lua=
-- defines a factorial function
function fact (n)
    if n == 0 then
        return 1
    else
        return n * fact(n-1)
    end
end
    
print("enter a number:")
a = io.read("*number")        -- read a number
print(fact(a))
```

----

<!-- .slide: data-transition="convex" -->

### Features

* Low usage of CPU and memory
  * Designed for embedded environment
* Table is everything in Lua data structure
  * Lua array is a kind of table

----

<!-- .slide: data-transition="convex" -->

### Applications

* Daily changed business logic and rule description
* Suitable for non-programmer

---

<!-- .slide: data-transition="convex" -->

## lupa

----

<!-- .slide: data-transition="convex" -->

### What is lupa?

* Integrates the runtimes of Lua or LuaJIT2 into CPython.
* A partial rewrite of LunaticPython in Cython
  * with some additional features such as proper coroutine support.
* Current version 1.10 published at 2 Sep 2021

----

<!-- .slide: data-transition="convex" -->

#### Features

Major features from official introduction

```
* separate Lua runtime states through a LuaRuntime class
* Python coroutine wrapper for Lua coroutines
* iteration support for Python objects in Lua and Lua objects in Python
* proper encoding and decoding of strings (configurable per runtime, UTF-8 by default)
* frees the GIL and supports threading in separate runtimes when calling into Lua
* tested with Python 2.7/3.5 and later
* written for LuaJIT2 (tested with LuaJIT 2.0.2), but also works with the normal Lua interpreter (5.1 and later)
* easy to hack on and extend as it is written in Cython, not C
```

----

<!-- .slide: data-transition="convex" -->

#### Highlights

* Overwhelming performance (compare to CPython)
* The unbelievable integration between Lua table :arrow_backward: :arrow_forward: Python dict

----

<!-- .slide: data-transition="convex" -->

### Overwhelming Performance

test code

```python=
#!/usr/bin/env python3
# -*- coding: utf-8 -*-


# from https://alchemypy.com/2020/05/13/speed-your-python-with-lua/


import sys
import time

from lupa import LuaRuntime


def sum_all_to_range_pure_python(end_range_number):
    start_time = time.time()
    the_sum = 0
    for number in range(1, end_range_number+1):
        the_sum += number
    stop_time = time.time()
    return (end_range_number, stop_time - start_time, the_sum)


def sum_all_to_range_with_lua_python(end_range_number):
    start_time = time.time()

    lua_code =  """
    function(n)
       sum = 0
       for i=1,n,1 do
           sum = sum + i
       end
       return sum
    end
    """
    lua_func = LuaRuntime(encoding=None).eval(lua_code)
    the_sum = lua_func(end_range_number)
    stop_time = time.time()
    return(end_range_number, stop_time - start_time, the_sum)


def main(argv=sys.argv[:]):
    end_range_number = 200000000

    python_result = sum_all_to_range_pure_python(end_range_number)
    lupa_result = sum_all_to_range_with_lua_python(end_range_number)

    print('python - range: %d, time: %g sec, sum: %d' % python_result)
    print('lupa - range: %d, time: %g sec, sum: %d' % lupa_result)


if __name__ == '__main__':
    sys.exit(main())
```

----

<!-- .slide: data-transition="convex" -->

### Overwhelming Performance

pypy3 > lupa > python3

```
python3 - range: 200000000, time: 7.86578 sec, sum: 20000000100000000
lupa - range: 200000000, time: 2.94546 sec, sum: 20000000100000000
pypy3 - range: 200000000, time: 0.270106 sec, sum: 20000000100000000
```

----

<!-- .slide: data-transition="convex" -->

### Unbelievable Integration

* Call Python built-in functions in Lua runtime
* Call Lua functions in Python runtime
* Access Python dict as Lua table in Lua runtime
* Access Lua table as Python dict in Python runtime

----

<!-- .slide: data-transition="convex" -->

### Unbelievable Integration

Python & Lua function

```python=
>>> import lupa
>>> from lupa import LuaRuntime
>>> lua = LuaRuntime(unpack_returned_tuples=True)

>>> lua.eval('1+1')
2

>>> lua_func = lua.eval('function(f, n) return f(n) end')

>>> def py_add1(n): return n+1
>>> lua_func(py_add1, 2)
3

>>> lua.eval('python.eval(" 2 ** 2 ")') == 4
True
>>> lua.eval('python.builtins.str(4)') == '4'
True
```

----

<!-- .slide: data-transition="convex" -->

### Unbelievable Integration

Python dict part 1

```python=
>>> lua_func = lua.eval('function(obj) return obj["get"] end')
>>> d = {'get' : 'value'}

>>> value = lua_func(d)
>>> value == d['get'] == 'value'
True
```

----

<!-- .slide: data-transition="convex" -->

### Unbelievable Integration

Python dict part 2

* when accessing lua table, **"obj[x] == obj.x"**
  * **"as_itemgetter()"**: access collection items using **"\_\_getitem\_\_()"**
  * **"as_attrgetter()"**: access object attributes

----

<!-- .slide: data-transition="convex" -->

### Unbelievable Integration

Python dict part 2

```python=
>>> value = lua_func( lupa.as_itemgetter(d) )
>>> value == d['get'] == 'value'
True

>>> dict_get = lua_func( lupa.as_attrgetter(d) )
>>> dict_get == d.get
True
>>> dict_get('get') == d.get('get') == 'value'
True

>>> lua_func = lua.eval(
...     'function(obj) return python.as_attrgetter(obj)["get"] end')
>>> dict_get = lua_func(d)
>>> dict_get('get') == d.get('get') == 'value'
True
```

----

<!-- .slide: data-transition="convex" -->

### Unbelievable Integration

Lua table part 1

```python=
>>> table = lua.eval('{10,20,30,40}')
>>> table[1]
10
>>> table[4]
40
>>> list(table)
[1, 2, 3, 4]
>>> list(table.values())
[10, 20, 30, 40]
>>> len(table)
4
```

----

<!-- .slide: data-transition="convex" -->

### Unbelievable Integration

Lua table part 2

```python=
>>> mapping = lua.eval('{ [1] = -1 }')
>>> list(mapping)
[1]

>>> mapping = lua.eval('{ [20] = -20; [3] = -3 }')
>>> mapping[20]
-20
>>> mapping[3]
-3
>>> sorted(mapping.values())
[-20, -3]
>>> sorted(mapping.items())
[(3, -3), (20, -20)]

>>> mapping[-3] = 3     # -3 used as key, not index!
>>> mapping[-3]
3
>>> sorted(mapping)
[-3, 3, 20]
>>> sorted(mapping.items())
[(-3, 3), (3, -3), (20, -20)]
```

---

<!-- .slide: data-transition="convex" -->

## Pitfalls

----

<!-- .slide: data-transition="convex" -->

### Variable Exposure of Python Runtime

* Bad guy can inject some code to Lua script to control Python runtime
  * **"python.func"** under lua runtime <!-- .element: class="fragment" data-fragment-index="1" -->

----

<!-- .slide: data-transition="convex" -->

### Variable Exposure of Python Runtime

```python=
>>> lua.execute('''
... for k, v in pairs(python) do
...     print(k)
... end
... ''')
eval
iter
iterex
enumerate
none
as_attrgetter
as_itemgetter
as_function
builtins
```

----

<!-- .slide: data-transition="convex" -->

### Variable Exposure of Python Runtime

```python=
>>> lua.execute('''
... for i in python.iter(python.builtins.dir(python.builtins)) do
...     print(i)
... end
... ''')
ArithmeticError
AssertionError
AttributeError
BaseException
.
.
.
map
max
memoryview
min
next
object
oct
open
ord
pow
print
property
quit
range
repr
reversed
round
set
setattr
slice
sorted
staticmethod
str
sum
super
tuple
type
vars
zip
```

----

<!-- .slide: data-transition="convex" -->

### Proxy Object of Lua table and Python dict

* Issues
  * Variable exposure of Python runtime again
  * Different behavior between native data and proxy object
* Solutions
  * Keep Lua runtime operate native data only

----

<!-- .slide: data-transition="convex" -->

### Proxy Object of Lua table and Python dict

Different behavior

```python=
>>> tbl = lua.eval('{a = 1, b = 2, c = 3}')
>>> len(tbl)
0
>>> len(list(tbl.items()))
3
>>> tbl
<Lua table at 0x1db1570>
>>> tbl.items()
LuaIter(<Lua table at 0x1db3570>)
```

----

<!-- .slide: data-transition="convex" -->

#### Utility functions

* Lua table has no corresponding utility methods that Python dict has, including **"len()"**, **"get()"**, **"setdefault()"**, etc

---

<!-- .slide: data-transition="convex" -->

## Workaround

----

<!-- .slide: data-transition="convex" -->

### Workaround

* How to implement and prepare isolated Lua runtime under Python runtime
  * Remove global data & functions in Lua runtime
  * Conversion of Python dict :arrow_backward: :arrow_forward: Lua table
    * Conversion of Python dict :arrow_right: Lua table: luadata
    * Conversion of Lua table :arrow_right: Python dict: DIY
  * Utility functions for Lua table: DIY

----

<!-- .slide: data-transition="convex" -->

### Isolated Lua Runtime

```python=
LUA_ISOLATED_RUNTIME_SCRIPT = '''
require = nil
package = nil

os = nil
io = nil

python = nil
'''

lua.execute(LUA_ISOLATED_RUNTIME_SCRIPT)
```

:warning: Execute it for every new **"LuaRuntime"**

----

<!-- .slide: data-transition="convex" -->

### Lua table :arrow_backward: :arrow_forward: Python dict

```python=
import luadata

def convert_table_to_dict(table):
    dict_data = {}
    for (k, v) in table.items():
        if lupa.lua_type(v) == 'table':
            vd = convert_table_to_dict(v)
        else:
            vd = v
        dict_data[k] = vd
    return dict_data

dict_data = {'a': 1, 'b': 2, 'c': {'d': 4, 'e': 5, 'f': 6}}
table = lua.eval(luadata.serialize(dict_data)) # convert dict to table
result_table = lua_func(table)
result_dict_data = convert_table_to_dict(result_table)
```

----

<!-- .slide: data-transition="convex" -->

###  Utility Functions

```python=
LUA_UTIL_FUNC_SCRIPT = '''
function get_table_size(tbl)
    local count = 0
    for _ in pairs(tbl) do
        count = count + 1
    end
    return count
end

function clear_table(tbl)
    for k, v in pairs(tbl) do
        tbl[k] = nil
    end
end

function get_table(tbl, k, default)
    local v = tbl[k]
    if v == nil then
        v = default
    end
    return v
end

function set_default_to_table(tbl, k, default)
    if tbl[k] == nil then
        tbl[k] = default
    end
    return tbl[k]
end
'''

lua.execute(LUA_UTIL_FUNC_SCRIPT)
```

---

<!-- .slide: data-transition="convex" -->

## Live demo

----

<!-- .slide: data-transition="convex" -->

#### Try It 

```python=
import lupa, luadata
from lupa import LuaRuntime
lua = LuaRuntime(unpack_returned_tuples=True)

def convert_table_to_dict(table):
    dict_data = {}
    for (k, v) in table.items():
        if lupa.lua_type(v) == 'table':
            vd = convert_table_to_dict(v)
        else:
            vd = v
        dict_data[k] = vd
    return dict_data

def convert_dict_to_table(dict_data):
    return lua.eval(luadata.serialize(dict_data))

LUA_ISOLATED_RUNTIME_SCRIPT = '''
require = nil
package = nil

os = nil
io = nil

python = nil
'''

lua.execute(LUA_ISOLATED_RUNTIME_SCRIPT)

LUA_UTIL_FUNC_SCRIPT = '''
function get_table_size(tbl)
    local count = 0
    for _ in pairs(tbl) do
        count = count + 1
    end
    return count
end

function clear_table(tbl)
    for k, v in pairs(tbl) do
        tbl[k] = nil
    end
end

function get_table(tbl, k, default)
    local v = tbl[k]
    if v == nil then
        v = default
    end
    return v
end

function set_default_to_table(tbl, k, default)
    if tbl[k] == nil then
        tbl[k] = default
    end
    return tbl[k]
end
'''

lua.execute(LUA_UTIL_FUNC_SCRIPT)

# test

dict_data = {'a': 1, 'b': 2, 'c': {'d': 4, 'e': 5, 'f': 6}}

tbl = lua.eval('{a = 1, b = 2, c = {d = 4, e = 5, f = 6}}')

lua_func = lua.eval('''
function(tbl)
    print(tbl)
    for k, v in pairs(tbl) do
        print(string.format("%s %s", k, v))
    end
    return tbl
end
''')
```

----

<!-- .slide: data-transition="convex" -->

#### Result

See demo

---

<!-- .slide: data-transition="convex" -->

## Conclusion

----

<!-- .slide: data-transition="convex" -->

### Experience

* Programming languages are just tools
  * Choose easy-to-use ones according to your application

----

<!-- .slide: data-transition="convex" -->

### Bless

:hash: {成為 Python 與 Lua 的橋梁|<big>Become bridge between Python and Lua</big>}

> [name=郭學聰 Hsueh-Tsung Kuo] [time=2021_10_02] [color=red] :notebook:

---

<!-- .slide: data-transition="convex" -->

## Resource 

----

<!-- .slide: data-transition="convex" -->

### Reference

* Lua
  * <small>https://www.lua.org/</small>
* lupa
  * <small>https://pypi.org/project/lupa/</small>
* luadata
  * <small>https://pypi.org/project/luadata/</small>

----

<!-- .slide: data-transition="convex" -->

### Utility

* Slides
  * Editor: HackMD
  * Presentation: Mozilla FireFox
* Broadcasting toolsets
  * OBS Studio
  * Panasonic DC-G100
    * LUMIX G VARIO 12-32mm / F3.5-F5.6

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

