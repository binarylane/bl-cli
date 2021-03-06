# Contributing

We love contributions! You are welcome to open a pull request, but it's a good idea to
open an issue and discuss your idea with us first.

Once you are ready to open a PR, please keep the following guidelines in mind:

1. Code should be `go fmt` compliant.
1. Types, structs and funcs should be documented.
1. Tests pass.

## Getting set up

`go-binarylane` uses go modules. Just fork this repo, clone your fork and off you go!

## Running tests

When working on code in this repository, tests can be run via:

```sh
go test -mod=vendor .
```

## Versioning

go-binarylane follows [semver](https://www.semver.org) versioning semantics.
New functionality should be accompanied by increment to the minor
version number. Any code merged to main is subject to release.

## Releasing

Releasing a new version of go-binarylane is currently a manual process.

Submit a separate pull request for the version change from the pull
request with your changes.

1. Update the `CHANGELOG.md` with your changes. If a version header
   for the next (unreleased) version does not exist, create one.
   Include one bullet point for each piece of new functionality in the
   release, including the pull request ID, description, and author(s).
   For example:

```
## [v1.8.0] - 2019-03-13

- #210 - @jcodybaker - Expose tags on storage volume create/list/get.
- #123 - @janedoe - Update test dependencies
```

2. Update the `libraryVersion` number in `binarylane.go`.
3. Make a pull request with these changes.  This PR should be separate from the PR containing the go-binarylane changes.
4. Once the pull request has been merged, [draft a new release](https://github.com/binarylane/go-binarylane/releases/new).
5. Update the `Tag version` and `Release title` field with the new go-binarylane version.  Be sure the version has a `v` prefixed in both places. Ex `v1.8.0`.
6. Copy the changelog bullet points to the description field.
7. Publish the release.
