// copied from https://github.com/RustCrypto/hashes/blob/master/sha2/src/consts.rs

// initial values for the digest limbs as big-endian integers
pub const HASH_IV: [u32; 8] = [
    0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
];

/// Constants necessary for SHA-256 family of digests.
pub const ROUND_CONSTANTS: [u32; 64] = [
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
    0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
    0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
    0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
    0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
    0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
];

/// Constants necessary for SHA-256 family of digests.
pub const ROUND_CONSTANTS_X4: [[u32; 4]; 16] = [
    [
        ROUND_CONSTANTS[3],
        ROUND_CONSTANTS[2],
        ROUND_CONSTANTS[1],
        ROUND_CONSTANTS[0],
    ],
    [
        ROUND_CONSTANTS[7],
        ROUND_CONSTANTS[6],
        ROUND_CONSTANTS[5],
        ROUND_CONSTANTS[4],
    ],
    [
        ROUND_CONSTANTS[11],
        ROUND_CONSTANTS[10],
        ROUND_CONSTANTS[9],
        ROUND_CONSTANTS[8],
    ],
    [
        ROUND_CONSTANTS[15],
        ROUND_CONSTANTS[14],
        ROUND_CONSTANTS[13],
        ROUND_CONSTANTS[12],
    ],
    [
        ROUND_CONSTANTS[19],
        ROUND_CONSTANTS[18],
        ROUND_CONSTANTS[17],
        ROUND_CONSTANTS[16],
    ],
    [
        ROUND_CONSTANTS[23],
        ROUND_CONSTANTS[22],
        ROUND_CONSTANTS[21],
        ROUND_CONSTANTS[20],
    ],
    [
        ROUND_CONSTANTS[27],
        ROUND_CONSTANTS[26],
        ROUND_CONSTANTS[25],
        ROUND_CONSTANTS[24],
    ],
    [
        ROUND_CONSTANTS[31],
        ROUND_CONSTANTS[30],
        ROUND_CONSTANTS[29],
        ROUND_CONSTANTS[28],
    ],
    [
        ROUND_CONSTANTS[35],
        ROUND_CONSTANTS[34],
        ROUND_CONSTANTS[33],
        ROUND_CONSTANTS[32],
    ],
    [
        ROUND_CONSTANTS[39],
        ROUND_CONSTANTS[38],
        ROUND_CONSTANTS[37],
        ROUND_CONSTANTS[36],
    ],
    [
        ROUND_CONSTANTS[43],
        ROUND_CONSTANTS[42],
        ROUND_CONSTANTS[41],
        ROUND_CONSTANTS[40],
    ],
    [
        ROUND_CONSTANTS[47],
        ROUND_CONSTANTS[46],
        ROUND_CONSTANTS[45],
        ROUND_CONSTANTS[44],
    ],
    [
        ROUND_CONSTANTS[51],
        ROUND_CONSTANTS[50],
        ROUND_CONSTANTS[49],
        ROUND_CONSTANTS[48],
    ],
    [
        ROUND_CONSTANTS[55],
        ROUND_CONSTANTS[54],
        ROUND_CONSTANTS[53],
        ROUND_CONSTANTS[52],
    ],
    [
        ROUND_CONSTANTS[59],
        ROUND_CONSTANTS[58],
        ROUND_CONSTANTS[57],
        ROUND_CONSTANTS[56],
    ],
    [
        ROUND_CONSTANTS[63],
        ROUND_CONSTANTS[62],
        ROUND_CONSTANTS[61],
        ROUND_CONSTANTS[60],
    ],
];

pub const K32X4: [[u32; 4]; 16] = [
    [
        ROUND_CONSTANTS[3],
        ROUND_CONSTANTS[2],
        ROUND_CONSTANTS[1],
        ROUND_CONSTANTS[0],
    ],
    [
        ROUND_CONSTANTS[7],
        ROUND_CONSTANTS[6],
        ROUND_CONSTANTS[5],
        ROUND_CONSTANTS[4],
    ],
    [
        ROUND_CONSTANTS[11],
        ROUND_CONSTANTS[10],
        ROUND_CONSTANTS[9],
        ROUND_CONSTANTS[8],
    ],
    [
        ROUND_CONSTANTS[15],
        ROUND_CONSTANTS[14],
        ROUND_CONSTANTS[13],
        ROUND_CONSTANTS[12],
    ],
    [
        ROUND_CONSTANTS[19],
        ROUND_CONSTANTS[18],
        ROUND_CONSTANTS[17],
        ROUND_CONSTANTS[16],
    ],
    [
        ROUND_CONSTANTS[23],
        ROUND_CONSTANTS[22],
        ROUND_CONSTANTS[21],
        ROUND_CONSTANTS[20],
    ],
    [
        ROUND_CONSTANTS[27],
        ROUND_CONSTANTS[26],
        ROUND_CONSTANTS[25],
        ROUND_CONSTANTS[24],
    ],
    [
        ROUND_CONSTANTS[31],
        ROUND_CONSTANTS[30],
        ROUND_CONSTANTS[29],
        ROUND_CONSTANTS[28],
    ],
    [
        ROUND_CONSTANTS[35],
        ROUND_CONSTANTS[34],
        ROUND_CONSTANTS[33],
        ROUND_CONSTANTS[32],
    ],
    [
        ROUND_CONSTANTS[39],
        ROUND_CONSTANTS[38],
        ROUND_CONSTANTS[37],
        ROUND_CONSTANTS[36],
    ],
    [
        ROUND_CONSTANTS[43],
        ROUND_CONSTANTS[42],
        ROUND_CONSTANTS[41],
        ROUND_CONSTANTS[40],
    ],
    [
        ROUND_CONSTANTS[47],
        ROUND_CONSTANTS[46],
        ROUND_CONSTANTS[45],
        ROUND_CONSTANTS[44],
    ],
    [
        ROUND_CONSTANTS[51],
        ROUND_CONSTANTS[50],
        ROUND_CONSTANTS[49],
        ROUND_CONSTANTS[48],
    ],
    [
        ROUND_CONSTANTS[55],
        ROUND_CONSTANTS[54],
        ROUND_CONSTANTS[53],
        ROUND_CONSTANTS[52],
    ],
    [
        ROUND_CONSTANTS[59],
        ROUND_CONSTANTS[58],
        ROUND_CONSTANTS[57],
        ROUND_CONSTANTS[56],
    ],
    [
        ROUND_CONSTANTS[63],
        ROUND_CONSTANTS[62],
        ROUND_CONSTANTS[61],
        ROUND_CONSTANTS[60],
    ],
];
