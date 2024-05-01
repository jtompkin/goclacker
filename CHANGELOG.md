# Changelog

- [1.0.0 - 2024-04-19](#100---2024-04-19)
- [1.1.0 - 2024-04-23](#110---2024-04-23)
- [1.1.1 - 2024-04-24](#111---2024-04-24)
- [1.1.2 - TBD](#112---TBD)

## TODO

- [ ] Command history

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

## [1.1.2](https://github.com/jtompkin/goclacker/releases/tag/v1.1.2) - TBD

Stack Ops: Cold War

### Added

- `swap` operator: swap top 2 stack values
- `froll` & `rroll` operators: roll stack forward or backward
