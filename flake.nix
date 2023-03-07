{ inputs =
    { nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
      flake-utils.url = "github:numtide/flake-utils";
    };

  nixConfig.bash-prompt = "[nix-develop-zk-Harness:] ";

  outputs = { nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system};
            python-pkgs = p: with p; [
              pip
            ];
        in
        {
          devShells.default =
            pkgs.mkShell {
              packages = [
                (pkgs.python3.withPackages python-pkgs)
              ];
              shellHook = ''
python -m venv pipenv
source ./pipenv/bin/activate
pip install -r requirements.txt
              '';
            };
        });
}
