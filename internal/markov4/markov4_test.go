package markov4

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestLearning(t *testing.T) {
	random.UseTestSeed()

	mt := MakeMarkov("")
	mt.LearnString("hello world")
	mt.LearnString("hello friend")
	mt.LearnString("hello people")
}

func TestLearningSubTokens(t *testing.T) {
	mt := MakeMarkov("")
	mt.LearnString("Hello :jiggled:")
	mt.LearnString("Hello world")
	mt.LearnString("Hello friend")
	mt.LearnString("Hello bob")

	fmt.Println(mt.Generate("Hello"))
	fmt.Println(mt.Generate("Hello"))
	fmt.Println(mt.Generate("Hello"))
	fmt.Println(mt.Generate("Hello"))
	fmt.Println(mt.Generate("Hello"))
	fmt.Println(mt.Generate("Hello"))
	fmt.Println(mt.Generate("Hello"))
}

func zTestDonaldTrump(t *testing.T) {
	read, err := os.OpenFile("../../infra/trump.txt", os.O_RDONLY, 0444)
	//write, err := os.OpenFile("../../infra/trump.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to open read: %+v\n", err)
	}

	defer read.Close()
	//defer write.Close()

	scanner := bufio.NewScanner(read)
	mt := MakeMarkov("")

	for scanner.Scan() {
		line := scanner.Text()
		//tw := CleanTweet(line, false)
		//write.WriteString(tw.Cleaned + "\n")
		mt.LearnStringVariant(line, cleaner.VariantTwitter)
	}

	random.UseTimeBasedSeed()

	for k := 0; k < 10; k++ {
		sentence := mt.Generate("Hillary")

		fmt.Println(sentence)
	}

	//words, total := mt.GenerateNextWord("will", "be")
	//
	//justWords := make([]string, 0, len(words))
	//for word, _ := range words {
	//	justWords = append(justWords, word)
	//}
	//
	//sort.Slice(justWords, func(i, j int) bool {
	//	iWord := justWords[i]
	//	jWord := justWords[j]
	//
	//	if words[iWord] < words[jWord] {
	//		return false
	//	}
	//
	//	return true
	//})
	//
	//for _, word := range justWords {
	//	stat := words[word]
	//
	//	if stat <= 1 {
	//		continue
	//	}
	//
	//	fmt.Printf("WORD '%s': %d out of %d (%f%%)\n", word, stat, total, (float64(stat)/float64(total))*100.0)
	//}

	//w, ok := mt.TokenWord["hillary"]
	//assert.True(t, ok)
	//
	//var ids []TokenID
	//for k := 0; k < 10; k++ {
	//	ids = mt.GenerateGraph(w)
	//	assert.True(t, len(ids) > 0)
	//	fmt.Println(strings.TrimSpace(strings.Join(mt.TokensToStrings(ids), " ")))
	//	//spew.Dump(strings.Join(mt.WordIDsToWords(ids), " "))
	//	//ids = mt.GenerateGraph( w )
	//	//spew.Dump(strings.Join(mt.WordIDsToWords(ids), " "))
	//}
}
