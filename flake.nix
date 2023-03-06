{ inputs =
    { nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
      flake-utils.url = "github:numtide/flake-utils";
    };

  outputs = { nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system};
            python-pkgs = p: with p; [
              # dash-bootstrap-components
              brotli
              click
              dash
              dash-core-components
              dash-html-components
              dash-table
              flask
              flask-compress
              itsdangerous
              jinja2
              markupsafe
              plotly
              six
              tenacity
              werkzeug
              pandas
            ];
        in
        {
          devShells.default =
            pkgs.mkShell {
              packages = [
                (pkgs.python3.withPackages python-pkgs)
              ];
            };
        });
}
