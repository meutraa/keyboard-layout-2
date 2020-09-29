package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	prand "math/rand"
	"runtime"
	"sort"
	"strings"

	"gitlab.com/meutraa/keyboard-gen/keyboard"
)

func binarySearch(a []Word, search float64, last Word) (result *Word, searchCount int) {
	mid := len(a) / 2
	switch {
	case len(a) == 0:
		result = &last
	case a[mid].cumulativePercentage > search:
		result, searchCount = binarySearch(a[:mid], search, a[mid])
	case a[mid].cumulativePercentage < search:
		result, searchCount = binarySearch(a[mid+1:], search, a[mid])
	default:
		result = &(a[mid])
	}
	searchCount++
	return
}

func createMessagesBook() string {
	bookBytes, err := ioutil.ReadFile("messages.txt")
	if nil != err {
		log.Fatalln("Unable to open messages.txt file", err)
		return ""
	}
	chars := map[byte]int64{}
	for _, b := range bookBytes {
		switch b {
		case '1', '0', '2', '3', '4', '5', '6', '7', '8', '9', 0:
		default:
			chars[[]byte(strings.ToLower(string(b)))[0]]++
		}
	}

	type kv struct {
		Key   byte
		Value int64
	}

	var ss []kv
	for k, v := range chars {
		ss = append(ss, kv{k, v})
		log.Printf("|%s| %d", string(k), v)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for i, c := range ss[:keyboard.NChars] {
		keyboard.Chars[i] = c.Key
		//log.Printf("|%s| %d", string(c.Key), c.Value)
	}

	noChars := []string{}
	for _, c := range ss[keyboard.NChars:] {
		noChars = append(noChars, string(c.Key))
	}

	words := strings.Fields(string(bookBytes))
	usedWords := []string{}
	bannedChars := strings.Join(noChars, "")
	for _, word := range words {
		if !strings.ContainsAny(word, bannedChars) {
			usedWords = append(usedWords, word)
		}
	}

	return string(bookBytes)
}

func createBook(words []Word, count int, pretty bool) string {
	lastComma := 20
	lastWord := ""
	var book strings.Builder
	for i := 0; i < count; i++ {
		ra, err := rand.Int(rand.Reader, big.NewInt(2147483647))
		if nil != err {
			log.Println(err)
			return ""
		}

		r := float64(ra.Int64()) / 2147483647.0
		word, _ := binarySearch(words, r, Word{})
		for lastWord == word.word {
			ra, err := rand.Int(rand.Reader, big.NewInt(2147483647))
			if nil != err {
				log.Println(err)
				return ""
			}

			r := float64(ra.Int64()) / 2147483647.0
			word, _ = binarySearch(words, r, Word{})
		}

		//fmt.Printf("%.8f %.8f %.8f %s\n", r, word.cumulativePercentage, word.percentage, word.word)

		switch word.word {
		case "-":
			if lastComma == 0 {
				continue
			}
			lastComma = 0
			book.WriteByte(' ')
			book.WriteString(word.word)
		case ":", ";", ",", "?", "!", ".":
			if lastComma == 0 {
				continue
			}
			lastComma = 0
			book.WriteString(word.word)
		default:
			if lastComma > prand.Intn(10)+5 {
				lastComma = 0
				book.WriteString(". ")
				if pretty {
					c := bytes.ToUpper([]byte{word.word[0]})
					book.Write(c)
					book.WriteString(word.word[1:])
				} else {
					book.WriteString(word.word)
				}
			} else {
				lastComma++
				book.WriteByte(' ')
				if pretty {
					switch word.word {
					case "i":
						book.WriteByte('I')
					default:
						book.WriteString(word.word)
					}
				} else {
					book.WriteString(word.word)
				}
			}
		}
		lastWord = word.word
	}
	b := strings.ReplaceAll(book.String()[2:], " 's", "'s")
	return b
}

func searchLoop(thread int, book string, distances *[4][12][4][12]float64, res chan keyboard.Keyboard) {
	mutationsStart := 3
	mutationsEnd := 0
	bestScore := initialScore

	sinceLast := 0
	mutations := mutationsStart
	total := 0
	gen := 1
	kb := keyboard.New()
	bkb := kb.Copy()
	for {
		total++
		if mutations == mutationsEnd {
			bkb = keyboard.New()
			mutations = mutationsStart
			gen++
			bestScore = initialScore
			sinceLast = 0
			// log.Println("Trying new keyboard layout")
		} else if (mutations == 2 && sinceLast > 500) ||
			(mutations == 3 && sinceLast > 500) ||
			(mutations == 1 && sinceLast > 1000) ||
			sinceLast > 100 {
			mutations--
			// log.Println("Switching at try", sinceLast, "mutation iteration to", mutations)
			sinceLast = 0
		}
		kb = bkb.Copy()
		kb.Book = &book
		for i := 0; i < mutations; i++ {
			kb.Mutate()
		}
		kb.FillScore(distances)
		score := kb.Score()
		sinceLast++
		if score < bestScore {
			kb.Gen = gen
			kb.Total = total
			kb.Thread = thread
			kb.Mutation = mutations
			kb.Iteration = sinceLast
			sinceLast = 0
			//mutations++
			bestScore = score
			res <- *kb
			bkb = kb.Copy()
		}
	}
}

const initialScore = 10000000.0

func main() {
	runtime.GOMAXPROCS(14)

	/*data, err := Parse()
	if nil != err {
		log.Fatalln("unable to parse data", err)
		return
	}*/

	//book := createBook(data, 10000, false)
	book := createMessagesBook()

	results := make(chan keyboard.Keyboard, 16)
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

	for i := 0; i < 14; i++ {
		go searchLoop(i, book, &distances, results)
	}

	topScore := initialScore

	for {
		select {
		case res := <-results:
			score := res.Score()
			if score < topScore {
				topScore = score
				fmt.Print((&res).DetailString(), &res, (&res).ScoreString())
			}
		}
	}
}
