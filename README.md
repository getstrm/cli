# Pace Command Line Interface

[![GitHub Actions](https://github.com/getstrm/cli/workflows/Build/badge.svg)](https://github.com/getstrm/cli/actions)
[![Latest Release](https://img.shields.io/github/v/release/getstrm/cli)](https://github.com/getstrm/cli/releases/latest)

This package contains a command line interface (CLI) for interacting with [Pace](https://pace.getstrm.com).

## Installation

### Builds
The Pace CLI is available for major OS platforms: Linux, Mac and Windows. Please note Windows builds are not tested by us, but should work properly.

### Manually

Download the latest release for your platform from
the [releases page](https://github.com/getstrm/cli/releases). Put the binary somewhere on your `$PATH`.

#### Shell Completion

In order to set up command completion, please follow the instructions below:

- for `bash` users \
  - add the following line to your `~/.bash_profile` or `~/.bashrc`
  `source <(pace completion bash)`
  - macOS users: `pace completion bash > /usr/local/etc/bash_completion.d/pace`
- for `zsh` users \
  ensure that shell completion is enabled, then run (only needs to be done once):
  `pace completion zsh > "${fpath[1]}/_pace"`
- for fish users \
  `pace completion fish > ~/.config/fish/completions/pace.fish` (or `$XDG_CONFIG_HOME` instead of `~/.config`)

### Homebrew

The CLI is also available through Homebrew. Install the formula as follows:

```
brew install pace/cli/pace
```

Setup command completion as described above.

Upgrades to the CLI can be done through `brew upgrade pace`.

### Other package managers

More package managers will be added in the future, so stay tuned.

## Configuration

The `pace` CLI can be configured using either the flags as specified by the help (as command line arguments), with
environment variables, or with a configuration file, named `config.yaml`, located in the Pace configuration directory
`~/.config/pace`

### Configuration directory

The Pace CLI stores it's information in a configuration directory, by default located in:
`$HOME/.config/pace/`. In this directory, the CLI looks for a file named: `config.yaml`, which is used for
setting global flags.

## Getting help

If you encounter an error, or you'd like a new feature, please create an
issue [here](https://github.com/getstrm/pace/cli/issues/new). Please be thorough in your description, as it helps us
to help you more quickly. At least include the version of the CLI, your OS. terminal and any custom Pace flags
that are present in your config or environment.

## More resources

See our [documentation](https://pace.getstrm.com/docs/readme/installation).
