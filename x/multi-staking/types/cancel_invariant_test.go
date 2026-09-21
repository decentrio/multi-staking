package types_test

import (
	"testing"

	"cosmossdk.io/math"
)

// invariant: after cancelling, unlocked-at-maturity must not exceed the remaining unlock amount.
func TestCancelRoundingInvariant(t *testing.T) {
	badOld, badNew, total := 0, 0, 0
	for _, ws := range []string{"0.1", "0.25", "0.3", "0.4", "0.5", "0.6", "0.75", "0.9", "0.333333333333333333", "1", "1.5", "2.7"} {
		w := math.LegacyMustNewDecFromStr(ws)
		for U := int64(1); U <= 120; U++ {
			B := w.MulInt64(U).TruncateInt() // staking entry after undelegate
			for x := int64(1); x <= U; x++ {
				c := w.MulInt64(x).TruncateInt() // bond cancelled
				if !c.IsPositive() {
					continue
				}
				total++
				rem := B.Sub(c)
				unlocked := math.LegacyNewDecFromInt(rem).Quo(w).TruncateInt()
				// old: unlock entry reduced by x
				if unlocked.GT(math.NewInt(U - x)) {
					badOld++
				}
				// new: unlock entry reduced by y = trunc(c / w)
				y := math.LegacyNewDecFromInt(c).Quo(w).TruncateInt()
				if y.IsPositive() && unlocked.GT(math.NewInt(U).Sub(y)) {
					badNew++
				}
			}
		}
	}
	t.Logf("cases=%d violations old=%d new=%d", total, badOld, badNew)
	if badNew != 0 {
		t.Fatalf("new formula still violates invariant")
	}
}
