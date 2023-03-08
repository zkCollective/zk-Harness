{ inputs =
    { nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
      flake-utils.url = "github:numtide/flake-utils";
      circom.url = "github:Polytopoi/circom/nix";
      gnark.url = "github:Polytopoi/gnark/nix";
    };

  nixConfig.bash-prompt = "[nix-develop-zk-Harness:] ";

  outputs = { nixpkgs, flake-utils, circom, gnark, ... }@inputs:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system};
            circom-out = circom.defaultPackage.${system};
            gnark-out = gnark.packages.${system}.default;
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
                gnark-out
                nodejs-19_x
                time
                zsh
              ];
              shellHook = ''
python -m venv pipenv
source ./pipenv/bin/activate
pip install -r requirements.txt
npm install
alias snarkjs="npx snarkjs"
              '';
            };
        });
}
