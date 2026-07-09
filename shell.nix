{ pkgs ? import (builtins.fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/refs/heads/nixos-26.05.tar.gz";
  }) {}
}:

pkgs.mkShell {
  packages = with pkgs; [
    go
    jq
    yq
    gnumake
    docker
    dive
    grype
    jdk
    graphviz
    plantuml
    goreleaser
    golangci-lint
  ];

  shellHook = ''
    echo "Go development shell"
    echo "nixpkgs channel: nixos-26.05"
    echo "Go: $(go version)"
  '';
}
