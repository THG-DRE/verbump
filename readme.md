# VerBump

## Introduction

A simple CLI tool that will analyse [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) since the last tag and bump a [Semantic Version](https://semver.org/) appropriately.

It is built to work with repositories that contain multiple applications that might share `/pkg`'s

## Usage

bump the version of `application1`.
```
./verbump bump \
    --repository "/path/to/repo" \
    --version-file "/path/to/repo/cmd/application1/version" \
    --include "pkg,cmd/application1"
```

bump the version of application2 independently of application1.
```
./verbump bump \
    --repository "/path/to/repo" \
    --version-file "/path/to/repo/cmd/application2/version" \
    --include "pkg,cmd/application2"
```

bump the version of anything that uses `pkg` if there were any changes.
```
./verbump bump \
    --repository "/path/to/repo" \
    --version-file "/path/to/repo/cmd/application2/version" \
    --include "pkg"
```
