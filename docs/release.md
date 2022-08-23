# How to create a release

- Ensure [CHANGELOG.md](/CHANGELOG.md) is up to date.
- Add heading with version and release date. e.g. `## 0.1.0 2022-06-21`
- Push a git tag with the version you wish to release.

```
git tag -a v0.1.0 -m "v0.1.0"
```

- Once the release is created edit the description to add the changelog text.
- Note: The release description will already have the commits added by GoReleaser.

## GoReleaser

This project uses [GoReleaser](https://github.com/goreleaser/goreleaser) to publish 
GitHub releases with binaries for Linux, Mac and Windows.

## Homebrew

GoReleaser also publishes `grid-intensity` to a [Homebrew Tap](https://docs.brew.sh/Taps) at https://github.com/thegreenwebfoundation/homebrew-carbon-aware-tools

See [homebrew.md](homebrew.md) for more details.
