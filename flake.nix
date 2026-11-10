{
  description = "MyTube: YouTube Music and your own music in one customizable desktop player";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { self, nixpkgs }:
    let
      inherit (nixpkgs) lib;
      # MyTube is built for Windows alone, cross-compiled on Linux: the
      # systems the shell runs on, not those MyTube runs on.
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      each = systems: f: lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      # The tools releases are built and published with, at the versions
      # flake.lock pins: Go, Bun for `mygo build`, NSIS for the Windows
      # installers, osslsigncode, which signs the Windows apps when
      # release.yml has a certificate, and GoReleaser; and Zola, which
      # builds the website.
      # release.yml, ci.yml and pages.yml run in it, as `nix develop` does
      # here.
      devShells = each systems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_27
            bun
            goreleaser
            jq
            zip
            zola
            nsis
            osslsigncode
            binutils
          ];
        };
      });

      formatter = each systems (pkgs: pkgs.nixfmt);
    };
}
