#!/usr/bin/env node

// Measuring arithmetics operations for ffiasm
// NOTE: if count < 10000000 then the result is not accurate.

var tmp = require("tmp-promise");
const path = require("path");
const util = require("util");
const { performance } = require("perf_hooks");
const exec = util.promisify(require("child_process").exec);
const fs = require('fs');
const os = require('os');
const bigInt = require("big-integer");
const buildZqField = require("ffiasm").buildZqField;

async function benchmarkMM(op, prime, count, x, y) {
    const dir = await tmp.dir({prefix: "circom_", unsafeCleanup: true });

    const source = await buildZqField(prime, "Fr");

    // console.log(dir.path);
    // Patch the generated file
    // Define the new line to be added
    let newLine = 'extern "C" bool Fr_init();\n';
    // Use regular expressions to find and modify the string
    source.hpp = source.hpp.replace(/(extern "C" void[^\n]+\n)/, newLine + '$1');

    await fs.promises.writeFile(path.join(dir.path, "fr.asm"), source.asm, "utf8");
    await fs.promises.writeFile(path.join(dir.path, "fr.hpp"), source.hpp, "utf8");
    await fs.promises.writeFile(path.join(dir.path, "fr.cpp"), source.cpp, "utf8");

    await exec(`cp  ${path.join(path.join(__dirname, "..", "src"),  `${op}.cpp`)} ${dir.path}`);

    if (process.platform === "darwin") {
        await exec("nasm -fmacho64 --prefix _ " +
            ` ${path.join(dir.path,  "fr.asm")}`
        );
    }  else if (process.platform === "linux") {
        await exec("nasm -felf64 " +
            ` ${path.join(dir.path,  "fr.asm")}`
        );
    } else throw("Unsupported platform");

    await exec("g++" +
       ` ${path.join(dir.path,  `${op}.cpp`)}` +
       ` ${path.join(dir.path,  "fr.o")}` +
       ` ${path.join(dir.path,  "fr.cpp")}` +
       ` -o ${path.join(dir.path, "benchmark")}` +
       " -lgmp -O3"
    );

    // ignore x and y for now
    let result = await exec(`${path.join(dir.path,  "benchmark")} ${count}`);

    if (result.stdout === '') {
          throw new Error("benchmark stdout is empty");
    }

    if (result.stderr !== '') {
          throw new Error("benchmark stderr is not empty");
    }


    return result.stdout;
}

async function run () {
    var result_string = "";
    // Read Arguments
    // The first two arguments are node and app.js
    if (process.argv.length != 8) {
        throw new Error(`Please provide all arguments: curve, field, operation, count, input, result`);
    }
    console.log("Process Arithmetics: " + process.argv[2] + " " + process.argv[3] + " " + process.argv[4] + " " + process.argv[5] + " " + process.argv[6] + " " + process.argv[7]);
    // Curve of which we should measure the native or scalar field operations
    const curve_name = process.argv[2];
    if (curve_name != "bn128") {
        throw new Error(`Only bn128 curve is currently supported`);
    }
    // Field scalar or base
    const field_name = process.argv[3];
    var field;
    // Those fields are for bn254
    if (field_name == "base") {
        field = bigInt("21888242871839275222246405745257275088696311157297823662689037894645226208583"); // Fq
    } else if (field_name == "scalar") {
        field = bigInt("21888242871839275222246405745257275088548364400416034343698204186575808495617"); // Fr
    } else {
        throw new Error(`Field should be base or scalar`);
    }
    // Operation to exute
    const operation = process.argv[4];
    if (!["add", "mul"].includes(operation)) {
        throw new Error(`Field should be: add or mul`);
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
    const y = parseInt(input['y'], 10);
    // Result csv file to write the result
    const path_name = process.argv[7];
    const dir_name = path.dirname(path_name);
    if (!fs.existsSync(dir_name)) {
        throw new Error(`Directory ${dir_name} does not exist`);
    }
    if (!fs.existsSync(path_name)) {
        result_string = "framework,category,curve,field,operation,input,ram,time,nbPhysicalCores,nbLogicalCores,count,cpu\n";
    }

    // Execute benchmark
    let elapsed = await benchmarkMM(operation, field, count, x, y);


    // Detect peripheral info
    const ram = process.memoryUsage().rss;
    const machine = os.cpus()[0].model;

    const input_path = input_file.substring(input_file.indexOf("input_file"));
    // FIXME cores and threads
    result_string += "rapidsnark,arithmetic," + curve_name + "," + field_name + "," + operation + "," + input_path + "," + ram + "," + elapsed + ",1,1," + count + "," + machine + "\n";

    fs.appendFileSync(path_name, result_string, function(err) {
        if(err) {
            return console.log(err);
        }
    });
}

run().then(() => {
    process.exit(0);
});
