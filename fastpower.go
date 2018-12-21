package numthy

import "math/big"
// import "fmt"

var intZero = big.NewInt(0)
var intOne  = big.NewInt(1)
var intTwo = big.NewInt(2) 


// Modular Arithmetic Powering, Right-Left Binary 
// Algorithm 1.2.1 (Cohen, p. 8)
// 
// Returns g^n (mod m). The result is placed in y 
// (if not nil) and returned by the function.
// 
// The resulting integer y will satisfy 
//   0 <= z < m
//
// Note that 'Right-Left' in the algorithm name refers to 
// a scan of the bits of the exponent from least significant to 
// most significant. 
//
// TODO: Summary of algorithm.
//
// This is a non-optimized reference implementation.

func PowerMod_RLB_Ref(y, g, n, m *big.Int) *big.Int {
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
			MultMod(y, z, y, m)   // y <- z * y
		}
		N.Rsh(N, 1)  // N <- floor(N/2)
		if 0 == N.Sign() {
			return y
		}
		MultMod(z, z, z, m) // z <- z * z 
	}
}

// Version(s) of PowerMod_RLB algorithm to test 
// effects of various implementation details.
func PowerMod_RLB(y, g, n, m *big.Int) *big.Int {
	// 1. Initialize
	if nil == y {
		y = new(big.Int)
	}
	y0 := y // save y
	y.SetUint64(1)
	nSign := n.Sign()
	if 0 == nSign {
		return y
	}
	var z *big.Int
	if 0 > nSign {
		// N = new(big.Int).Set(n).Neg(n)
		// TODO: Write my own ModInverse
		z = new(big.Int).ModInverse(g, m)
	} else {
		// N = new(big.Int).Set(n)
		z = new(big.Int).Set(g)
	}

	nLen := n.BitLen()
	rr := new(big.Int)  // result of MultModAlt
	qq := new(big.Int)  // extra int for MultModAlt
	pp := new(big.Int)  // extra int for MultModAlt
	for i := 0; i < nLen; i++ {
		if (1 == n.Bit(i)) {
			MultModAlt(rr, qq, pp, z, y, m)
			// *y, *rr = *rr, *y
			y, rr = rr, y

		}
		MultModAlt(rr, qq, pp, z, z, m)
		// *z, *rr = *rr, *z
		z, rr = rr, z

	}

	*y0 = *y 
	return y0
}


// Modular Arithmetic Powering, Left-Right Binary
// Algorithm 1.2.2 (Cohen, p.9)
//
// Returns g^n (mod m). The result is placed in y 
// (if not nil) and returned by the function.
// 
// The resulting integer y will satisfy 
//   0 <= z < m
//
// TODO: Summary of algorithm.
//
// This is a non-optimized reference implementation.

func PowerMod_LRB_Ref(y, g, n, m *big.Int) *big.Int {
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
		z = g  					// Note: z is not mutated beyond this point
	}
	y.Set(z)
	e := N.BitLen()-1    // 2^e <= N < 2^(e+1)
	E := big.NewInt(1) 
	E.Lsh(E, uint(e))          // E = 2^e
	N.Sub(N, E)          // remove highest order bit from N

	for {
		if 0 == E.Cmp(intOne) {  // 2. [Finished?]
			return y
		}
		E.Rsh(E, 1)
		MultMod(y, y, y, m)
		if 0 <= N.Cmp(E) {  // N >= E   (3. [Multiply?])
			N.Sub(N, E) 
			MultMod(y, y, z, m)
		}
	}
}

// Modular Arithmetic Powering, Left-Right Binary, With Bit Operations
// Algorithm 1.2.3 (Cohen, p.9)
//
// Returns g^n (mod m). The result is placed in y 
// (if not nil) and returned by the function.
// 
// The resulting integer y will satisfy 
//   0 <= z < m
//
// TODO: Summary of algorithm.
//
// This is a non-optimized reference implementation.
func PowerMod_LRB_BitOp_Ref(y, g, n, m *big.Int) *big.Int {
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
		z = g  					// Note: z is not mutated beyond this point
	}
	y.Set(z)
	f := N.BitLen()-1    // 2^f <= N < 2^(f+1)
	for {
		if 0 == f {
			return y
		}
		f--
		MultMod(y, y, y, m)
		if 1 == N.Bit(f) {
			MultMod(y, y, z, m)			
		}
	}
}

// Version(s) of PowerMod_LRB_BitOp algorithm to test 
// effects of various implementation details.
func PowerMod_LRB_BitOp(y, g, n, m *big.Int) *big.Int {
	// 1. Initialize
	if nil == y {
		y = new(big.Int)
	}
	y0 := y // save y
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
		N = n
		z = g  					// Note: z is not mutated beyond this point
	}
	y.Set(z)
	f := N.BitLen()-1    // 2^f <= N < 2^(f+1)
	rr := new(big.Int)  // result of MultModAlt
	qq := new(big.Int)  // extra int for MultModAlt
	pp := new(big.Int)  // extra int for MultModAlt
	for {
		if 0 == f {
			*y0 = *y
			return y0
		}
		f--
		MultModAlt(rr, qq, pp, y, y, m)
		y, rr = rr, y

		if 1 == N.Bit(f) {
			MultModAlt(rr, qq, pp, y, z, m)
			y, rr = rr, y			
		}
	}
}

