# Themes and plugins
<!-- docs: group=using order=80 -->

Themes and plugins are add-ons: a theme sets MyTube's colours, typefaces
and corners, and a plugin adds pages, commands and more. MyTube loads
those it ships with and those you put in your add-ons folders.

## Add a theme or a plugin

The folders are in MyTube's
[data folder](data-and-privacy.md#the-data-folder), `addons/themes` and
`addons/plugins`: **open themes folder** in
[Settings › Appearance](look-and-feel.md#themes) › Theme and **open plugins
folder** in [Settings › Plugins](settings.md#plugins) open them.

Put the add-on's folder there, named after its `id`, then
[reload add-ons](#reload-add-ons). A theme joins the theme cards of
Settings › Appearance, and works with every
[layout](look-and-feel.md#layouts) and option there. The repository's
[examples](../examples/README.md) are themes and plugins to start from.

## Enable a plugin

A plugin doesn't run until you enable it, with its switch in Settings ›
Plugins › Installed plugins. MyTube asks first:

> Plugins can control the player, change your library and reach the
> internet, whatever they ask for. Enable plugins from authors you trust.

followed by what the plugin asks for. Only add plugins from authors you
trust.

A plugin that fails to start says why, in a notice and on its row in
Settings › Plugins ("failed:" and the reason); **reload** tries it again.

## What a plugin can do

A plugin runs in a JavaScript engine of MyTube's own, inside MyTube,
not in a browser. Through the [plugin API](plugins.md#the-api), it can:

- add **pages**, with an item of their own in the sidebar;
- add **commands** to the command palette, with keys you can change in
  Settings › Shortcuts, under "plugins";
- add **actions** to the menu of every track;
- add **widgets** to the now playing page, under the lyrics;
- control the **player** and the queue, and read what plays, and the
  spectrum of its sound;
- search YouTube Music and Spotify, and read albums, artists, playlists,
  radios and lyrics;
- read and change your **liked and disliked songs** and your mixes;
- read your **settings**, but not change them;
- show toasts and **desktop notifications**, open web pages in your browser,
  reach the **internet**, and keep up to 1 MB of its own **storage**.

What it draws, it draws in your theme's colours and type. Everything it
added goes when you turn it off.
[How MyTube works](how-it-works.md#themes-and-plugins) tells how plugins
run, and [the add-ons' part of the architecture](architecture/addons.md)
shows the code.

### What it asks for

The question before you enable a plugin lists what it asks for:

| It asks for | What it does |
|---|---|
| `playback` | controls the player |
| `network` | talks to the internet |
| `notifications` | shows desktop notifications |
| `storage` | keeps values between runs |

These tell you what the plugin does. They aren't enforced: every plugin can
use the whole API, which is why MyTube asks you to trust it.

## Remove

- In Settings › Plugins, the trash button beside a plugin removes it, after
  asking. The plugin and its files are deleted; what it stored stays.
- To remove a theme, delete its folder from the themes folder and reload
  add-ons.
- Built-in themes and plugins can't be removed.

## Reload add-ons

**Reload add-ons**, in the command palette, finds the themes and plugins
in your add-ons folders again, reads your
[colour file](look-and-feel.md#colours-from-a-file) again, and starts
again the plugins that failed. **reload** in Settings › Plugins does the
same. Run it after you put an add-on in a folder yourself, or change one.

## Make your own

- [Themes](themes.md) lists every colour, font and corner a theme sets,
  and shows a complete one.
- [Plugins](plugins.md) documents the plugin API, with its types, and how
  to try a plugin from the plugins folder.
