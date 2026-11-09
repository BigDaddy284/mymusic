# MyMusic for Linux, built from source. MyGo needs no C toolchain: it
# loads GTK, WebKitGTK, the tray's AppIndicator, the keyring and the sound
# libraries as it runs, which the wrapper finds for it.
{
  lib,
  buildGoModule,
  go_1_27,
  wrapGAppsHook3,
  copyDesktopItems,
  makeDesktopItem,
  gtk3,
  glib,
  glib-networking,
  gsettings-desktop-schemas,
  cairo,
  pango,
  at-spi2-core,
  gdk-pixbuf,
  libepoxy,
  webkitgtk_4_1,
  libsoup_3,
  libayatana-appindicator,
  libdbusmenu-gtk3,
  gtk-layer-shell,
  libsecret,
  libpulseaudio,
  alsa-lib,
  version,
}:

let
  # What MyMusic loads as it runs, by these libraries' names.
  runtimeLibraries = [
    gtk3
    glib
    cairo
    pango
    at-spi2-core
    gdk-pixbuf
    libepoxy
    webkitgtk_4_1
    libsoup_3
    libayatana-appindicator
    libdbusmenu-gtk3
    gtk-layer-shell
    libsecret
    libpulseaudio
    alsa-lib
  ];
in
(buildGoModule.override { go = go_1_27; }) {
  pname = "mymusic";
  inherit version;

  src = lib.cleanSource ../.;
  # The hash of the Go modules go.sum names: a change of go.sum changes it,
  # and `nix build` then tells the new one.
  vendorHash = "sha256-Cwk5Ebjti0k+LZ2F9v7pNmjbYDdU6CMQu2tMrr5S5aU=";

  subPackages = [ "." ];
  env.CGO_ENABLED = "0";
  ldflags = [
    "-s"
    "-w"
  ];
  # CI runs the tests; some of them need the network or a sound card.
  doCheck = false;

  nativeBuildInputs = [
    wrapGAppsHook3
    copyDesktopItems
  ];
  # For the wrapper: GSettings' schemas, and TLS for WebKitGTK's sign-in.
  buildInputs = runtimeLibraries ++ [
    glib-networking
    gsettings-desktop-schemas
  ];

  preFixup = ''
    gappsWrapperArgs+=(--prefix LD_LIBRARY_PATH : "${lib.makeLibraryPath runtimeLibraries}")
  '';

  postInstall = ''
    install -Dm644 brand/icon.png $out/share/pixmaps/mymusic.png
  '';

  desktopItems = [
    (makeDesktopItem {
      name = "mymusic";
      desktopName = "MyMusic";
      comment = "Spotify and YouTube Music in one customizable player";
      exec = "mymusic %U";
      icon = "mymusic";
      categories = [
        "AudioVideo"
        "Audio"
        "Player"
      ];
      mimeTypes = [ "x-scheme-handler/mymusic" ];
    })
  ];

  meta = {
    description = "Spotify and YouTube Music in one customizable desktop player";
    homepage = "https://github.com/bigdaddy284/mymusic";
    license = with lib.licenses; [
      mit
      gpl3Only
    ];
    mainProgram = "mymusic";
    platforms = lib.platforms.linux;
  };
}
