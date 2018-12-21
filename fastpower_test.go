package numthy

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"testing"
)

type funcPowerMod func(out, g, n, m *big.Int) *big.Int
type argPowerMod struct {
	g, n, M, expected string
}

func BuiltinExpZZ(out, g, n, m *big.Int) *big.Int { 
	return out.Exp(g, n, m) 
}

var powerModCases = []argPowerMod {
	argPowerMod{"3", "4", "25", "6"},
	argPowerMod{"2", "20", "7891", "6964"},
	argPowerMod{"0", "9483", "3493841", "0"},
	argPowerMod{"1", "9483", "3493841", "1"}, 
	argPowerMod{"13", "17534536", "17534537", "1"},
	argPowerMod{"12763457623", "17172534536", "99349283747534537", "60054740417897460"},
	argPowerMod{"283746589273648629384756298347663457623",
				"972987579238746834", 
				"639173479038475928374576289347685902374568234",
				 "372712510135067550268885892131308456226037873"},
	argPowerMod{"20398475902837495287438957693847569817689176982387495762893476582374658923746950287986982374658937409576398652834765",
				"92837490852790348759028374589762938238470560918690197639109457692736157384697861726894756198",
				"9876987283746598267187236487698726309581782746189237648912376487943659873649857618475693847561827364712637481568396418923",
				"4227851657401596239805556221161222815611188401865721801371457036594628468848826605639560135706311432453173172723268210052"},
}

func TestBuiltinExpZZ(t *testing.T) {
	for _, arg := range powerModCases {
		checkPowerMod(t, BuiltinExpZZ, arg, "TestBuiltinExpZZ")
	}
}

func TestPowerMod_RLB_Ref(t *testing.T) {
	for _, arg := range powerModCases {
		checkPowerMod(t, PowerMod_RLB_Ref, arg, "TestPowerMod_RLB_Ref")
	}
}

func TestPowerMod_RLB(t *testing.T) {
	for _, arg := range powerModCases {
		checkPowerMod(t, PowerMod_RLB, arg, "TestPowerMod_RLB")
	}
}

func TestPowerMod_LRB_Ref(t *testing.T) {
	for _, arg := range powerModCases {
		checkPowerMod(t, PowerMod_LRB_Ref, arg, "TestPowerMod_LRB_Ref")
	}
}

func TestPowerMod_LRB_BitOp_Ref(t *testing.T) {
	for _, arg := range powerModCases {
		checkPowerMod(t, PowerMod_LRB_BitOp_Ref, arg, "TestPowerMod_LRB_BitOp_Ref")
	}
}

func TestPowerMod_LRB_BitOp(t *testing.T) {
	for _, arg := range powerModCases {
		checkPowerMod(t, PowerMod_LRB_BitOp, arg, "TestPowerMod_LRB_BitOp")
	}
}

func TestPowerMod_LRWindowed(t *testing.T) {
	for _, arg := range powerModCases {
		checkPowerMod(t, PowerMod_LRWindowed, arg, "TestPowerMod_LRWindowed")
	}
}

func checkPowerMod(t *testing.T, f funcPowerMod, a argPowerMod, tag string) {
	g, ok1 := new(big.Int).SetString(a.g, 10)
	n, ok2 := new(big.Int).SetString(a.n, 10)
	M, ok3 := new(big.Int).SetString(a.M, 10)
	out := new(big.Int)
	expected, ok4 := new(big.Int).SetString(a.expected, 10)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		t.Errorf("Input Error")
	}

	f(out, g, n, M)
	if expected.Cmp(out) != 0 {
		t.Errorf("%v: Computing %d**%d (mod %d): expected %d, got %d", tag, g, n, M, expected, out)
	}
}

// ******************  Benchmarks ****************************


func BenchmarkPowerMod(b *testing.B) {
	testFuncs := []struct {
		name string
		fun funcPowerMod
	}{
		{"BuiltinExpZZ", BuiltinExpZZ},
		// {"PowerMod_RLB_Ref", PowerMod_RLB_Ref},
		{"PowerMod_RLB", PowerMod_RLB},
		// {"PowerMod_LRB_Ref", PowerMod_LRB_Ref},
		// {"PowerMod_LRB_BitOp_Ref", PowerMod_LRB_BitOp_Ref},
		{"PowerMod_LRB_BitOp", PowerMod_LRB_BitOp},
		{"PowerMod_LRWindowed", PowerMod_LRWindowed},

	}
	for _, testFunc := range testFuncs {
		for k := 6.; k <= 16; k++ {
			bitLength := uint(math.Pow(2, k))
			b.Run(fmt.Sprintf("%s/%d", testFunc.name, bitLength), func(b *testing.B) {
				// Always the same seed for reproducible results.
				rnd := rand.New(rand.NewSource(0))
				out := new(big.Int)
				// out2 := new(big.Int)
				g := new(big.Int)
				n := new(big.Int)
				M := new(big.Int)

				ub := new(big.Int) // extra arg for RandInt

				for i := 0; i < b.N; i++ {
					g = RandInt(g, ub, rnd, bitLength)
					n = RandInt(n, ub, rnd, 512)
					M = RandInt(M, ub, rnd, bitLength)
					M.Add(M, intTwo)
					testFunc.fun(out, g, n, M)
					// BuiltinExpZZ(out2, g, n, M)
					// if 0 != out.Cmp(out2) {
					// 	panic("mismatched results")
					// }
				}
			})
		}
	}
}

// ******************  Utilities ****************************

// Returns random integer in [0, 2^n)
func RandInt(out, upperBound *big.Int, rnd *rand.Rand, bitLength uint) *big.Int {
	upperBound.Lsh(intOne, bitLength)
	return out.Rand(rnd, upperBound)
}









