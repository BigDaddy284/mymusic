# Contributing
<!-- docs: group=project order=10 -->

Thank you for helping. There are three ways in: themes and plugins for the
community registry, fixes and features of the app, and reports of what does
not work.

## Themes and plugins

[community/README.md](community/README.md) tells how to share a theme or a
plugin: a pull request adding one folder.

## The app

You need [Go](https://go.dev/dl/) 1.27.3 or later and [Bun](https://bun.sh) 1.2 or later.

```sh
git clone https://github.com/bigdaddy284/mymusic
cd mymusic
bun install
bun run dev
```

`bun run dev` builds the app, starts it, and rebuilds it as you edit it.
Development builds keep their data apart from an installed MyMusic, under
the name `MyMusic Dev`.

Before you open a pull request:

```sh
gofmt -l .                      # prints nothing
GOOS=windows go vet ./...
bun run build:plugins           # the community plugins' JavaScript, which tests run
go test ./...
bun run typecheck
go run ./tools/registry -check
```

[docs/architecture](docs/architecture/README.md) tells where things live. A few
rules keep the project healthy:

- **Windows only.** MyMusic is built for Windows alone. Linux still compiles the code, with stubs that do nothing, because `mygo build` runs it there to write its bindings and releases are cross-compiled on Linux: `GOOS=windows go vet ./...` checks a change from Linux.
- **The theme contract is public.** Community themes set the tokens of
  [docs/themes.md](docs/themes.md): colors, fonts and a radius. Rename or
  remove one only with a plan for the themes that set it, and document new
  ones there.
- **The plugin API is public.** Change `internal/plugins/api.js`,
  `sdk/index.d.ts` and [docs/plugins.md](docs/plugins.md) together; the
  typecheck checks the community's plugins against the types.
- **Tests do not need the network.** Parsers are tested on recorded
  responses in `testdata/`; live tests run only with `MYMUSIC_LIVE=1`.
  The Opus decoder's conformance test runs when `OPUS_TESTVECTORS` names
  the folder of the
  [RFC 8251 test vectors](https://opus-codec.org/static/testvectors/opus_testvectors-rfc8251.tar.gz),
  which CI downloads. The tests that open the audio device, silently, run
  only with `MYMUSIC_AUDIO=1`, and those that compare the native engine's
  sound with Chromium's Web Audio when `MYMUSIC_BROWSER` names a Chromium
  browser, as `msedge.exe`.
- **Honest errors.** When a service refuses something, say what and how to
  get past it, in the words of the interface.

## Keeping changes apart

Many pull requests are open at once, and one that merges shouldn't make the
others conflict. So a feature lives in files of its own: adding one touches
those files and a few new ones, not a shared list or a shared struct.

Each pull request branches off main, and only main merges into it. One
built on another's branch carries the other's changes a second time once
that one is squashed into main, and conflicts: a pull request that needs
another waits for it to merge. CI's Overlap check merges each pull request
with every other open one ([tools/overlap](tools/overlap/main.go)). When it
fails, the pull request pushed last changes so that the two don't touch the
same lines, as below, or waits for the other to merge.

- **Settings.** A setting lives in its section's struct in
  `internal/settings`, in that section's file, such as `Playback` in
  `internal/settings/playback.go`, with its default in `defaultPlayback()`
  and its clean-up in the section's `clean` beside it. Each field has a doc
  comment of its own, so a new field never makes gofmt re-align the ones
  beside it. `settings.json` stays flat: the sections are embedded. A
  literal names the section, as in
  `settings.Settings{Playback: settings.Playback{Volume: 1}}`: editors'
  gopls doesn't read a section's field set on `Settings` itself, which
  Go 1.27 compiles, and a test keeps it so. Add the key to
  `testdata/full.json` with a value that isn't its default, and rewrite
  `testdata/default.json` with
  `go test ./internal/settings -run TestDefaultFile -update`.
- **The settings page.** A setting's row lives in its feature's file, as
  `settings_surround.go` does, and registers itself from `init()` with
  `addSettings(place, order, build)`: a place is a tab, as `"playback"`,
  or a group, as `soundRows`. Orders go up in tens, so a new row slots in
  between. Two rows at the same order stop MyMusic at start-up, so the
  tests catch it; the second to merge takes another number. Nothing edits
  a central list or switch.
- **Options, not arguments.** When a function or an interface method would
  take one more argument, give it an options struct whose zero value keeps
  what it does today, and add a field, as `audio.Cue` does for how a track
  starts. Don't add bool or positional arguments to interfaces that tests
  fake.
- **Test fakes.** Fakes embed the do-nothing doubles of
  `internal/player/playertest` (`Services`, `Stream` and `Device`), and
  override only what the test needs. A new method on one of the player's
  interfaces goes there too, doing nothing.
- **Tests.** A feature's tests go in a test file of their own, named after
  the feature, such as `speed_test.go`, not at the end of a shared one.
- **Docs.** Each part's docs are a file of their own in
  [docs/architecture](docs/architecture/README.md), which nothing lists by
  hand. A feature that is a part of its own, such as podcasts, gets a new
  file there; one that adds to a part goes in that part's file, beside
  what it's nearest to.
- **Beside, not at the end.** Where a change has to add to something
  shared, such as a field of `app.Options`, a check of a condition or an
  item of a list, it goes beside the ones it's nearest to. The end, and
  the top, are where every other change adds too.
- **The changelog.** A change people will notice adds a file of its own to
  changelog.d, as [Commits](#commits) tells.
- **Small files.** Keep files small, and name them after what they hold.
  Split a file when two features that have nothing to do with each other
  keep editing it.

## Commits

Commits follow [Conventional Commits](https://www.conventionalcommits.org),
one change per commit:

```
feat(player): crossfade between tracks
fix(youtube): read albums without a year
docs: explain the plugin storage
```

Types: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `build`, `ci`,
`chore`. Scopes name the part changed: `app`, `player`, `ui`, `theme`,
`themes`, `plugins`, `youtube`, `spotify`, `stream`, `match`, `lyrics`,
`library`, `addons`, `community`, `sdk`, `tools`.

A change people will notice also gets an entry in the changelog: a file in
[changelog.d](changelog.d/README.md), in the folder of its section as
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) has them (`added`,
`changed`, `removed`, `fixed` and so on), named after its issue, such as
`changelog.d/fixed/101-stream-retry.md`. It holds the entry as
[CHANGELOG.md](CHANGELOG.md) has them, a Markdown list item. Pull requests
don't edit CHANGELOG.md itself, so they never conflict over it;
`go run ./tools/changelog -print` shows the entries of the next release.

## Releases

1. Run `go run ./tools/changelog`, which folds the entries of changelog.d
   into `## [Unreleased]` in CHANGELOG.md and deletes them. Rename
   `## [Unreleased]` to `## [x.y.z] - YYYY-MM-DD`, add an empty
   `## [Unreleased]` above it, and update the links at the foot of the file.
2. Set `version` in package.json to `x.y.z`, and commit both as
   `chore: release vx.y.z`.
3. Tag the commit and push it with the tag:

   ```sh
   git tag -a vx.y.z -m "MyMusic vx.y.z"
   git push origin main vx.y.z
   ```

The tag runs [the release workflow](.github/workflows/release.yml). It
checks that the tag, package.json and CHANGELOG.md agree, builds the Windows
installers and portable zips for x64 and Arm64, and publishes the release with the version's section of
CHANGELOG.md, the downloads of
[.github/release-notes.md](.github/release-notes.md) and `SHA256SUMS`. To
build a tag again, run the workflow from the Actions tab with the tag.

The [website](site/templates/index.html)'s download buttons open the latest
release, and the page reads the latest version from GitHub's API. Once the
release is published, set the two `data-version` spans to it, for when the
API does not answer.

Installed apps update themselves from the latest release. Its update files are
signed with the key whose public half is `updates.publicKey` in
[mygo.config.ts](mygo.config.ts); the secret half is the repository's
`MYGO_UPDATER_PRIVATE_KEY` secret, which the workflow signs with. Without it, the
release has no update files and apps stay on the version they have.

### Signing

The workflow signs the Windows apps and installers once their secrets are
in the repository; without them it says so in a notice and skips that
part, and the release goes out as before. `mygo build` signs, before it
makes the update archives, so that updates install the signed app too:
see the comments of [mygo.config.ts](mygo.config.ts) and of the workflow.

| Secret | What it holds |
|---|---|
| `WINDOWS_CERTIFICATE` | a code signing certificate and its private key, as a `.pfx`, in base64 |
| `WINDOWS_CERTIFICATE_PASSWORD` | its password, if it has one |

Get a code signing certificate that comes as a `.pfx`, then:

```sh
base64 -w0 code-signing.pfx | gh secret set WINDOWS_CERTIFICATE
gh secret set WINDOWS_CERTIFICATE_PASSWORD
```

Most certificates sold today keep their key on a hardware token or in a
signing service, and never come as a `.pfx`: those sign with MyGo's
`windows.signCommand` instead, which takes a step of its own. SmartScreen
may still warn about a new certificate until it has a reputation.

The first signed release: check on Windows that it opens without a
warning. Then take "not signed yet" out of the README's first launch,
docs/troubleshooting.md and the website, its FAQ in site/data/faq.toml
too. The release notes leave it out by themselves.

## The website

The [website](site/templates/index.html) and the docs' pages are built with
[Zola](https://www.getzola.org) from `site/` and `docs/` (see
[site/zola.toml](site/zola.toml)), and
[pages.yml](.github/workflows/pages.yml) publishes them to GitHub Pages, at
https://bigdaddy284.github.io/mymusic/, whenever what they're built from
changes on main.

## Translations of the README

[README.zh-CN.md](README.zh-CN.md) (Simplified Chinese) and
[README.id.md](README.id.md) (Indonesian) translate
[README.md](README.md), which stays the one that's current. A pull request
that changes the README needn't change them: say in it which sections
changed, so that a translator can follow. Corrections and new languages are
welcome, as `README.<language code>.md` with a link in each README's
language line.

## Reports

Open an [issue](https://github.com/bigdaddy284/mymusic/issues) with what you did,
what happened, and what you expected, with your system and MyMusic's
version (Settings › About). For security problems, see
[SECURITY.md](SECURITY.md) instead.

## Conduct

Everyone taking part follows the [code of conduct](CODE_OF_CONDUCT.md).
