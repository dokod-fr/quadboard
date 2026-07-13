{
  description = "QuadBoard development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gotools
            golangci-lint
            gofumpt

            templ
            air
            go-task # Task runner for Go
            cobra-cli # subcommand generator for Go applications
            git

            nil
            nixpkgs-fmt
          ];

          shellHook = ''
            echo ""
            echo "🚀 Welcome to the QuadBoard development shell"
            echo ""
            echo "Go      : $(go version)"
            echo "Task    : $(task --version)"
            echo "Templ   : $(templ version || true)"
            echo ""
          '';
        };
      });
}