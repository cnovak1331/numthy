package numthy

import "math/big"

var intZero = big.NewInt(0)

// Modular Arithmetic Powering, Right-Left Binary 
// Algorithm 1.2.1 (Cohen, p. 8)
// 
// Returns g^n (mod m). The result is placed in t 
// (if not nil) and returned by the function.
// 
// The resulting integer y will satisfy 
//   0 <= z < m
func PowerMod_RLB(g, n, y, m *big.Int) *big.Int {
	// 1. Initialize
	if nil == y {
		y = new(big.Int)
	}
	y.SetUint64(1)
	nSign := n.Sign()
	if 0 == nSign {
		return y
	}
	var N, z *big.Int
	if 0 > nSign {
		N = new(big.Int).Set(n).Neg(n)
		// TODO: Write my own ModInverse
		z = new(big.Int).ModInverse(g, m)
	} else {
		N = new(big.Int).Set(n)
		z = new(big.Int).Set(g)
	}

	for {
		if (1 == N.Bit(0)) {  // if N is odd
			y = MultMod(z, y, y, m)   // y <- z * y
		}
		N.Rsh(N, 1)  // N <- floor(N/2)
		if 0 == N.Sign() {
			return y
		}
		z = MultMod(z, z, z, m) // z <- z * z 
	}
}

// Multiplies Integers x and y and 
// reduces the resulting value mod m.
// The result is placed in z (if not nil)
// and also returned by the function. 
//
// The resulting integer z will satisfy 
//   0 <= z < m
func MultMod(x, y, z, m *big.Int) *big.Int {
	// TODO: check m >= 2
	if nil == z {
		z = new(big.Int)
	}
	z.Mul(x, y).Mod(z, m)
	return z
}
