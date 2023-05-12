pragma circom 2.0.0;
include "../circomlib/comparators.circom";
include "../circomlib/compconstant.circom";
include "../circomlib/poseidon.circom";
include "../circomlib/bitify.circom";
include "../circomlib/escalarmulany.circom";
include "../circomlib/escalarmulfix.circom";
//This circom code is used to constrain the schnorr signature.
// Based off of the wikipedia article: https://en.wikipedia.org/wiki/Schnorr_signature
// Inspired by Circom's EdDSA Poseidon: https://github.com/iden3/circomlib/blob/master/circuits/eddsaposeidon.circom
template SchnorrPosedion(q, gx, gy){ //q is the subgroup order, g = (g_x,g_y) is the base boint of Baby Jub
    signal input enabled; 
    signal input M; //message
    //y = (yx, yy) = g^x is a point
    signal input yx;
    signal input yy; 
    signal input S; //S = k - xe is apart of our signature
    signal input e; //e is apart of our signature

    //need to ensure that S is in our subgroup (from eddsa code)
    component snum2bits = Num2Bits(253);
    snum2bits.in <== S;
    component compConstant = CompConstant(q);
    for (var i=0; i<253; i++) {
        snum2bits.out[i] ==> compConstant.in[i];
    }
    compConstant.in[253] <== 0;
    compConstant.out*enabled === 0;

    //calculate g^s
    component mulAny = EscalarMulAny(253);
    for(var i = 0; i<253; i++){
        mulAny.e[i] <== snum2bits.out[i];
    }
    mulAny.p[0] <== gx;
    mulAny.p[1] <== gy;

    //calculate y^e
    component enum2bits = Num2Bits(254);
    enum2bits.in <== e;

    component mulAny1 = EscalarMulAny(254);
    for(var i = 0; i<254; i++){
        mulAny1.e[i] <== enum2bits.out[i];
    }
    mulAny1.p[0] <== yx;
    mulAny1.p[1] <== yy;

    //rv = g^sy^e (which is just adding g^s and y^e)
    component add1 = BabyAdd();
    add1.x1 <== mulAny.out[0];
    add1.y1 <== mulAny.out[1];
    add1.x2 <== mulAny1.out[0];
    add1.y2 <== mulAny1.out[1];

    //hash H(rv || M)
    component ev = Poseidon(3);
    ev.inputs[0] <== add1.xout;
    ev.inputs[1] <== add1.yout;
    ev.inputs[2] <== M;


    //check if e == ev
    component eqCheck = ForceEqualIfEnabled();
    eqCheck.enabled <== enabled;
    eqCheck.in[0] <== e;
    eqCheck.in[1] <== ev.out;
}