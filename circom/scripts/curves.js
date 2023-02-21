#!/usr/bin/env node

// Measuring elliptic curve operations for ffjavascript

const fs = require('fs');
const os = require('os');
const path = require('path');
const ffjs = require('ffjavascript');
const buildBn128 = ffjs.buildBn128;
const buildBls12381 = ffjs.buildBls12381;
const BigBuffer = ffjs.BigBuffer;
const F1Field = ffjs.F1Field;

async function getCurve(curve_name, singleThread) {
    if (curve_name == "bn128") {
        curve = await buildBn128(singleThread);
    } else if (curve_name == "bls12_381") {
        curve = await buildBls12381(singleThread);
    } else {
        throw new Error(`Curve not supported: ${curve}`);
    }
    return curve;
}

// Note that if we try to move the measuring in another function by passing
// field.add (or any other operation) then we get the following error:
// ffjavascript/build/main.cjs:1534
//         return res >= this.p ? res-this.p : res;
async function benchmark(curve, G, operation, x, y, count) {
    var start;
    switch (operation) {
        case "scalar-multiplication":
            start = process.hrtime();
            for (let step = 0; step < count; step++) {
                var r = G.timesScalar(G.g, x)
            }
            return process.hrtime(start)[1] / count / 1024; // milli seconds
        case "multi-scalar-multiplication":
            start = process.hrtime();
            const N = x;
            const scalars = new BigBuffer(N*curve.Fr.n8);
            const bases = new BigBuffer(N*G.F.n8*2);
            for (let step = 0; step < count; step++) {
                var r = await G.multiExpAffine(bases, scalars, false, "");
            }
            return process.hrtime(start)[1] / count / 1024; // milli seconds
        case "pairing":
            start = process.hrtime();
            const g1 = curve.G1.timesScalar(curve.G1.g, x);
            const g2 = curve.G2.timesScalar(curve.G2.g, y);
            for (let step = 0; step < count; step++) {
                const pre1 = curve.prepareG1(g1);
                const pre2 = curve.prepareG2(g2);
                const r1 = curve.millerLoop(pre1, pre2);
                const r2 = curve.finalExponentiation(r1);
            }
            return process.hrtime(start)[1] / count / 1024; // milli seconds
        default:
            throw new Error(`Operation not supported: ${operation}`);
    }
}


async function run () {
    const singleThread = true;
    var result_string = "";
    // Read Arguments
    // The first two arguments are node and app.js
    if (process.argv.length != 8) {
        throw new Error(`Please provide all arguments: curve, g, operation, count, input, result`);
    }
    console.log("Process Arithmetics: " + process.argv[2] + " " + process.argv[3] + " " + process.argv[4] + " " + process.argv[5] + " " + process.argv[6] + " " + process.argv[7]);
    // Curve of which we should measure the native or scalar field operations
    const curve_name = process.argv[2];
    const curve = await getCurve(curve_name, singleThread);
    // Group
    // For pairing it does not matter if it is G1 or G2.
    const group_name = process.argv[3];
    var G;
    if (group_name == "g1") {
        G = curve.G1;
    } else if (group_name == "g2") {
        G = curve.G2;
    } else {
        throw new Error(`G should be: g1 or g2`);
    }
    // Operation to exute
    var operation = process.argv[4];
    if (!["scalar-multiplication", "multi-scalar-multiplication", "pairing"].includes(operation)) {
        throw new Error(`Field should be: scalar-multiplication, multi-scalar-multiplication, or pairing"`);
    }
    // Counter of how many times to run the operation
    const count = parseInt(process.argv[5], 10);
    // JSON Input file
    const input_file = process.argv[6];
    const input = JSON.parse(fs.readFileSync(input_file, 'utf8'));
    if (!input.hasOwnProperty("x")) {
        throw new Error(`Input x is missing from the provided input.`);
    }
    const x = parseInt(input['x'], 10);
    var y;
    if (operation == "pairing") {
        if (!input.hasOwnProperty("y")) {
            throw new Error(`Input y is missing from the provided input.`);
        }
        y = parseInt(input['y'], 10);
    } else {
        // Give a dummy value since it's not going to be used.
        y = 0;
    }
    // Result csv file to write the result
    const path_name = process.argv[7];
    const dir_name = path.dirname(path_name);
    if (!fs.existsSync(dir_name)) {
        throw new Error(`Directory ${dir_name} does not exist`);
    }
    if (!fs.existsSync(path_name)) {
        result_string = "framework,category,curve,operation,input,ram,time,nbPhysicalCores,nbLogicalCores,cpu\n";
    }

    // Execute benchmark
    elapsed = Math.floor(await benchmark(curve, G, operation, x, y, count));

    // Detect peripheral info
    const ram = process.memoryUsage().heapUsed;
    const machine = os.cpus()[0].model;

    // Prepend g1 or g2
    if (["scalar-multiplication", "multi-scalar-multiplication"].includes(operation)) {
        if (group_name == "g1") {
            operation = "g1-" + operation;
        } else {
            operation = "g2-" + operation;
        }
    }

    result_string += "circom,ec," + curve_name + "," + operation + "," + path_name + "," + ram + "," + elapsed + ",1,1," + machine + "\n";

    fs.appendFileSync(path_name, result_string, function(err) {
        if(err) {
            return console.log(err);
        }
    }); 
}

run().then(() => {
    process.exit(0);
});
