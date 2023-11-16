# VerBump

## Introduction

A simple cli tool that will analyse [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) since the last tag and bump a [Semantic Version](https://semver.org/) appropriately.

It is built to work with repositories that contain multiple applications that might share `/pkg`'s

## Usage

bump the version of application1.
```
./verbump bump --repository "/path/to/repo" --current-version "0.4.8" --include "pkg,cmd/application1"
```

bump the version of application2 independently of application1.
```
./verbump bump --repository "/path/to/repo" --current-version "0.4.8" --include "pkg,cmd/application1"
```

bump the version of anything that uses `pkg` if there were any changes.
```
./verbump bump --repository "/path/to/repo" --current-version "0.4.8" --include "pkg"
```