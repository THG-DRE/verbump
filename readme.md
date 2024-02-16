# VerBump

## Introduction

A simple CLI tool that will analyse [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) since the last tag and bump a [Semantic Version](https://semver.org/) appropriately.

It is built to work with repositories that contain multiple applications that might share `/pkg`'s

## Usage

### version multiple applications in the same repository

Bump the version of `app1` taking into account any commits to shared components in `pkg`.
```
./verbump bump --include "pkg,cmd/app1" --version-file "cmd/app1/version"
```

Commit `fix: ...` to `cmd/app1`: `0.0.1` -> `0.0.2`

Bump the version of `app2` taking into account any commits to shared components in `pkg`.
```
./verbump bump --include "pkg,cmd/app2" --version-file "cmd/app2/version"
```

Commit `feat: ...` to `pkg`: `0.0.2` -> `0.1.0`

---

### version an application with a pre release flag

Bump the version of `app1` and add a pre release label of `alpha`

```
./verbump bump -i "pkg,cmd/app2" -v "cmd/app2/version" --pre-release "alpha"
```

commit `fix!: ...` to `pkg`: `0.1.0` -> `1.0.0-alpha.0`

Running the same command again: `1.0.0-alpha.0` -> `1.0.0-alpha.1`

Change the pre release label to `rc`: `1.0.0-alpha.1` -> `1.0.0-rc.0`

```
./verbump bump -i "pkg,cmd/app2" -v "cmd/app2/version" -p "alpha"
```

running the same command: `1.0.0-rc.0` -> `1.0.0`

```
./verbump bump -i "pkg,cmd/app2" -v "cmd/app2/version"
```
