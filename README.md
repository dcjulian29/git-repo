# git-repo

git-repo is a CLI binary that combines multiple repository git workflows, built with
[Cobra](https://github.com/spf13/cobra) and
[go-toolbox](https://github.com/dcjulian29/go-toolbox).

## Features

| Command | Description |
|---|---|
| `git-repo config show` | Print the full configuration |
| `git-repo config directory <path>` | Set or change the managed directory |
| `git-repo config add <name> <url>` | Add a repository to the config |
| `git-repo config remove <name>` | Remove a repository from the config |
| `git-repo config list` | List all configured repositories |
| `git-repo init` | Clone missing repos defined in config; skip those that exist |
| `git-repo status` | Colour-coded table: dirty / push / pull / diverged / untracked |
| `git-repo sync` | Fetch → pull (rebase + prune + submodules) → push |

## Configuration

Default path: `~/.config/git-repo.yml`

```yaml
directory: ~/src

repositories:
  - name: my-project
    url: https://github.com/example/my-project.git
  - name: another-repo
    url: git@github.com:example/another-repo.git
```

Copy `git-repo.yml.example` as a starting point:

```bash
cp git-repo.yml.example ~/.config/git-repo.yml
```

## Quick start

```bash
# 1. Download dependencies
go mod tidy

# 2. Build
go build          # → bin/git-repo

# 3. Bootstrap from config
git-repo config directory ~/src
git-repo config add my-lib https://github.com/example/my-lib.git
git-repo init

# 4. Daily use
git-repo status
git-repo sync
```

## Status flags

```
--actions     Show only repos that require any action
--dirty       Show only repos with uncommitted changes
--push        Show only repos that need a push
--pull        Show only repos that need a pull
--diverged    Show only repos that have diverged from remote
--untracked   Show only repos with untracked files
```
