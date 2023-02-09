pragma circom 2.1.2;

template Cubic() {
    signal input x;
    signal input y;
    
    //x**3 + x + 5 == y
    var n = 3;
    signal xs[n];
    xs[0] <== x;
    for (var i=1; i<n; i++) {
        xs[i] <== xs[i-1] * xs[0];
    }
    y === xs[n-1] + x + 5;
}

component main {public [y]} = Cubic();
