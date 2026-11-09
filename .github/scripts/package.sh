#!/usr/bin/env bash
# Makes the download mygo build does not: a portable zip of each Windows
# build, from what mygo build left in build/<target>/.
#
#   .github/scripts/package.sh <version> [out]
set -euo pipefail
shopt -s nullglob

v="${1:?usage: package.sh <version> [out]}"
out="${2:-release}"
mkdir -p "$out"
out="$(cd "$out" && pwd)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

# Windows: the app and the files beside it, in a MyMusic folder that runs
# from wherever it is unpacked.
for dir in build/windows-*/; do
	arch="$(basename "$dir")"
	arch="${arch#windows-}"
	stage="$work/portable-$arch/MyMusic"
	mkdir -p "$stage"
	cp "$dir"MyMusic.exe "$stage/"
	for f in "$dir"*.png; do cp "$f" "$stage/"; done
	(cd "$work/portable-$arch" && zip -qr9 "$out/mymusic-$v-windows-$arch-portable.zip" MyMusic)
	echo "made mymusic-$v-windows-$arch-portable.zip"
done
