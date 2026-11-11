# Lyrics
<!-- docs: group=using order=30 -->

MyTube shows synced and word-by-word lyrics for YouTube Music tracks and the music on your computer, and can translate them.

Lyrics come from Unison, LRCLIB and YouTube Music, for the track that plays, whichever service it's from. Nothing needs setting up. The lyrics themselves go to Unison, for their translation, only while you turn it on.

## Where lyrics show

- **The lyrics panel** opens at the right of the window, beside the page you're on, from the Lyrics button by the player (beside Queue), or with Ctrl+Shift+L. The same keys or button close it.
- **The now playing page** shows them larger, under the player: open it from the cover in the player bar, the card at the foot of the sidebar, or Alt+4.

![The lyrics panel at the right of the window, beside Library › Liked: its Synced and Plain text choice, the Translate and Romaji switches, and the lyrics of アイドル by YOASOBI, the line sung now lighting up in the accent colour](shots/lyrics.webp)

The panel and the page show the same lyrics, with the same switches above them, and each keeps where it was scrolled. Synced lyrics follow the song; scroll them to read back, and they stay where you put them for a few seconds before they follow again. Click a line of synced lyrics to play the song from there, once it plays.

The panel says what it found:

- "Looking for the lyrics…" while it looks;
- "No lyrics found for this track." when none of the sources has any;
- "Lyrics can't be loaded right now. Check your connection, or try again in a moment.", with try again, when a source is out of reach;
- "This track is instrumental." for a track whose lyrics source says so.

Refresh this page asks again. For a podcast episode, the panel says Details, and shows the episode's description and chapters in place of lyrics: see [Podcasts](podcasts.md#episode-details-and-chapters).

## Where they come from

MyTube asks, in order:

1. **[Unison](https://unison.boidu.dev),** an open, crowdsourced database of lyrics, by the track's video ID, or by its title, artists, album and length. Unison has many lyrics timed word by word. While it's out of reach, MyTube goes on to LRCLIB after a wait of 3 seconds at most.
2. **[LRCLIB](https://lrclib.net),** the open database of synced lyrics, by the track's title, artists, album and length: when Unison has no lyrics, can't be reached, or has only plain lyrics where LRCLIB has synced ones.
3. **YouTube Music:** when neither has any lyrics, or both have only plain ones and YouTube Music's are longer.

A track on your computer, which has no video ID, is looked up by its title, artists, album and length. The line at the end of the lyrics names where they came from.

## Synced and plain lyrics

Synced lyrics have a time for each line. The line sung now lights in the accent colour, the lines sung before it dim, and the view keeps it in the middle.

- **Synced or Plain text:** above synced lyrics, this choice shows them as plain text instead, without following the song. MyTube remembers it for every track.
- **When "synced" is offered:** only for lyrics whose lines can be followed in time. Lyrics of one timed line, or whose lines all have one time, count as plain, and another source's lyrics take their place when it has more of them.
- **At another [speed](playing-music.md#speed),** synced lyrics keep the track's own time, so they still light as the words are sung.

Plain lyrics show as they are, with no choice above them.

## Word by word

Lyrics timed word by word, as Unison has many, light each word of the line sung now as it's sung, in the panel and on the now playing page:

- each word fills with the accent from its start to its end;
- words sung stay lit, and words to come keep the colour of the lines to come;
- background vocals light the same way.

While your system asks for [less motion](look-and-feel.md#reduced-motion), each word lights whole as it starts.

## Translate

Lyrics in another language can show their translation: turn on Translate above them, in the panel or on the now playing page. Each line then shows its translation under it, smaller, moving with it in synced lyrics.

- **When the switch shows:** for lyrics in another language than the one they translate to. That's your system's language by default; "Translate lyrics to", in [Settings › Playback](settings.md#playback) › Lyrics, chooses another, among 47.
- **It's remembered** for the next song in the same language: turned on for Korean lyrics, it stays on for Korean lyrics.
- **What is sent:** while Translate is on, the lyrics of the track playing and its video ID go to Unison, which has Google translate them. Nothing is sent while it's off, as it is by default. [Your data and privacy](data-and-privacy.md#what-each-feature-sends-where) lists what every feature sends.
- **Once per song:** each translation is kept on disk, by track and language, so a song is translated once.
- **While it comes,** the panel says "Translation is on its way…". When the service is busy, the lyrics show untranslated, with "Translation isn't available right now." and try again.

## Credits

Lyrics from [Unison](https://unison.boidu.dev), an open, crowdsourced database of lyrics, under the Open Database License (ODbL). The lyrics panel ends Unison's lyrics with "lyrics from Unison (https://unison.boidu.dev)", as their licence asks. Synced lyrics also come from [LRCLIB](https://lrclib.net), the open database of synced lyrics.
