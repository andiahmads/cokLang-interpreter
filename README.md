# interpreter-Cok-Lang

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
# monkey language sightings
``` console
let age = 1;
let name = "Monkey";

let result = 10 * (20 / 2);

And here is a hash, where values are associated with keys:
let thorsten = {"name": "Thorsten", "age": 28};
Accessing the elements in arrays and hashes is done with index expressions:
myArray[0] // => 1
thorsten["name"] // => "Thorsten"
The let statements can also be used to bind functions to names. Here’s a small function that
adds two numbers:
let add = fn(a, b) { return a + b; };
But Monkey not only supports return statements. Implicit return values are also possible,
which means we can leave out the return if we want to:
let add = fn(a, b) { a + b; };
And calling a function is as easy as you’d expect:
add(1, 2);
A more complex f
```

# support arrays and hashes. Here’s what binding an array of integers to a name looks like:
``` console
let myArray = [1, 2, 3, 4, 5];
```


# quick start
``` console
go test ./lexer
```

