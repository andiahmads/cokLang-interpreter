# The CokLang Programming Language & Interpreter
Each interpreter is built to interpret a specific programming language. That's how you "implement" a programming language. 
Without a compiler or interpreter, a programming language is nothing more than an idea or specification.  
We will parse and evaluate our own language called cokLang.

# COK lang specification
```
- C-like syntax
- variable bindings
- integers and booleans
- arithmetic expressions
- built-in functions
- first-class and higher-order functions
- closures
- a string data structure
- an array data structure
- a hash data structure
```

# COK language sightings
``` env
let age = 1;
let name = "Monkey";

let result = 10 * (20 / 2);

```

Besides integers, booleans and strings, the cok interpreter we’re going to build will also
support arrays and hashes. Here’s what binding an array of integers to a name looks like:
``` env
let myArray = [1, 2, 3, 4, 5];
```

And here is a hash, where values are associated with keys:
``` env
let thorsten = {"name": "Thorsten", "age": 28};
```

Accessing the elements in arrays and hashes is done with index expressions:
``` env
myArray[0] // => 1
thorsten["name"] // => "Thorsten"
```

The let statements can also be used to bind functions to names. Here’s a small function that
adds two numbers:
``` env
let add = fn(a, b) { return a + b; };
```

But cok not only supports return statements. Implicit return values are also possible, which means we can leave out the return if we want to:
``` env
let add = fn(a, b) { a + b; };
```

And calling a function is as easy as you’d expect:
``` env
add(1, 2);
```


A more complex function, such as a fibonacci function that returns the Nth Fibonacci number,
might look like this:
``` env
let fibonacci = fn(x) {
if (x == 0) {
   0
} else {
if (x == 1) {
    1
} else {
    fibonacci(x - 1) + fibonacci(x - 2);
 }
}
```

# quick start
``` console
go test ./lexer
```

