# Changelog

- [1.0.0 - 2024-04-19](#100---2024-04-19)
- [1.1.0 - 2024-04-23](#110---2024-04-23)
- [1.1.1 - 2024-04-24](#111---2024-04-24)
- [1.2.0 - 2024-04-30](#120---2024-04-30)
- [1.3.0 - 2024-05-08](#130---2024-05-08)
- [1.3.1 - TBD](#131---TBD)

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

## [1.3.1](https://github.com/jtompkin/goclacker/releases/tag/v1.3.1) - TBD

Topper

### Added

- Debug operators

### Changed

- `&t` now `&Nt` format specifier: display top N values of the stack. Probably buggy
