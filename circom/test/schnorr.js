const path = require("path");
const assert = require("chai").assert;
const wasm_tester = require("circom_tester").wasm;
const { buildSchnorr } = require("./schnorrhelper.js");
const { Scalar } = require("ffjavascript");
const buildBabyjub = require("circomlibjs").buildBabyjub;

// A helper function to generate a random integer within our base field.
// Inspired by Circom's EdDSA Poseidon Test: https://github.com/iden3/circomlib/blob/master/test/eddsa.js
function generateRandomBigInt(maxBits) {
    const words = Math.ceil(maxBits / 32);
    const arr = new Uint32Array(words);
    crypto.getRandomValues(arr);
    // clear any excess bits in the last word
    const excessBits = (words * 32) - maxBits;
    if (excessBits > 0) {
        const mask = (1 << (32 - excessBits)) - 1;
        arr[words - 1] &= mask;
    }
    return BigInt(arr.join(''));
}

describe("Schnorr Test", () => {
    var circ_file = path.join(__dirname, "circuits", "schnorr_test.circom");
    var circ, num_constraints;
    let schnorr;

    before(async () => {
        schnorr = await buildSchnorr();
        babyJub = await buildBabyjub();
        F = babyJub.F;
        circuit = await wasm_tester(circ_file);


        await circuit.loadConstraints();
        num_constraints = circuit.constraints.length;
        console.log("Schnorr #Constraints:", num_constraints);

    });

    it("Sign a single number", async () => {
        //message
        let msg = F.e(1234);
        msg = F.toObject(msg);

        //random private key
        const prvKey = generateRandomBigInt(253);

        //random integer selected by the verifier
        const k = generateRandomBigInt(253);

        //obtain private key from public key
        const pubKey = schnorr.prv2pub(prvKey);

        //obtain the signature
        const signature = schnorr.signPoseidon(prvKey, msg, k);
        //verify that everything is correct
        assert(schnorr.verifyPoseidon(signature, pubKey, msg));

        let input = {
            "enabled": "1",
            "M": msg.toString(),
            "yx": F.toObject(pubKey[0]).toString(),
            "yy": F.toObject(pubKey[1]).toString(),
            "S": signature.s.toString(),
            "e": signature.e.toString()
        };
        const w = await circuit.calculateWitness(input, true);
        await circuit.checkConstraints(w);
    });

    it("Detect Invalid signature", async () => {
        let msg = F.e(1234);
        msg = F.toObject(msg);
        const prvKey = generateRandomBigInt(253);
        const k = generateRandomBigInt(253);
        const pubKey = schnorr.prv2pub(prvKey);
        const signature = schnorr.signPoseidon(prvKey, msg, k);
        assert(schnorr.verifyPoseidon(signature, pubKey, msg));
        try {
            let yy = F.toObject(pubKey[1]);
            yy = Scalar.add(yy, "1"); //this makes the signature invalid
            await circuit.calculateWitness({
                "enabled": "1",
                "M": msg.toString(),
                "yx": F.toObject(pubKey[0]).toString(),
                "yy": yy.toString(),
                "S": signature.s.toString(),
                "e": signature.e.toString()

            }, true);
            assert(false);
        } catch (err) {
            console.log("err ", err);
            assert(err.message.includes("Assert Failed"));
        }
    });

    it("Test a disabled circuit with a bad signature", async () => {
        let msg = F.e(1234);
        msg = F.toObject(msg);
        const prvKey = generateRandomBigInt(253);
        const k = generateRandomBigInt(253);
        const pubKey = schnorr.prv2pub(prvKey);
        const signature = schnorr.signPoseidon(prvKey, msg, k);
        assert(schnorr.verifyPoseidon(signature, pubKey, msg));
        let yy = F.toObject(pubKey[1]);
        yy = Scalar.add(yy, "1");
        const w = await circuit.calculateWitness({
            "enabled": "0", //setting enabled to 0 make circuit dissabled
            "M": msg.toString(),
            "yx": F.toObject(pubKey[0]).toString(),
            "yy": yy.toString(),
            "S": signature.s.toString(),
            "e": signature.e.toString()

        }, true);
        await circuit.checkConstraints(w);
    });
});