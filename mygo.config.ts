import { defineConfig } from "mygo-cli";
import pkg from "./package.json" with { type: "json" };

export default defineConfig({
  name: "MyMusic",
  // The slug names the update archives and the Windows installer, which
  // release.yml and .github/scripts/package.sh look for by name. It is
  // MyGo's default for "MyMusic": pinned so it can't move.
  slug: "mymusic",
  identifier: "app.mymusic.desktop",
  version: pkg.version,
  copyright: "© 2026 MyMusic contributors, based on MixMusic, © 2026 MixMusic contributors",
  // The icon that `mygo build` makes the Windows icons of. It isn't
  // in resources/, MyGo's place for it by default: `mygo build` ships every
  // file there beside the app, which reads none as it runs (#303).
  icon: "brand/icon.png",
  // MyGo draws the windows itself, in Go (the root package): the only page
  // is that of the picture of music videos, public/video.html with
  // video.js, which `mygo build` embeds as it is. They call Go without a generated
  // client, so `mygo dev` writes it out of the way.
  frontendDist: "public",
  bindings: ".mygo/bindings.ts",
  out: "build",
  urlSchemes: ["mymusic"],
  // Apps update themselves from the releases of the repository, signed with
  // the key whose public half is here (the secret one is the
  // MYGO_UPDATER_PRIVATE_KEY secret of the release workflow).
  updates: {
    publicKey: "y/g7m+MircfXqdGupUFNygznyyR/Y41D1YAOYw41ISU=",
    github: "bigdaddy284/mymusic",
  },
  // The workflow's code signing certificate, whose password is
  // MYGO_WINDOWS_CERTIFICATE_PASSWORD: osslsigncode signs the app, the
  // installer and its uninstaller with it, timestamped. Unsigned without it.
  windows: {
    certificate: process.env.WINDOWS_CERTIFICATE_FILE || undefined,
  },
});
