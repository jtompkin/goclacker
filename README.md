# goclacker

Command line reverse Polish notation (RPN) calculator. This stack is ready to
Go.

Josh Tompkin

jtompkin-dev@pm.me

https://github.com/jtompkin/goclacker

## Installation

Install with Go.

```
go install github.com/jtompkin/goclacker@latest
```

Build with Go.

```
git clone https://github.com/jtompkin/goclacker.git
cd goclacker
go build .
```

Live, laugh, love with Go.

If you are not familiar with [Go](https://go.dev), the binary will be in
`~/go/bin` for linux, `C:\users\<USER>\go\bin` for windows (probably).

Binaries are available on the
[release](https://github.com/jtompkin/goclacker/releases/latest) page.

## Usage

```
goclacker [-V] [-h] [-s] [-l] int [-w] string [-p] string [program...]
```

If any positional arguments (`program...`) are supplied, they will be
interpreted and executed by the calculator. To enter interactive mode, do not
provide any arguments to `program...`. Run `goclacker -h` to see information
about the other command line arguments. If the program does not start, you must
first denounce [infix notation](https://en.wikipedia.org/wiki/Satan) and your
god and it will then work as intended.

## Interactive mode

Type a number and press enter to push it to the stack. Type an operator and
press enter to execute that operator on the stack. Enter `help` to see available
operators. Enter mutliple commands separated by a space and press enter to
execute them in order.

## Prompt

If you're into accessorizing your command line RPN calculators (I know you are),
you can create your own custom prompt with the `-p` flag. Just provide a single
string that defines what you want the prompt to look like. If you wanna go
really crazy, you can include format specifyers that will print information
about your current calculating environment! All format specifyers are prefixed
by a `&` character. Some examples:

`goclacker -p ' &c > '` would make a prompt that prints the current stack size
and a greater than character. All spaces are preserved; no extra whitespace is
ever added. This happens to be the default prompt.

`goclacker -p '-&t-&l- <3 '` would make a prompt that prints the top value in
the stack and the stack size limit surrounded by `-` characters and a heart. For
when you're in the *mood* for that reverse Polish goodness.

| Specifyer | Value               |
|----------:|---------------------|
|         l | stack size limit    |
|         c | current stack size  |
|         t | top stack value     |
|         s | current stash value |

You can probably break this if you try hard enough, so please do.

## Words

Custom commands (called words) can be defined in a config file (see [config
file](#configuration) if you wanna know how).

You can also define words in interactive mode!! To do so, start your command
with `=` and then type the word and the program you want to run when you enter
the word (the `>` is not typed, ya dingus). Now, when `sqrt` is entered at the
prompt, 0.5 is pushed to the stack, and the exponentiation operator is called.
That is apparently the same thing as taking the square root. Crazy. Entering
`pi` would simply push the value of pi to the stack.

```
  > = sqrt 0.5 ^
  > = pi 3.14159265358979323846
```

These two words happen to be automagically defined whenever you start the
program. If you hate them (or any other words you define) you can delete a
defined word by providing its name after `=` without any definition. All
currently defined words can be viewed by entering `words`.

```
  > = sqrt
```

## Configuration

If you have crafted a beautiful prompt or have a list of words that you can't
live without, a config file is what you need. Provide the path to this text file
with the `-c` flag, and it will set the prompt format and define any words
inside every time you start goclacker.

The format is as follows:

- First line is the prompt format
- Any other lines are word definitions

The first line is **always** interpreted as the prompt format. Leave it blank if
you want the default prompt. You can surround your format with `"` on either
side if you would like (not `'`!!); surrounding double quotes will not be
included in the prompt. If you want a blank prompt (because you are boring), use
a single `"` as your first line. Any format provided with `-p` will ovveride
whatever is in the config file.

Word definitions are the same as in interactive mode, except the `=` is not
included - i.e. The first word per line is the word itself.

A configuration file containing the following three lines would set the prompt
to look like `------> ` (notice the lack of `"` and the preserved whitespace),
and define the same words as in the [interactive](#words) example.

```
------> "
sqrt 0.5 ^
pi 3.14159265358979323846
```
