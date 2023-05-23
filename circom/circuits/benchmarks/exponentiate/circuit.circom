pragma circom 2.1.2;

template Exponentiate(E) {
    signal input X;
    signal input Y;

    signal res[E];

    //X**E == Y
    res[0] <== X;
    for (var i=1; i<E; i++) {
        res[i] <== res[i-1] * X;
    }
    res[E-1] === Y;
}

component main {public [X, Y]} = Exponentiate({TEMPLATE_VARS});
