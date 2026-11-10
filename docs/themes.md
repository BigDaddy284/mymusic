# Themes
<!-- docs: group=make order=10 -->

A theme sets the colors, the typefaces and the corners of MyTube, which
draws everything with them. A theme is a folder:

```
my-theme/
├── theme.json   the manifest: colors, fonts, radius
└── preview.png  optional: an image of the theme
```

## Make one

1. In Settings › Appearance, choose **open themes folder**.
2. Copy a theme there, for example
   [examples/themes/paper-mint](../examples/themes/paper-mint), into a folder
   named after its new `id`.
3. Change `id`, `name` and `author`, then the colors.
4. Run **[Reload add-ons](addons.md#reload-add-ons)** from the command
   palette (<kbd>Ctrl</kbd> <kbd>K</kbd>) and pick the theme in
   [Settings › Appearance](look-and-feel.md#themes). Themes that fail to load
   are listed there, with the reason.

Every change to `theme.json` shows after a reload.

## theme.json

```json
{
  "id": "my-theme",
  "name": "My Theme",
  "version": "1.0.0",
  "author": "Your Name",
  "description": "What it looks like, in one sentence.",
  "scheme": "dark",
  "colors": {
    "bg": "#000000",
    "surface": "#0b0b0b",
    "surface-2": "#181818",
    "surface-3": "#272727",
    "border": "rgba(255, 255, 255, 0.16)",
    "text": "#ffffff",
    "text-muted": "#a3a3a3",
    "accent": "#00bde3",
    "accent-text": "#000000",
    "danger": "#ff5a4f",
    "success": "#a4e600",
    "youtube": "#ff4040",
    "tile-1": "#00bde3",
    "tile-2": "#e6007e",
    "tile-3": "#a4e600",
    "tile-4": "#ffb000",
    "tile-text": "#000000"
  },
  "fonts": {
    "body": "Selawik, 'Segoe UI', system-ui, sans-serif",
    "display": "Selawik, 'Segoe UI', system-ui, sans-serif"
  },
  "radius": "0px",
  "files": ["theme.json"]
}
```

| Field | |
|---|---|
| `id` | lowercase letters, digits and dashes; the folder's name |
| `version` | a semantic version, as `1.2.0`; raise it with every change you publish |
| `scheme` | `dark` or `light`: the controls MyTube draws, such as switches and sliders, follow it |
| `colors` | every color below, as `#rgb`, `#rrggbb`, `#rrggbbaa`, `rgb()`, `rgba()`, `hsl()`, `oklch()` or `transparent` |
| `fonts` | font-family lists, of which MyTube uses the first the computer has: `body` for text, `display` for titles and the large lowercase words (the body's when left out) |
| `radius` | the corner radius of tiles and cards in the Tiles layout, as `12px`; smaller things round by 60% of it |
| `preview` | an image of the theme |
| `files` | every file of the theme |

### Colors

| Color | Used for |
|---|---|
| `bg` | the window's ground |
| `surface`, `surface-2`, `surface-3` | panels, then hovered and pressed things on them |
| `border` | hairlines between rows and panels |
| `text`, `text-muted` | ink, and secondary ink; outlined controls are drawn in `text` |
| `accent`, `accent-text` | everything live: the playing row, the selection, the active page, primary buttons, and ink on them |
| `danger`, `success` | errors and confirmations |
| `youtube` | the mark of YouTube Music's tracks |
| `tile-1` … `tile-4`, `tile-text` | the flat tiles of the home page and the stats, and ink on them |

A theme may also set these, as the built-in MyTube themes do:

| Token | Used for |
|---|---|
| `accent-2` | the far end of a gradient from `accent`, left to right, which fills the primary buttons, the play buttons and the seek bar, and, faintly, what is chosen; the user's accent or the artwork's turns it off |
| `hover`, `pressed` | solid fills of what the pointer is over and presses, in place of tints of `text` |

Keep text readable: `text-muted` on `bg`, `accent-text` on `accent` and
`tile-text` on every tile need a contrast of at least 4.5:1. When the user
picks an accent of their own, or lets it follow the artwork, MyTube picks
the accent's ink itself, and keeps the accent readable on your `bg`.

The Studio, Blend, Expressive and Crate layouts fill what is chosen, as the
current page, and their secondary buttons with `surface-3`: `text` on it
needs 4.5:1 as well. They set the accent as small text on the panels, so
there MyTube darkens or lightens yours until it reads at 4.5:1 on
`surface`.

### Fonts

MyTube ships Selawik, Quicksand, Mochi's round face, Roboto Flex, the
Material themes' face, at weights 400, 500 and 700, and Archivo
SemiCondensed at 800, the Walnut themes' headings'; other families must be
on the user's computer, so end each list with a family Windows has:
`system-ui` or `sans-serif` for the system's font, `monospace` for its
monospaced one.

### Layouts

A theme sets the colors, the faces and the corners; the layout is the
user's to choose, in [Settings › Appearance](look-and-feel.md#layouts),
with any theme:

- **Blend**, the default: Expressive's rounded shapes and springs, on
  Crate's dense rows, big covers in cards and bold headings.
- **Expressive**, after Material 3 Expressive: tonal fills, big rounded
  shapes, pill buttons, a floating player bar and springy motion.
- **Crate**: dense rows, big covers in cards and bold, tight headings.
- **Tiles**, MyTube's first: flat tiles, giant light lowercase words
  and outlined boxes. Only Tiles rounds with your `radius`; the others
  have corners of their own.

The layouts set your `display` face at their own weights: light in Tiles,
regular in Expressive, and heavy in Blend and Crate. A face with those
weights looks its best in all of them.

## Effects

The themes MyTube ships with add effects of their own, which MyTube
draws and the user tunes in
[Settings › Appearance](look-and-feel.md#theme-options): Mix's tile textures,
Aurora's drifting lights and frosted glass, Sakura's petals, Terminal's
scanlines and glow, Greenroom's and Glasshouse's grounds in the colour of
the cover playing, Anime's speed lines and screentone, and Mochi's rising
sparkles, and the Material themes' colours, which MyTube makes of a
seed colour in Material's HCT color space, as Material's dynamic color
does: the cover's while it plays and "Follow the artwork" is on, the
accent you chose otherwise, or MyTube's pink. A theme of yours sets the
tokens above, and MyTube draws
everything else with them.

Themes made for the web interface of MyTube 0.4 and before still load:
MyTube reads their colors, fonts and radius, and ignores their stylesheet
and their options.
