#!/usr/bin/env bash
# Builds MyMusic for each Windows platform given, as `bun run build --
# -platform` does, one platform at a time.
#
#   .github/scripts/build.sh windows/amd64,windows/arm64
#
# After each, the log lists what its archive holds, and the script stops
# if it ships Go source (#303).
set -euo pipefail

# shipped lists what the archive of the build in dir holds, and fails if
# the build ships Go source, beside the app or in its archive: `mygo
# build` ships every file of resources/ with the app, as 0.5.9's did
# resources.go (#303).
shipped() {
	local dir="$1" archive found
	for archive in "$dir"/*.tar.gz; do
		if [ -f "$archive" ]; then
			echo "$archive holds:"
			tar -tzf "$archive" | sed 's/^/  /'
		fi
	done
	found="$(
		find "$dir" -name '*.go'
		for archive in "$dir"/*.tar.gz; do
			if [ -f "$archive" ]; then
				tar -tzf "$archive" | awk -v a="$archive" '/\.go$/ { print a ": " $0 }'
			fi
		done
	)"
	if [ -n "$found" ]; then
		echo "build.sh: the build ships Go source, which the app never reads:" >&2
		echo "$found" >&2
		echo "mygo build ships every file of resources/ with the app: keep files the app doesn't read as it runs out of it." >&2
		exit 1
	fi
}

IFS=, read -r -a platforms <<<"${1:?usage: build.sh <platform>[,<platform>...]}"
for p in "${platforms[@]}"; do
	case "$p" in
	windows/*) ;;
	*)
		echo "build.sh: MyMusic builds for Windows alone, not $p" >&2
		exit 1
		;;
	esac
	bun run build -- -platform "$p"
	shipped "build/${p/\//-}"
done
