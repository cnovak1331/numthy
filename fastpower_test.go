package numthy

import (
	"math/big"
	"testing"
)

type funcPowerMod func(g, n, out, m *big.Int) *big.Int
type argPowerMod struct {
	g, n, M, expected string
}

func TestPowerMod_RLB(t *testing.T) {
	arg1 := argPowerMod{"2", "20", "7891", "6964"}
	checkPowerMod(t, PowerMod_RLB, arg1)

}

func checkPowerMod(t *testing.T, f funcPowerMod, a argPowerMod) {
	g, ok1 := new(big.Int).SetString(a.g, 10)
	n, ok2 := new(big.Int).SetString(a.n, 10)
	M, ok3 := new(big.Int).SetString(a.M, 10)
	out := new(big.Int)
	expected, ok4 := new(big.Int).SetString(a.expected, 10)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		t.Errorf("Input Error")
	}

	f(g, n, out, M)
	if expected.Cmp(out) != 0 {
		t.Errorf("Computing %d**%d (mod %d): expected %d, got %d", g, n, M, expected, out)
	}
}
