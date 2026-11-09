# Spotify in MyMusic

"In MyMusic, from Spotify" plays Spotify tracks with Spotify's Web
Playback SDK, in the hidden window of `public/spotify.html`
([its web view](webviews.md)). Spotify encrypts its audio, and the SDK
plays only where the web view has a decryption module (a CDM) it takes.
On Windows, WebView2 is Chromium with Widevine, so MyMusic plays Spotify
itself, and does so by default (as tried for #295 in November 2026).

## What the SDK asks for

- **Its frame.** The SDK plays in an iframe of
  `https://sdk.scdn.co/embedded/index.html` (`allow="encrypted-media;
  autoplay"`). It looks for Widevine (`com.widevine.alpha`), then
  PlayReady, then FairPlay (`com.apple.fps.1_0`), through EME or WebKit's
  older `WebKitMediaKeys`.
- **Where none answers**, it fails with `initialization_error` ("Failed
  to initialize player") at once, before it asks for a token.
- **Widevine at software robustness** (`SW_SECURE_CRYPTO`, level L3) is
  enough.
