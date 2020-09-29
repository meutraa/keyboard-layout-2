package main

import (
	"log"
	"math"
	"testing"

	"gitlab.com/meutraa/keyboard-gen/keyboard"
)

func BenchmarkFillScore(b *testing.B) {
	kb := keyboard.New()
	kb.Fill(0)
	data, err := Parse()
	if nil != err {
		log.Fatalln("unable to parse data", err)
		return
	}

	book := createBook(data, 10000, false)
	kb.Book = &book

	distances := [4][12][4][12]float64{}
	for a, va := range distances {
		for b, vb := range va {
			for c, vc := range vb {
				for d, _ := range vc {
					h := math.Abs(float64(c - a))
					w := math.Abs(float64(d - b))
					distances[a][b][c][d] = math.Sqrt(h*h + w*w)
				}
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kb.FillScore(&distances)
	}
}
