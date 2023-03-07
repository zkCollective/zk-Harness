{ inputs =
    { nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
      flake-utils.url = "github:numtide/flake-utils";
      circom.url = "github:Polytopoi/circom/nix";
    };

  nixConfig.bash-prompt = "[nix-develop-zk-Harness:] ";

  outputs = { nixpkgs, flake-utils, circom, ... }@inputs:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system};
            circom-out = circom.defaultPackage.${system};
            python-pkgs = p: with p; [
              pip
            ];
        in
        {
          devShells.default =
            pkgs.mkShell {
              packages = with pkgs; [
                (python3.withPackages python-pkgs)
                gnumake
                circom-out
              ];
              shellHook = ''
python -m venv pipenv
source ./pipenv/bin/activate
pip install -r requirements.txt
              '';
            };
        });
}
