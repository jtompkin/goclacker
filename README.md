# goclacker
Command line reverse Polish notation (RPN) calculator. This stack is ready to Go.

Josh Tompkin

jtompkin-dev@pm.me

https://github.com/jtompkin/goclacker

## Installation

Install with Go

```
go install github.com/jtompkin/goclacker@latest
```

Build with Go

```
git clone https://github.com/jtompkin/goclacker.git
cd goclacker
go build .
```

Live, laugh, love with Go.

## Usage

```
goclacker [-V] [-h] [-s] [-l] int [-w] string [-p] string [program...]
```

If any positional arguments (`program...`) are supplied, they will be interpreted and executed by the calculator. To enter interactive mode, do not provide any arguments to `program...`. Run `goclacker -h` to see information about the other command line arguments. If the program does not start, you must first denounce [infix notation](https://en.wikipedia.org/wiki/Satan) and your god and it will then work as intended.

## Interactive mode

Type a number and press enter to push it to the stack. Type an operator and press enter to execute that operator on the stack. Enter `help` to see available operators. Enter mutliple commands separated by a space and press enter to execute them in order. 

## Words

Custom commands (called words) can be defined in interactive mode or in a file. To define words in a file, provide one word statement on each line. A word statement consists of the word itself and its definition. The word and definition are separated by a space, and all operations in the definition are separated by a space. Provide the path to this file when calling the program with `-w`. A file containing the following two lines would define two words: `sqrt`, which pushes 0.5 to the stack and then calls the exponent operator (which happens to take the square root (which I totally knew before making this thing)), and `pi`, which pushes the value of pi to the stack:

```
sqrt 0.5 ^
pi 3.14159265358979323846
```

You can also define words in interactive mode!!!! To do so, start your command with `=` and then enter the word and its definition just like you would in a file. The following two commands would accomplish the same word definitions as the file above (the `>` is not typed, ya dingus).

```
  > = sqrt 0.5 ^
  > = pi 3.14159265358979323846
```

These two words happen to be automagically defined whenever you start the program. If you hate them (or any other words you define) you can delete a defined word by providing its name after `=` without any definition. All currently defined words can be viewed by entering `words`.

```
  > = sqrt
```

## Prompt

If you're into accessorizing your command line RPN calculators (I know you are), you can create your own custom prompt with the `-p` flag. Just provide a single string that defines what you want the prompt to look like. If you wanna go really crazy, you can include format specifyers that will print information about your current calculating environment! All format specifyers are prefixed by a `&` character. Some examples:

`goclacker -p ' &c > '` would make a prompt that prints the current stack size and a greater than character. All spaces are preserved; no extra whitespace is ever added. This happens to be the default prompt.

`goclacker -p '-&t-&l- <3 '` would make a prompt that prints the top value in the stack and the stack size limit surrounded by `-` characters and a heart. For when you're in the *mood* for that reverse Polish goodness.

| Specifyer |        Value        |
|----------:|---------------------|
|          l| stack size limit    |
|          c| current stack size  |
|          t| top stack value     |
|          s| current stash value |

You can probably break this if you try hard enough, so please do.
