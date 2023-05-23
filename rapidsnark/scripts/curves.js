#!/usr/bin/env node

// Measuring arithmetics operations for ffiasm

var tmp = require("tmp-promise");
const path = require("path");
const util = require("util");
const exec = util.promisify(require("child_process").exec);
const fs = require('fs');
const os = require('os');
const bigInt = require("big-integer");
const buildZqField = require("ffiasm").buildZqField;


async function generate_field(dir, prime, name) {
    const source = await buildZqField(prime, name);

    // Patch the generated file
    // Define the new line to be added
    let newLine = 'extern "C" bool ' + name + '_init();\n';
    // Use regular expressions to find and modify the string
    source.hpp = source.hpp.replace(/(extern "C" void[^\n]+\n)/, newLine + '$1');

    let lower_name = name.toLowerCase();

    await fs.promises.writeFile(path.join(dir.path, lower_name + ".asm"), source.asm, "utf8");
    await fs.promises.writeFile(path.join(dir.path, lower_name + ".hpp"), source.hpp, "utf8");
    await fs.promises.writeFile(path.join(dir.path, lower_name + ".cpp"), source.cpp, "utf8");

    if (process.platform === "darwin") {
        await exec("nasm -fmacho64 --prefix _ " +
            ` ${path.join(dir.path, lower_name + ".asm")}`
        );
    }  else if (process.platform === "linux") {
        await exec("nasm -felf64 " +
            ` ${path.join(dir.path, lower_name + ".asm")}`
        );
    } else throw("Unsupported platform");
}

async function benchmarkMM(g_file, count, x) {
    const dir = await tmp.dir({prefix: "circom_", unsafeCleanup: true });

    // Those are hardcoded for bn254 rn
    await generate_field(dir, bigInt("21888242871839275222246405745257275088548364400416034343698204186575808495617"), "Fr");
    await generate_field(dir, bigInt("21888242871839275222246405745257275088696311157297823662689037894645226208583"), "Fq");

    await exec(`cp  ${path.join(path.join(__dirname, "..", "src"),  `${g_file}`)} ${dir.path}`);

    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "alt_bn128.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "alt_bn128.cpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "f2field.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "f2field.cpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "splitparstr.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "splitparstr.cpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "curve.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "curve.cpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "exp.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "naf.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "naf.cpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "multiexp.cpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "multiexp.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "misc.hpp")} ${dir.path}`);
    await exec(`cp  ${path.join(__dirname, "..", "..", "node_modules", "ffiasm", "c", "misc.cpp")} ${dir.path}`);

    await exec("g++" +
       ` -I.${dir.path}` +
       ` ${path.join(dir.path,  `${g_file}`)}` +
       ` ${path.join(dir.path,  "alt_bn128.cpp")}` +
       ` ${path.join(dir.path,  "splitparstr.cpp")}` +
       ` ${path.join(dir.path,  "misc.cpp")}` +
       ` ${path.join(dir.path,  "naf.cpp")}` +
       ` ${path.join(dir.path,  "fr.o")}` +
       ` ${path.join(dir.path,  "fr.cpp")}` +
       ` ${path.join(dir.path,  "fq.o")}` +
       ` ${path.join(dir.path,  "fq.cpp")}` +
       ` -o ${path.join(dir.path, "benchmark")}` +
       " -lgmp -fopenmp -O3"
    );

    let result = await exec(`${path.join(dir.path,  "benchmark")} ${x} ${count}`);

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
        throw new Error(`Please provide all arguments: curve, g, operation, count, input, result`);
    }
    console.log("Process Arithmetics: " + process.argv[2] + " " + process.argv[3] + " " + process.argv[4] + " " + process.argv[5] + " " + process.argv[6] + " " + process.argv[7]);
    // Curve of which we should measure the native or scalar field operations
    const curve_name = process.argv[2];
    if (curve_name != "bn128") {
        throw new Error(`Only bn128 curve is currently supported`);
    }
    // Group
    const group_name = process.argv[3];
    var g_file;
    if (group_name == "g1") {
        g_file = "multiexp_g1.cpp";
    } else if (group_name == "g2") {
        g_file = "multiexp_g2.cpp";
    } else {
        throw new Error(`G should be: g1 or g2`);
    }
    // Operation to exute
    let operation = process.argv[4];
    if (!["multi-scalar-multiplication"].includes(operation)) {
        throw new Error(`Field should be: multi-scalar-multiplication"`);
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
    // Result csv file to write the result
    const path_name = process.argv[7];
    const dir_name = path.dirname(path_name);
    if (!fs.existsSync(dir_name)) {
        throw new Error(`Directory ${dir_name} does not exist`);
    }
    if (!fs.existsSync(path_name)) {
        result_string = "framework,category,curve,operation,input,ram,time,nbPhysicalCores,nbLogicalCores,count,cpu\n";
    }

    // Execute benchmark
    let elapsed = await benchmarkMM(g_file, count, x);


    // Detect peripheral info
    const ram = process.memoryUsage().rss;
    const machine = os.cpus()[0].model;

    // Prepend g1 or g2
    if (["scalar-multiplication", "multi-scalar-multiplication"].includes(operation)) {
        if (group_name == "g1") {
            operation = "g1-" + operation;
        } else {
            operation = "g2-" + operation;
        }
    }

    const input_path = input_file.substring(input_file.indexOf("input_file"));
    // FIXME cores and threads
    result_string += "rapidsnark,ec," + curve_name + "," + operation + "," + input_path + "," + ram + "," + elapsed + ",1,1," + count + "," + machine + "\n";

    fs.appendFileSync(path_name, result_string, function(err) {
        if(err) {
            return console.log(err);
        }
    });
}

run().then(() => {
    process.exit(0);
});
