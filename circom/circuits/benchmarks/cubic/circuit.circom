pragma circom 2.1.2;

template Cubic() {
    signal input X;
    signal input Y;
    
    //x**3 + x + 5 == y
    var n = 3;
    signal xs[n];
    xs[0] <== X;
    for (var i=1; i<n; i++) {
        xs[i] <== xs[i-1] * xs[0];
    }
    Y === xs[n-1] + X + 5;
}

component main {public [Y]} = Cubic();
