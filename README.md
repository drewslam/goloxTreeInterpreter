# Golox

A tree-walk interpreter for the Lox programming language, implemented in Go. Based on the Java implementation from [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom.

## About

Lox is a dynamically-typed, object-oriented scripting language designed for learning interpreter implementation. This Go implementation follows the book's architecture while adapting to Go's idioms and conventions.

## Features

- Full lexical scanning and parsing
- Object-oriented programming with classes and inheritance
- First-class functions and closures
- Lexical scoping with block-level resolution
- Tree-walk interpreter with environment-based variable binding

## Requirements

- Go 1.22.5 or later

## Installation

Clone the repository and build:

```bash
git clone https://github.com/drewslam/goloxTreeInterpreter.git
cd goloxTreeInterpreter
go build Lox.go
```

## Usage

Run a Lox script:
```bash
./Lox script.lox
```

Start the REPL:
```bash
./Lox
```

## Project Structure

```
.
├── ast/             # Abstract syntax tree definitions
├── scanner/         # Lexical analysis
├── parser/          # Syntax analysis
├── resolver/        # Semantic analysis and variable resolution
├── interpreter/     # Tree-walk interpreter
├── environment/     # Variable binding and scopes
├── object/          # Runtime object representations
├── loxCallable/     # Function and class interfaces
├── loxError/        # Error handling
└── token/           # Token definitions
```

## Language Examples

Classes and inheritance:
```lox
class Animal {
  init(name) {
    this.name = name;
  }

  speak() {
    print this.name + " makes a sound";
  }
}

class Dog < Animal {
  speak() {
    print this.name + " barks";
  }
}

var dog = Dog("Buddy");
dog.speak();  // Buddy barks
```

Functions and closures:
```lox
fun makeCounter() {
  var count = 0;
  fun increment() {
    count = count + 1;
    return count;
  }
  return increment;
}

var counter = makeCounter();
print counter();  // 1
print counter();  // 2
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

The original Lox language specification and Java implementation from Crafting Interpreters are licensed under the MIT License. Copyright (c) 2015 Robert Nystrom.

## Acknowledgments

- Robert Nystrom for the excellent [Crafting Interpreters](https://craftinginterpreters.com/) book
- The original Lox language design and implementation
