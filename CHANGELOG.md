# Changelog

- [1.0.0 - 2024-04-19](#100---2024-04-19)
- [1.1.0 - 2024-04-23](#110---2024-04-23)
- [1.1.1 - 2024-04-24](#111---2024-04-24)
- [1.2.0 - 2024-04-30](#120---2024-04-30)
- [1.3.0 - 2024-05-08](#130---2024-05-08)
- [1.3.1 - 2024-05-30](#131---2024-05-30)
- [1.3.2 - 2024-06-14](#132---2024-06-14)
- [1.4.0 - 2024-06-26](#140---2024-06-26)

## TODO

- [ ] better prompt format parsing

## [1.0.0](https://github.com/jtompkin/goclacker/releases/tag/v1.0.0) - 2024-04-19

It's Go time (I'm sorry).

### Added

- EVERYTHING

## [1.1.0](https://github.com/jtompkin/goclacker/releases/tag/v1.1.0) - 2024-04-23

Clearing and cleaning and configing

### Added

- `cls` operator: clear the terminal screen (actual magic don't ask me how)
  (also not at all guaranteed to work in all terminals)
- Define custom prompt and words in one config file! Provide path to config file
  with `-c` flag.

### Fixed

- Blank lines in word definition file (now config file) no longer destroy the
world.

### Changed

- Removed external dependency for ordered map. All vanilla now.

## [1.1.1](https://github.com/jtompkin/goclacker/releases/tag/v1.1.1) - 2024-04-24

Whoopsies

### Fixed

- `help` can be ran multiple times now.

## [1.2.0](https://github.com/jtompkin/goclacker/releases/tag/v1.2.0) - 2024-04-30

Stack Ops 2: Electric Boogaloo

### Added

- `swap` operator: swap top 2 stack values
- `froll` & `rroll` operators: roll stack forward or backward

### Changed

- no stack limit if `-l` is negative

## [1.3.0](https://github.com/jtompkin/goclacker/releases/tag/v1.3.0) - 2024-05-08

Terminal woes

### Added

- Arrow keys scroll through history in interactive mode.

### Changed

- Interactive mode generally better

## [1.3.1](https://github.com/jtompkin/goclacker/releases/tag/v1.3.1) - 2024-05-30

Topper

### Added

- Debug operators

### Changed

- `&t` now `&Nt` format specifier: display top N values of the stack. Probably
buggy
- Config file now interprets any additional lines as regular programs

## [1.3.2](https://github.com/jtompkin/goclacker/releases/tag/v1.3.2) - 2024-06-14

Idunno

### Changed

- Debug operators no longer show in help.
- Other stuff probably.

## [1.4.0](https://github.com/jtompkin/goclacker/releases/tag/v1.4.0) - 2024-06-26

Double value

### Added

- Value words: use `==` to define value words. They execute their definition and
  alias the top value in the resulting stack to the word, as opposed to aliasing
  the definition itself.
- `e` value word: alias for the value of e
- `asin`, `acos`, `atan` operators for arc trig functions.

### Changed

- `pi` is now a value word.
