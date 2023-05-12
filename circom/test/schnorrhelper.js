const { Scalar } = require("ffjavascript");
const { buildBabyjub, buildPoseidon } = require("circomlibjs");

async function buildSchnorr() {
    const babyJub = await buildBabyjub("bn128");
    const poseidon = await buildPoseidon();
    return new Schnorr(babyJub, poseidon);
}
// base field: 21888242871839275222246405745257275088548364400416034343698204186575808495617
class Schnorr {

    constructor(babyJub, poseidon) {
        this.babyJub = babyJub;
        this.poseidon = poseidon;
        this.F = babyJub.F;
    }
    //converts our private key to our public key
    prv2pub(prv) {
        let pub = this.babyJub.mulPointEscalar(this.babyJub.Base8, prv); //pub  = g^prv where g is the base point of babyjub
        return pub;
    }

    //calculates the signature for schnorr
    signPoseidon(prv, msg, k) {
        //calulcate r = g^k
        const F = this.babyJub.F;
        const r = this.babyJub.mulPointEscalar(this.babyJub.Base8, k);
        //calculate H(r||M) = e where H is the poseidon has function
        let e = this.poseidon([F.toObject(r[0]), F.toObject(r[1]), msg]);
        e = F.toObject(e);
        //calculate s = k - prv*e
        let s = Scalar.sub(k, Scalar.mul(prv, e)); 
        s = Scalar.mod(s, this.babyJub.subOrder); //we must ensure that s is positive and within our subgroup order
        s = Scalar.add(s, this.babyJub.subOrder);
        s = Scalar.mod(s, this.babyJub.subOrder);
        //return signature scheme
        return {
            e: e,
            s: s
        };
    }
    //signature = (s,e) 
    //verifies that e_v and e are the same
    verifyPoseidon(sig, y, msg) {

        // Check parameters for schnorr
        if (typeof sig != "object") return false;
        if (!Array.isArray(y)) return false; // making sure that y
        if (y.length != 2) return false;
        if (!this.babyJub.inCurve(y)) return false; //making sure that y is on the baby jub curve
        if (sig.s >= this.babyJub.subOrder) return false;
        const e = sig.e;
        const gs = this.babyJub.mulPointEscalar(this.babyJub.Base8, sig.s); //calculates g^s
        const ye = this.babyJub.mulPointEscalar(y, sig.e); //calculates y^e
        let rv = this.babyJub.addPoint(gs, ye); //adds g^s and y^e 
        let ev = this.poseidon([F.toObject(rv[0]), F.toObject(rv[1]), msg]); //H(r_v || M)
        ev = F.toObject(ev);
        if (!Scalar.eq(e, ev)) return false; //checks if e == e_v
        return true;
    }
}

module.exports = { buildSchnorr };
