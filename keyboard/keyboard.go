package keyboard

import (
	"fmt"
	"log"
	"math"
	prand "math/rand"
	"strings"
	"time"
)

const NChars = 35

var Chars = [NChars]byte{}

var reserved = [4][12]byte{
	{'E', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x'},
	{'B', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x'},
	{'C', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x', 'S'},
	{'T', 'x', 'A', 'M', 'X', 'x', 'x', 'H', 'L', 'D', 'U', 'R'},
}

var keyPrintingMap = map[byte]rune{
	'E': '⎋',
	'C': '⎈',
	'B': '←',
	'T': '↹',
	'A': '⎇',
	'M': '◆',
	'X': '⌘',
	'H': '⊞',
	'L': '←',
	'D': '↓',
	'U': '↑',
	'R': '→',
	'S': '⇧',
	// 'Y': '↩',
}

var absFinger = [4][12]int{
	{1, 1, 1, 2, 3, 3, 6, 6, 7, 8, 8, 8},
	{0, 0, 1, 2, 3, 3, 6, 6, 7, 8, 9, 9},
	{0, 0, 1, 2, 3, 3, 6, 6, 7, 8, 9, 9},
	{0, 4, 4, 4, 4, 4, 5, 5, 6, 7, 8, 9},
}

var handFinger = [4][12]int{
	{1, 1, 1, 2, 3, 3, 3, 3, 2, 1, 1, 1},
	{0, 0, 1, 2, 3, 3, 3, 3, 2, 1, 0, 0},
	{0, 0, 1, 2, 3, 3, 3, 3, 2, 1, 0, 0},
	{0, 4, 4, 4, 4, 4, 4, 4, 3, 2, 1, 0},
}

var handColumn = [4][12]int{
	{0, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1, 0},
	{0, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1, 0},
	{0, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1, 0},
	{0, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1, 0},
}

var hand = [4][12]int{
	{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1},
}

var effort = [4][12]int64{
	{7, 4, 1, 1, 4, 7, 5, 4, 1, 1, 3, 5},
	{3, 1, 0, 0, 0, 3, 3, 0, 0, 0, 1, 3},
	{5, 5, 5, 5, 2, 4, 4, 2, 4, 4, 4, 5},
	{7, 9, 9, 7, 1, 0, 0, 1, 0, 0, 0, 0},
}

var fingerPosition = [10]KeyPosition{
	{1, 1}, {1, 2}, {1, 3}, {1, 4}, {3, 5},
	{3, 6}, {1, 7}, {1, 8}, {1, 9}, {1, 10},
}

var targetFingerUsage = [10]float64{
	0.06, 0.10, 0.11, 0.12, 0.11,
	0.11, 0.12, 0.11, 0.10, 0.06,
}

const reset = "\x1b[0m"
const grey = "\x1b[1;30m"

var effortColor = map[int64]string{
	0: "\x1b[38;2;33;150;243m",
	1: "\x1b[38;2;53;143;226m",
	2: "\x1b[38;2;74;136;210m",
	3: "\x1b[38;2;94;129;194m",
	4: "\x1b[38;2;115;123;177m",
	5: "\x1b[38;2;136;116;161m",
	6: "\x1b[38;2;156;109;145m",
	7: "\x1b[38;2;177;103;128m",
	8: "\x1b[38;2;197;96;112m",
	9: "\x1b[38;2;218;89;96m",
}

type KeyPosition struct {
	i, j int
}

type Keyboard struct {
	Book              *string
	layout            [4][12]byte
	keyPositionLookup [128]KeyPosition
	handOverUse       int64
	repeatedPresses   int64
	repeatFinger1Gap  int64
	effort            int64
	inward            int64
	comfyInward       int64
	outward           int64
	comfyOutward      int64
	rowjump           int64
	distance          float64
	fingers           []float64
	hands             []float64
	fingerInequality  float64
	handInequality    float64
	Thread            int
	Gen               int
	Total             int
	Mutation          int
	Iteration         int
}

func New() *Keyboard {
	kb := &Keyboard{}
	kb.layout = [4][12]byte{}
	kb.keyPositionLookup = [128]KeyPosition{}
	kb.Fill(time.Now().Unix())
	return kb
}

func (kb *Keyboard) Fill(seed int64) {
	shars := Chars
	prand.Seed(seed)
	prand.Shuffle(len(Chars), func(i, j int) {
		shars[i], shars[j] = shars[j], shars[i]
	})
	for i, j := 0, 0; i < 48; i++ {
		p := i / 12
		q := i % 12
		if reserved[p][q] != 'x' {
			kb.layout[p][q] = reserved[p][q]
			continue
		}
		kb.layout[p][q] = shars[j]
		kb.keyPositionLookup[kb.layout[p][q]] = KeyPosition{p, q}
		j++
	}
}

func (kb *Keyboard) FillScore(distances *[4][12][4][12]float64) {
	keyZ := kb.keyPositionLookup[' ']
	keyA := kb.keyPositionLookup[' ']

	fingerUsage := [10]int64{}
	handUsage := [2]int64{}
	wasAnInroll := false
	wasAnOutroll := false

	for _, b := range []byte(*kb.Book) {
		keyB := kb.keyPositionLookup[b]
		ai, aj, bi, bj, zi, zj := keyA.i, keyA.j, keyB.i, keyB.j, keyZ.i, keyZ.j
		afa, afb, afz := absFinger[ai][aj], absFinger[bi][bj], absFinger[zi][zj]

		if afa == afb {
			kb.repeatedPresses++
		}

		if afz == afb {
			kb.repeatFinger1Gap++
		}

		// Add finger usage
		fingerUsage[afb]++

		// This is a same hand movement
		if hand[ai][aj] == hand[bi][bj] {
			if handFinger[ai][aj] < handFinger[bi][bj] && handColumn[ai][aj] < handColumn[bi][bj] {
				if effort[ai][aj] <= 2 && effort[bi][bj] <= 2 {
					kb.comfyInward++
				} else {
					kb.inward++
				}
				if (wasAnOutroll) && hand[bi][bj] == hand[ai][aj] {
					kb.handOverUse++
				}
				wasAnInroll = true
				wasAnOutroll = false
			} else if handFinger[ai][aj] > handFinger[bi][bj] && handColumn[ai][aj] > handColumn[bi][bj] {
				if effort[ai][aj] <= 2 && effort[bi][bj] <= 2 {
					kb.comfyOutward++
				} else {
					kb.outward++
				}
				if (wasAnInroll) && hand[bi][bj] == hand[ai][aj] {
					kb.handOverUse++
				}
				wasAnInroll = false
				wasAnOutroll = true
			} else if handFinger[ai][aj] != 4 && handFinger[bi][bj] != 4 && math.Abs(float64(bi)-float64(ai)) > 1.0 {
				if (wasAnInroll || wasAnOutroll) && hand[bi][bj] == hand[ai][aj] {
					kb.handOverUse++
				}
				kb.rowjump++
				wasAnInroll = false
				wasAnOutroll = false
			} else {
				if (wasAnInroll || wasAnOutroll) && hand[bi][bj] == hand[ai][aj] {
					kb.handOverUse++
				}
				wasAnInroll = false
				wasAnOutroll = false
			}
		}

		q := fingerPosition[afb]
		kb.distance += distances[bi][bj][q.i][q.j]

		// kb.distance += math.Sqrt(h*h + w*w)
		fingerPosition[afb] = KeyPosition{bi, bj}

		kb.effort += effort[bi][bj]

		keyZ, keyA = keyA, keyB
	}

	inequality := 0.0
	handInequality := 0.0
	kb.fingers = make([]float64, 10)
	for i, fu := range fingerUsage {
		usage := (float64(fu) / float64(len(*(kb.Book))))
		kb.fingers[i] = usage
		inequality += math.Abs(targetFingerUsage[i] - usage)
		if i <= 4 {
			handUsage[0] += fu
		} else {
			handUsage[1] += fu
		}
	}

	kb.hands = make([]float64, 2)
	for i, hu := range handUsage {
		usage := (float64(hu) / float64(len(*(kb.Book))))
		kb.hands[i] = usage
		handInequality += math.Abs(0.5 - usage)
	}

	kb.fingerInequality = inequality
	kb.handInequality = handInequality
}

func (kb *Keyboard) Score() float64 {
	return 100000 + ((kb.distance/4)+
		float64(kb.repeatedPresses*2)+
		float64(kb.repeatFinger1Gap)+
		float64(kb.effort)*0.04-
		float64(kb.comfyInward)-
		(float64(kb.inward)/4)-
		(float64(kb.comfyOutward)/2)-
		(float64(kb.outward)/8)+
		float64(kb.handOverUse/2))*
		(1+(kb.fingerInequality/4))
}

func (kb *Keyboard) Copy() *Keyboard {
	newLookup := [128]KeyPosition{}
	for k, v := range kb.keyPositionLookup {
		newLookup[k] = v
	}
	newLayout := [4][12]byte{}
	for i, v := range kb.layout {
		for j, w := range v {
			newLayout[i][j] = w
		}
	}
	return &Keyboard{
		layout:            newLayout,
		keyPositionLookup: newLookup,
	}
}

func (kb *Keyboard) Mutate() (s, t int) {
	a := prand.Int31n(NChars) // 16
	b := prand.Int31n(NChars) // 3

	ca := Chars[a]
	cb := Chars[b]

	kpa := kb.keyPositionLookup[ca] // {0, 3}
	kpb := kb.keyPositionLookup[cb] // {1, 3}

	kb.layout[kpa.i][kpa.j], kb.layout[kpb.i][kpb.j] = kb.layout[kpb.i][kpb.j], kb.layout[kpa.i][kpa.j]
	kb.keyPositionLookup[ca], kb.keyPositionLookup[cb] = kb.keyPositionLookup[cb], kb.keyPositionLookup[ca]

	return int(a), int(b)
}

func (kb *Keyboard) DetailString() string {
	return fmt.Sprintf(`
    Thread: %v  Gen: %v  Total: %v  Mutation: %v  Iter: %v

`, kb.Thread, kb.Gen, kb.Total, kb.Mutation, kb.Iteration)
}

func (kb *Keyboard) ScoreString() string {
	score := kb.Score()
	perc := func(v int64, mul float64) float64 {
		return float64(v) * mul / score * 100
	}
	return fmt.Sprintf(`
    Score:                  %.0f
    Repeated Finger 0 Gap: %5.1f%%     %v
    Repeated Finger 1 Gap: %5.1f%%     %v
    Comfortableness:       %5.1f%%     %v
    Comfy Inward Rolls:    %5.1f%%    %.0v
    Other Inward Rolls:    %5.1f%%    %.0v
    Comfy Outward Rolls:   %5.1f%%    %.0v
    Other Outward Rolls:   %5.1f%%    %.0v
    Hand Overuse:          %5.1f%%     %.0v
    Rowjumps:              %5.1f%%     %.0v
    Distance:              %5.1f%%     %.0f
    Hand Inequality:        %.3f     %.3f
    Finger Inequality:      %.3f     %.3f

    Left:     %4.1f               Right:   %4.1f
    %4.1f %4.1f %4.1f %4.1f %4.1f    %4.1f %4.1f %4.1f %4.1f %4.1f


`,
		score,
		perc(kb.repeatedPresses*2, 1),
		kb.repeatedPresses,
		perc(kb.repeatFinger1Gap, 1),
		kb.repeatFinger1Gap,
		perc(kb.effort, 0.04),
		kb.effort,
		perc(-kb.comfyInward, 1),
		-kb.comfyInward,
		perc(-kb.inward, 0.25),
		-kb.inward,
		perc(-kb.comfyOutward, 0.5),
		-kb.comfyOutward,
		perc(-kb.outward, 0.125),
		-kb.outward,
		perc(kb.handOverUse, 0.5),
		kb.handOverUse,
		perc(kb.rowjump, 0),
		kb.rowjump,
		kb.distance*0.25*100/score,
		kb.distance,
		kb.handInequality,
		kb.handInequality,
		kb.fingerInequality*0.25,
		kb.fingerInequality,
		kb.hands[0]*100,
		kb.hands[1]*100,
		kb.fingers[0]*100,
		kb.fingers[1]*100,
		kb.fingers[2]*100,
		kb.fingers[3]*100,
		kb.fingers[4]*100,
		kb.fingers[5]*100,
		kb.fingers[6]*100,
		kb.fingers[7]*100,
		kb.fingers[8]*100,
		kb.fingers[9]*100,
	)
}

func (kb *Keyboard) String() string {
	var str strings.Builder
	for i, row := range kb.layout {
		str.WriteString("    ")
		for j, ch := range row {
			switch reserved[i][j] {
			case 'x':
				color, ok := effortColor[effort[i][j]]
				if !ok {
					log.Println("unable to get effort color for", ch)
				} else {
					str.WriteString(color)
				}
				if ch == '\n' {
					str.WriteRune('↩')
				} else {
					str.WriteByte(ch)
				}
				if ok {
					str.WriteString(reset)
				}
			default:
				str.WriteString(grey)
				r, ok := keyPrintingMap[ch]
				if !ok {
					log.Println(ch, "not mapped")
					str.WriteByte(' ')
				} else {
					str.WriteRune(r)
				}
				str.WriteString(reset)
			}
			str.WriteString("  ")
		}
		str.WriteByte('\n')
	}
	return str.String()
}
