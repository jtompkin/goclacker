# Changelog

- [1.0.0 - 2024-04-19](#100---2024-04-19)
- [1.1.0 - 2024-04-23](#110---2024-04-23)

## TODO

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
  with `-c` flag. First line must be the prompt definition; you can surround it
  in `"` if you wish. Leave the first line blank if you want the default prompt.
  Every other line is a word definition. `-p` flag overrides the prompt
  definition from this file.

### Fixed

- Blank lines in word definition file (now config file) no longer destroy the
world.

### Changed

- Removed external dependency for ordered map. All vanilla now.