// Modular Arithmetic Powering, Left-Right Base 2^k
// Algorithm 1.2.4 (Cohen, p.10)
//
// Returns g^n (mod m). The result is placed in y 
// (if not nil) and returned by the function.
// 
// The resulting integer y will satisfy 
//   0 <= y < m
//
// TODO: Summary of algorithm.

var powerMod_LRW_windowSize = 4

func PowerMod_LRWindowed(y, g, n, m *big.Int) *big.Int {
// 1. Initialize
	if nil == y {
		y = new(big.Int)
	}
	y0 := y // save y
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
		N = n
		z = g  					// Note: z is not mutated beyond this point
	}
	k := uint(powerMod_LRW_windowSize) // For notation of window size
	e := (uint(N.BitLen()-1) / k)    // 2^kf <= N < 2^k(f+1) 
	f := e

	rr := new(big.Int)  // result of MultModAlt
	qq := new(big.Int)  // extra int for MultModAlt
	pp := new(big.Int)  // extra int for MultModAlt

	// Precompute odd powers of z  (powers[i] = z^(2i+1), for 0 <= i < 2^(k-1)
	powers := make([]*big.Int, 1<<(k-1))   
	powers[0] = z
	// fmt.Printf("\npowers[0] = %d\n", powers[0])
	MultModAlt(rr, qq, pp, z, z, m)
	zSqr := new(big.Int).Set(rr)
	for i := 1; i < 1<<(k-1); i++ {
		MultModAlt(rr, qq, pp, powers[i-1], zSqr, m)
		powers[i] = new(big.Int).Set(rr)
		// fmt.Printf("powers[%d] = %d\n", i, powers[i])
	}

	for {
		a := DigitInPowerOfTwoBase(k, f, N)
		if (0 == a) {
			for i := 0; i < int(k); i++ {
				MultModAlt(rr, qq, pp, y, y, m)
				y, rr = rr, y
			}
		} else {
			t, b := PowerOfTwoAndOddDivisor(a)
			if e != f {
				for i := 0; i < int(k-t); i++ {
					MultModAlt(rr, qq, pp, y, y, m)
					y, rr = rr, y
				}

				MultModAlt(rr, qq, pp, y, powers[(b-1)>>1], m)   // b = 2i+1 => i = (b-1)/2
				y, rr = rr, y
			} else { // e == f (first loop)
				y.Set(powers[(b-1)>>1])
			}

			for i := 0; i < int(t); i++ {
				MultModAlt(rr, qq, pp, y, y, m)
				y, rr = rr, y
			}
		}

		// 4. Finished?
		if 0 == f {
			*y0 = *y
			return y0
		}
		f--
	}
}

// Flexible Base 2^k Digits
// Sub-Algorithm 1.2.4.1 (Cohen, p.11/45)
//
// For positive integer N and k >= 1, this algorithm computes
// integers t_i and a_i, 0 <= i <= e such that 
//
//   N = 2^t_0 * (a_0 + 2^t_1 * (a_1 + ... + 2^t_e * a_e))
//  
// where t_i >= k for i >=1 and a_i is odd such that 1 <= a_i <= 2^k - 1
func FlexibleDigitsInPowerOfTwoBase(k uint, N *big.Int) ([]uint, []uint) {
	if N.Sign() <= 0 {
		panic("FlexibleDigitsInPowerOfTwoBase called with non-positive argument N")
	}
	if k <= 0 {
		panic("FlexibleDigitsInPowerOfTwoBase called with non-positive window size k")
	}
	capacity := N.BitLen() / int(k) + 1

	t := make([]uint, 0, capacity)
	a := make([]uint, 0, capacity)

	return t, a
}

// **************** Utilities *************************************

// Multiplies Integers x and y and 
// reduces the resulting value mod m.
// The result is placed in z (if not nil)
// and also returned by the function. 
//
// The resulting integer z will satisfy 
//   0 <= z < m

func MultMod(z, x, y, m *big.Int) *big.Int {
	z.Mul(x, y)
	z.Mod(z, m)
	return z
}

func MultModAlt(r, q, p, x, y, m *big.Int) *big.Int {
	p.Mul(x, y)
	q.DivMod(p, m, r)
	return r
}

func DigitInPowerOfTwoBase(k, f uint, N *big.Int) uint {
	var res uint

	for i := int(k*(f+1))-1 ; i >= int(k*f); i-- {
		res <<= 1
		res += N.Bit(i)
	} 
	// fmt.Printf("Digit(%d, %d, %d (%X)) = %d\n", k, f, N, N, res)

	return res
}

// Given a > 0, returns t >= 0 and odd b > 0 such that a = 2^t * b
func PowerOfTwoAndOddDivisor(a uint) (uint, uint) {
	b := a
	var t uint
	for t = 0; b > 0; t++ {
		// fmt.Printf("  a=%d, (t, b) = (%d, %d)\n", a, t, b)
		if (b & 1 == 1) {
			// fmt.Printf("PowerOfTwoAndOddDivisor(%d) = (%d, %d)\n", a, t, b)
			return t, b
		}
		b = b>>1
	}
	panic("PowerOfTwoAndOddDivisor did not terminate") 
}

