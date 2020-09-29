package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Word struct {
	word                 string
	count                int64
	percentage           float64
	cumulativePercentage float64
}

type WordList []Word

func (s WordList) Len() int {
	return len(s)
}
func (s WordList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s WordList) Less(i, j int) bool {
	return s[i].count > s[j].count
}

func Parse() ([]Word, error) {
	file, err := os.Open("simple.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	words := make(map[string]int64, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		word, countStr := parts[0], parts[1]
		count, err := strconv.ParseInt(countStr, 10, 64)
		if nil != err {
			return nil, err
		}
		lo := strings.ToLower(word)
		li := strconv.QuoteToASCII(lo)
		if strings.ContainsAny(li, "0123456789&=`@()_:!{}></\\$*#][") {
			continue
		}
		l := li[1 : len(li)-1]
		if len(l) == 1 && (l != "a" && l != "i" && l != "-" && l != ":" && l != ";" && l != "!" && l != "," && l != "?" && l != ".") {
			continue
		}
		if len(l) > 1 && l[len(l)-1] == '.' {
			continue
		}
		if _, ok := words[l]; ok {
			log.Println("dupe word")
			words[l] += count
		} else {
			words[l] = count
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	unsorted := make([]Word, 0)

	var total float64

	for word, count := range words {
		if count < 5000 {
			continue
		}
		total += float64(count)
		unsorted = append(unsorted, Word{
			word:  word,
			count: count,
		})
	}

	sort.Sort(WordList(unsorted))

	sorted := make([]Word, 0)

	var last float64
	for _, word := range unsorted {

		word.percentage = float64(word.count) / total
		word.cumulativePercentage = word.percentage + last
		last = word.cumulativePercentage
		sorted = append(sorted, word)
	}

	/*for _, word := range sorted {
		log.Println(word.cumulativePercentage, word.word)
	}*/

	return sorted, nil

	/*	out, err := os.Create("best.txt")
		if err != nil {
			log.Fatalln(err)
		}

		defer out.Close()

		w := bufio.NewWriter(out)

		for _, word := range sorted {
			w.WriteString(word.word + "\t" + strconv.FormatInt(word.count, 10) + "\t" + strconv.FormatFloat(word.percentage, 'f', -1, 64) + "\t" + strconv.FormatFloat(word.cumulativePercentage, 'f', -1, 64) + "\n")
		}
		w.Flush()*/
}
