"""Draws MyMusic's icon, five sound bars, and the Music mark, a music note, in the theme's accent
gradient, the icon on a dark rounded tile.

The same style as StackSpace's and Stills & Reels' icons (StackSpace tool/make_icon.py): a dark
tile with a faint gradient rim, a soft glow, and a two-tone gradient mark, blue #2668FF to
purple #B048FF.

Run from the repository root with Pillow installed:
    python3 brand/make_icon.py
Writes the app's icon, the bars, to brand/icon.png (1024 px), brand/tray.png (32 px),
brand/icon.svg and site/static/assets/icon.png, and the Music mark, the note the sidebar shows
beside its name, to brand/mark.svg.
"""
from PIL import Image, ImageDraw, ImageFilter, ImageChops

S = 1024
SS = 4  # supersampling for the glyph masks
A = (38, 104, 255)    # theme accentA
B = (176, 72, 255)    # theme accentB
LIGHT_TOP, LIGHT_BOTTOM = (70, 140, 255), (190, 100, 255)
DARK_TOP, DARK_BOTTOM = (40, 90, 220), (130, 60, 220)

# The note, in the tile's 1024 units: a slanted beam, two stems and two heads.
BEAM = [(395, 250), (735, 180), (735, 285), (395, 355)]
STEMS = [(395, 300, 455, 700), (675, 230, 735, 630)]
HEADS = [(350, 705, 105, 80), (630, 635, 105, 80)]

# The bars, in the tile's 1024 units: their heights, around the middle, and which are in the
# lighter tone.
BAR_HEIGHTS = [260, 470, 640, 400, 560]
BAR_W, BAR_GAP = 96, 44
BAR_LIGHT = (1, 2, 4)


def bar_rects():
    x = 512 - (5 * BAR_W + 4 * BAR_GAP) / 2
    for i, h in enumerate(BAR_HEIGHTS):
        yield i in BAR_LIGHT, (x, 512 - h / 2, x + BAR_W, 512 + h / 2)
        x += BAR_W + BAR_GAP


def gradient(size, c1, c2):
    w, h = size
    g = Image.new('RGB', (w, h))
    px = g.load()
    for y in range(h):
        for x in range(w):
            t = (x + y) / (w + h - 2)
            px[x, y] = tuple(int(c1[i] + (c2[i] - c1[i]) * t) for i in range(3))
    return g


def vertical(c1, c2, y0, y1):
    g = Image.new('RGB', (1, S))
    for y in range(S):
        t = min(1, max(0, (y - y0) / (y1 - y0)))
        g.putpixel((0, y), tuple(int(c1[i] + (c2[i] - c1[i]) * t) for i in range(3)))
    return g.resize((S, S))


def tile_mask(inset=0):
    m = Image.new('L', (S, S), 0)
    r = int(S * 0.225)
    ImageDraw.Draw(m).rounded_rectangle((inset, inset, S - 1 - inset, S - 1 - inset), radius=r - inset, fill=255)
    return m


def bar_masks():
    k = SS
    light = Image.new('L', (S * k, S * k), 0)
    dark = Image.new('L', (S * k, S * k), 0)
    for lit, (x0, y0, x1, y1) in bar_rects():
        ImageDraw.Draw(light if lit else dark).rounded_rectangle(
            (x0 * k, y0 * k, x1 * k, y1 * k), radius=BAR_W * k / 2, fill=255)
    return light.resize((S, S), Image.LANCZOS), dark.resize((S, S), Image.LANCZOS)


def glyph_masks():
    k = SS
    light = Image.new('L', (S * k, S * k), 0)
    dark = Image.new('L', (S * k, S * k), 0)
    dl, dd = ImageDraw.Draw(light), ImageDraw.Draw(dark)
    dl.polygon([(x * k, y * k) for x, y in BEAM], fill=255)
    for x0, y0, x1, y1 in STEMS:
        dl.rectangle((x0 * k, y0 * k, x1 * k, y1 * k), fill=255)
    for cx, cy, rx, ry in HEADS:
        dd.ellipse(((cx - rx) * k, (cy - ry) * k, (cx + rx) * k, (cy + ry) * k), fill=255)
    light = light.resize((S, S), Image.LANCZOS)
    dark = dark.resize((S, S), Image.LANCZOS)
    return light, ImageChops.subtract(dark, light)


