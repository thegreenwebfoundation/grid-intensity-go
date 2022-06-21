# How to create a release

Push a git tag with the version you wish to release.

```
git tag -a v0.1.0 -m "Initial release"
```

## GoReleaser

This project uses [GoReleaser](https://github.com/goreleaser/goreleaser) to publish 
GitHub releases with binaries for Linux, Mac and Windows.

## Homebrew

GoReleaser also publishes `grid-intensity` to a [Homebrew Tap](https://docs.brew.sh/Taps) at https://github.com/thegreenwebfoundation/homebrew-carbon-aware-tools

See [homebrew.md](homebrew.md) for more details.
