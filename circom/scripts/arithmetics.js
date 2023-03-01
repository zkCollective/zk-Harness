#!/usr/bin/env node

// Measuring arithmetics operations for ffjavascript

const fs = require('fs');
const os = require('os');
const path = require('path');
const ffjs = require('ffjavascript');
const buildBn128 = ffjs.buildBn128;
const buildBls12381 = ffjs.buildBls12381;
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
function benchmark(field, operation, x, y, count) {
    x = field.e(x);
    y = field.e(y);
    var start;
    var time;
    var hrTime;
    switch (operation) {
        case "add":
            start = process.hrtime();
            for (let step = 0; step < count; step++) {
                field.add(x, y);
            }
            hrTime = process.hrtime(start)
            time = hrTime[0] * 1000000000 + hrTime[1];
            return time / count; // nano seconds
        case "sub":
            start = process.hrtime();
            for (let step = 0; step < count; step++) {
                field.sub(x, y);
            }
            hrTime = process.hrtime(start)
            time = hrTime[0] * 1000000000 + hrTime[1];
            return time / count; // nano seconds
        case "mul":
            start = process.hrtime();
            for (let step = 0; step < count; step++) {
                field.mul(x, y);
            }
            hrTime = process.hrtime(start)
            time = hrTime[0] * 1000000000 + hrTime[1];
            return time / count; // nano seconds
        case "inv":
            start = process.hrtime();
            for (let step = 0; step < count; step++) {
                field.inv(x);
            }
            hrTime = process.hrtime(start)
            time = hrTime[0] * 1000000000 + hrTime[1];
            return time / count; // nano seconds
        case "exp":
            start = process.hrtime();
            for (let step = 0; step < count; step++) {
                field.exp(x, y);
            }
            hrTime = process.hrtime(start)
            time = hrTime[0] * 1000000000 + hrTime[1];
            return time / count; // nano seconds
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
        throw new Error(`Please provide all arguments: curve, field, operation, count, input, result`);
    }
    console.log("Process Arithmetics: " + process.argv[2] + " " + process.argv[3] + " " + process.argv[4] + " " + process.argv[5] + " " + process.argv[6] + " " + process.argv[7]);
    // Curve of which we should measure the native or scalar field operations
    const curve_name = process.argv[2];
    const curve = await getCurve(curve_name, singleThread);
    // Field scalr or base 
    const field_name = process.argv[3];
    var field;
    if (field_name == "base") {
        field = new F1Field(curve.q);
    } else if (field_name == "scalar") {
        field = new F1Field(curve.r);
    } else {
        throw new Error(`Field should be base or scalar`);
    }
    // Operation to exute
    const operation = process.argv[4];
    if (!["add", "sub", "mul", "inv", "exp"].includes(operation)) {
        throw new Error(`Field should be: add, sub, mul, inv, or exp`);
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
    if (operation != "inv") {
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
        result_string = "framework,category,curve,field,operation,input,ram,time,nbPhysicalCores,nbLogicalCores,cpu\n";
    }

    // Execute benchmark
    elapsed = Math.floor(benchmark(field, operation, x, y, count));

    // Detect peripheral info
    const ram = process.memoryUsage().rss;
    const machine = os.cpus()[0].model;

    const input_path = input_file.substring(input_file.indexOf("input_file"));
    result_string += "snarkjs,arithmetic," + curve_name + "," + field_name + "," + operation + "," + input_path + "," + ram + "," + elapsed + ",1,1," + machine + "\n";

    fs.appendFileSync(path_name, result_string, function(err) {
        if(err) {
            return console.log(err);
        }
    }); 
}

run().then(() => {
    process.exit(0);
});