def build():
    out = Image.new('RGBA', (S, S), (0, 0, 0, 0))
    bg = gradient((64, 64), (30, 30, 44), (12, 12, 18)).resize((S, S), Image.BILINEAR)
    tile = tile_mask()
    out.paste(bg, (0, 0), tile)
    rim = ImageChops.subtract(tile, tile_mask(10))
    acc = gradient((256, 256), A, B).resize((S, S), Image.BILINEAR)
    rim_layer = Image.new('RGBA', (S, S), (0, 0, 0, 0))
    rim_layer.paste(acc, (0, 0), rim.point(lambda v: v * 150 // 255))
    out = Image.alpha_composite(out, rim_layer)
    light, dark = bar_masks()
    g = ImageChops.lighter(light, dark)
    y0, y1 = g.getbbox()[1], g.getbbox()[3]
    glow_a = ImageChops.multiply(g.filter(ImageFilter.GaussianBlur(48)).point(lambda v: min(255, v * 2)), tile)
    glow = Image.new('RGBA', (S, S), (0, 0, 0, 0))
    glow.paste(vertical(A, B, y0, y1), (0, 0), glow_a.point(lambda v: v * 200 // 255))
    out = Image.alpha_composite(out, glow)
    glyph = Image.new('RGBA', (S, S), (0, 0, 0, 0))
    glyph.paste(vertical(DARK_TOP, DARK_BOTTOM, y0, y1), (0, 0), dark)
    glyph.paste(vertical(LIGHT_TOP, LIGHT_BOTTOM, y0, y1), (0, 0), light)
    return Image.alpha_composite(out, glyph)


def hexc(c):
    return '#%02x%02x%02x' % c


def glyph_svg():
    """The note's shapes, with gradients across its height in user units."""
    y0 = min(y for _, y in BEAM)
    y1 = max(cy + ry for _, cy, _, ry in HEADS)
    beam = ' '.join('%d,%d' % p for p in BEAM)
    stems = ''.join('<rect x="%d" y="%d" width="%d" height="%d"/>' % (x0, a, x1 - x0, b - a) for x0, a, x1, b in STEMS)
    heads = ''.join('<ellipse cx="%d" cy="%d" rx="%d" ry="%d"/>' % h for h in HEADS)
    defs = (
        f'<linearGradient id="light" gradientUnits="userSpaceOnUse" x1="0" y1="{y0}" x2="0" y2="{y1}">'
        f'<stop offset="0" stop-color="{hexc(LIGHT_TOP)}"/><stop offset="1" stop-color="{hexc(LIGHT_BOTTOM)}"/></linearGradient>'
        f'<linearGradient id="dark" gradientUnits="userSpaceOnUse" x1="0" y1="{y0}" x2="0" y2="{y1}">'
        f'<stop offset="0" stop-color="{hexc(DARK_TOP)}"/><stop offset="1" stop-color="{hexc(DARK_BOTTOM)}"/></linearGradient>'
    )
    body = f'<g fill="url(#dark)">{heads}</g><g fill="url(#light)"><polygon points="{beam}"/>{stems}</g>'
    return defs, body


def bars_svg():
    """The bars' shapes, with gradients across their height in user units."""
    rects = list(bar_rects())
    y0 = min(r[1] for _, r in rects)
    y1 = max(r[3] for _, r in rects)
    defs = (
        f'<linearGradient id="light" gradientUnits="userSpaceOnUse" x1="0" y1="{y0:g}" x2="0" y2="{y1:g}">'
        f'<stop offset="0" stop-color="{hexc(LIGHT_TOP)}"/><stop offset="1" stop-color="{hexc(LIGHT_BOTTOM)}"/></linearGradient>'
        f'<linearGradient id="dark" gradientUnits="userSpaceOnUse" x1="0" y1="{y0:g}" x2="0" y2="{y1:g}">'
        f'<stop offset="0" stop-color="{hexc(DARK_TOP)}"/><stop offset="1" stop-color="{hexc(DARK_BOTTOM)}"/></linearGradient>'
    )
    def group(lit):
        return ''.join('<rect x="%g" y="%g" width="%g" height="%g" rx="%g"/>' % (x0, a, x1 - x0, b - a, BAR_W / 2)
                       for l, (x0, a, x1, b) in rects if l == lit)
    body = f'<g fill="url(#dark)">{group(False)}</g><g fill="url(#light)">{group(True)}</g>'
    return defs, body


def write_svgs():
    defs, body = glyph_svg()
    note = '<!-- MyMusic\'s Music mark: drawn by brand/make_icon.py; edit that, not this. -->\n'
    with open('brand/mark.svg', 'w') as f:
        f.write('<svg xmlns="http://www.w3.org/2000/svg" viewBox="180 172 620 620">\n' + note)
        f.write(f'<defs>{defs}</defs>{body}\n</svg>\n')
    defs, body = bars_svg()
    note = '<!-- MyMusic\'s icon: drawn by brand/make_icon.py; edit that, not this. -->\n'
    tile = (
        '<linearGradient id="tile" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#1e1e2c"/><stop offset="1" stop-color="#0c0c12"/></linearGradient>'
        f'<linearGradient id="rim" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="{hexc(A)}"/><stop offset="1" stop-color="{hexc(B)}"/></linearGradient>'
        '<filter id="glow" x="-50%" y="-50%" width="200%" height="200%"><feGaussianBlur stdDeviation="40"/></filter>'
        '<clipPath id="clip"><rect width="1024" height="1024" rx="230"/></clipPath>'
    )
    with open('brand/icon.svg', 'w') as f:
        f.write('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" width="1024" height="1024">\n' + note)
        f.write(f'<defs>{defs}{tile}</defs>')
        f.write('<rect width="1024" height="1024" rx="230" fill="url(#tile)"/>')
        f.write('<rect x="5" y="5" width="1014" height="1014" rx="225" fill="none" stroke="url(#rim)" stroke-width="10" stroke-opacity="0.6"/>')
        f.write(f'<g clip-path="url(#clip)"><g filter="url(#glow)" opacity="0.8">{body}</g></g>{body}\n</svg>\n')


if __name__ == '__main__':
    img = build()
    img.save('brand/icon.png', optimize=True)
    img.resize((32, 32), Image.LANCZOS).save('brand/tray.png', optimize=True)
    img.save('site/static/assets/icon.png', optimize=True)
    write_svgs()
