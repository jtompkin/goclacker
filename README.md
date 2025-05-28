# goclacker

<!--toc:start-->
- [goclacker](#goclacker)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Interactive mode](#interactive-mode)
  - [Prompt](#prompt)
  - [Words](#words)
    - [Value words](#value-words)
  - [Configuration](#configuration)
  - [License](#license)
<!--toc:end-->

Command line reverse Polish notation (RPN) calculator. This stack is ready to
Go.

By Josh Tompkin

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
go get -u golang.org/x/term
go build .
```

Live, laugh, love with Go.

If you are not familiar with [Go](https://go.dev), the binary will be in
`~/go/bin` for Linux, `C:\users\%username%\go\bin` for Windows (probably).

Pre-built binaries are available on the
[release](https://github.com/jtompkin/goclacker/releases/latest) page.

A Nix flake is also available for all you functional bros out there. To align
myself with Nix's goal of not documenting anything, I will not provide
instructions on how to use it.

## Usage

```
goclacker [-V] [-h] [-s] [-d] [-r] [-l] int [-c] string [-p] string [program]...
```

If any positional arguments (`program...`) are supplied, they will be
interpreted and executed by the calculator. To enter interactive mode, do not
provide any arguments to `program...`. Run `goclacker -h` to see information
about the other command line arguments. If the program does not work, you must
first denounce [infix notation](https://en.wikipedia.org/wiki/Satan) and your
god, and it will then work as intended.

## Interactive mode

Type a number and press enter to push it to the stack. Type an operator and
press enter to execute that operator on the stack. Enter `help` to see available
operators. Enter multiple commands separated by a space and press enter to
execute them in order.

## Prompt

If you're in to accessorizing your command line RPN calculators (I know you
are), you can create your own custom prompt with the `-p` flag. Just provide a
single string that defines what you want the prompt to look like. If you wanna
go really crazy, you can include format specifiers that will print information
about your current calculating environment! All format specifiers are prefixed
by an `&` character; some use parameters that come after the `&` and before the
specifier. Some examples:

 - `goclacker -p ' &c > '` would make a prompt that prints the current stack
 size and a greater than character. All spaces are preserved; no extra
 whitespace is ever added. This happens to be the default prompt.

 - `goclacker -p '-&3t-&l- <3 '` would make a prompt that prints the top 3
 values in the stack and the stack size limit surrounded by `-` characters, and
 a heart. For when you're in the *mood* for that reverse Polish goodness.

| specifier | value               |
|----------:|---------------------|
|         l | stack size limit    |
|         c | current stack size  |
|        Nt | top N stack values  |
|         s | current stash value |

You can probably break this if you try hard enough, so please do.

## Words

Custom commands (called words) can be defined in a config file (see [config
file](#configuration) if you wanna know how).

You can also define words in interactive mode!! To do so, start your command
with `=` and then type the word and the program you want to run when you enter
the word (the `>` is not typed, ya dingus). Consider entering the following
lines at the interactive prompt:

```
  > = sqrt 0.5 ^
  > = logb log swap log / -1 ^
```

Now, when `sqrt` is entered at the prompt, 0.5 is pushed to the stack, and the
exponentiation operator is called. That is apparently the same thing as taking
the square root. Math is crazy. Similarly, when `logb` is entered at the prompt,
all of the commands in its definition will be executed, effectively popping `a`
and `b` and pushing the logarithm base `b` of `a`.

These two words happen to be automagically defined whenever you start the
program. If you hate them (or any other words you define) you can delete a
defined word by providing its name after `=` without any definition. You can
also freely redefine any currently defined word.

```
  > = sqrt
  > = logb
```

All currently defined words can be viewed by entering `words`. Words
will be separated from their definition by a `:`.

### Value words

You can also define value words by beginning your command with `==`. Value words
are essentially aliases for numerical values, so instead of executing the
commands in the definition, they just push the defined value straight into the
stack when entered after they are defined.

```
  > == pi 22 7 /
```

This would start a sub-stack, push 22 and 7 to it and call the division
operator. It would then set the value of `pi` to the value at the top of this
sub-stack---in this case the only value in the sub-stack.

Value words are separated from their value in the `words` screen by an `=`.

## Configuration

If you have crafted a beautiful prompt or have a list of words that you can't
live without, a config file is what you need. Provide the path to this text file
with the `-c` flag, and it will set the prompt format and execute any additional
programs you supply. Goclacker looks for default config files in the following
locations and loads the first one it finds:

- `./.goclacker`
- `~/.goclacker`
- `~/.config/goclacker/config`

Passing anything---including an empty string---to `-c` will disable default 
config files.

The format is as follows:

- First line is the prompt format.
- All other lines are programs to execute.

The first line is **always** interpreted as the prompt format. Leave it blank if
you want the default prompt. You can surround your format with `"` on either
side if you would like (not `'`!!); a single pair of surrounding double quotes
will not be included in the prompt. If you want a blank prompt (because you are
boring), place a single `"` in the first line. Any format provided with `-p`
will override whatever is in the config file.

Any other lines are interpreted just as if they were entered in interactive
mode. Defining words that will be present each time you start the program is the
most useful use of this, but any regular calculations can also be done here, if
you want certain values to be in your stack at start-up.

A configuration file containing the following lines would set the prompt to look
like `------> ` (notice the lack of `"` and the preserved whitespace), and
define the a word and a value word. It would then push the square root of pi,
push the value 2, and call the multiplication operator. These last three lines
could all be put on the same line, just like in interactive mode.

```
------> "
= sqrt 0.5 ^
== pi 3.14159265358979323846
pi sqrt
2
*
```

## License

Licensed under the [MIT](https://spdx.org/licenses/MIT.html) license. See 
LICENSE file.

SPDX-License-Identifier: MIT
