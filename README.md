# artifact-selector

A CLI tool that fetches GitHub release artifacts and ranks them by best match for your system. The first result is always the best possible match.

## How it works

`artifact-selector` queries a GitHub project's releases, filters the available artifacts based on your system specifications (architecture, OS, extension preference, etc.), and outputs a JSON-ranked list. Artifacts are scored across multiple dimensions — extension preference, architecture, OS, OS version, content type, and musl preference — then sorted so the top result is the best match.

## Installation

Requires Go 1.23+.

```sh
git clone https://github.com/darthbanana13/artifact-selector.git
cd artifact-selector
make build
```

The binary will be at `bin/artifact-selector`.

## Usage

```sh
artifact-selector -g user/project -r latest -e deb,appimage,tar.gz -a x86_64 -o ubuntu -O 24.04
```

Output is JSON:

```json
{
  "version": "v0.10.4",
  "artifacts": [
    { "name": "project-0.10.4-linux-amd64.deb", "..." : "..." },
    ...
  ]
}
```

Pipe through `jq` for pretty-printing:

```sh
artifact-selector -g neovim/neovim | jq '.'
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--github` | `-g` | `neovim/neovim` | The `user/project_name` to look up on GitHub |
| `--release` | `-r` | `latest` | Release version to fetch (e.g. `latest`, `v1.2`) |
| `--extension` | `-e` | `deb,,appimage,tar.zst,tbz,tar.gz,tar.xz` | Comma-separated extension preference list. Earlier entries are preferred. Use empty string for no-extension (bare binary) |
| `--arch` | `-a` | `x86_64` | Target architecture (e.g. `amd64`, `arm64`, `x86`) |
| `--os` | `-o` | `ubuntu` | Target OS or distro (e.g. `ubuntu`, `linux`, `macos`) |
| `--os-version` | `-O` | `24.04` | Target OS/distro version |
| `--musl` | `-m` | `false` | Exclude musl artifacts |
| `--musl-prefer` | `-M` | `false` | Prefer musl over non-musl artifacts |
| `--verbose` | `-v` | off | Increase verbosity. Use `-v` for verbose, `-vv` for very verbose |
| `--regex` | `-X` | — | Additional regex filter(s) to apply to artifact names |
| `--regex-meta` | `-K` | — | Metadata key name for the regex match |
| `--regex-lower` | `-L` | `no` | Lowercase artifact name before applying regex. Values: `yes`, `no`, `y`, `n` |
| `--regex-filter` | `-F` | `no` | Regex filter mode. `yes`/`y` = keep only matches, `no`/`n` = add metadata only, `exclude`/`e` = exclude matches |
| `--token` | `-t` | — | GitHub token (classic or fine-grained with `repo` scope). Also reads from `GITHUB_TOKEN` env var or `~/.config/artifact-selector/config.json` |

## GitHub Token

The token is resolved in this order:

1. `--token` / `-t` flag
2. `GITHUB_TOKEN` environment variable
3. `$XDG_CONFIG_HOME/artifact-selector/config.json` (`github_token` key)
4. `$HOME/.config/artifact-selector/config.json` (`github_token` key)

## Examples

Find the best Neovim artifact for Ubuntu 24.04 on x86_64, preferring `.deb`:

```sh
artifact-selector -g neovim/neovim -e deb,appimage,tar.gz -a x86_64 -o ubuntu -O 24.04
```

Find a musl-preferred static binary for Alpine:

```sh
artifact-selector -g some/project -M -e tar.gz, -a x86_64 -o linux
```

Use regex to filter artifacts containing "static":

```sh
artifact-selector -g some/project -X "static" -F yes
```

## License

GPLv3 — see [LICENSE](LICENSE).
