/*
    Copyright 2018 0KIMS association.

    This file is part of circom (Zero Knowledge Circuit Compiler).

    circom is a free software: you can redistribute it and/or modify it
    under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    circom is distributed in the hope that it will be useful, but WITHOUT
    ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
    or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.

    You should have received a copy of the GNU General Public License
    along with circom. If not, see <https://www.gnu.org/licenses/>.
*/
pragma circom 2.0.3;

include "./sha256.circom";

template Main(N) {
    signal input in[N];
    signal input hash[32];
    signal output out[32];

    component sha256 = Sha256Bytes(N);
    sha256.in <== in;
    out <== sha256.out;

    for (var i = 0; i < 32; i++) {
        out[i] === hash[i];
    }

}

// render this file before compilation
component main = Main({TEMPLATE_VARS});
