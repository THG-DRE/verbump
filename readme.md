# VerBump

## Introduction

A simple CLI tool that will analyse [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) since the last tag and bump a [Semantic Version](https://semver.org/) appropriately.

It is built to work with repositories that contain multiple applications that might share `/pkg`'s

## Usage

### version multiple applications in the same repository

bump the version of `app1` taking into account any commits to shared components in `pkg`.
```
./verbump bump  --repository "." --include "pkg,cmd/app1" --version-file "cmd/app1/version"
```

commit `fix: ...` to `cmd/app1` results in version `0.0.1` being bumped to `0.0.2`

bump the version of `app2` taking into account any commits to shared components in `pkg`.
```
./verbump bump --repository "." --include "pkg,cmd/app2" --version-file "cmd/app2/version"
```

commit `feat: ...` to `pkg` results in version `0.0.2` being bumped to `0.1.0`

---

### version an application with a pre release flag

bump the version of `app1` and add a pre release label of `alpha`

```
./verbump bump --repository "." --include "pkg,cmd/app2" --version-file "cmd/app2/version" --pre-release "alpha"
```

commit `fix!: ...` to `pkg` results in version `0.1.0` being bumped to `1.0.0-alpha.0`

running the same command again will bump the version from `1.0.0-alpha.0` to `1.0.0-alpha.1`

running the same command with a new pre release label of `rc` will bump the version from `1.0.0-alpha.1` to `1.0.0-rc.0`

running the same command without pre release label will bump the version from `1.0.0-rc.0` to `1.0.0`
