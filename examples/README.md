# Example themes and plugins

Working examples of a theme and a plugin, to read, copy and change:

- [themes/neon-nights](themes/neon-nights) and
  [themes/paper-mint](themes/paper-mint), made as
  [docs/themes.md](../docs/themes.md) tells;
- [plugins/sleep-timer](plugins/sleep-timer) and
  [plugins/spectrum](plugins/spectrum), written in TypeScript against the
  plugin API of [docs/plugins.md](../docs/plugins.md).

To use one, copy its folder into your add-ons folder: Settings ›
Appearance › open themes folder for a theme, Settings › Plugins › open
plugins folder for a plugin. A plugin runs as the JavaScript its
TypeScript compiles to, which Git leaves out:

```sh
bun install
bun run build:plugins
```

then copy `plugin.json` and `index.js` of its folder.
