# Add a theme
<!-- docs: group=using order=80 -->

Themes are add-ons: a theme sets MyTube's colours, typefaces and corners.
MyTube loads those it ships with and those you put in your themes folder.

## Where themes go

The folder is in MyTube's
[data folder](data-and-privacy.md#the-data-folder), `addons/themes`:
**open themes folder** in
[Settings › Appearance](look-and-feel.md#themes) › Theme opens it.

Put the theme's folder there, named after its `id`, then
[reload add-ons](#reload-add-ons). A theme joins the theme cards of
Settings › Appearance, and works with every
[layout](look-and-feel.md#layouts) and option there. The repository's
[examples](../examples/README.md) are themes to start from.

## Remove

- To remove a theme, delete its folder from the themes folder and reload
  add-ons.
- Built-in themes can't be removed.

## Reload add-ons

**Reload add-ons**, in the command palette, finds the themes in your
themes folder again, and reads your
[colour file](look-and-feel.md#colours-from-a-file) again. Run it after
you put a theme in the folder yourself, or change one.

## Make your own

[Themes](themes.md) lists every colour, font and corner a theme sets,
and shows a complete one.
