# Add-ons

- **`internal/addons`** validates theme manifests, lists the
  built-in add-ons and those of the user's add-ons folders, and removes
  the user's.
- **`internal/colorfile`** reads the colours of pywal's `colors.json`,
  matugen's JSON, MyTube's matugen template and a `theme.json`, telling
  them apart by what they hold; `internal/app/colorfile.go` follows the
  file the settings name with MyGo's file watch, and keeps the colours it
  read last while the file is half written.
