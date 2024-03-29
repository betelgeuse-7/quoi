### Quoi Programming Language

[Pronunciation](https://forvo.com/word/quoi/) (like 'kwa')

Quoi is a simple programming language. This repository is an implementation of this language that compiles Quoi to Go.

Quoi is an explicitly, and statically typed programming language.

##### Some code samples

```cpp
Stdout::println("Hello world")
```

```rust
fun factorial(int n) -> int {
    int product = 1.
    int j = 1.
    loop (lte j n) {
      j = (+ j 1).
      product = (* product j).
    }
    return product.
}
```

Quoi does not have a lot of features, or syntactic sugar.

There are 3 primitive data types: 
```lisp
int         ; 64-bit
bool
string      ; utf-8 encoded strings
```

##### Operators
- Operators in Quoi are prefix operators (like in Lisp).
- They are enclosed in parenthesis ("like in Lisp" 2).
```python
+ - * / = ' lt lte gt gte 
and or not
```
```lisp
(* (+ 1 2) (/ 6 2))         ; result is 9
(and true true)             ; true
(not (and true false))      ; true 
(lt 5 4)                   ; false
(not (gte 5 5))            ; false
```

- There are lists.

  - List literals start with an opening square bracket, and end with a closing one.
  - List types are in the form of ```listof <type>```.
  - There is a list indexing operator. (```(' list index)```)
    - This operator returns the value stored at that index. To place a new value at that index use ```List::replace_<typeof_list>(list, index, new_value)```

```lisp
listof string names = ["Jennifer", "Hasan"].
listof int nx = [1, 2, 56, 9910].

Stdout::println(String::from_int((' nx 2))) ; prints 56
```

<a id="datatypes"></a>
There are user-defined data types (```datatype```).

```lisp
; declaration
datatype User {
    string name
    int age
}

; initialization
User u = User {name="Jennifer" age=34}.
User u2 = User{}. ; name, and age are set to their zero values.
User u3 = User{name="Jennifer"}. ; age is set to its zero value, which is 0.

; getting field values
string jennifer = (get u name).

; setting field values
; set operator returns back 'u' with name as "Hasan".

print_User(u). ; Jennifer, 34
u = (set u name "Hasan").
; don't forget to reassign u; otherwise, u is not changed (well, before that, it is a compilation error; because, set expression is not used.).
print_User(u). ; Hasan, 34
```

##### Zero values

- "" for strings
- 0 for ints
- false for bools

Functions: 

```rust
fun hello_world() { }
fun hello_world() -> { } 

fun greet(string name) { }

fun some_func(int a, int b) -> string, bool {
  return "Hey", true.
}
```

Loops:

```rust
loop <condition> {

}
```
As long as condition is true, call statements inside the block ({}).

```c
for (int i = 0; i < 10; i++) {
    printf("#%d\n", i);
}
```
Equivalent of this classical loop above: 
```rust
int i = 0.
loop (lt i 10) {
  string msg = String::concat("#", String::from_int(i)).
  Stdout::println(msg).
  Stdout::print("\n").
  i = (+ i 1).
}
```

Branching:
```rust
if <condition> {

} elseif <condition> {

} else {

}
```

See [datatype](#datatypes)   

##### Keywords

List of all keywords: 

``` 
datatype, fun, int, string, bool, listof, block, end, if, elseif, else, loop, return, break, continue
```

--- 
##### Some notes about the syntax

- Statements end with dots.
- Spacing is not strict. As long as you separate keywords with at least one whitespace character, the rest doesn't matter.
- Escape sequences (TODO)
- Newlines are required after every field in ```datatype``` declarations.

##### Some notes about the semantics

- PARADIGM

- It is explicitly, and statically typed.
- It does not allow function overloading.
- Functions can only be declared globally. No function declarations in other functions' bodies, or in any other block (ifs, loops, arbitrary blocks, etc.). 
- No way to make a variable constant, but there is a convention that ```ALL_UPPERCASE``` variables are meant to be constants (like in Python).
- You can create new blocks that have their own scopes, using ```block```, and ```end``` keywords.
- Variables in a scope, cannot be accessed outside of said scope. It will raise some kind of a ```ReferenceError``` (like in Javascript).
  - ```lisp
    int day = 15.
    block 
        int day = 30.
        int age = 15.
        Stdout::println(day)      ; prints 30
                                  ; if a global variable and a variable in a scope has the same name,
                                  ; and a (pseudo-)function references that name, then the function will
                                  ; use the one which is in the same block as it is. 
    end
    Stdout::println(day).         ; prints 15
    Stdout::println(age).         ; reference error
    ```
- Ability to compose different types to create a compound data type, using the ```datatype``` keyword.
- There are no methods attached to a data type, but you can just create functions that take in any data type.
  - ```sml
    datatype City {
        string name
        int founded_in
    }

    fun introduce_city(City c) -> string {
        string res = "City".
        return res.
    }
    ```
- No module system; but there are namespaces that you can use. They form the standard library.
  - When you use a namespace, the code necessary to provide that service is injected in the compiled Go code.
- No floats.
- Only one looping construct (```loop``` keyword).
- No manual memory management. Quoi programs are compiled to Go, and the Go runtime handles all the memory management using a garbage collector.
- All the code is written in one file (this may change). There is no entry point to the program (a main function), so the instructions just run sequentially, top to bottom. Compiled Go code is in one file.
- Global variables can be accessed anywhere in the program.
- No pointers, but values are pass-by-reference; meaning when you pass an argument to a function, you basically pass a pointer to that argument, so the callee can change the argument's value.
  - ```lisp
    int age = 30.

    fun celebrate_birthday(int age) {
        Stdout::println("Happy birthday").
        age = (+ age 1).
    }
    ;; we can put '->' after parameter list, even if there are no return types.
    ;; so, the below function declaration is also valid:
    ;;
    ;; func celebrate_birthday(int age) -> {
    ;;  Stdout::println("Happy birthday").
    ;;  age = (+ age 1).
    ;; }
    celebrate_birthday(age).
    Stdout::println(age).       ; 31
    ```
- No function signatures.
- Functions are not values. They cannot be assigned to variables.
- We can reference functions before their declarations.

##### Namespaces

- Namespaces form the standard library.
- They provide functions to manipulate built-in data types, or print to the console, etc.
c
Syntax:

```cpp
<namespace>::<function>().
```

```lisp
; get the index of the first occurence of character 'e' in string "Hello"
int idx = String::index("Hello", "e").
Stdout::print("Index of 'e': ").
Stdout::println(idx).
```