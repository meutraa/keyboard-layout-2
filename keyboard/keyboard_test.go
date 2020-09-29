package keyboard

import (
	"testing"
)

func NewTestKeyboard() *Keyboard {
	kb := &Keyboard{}
	kb.layout = [4][12]byte{}
	kb.keyPositionLookup = [128]KeyPosition{}
	kb.Fill(0)
	return kb
}

/*func BenchmarkDistance(b *testing.B) {
	x, y := rand.Int63n(4), rand.Int63n(12)
	var a float64
	for i := 0; i < b.N; i++ {
		x <<= 4
		y <<= 4
		a += math.Sqrt(float64(x*x) + float64(y*y))
	}
}*/

func TestMutate(t *testing.T) {
	kb := NewTestKeyboard()

	for i := 0; i < 0; i++ {
		_testMutate(t, kb)
	}
}

/*func BenchmarkMutate(b *testing.B) {
	kb := NewTestKeyboard()

	for i := 0; i < b.N; i++ {
		kb.Mutate()
	}
}*/

func _testMutate(t *testing.T, kb *Keyboard) {
	// Do the mutation
	okb := kb.Copy()
	x, y := kb.Mutate()
	t.Logf("\n%v", okb)
	t.Logf("\n%v", kb)

	la := Chars[x] // b
	lb := Chars[y] // u

	ra := okb.keyPositionLookup[la] // b (2, 2)
	rb := okb.keyPositionLookup[lb] // u (3, 1)

	t.Logf("%c (%d, %d) <-> %c (%d, %d)", la, ra.i, ra.j, lb, rb.i, rb.j)

	p := okb.layout[ra.i][ra.j] // u
	q := okb.layout[rb.i][rb.j] // b
	t.Log(p, la)
	if p != la {
		t.Errorf("key %c before mutate meant to be at (%d, %d) but %c is", la, ra.i, ra.j, p)
	}
	if q != lb {
		t.Errorf("key %c before mutate meant to be at (%d, %d) but %c is", lb, rb.i, rb.j, q)
	}

	u := okb.keyPositionLookup[la]
	b := okb.keyPositionLookup[lb]
	if u.i != ra.i || u.j != ra.j {
		t.Errorf("key %c before mutate meant to be at (%d, %d) but at (%d, %d)", la, ra.i, ra.j, u.i, u.j)
	}
	if b.i != rb.i || b.j != rb.j {
		t.Errorf("key %c before mutate meant to be at (%d, %d) but at (%d, %d)", lb, rb.i, rb.j, b.i, b.j)
	}

	p = kb.layout[ra.i][ra.j]
	q = kb.layout[rb.i][rb.j]
	// Check the results
	if p != lb {
		t.Errorf("key %c after mutate meant to be at (%d, %d) but %c is", lb, ra.i, ra.j, p)
	}
	if q != la {
		t.Errorf("key %c after mutate meant to be at (%d, %d) but %c is", la, rb.i, rb.j, q)
	}

	u = kb.keyPositionLookup[la]
	b = kb.keyPositionLookup[lb]
	if b.i != ra.i || b.j != ra.j {
		t.Errorf("key %c after mutate meant to be at (%d, %d) but at (%d, %d)", la, ra.i, ra.j, b.i, b.j)
	}
	if u.i != rb.i || u.j != rb.j {
		t.Errorf("key %c after mutate meant to be at (%d, %d) but at (%d, %d)", lb, rb.i, rb.j, u.i, u.j)
	}
}
