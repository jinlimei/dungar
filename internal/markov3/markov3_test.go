package markov3

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestLearning(t *testing.T) {
	random.UseTestSeed()

	mt := MakeMarkov("")
	mt.LearnString("hello world", cleaner.VariantPlain)
	mt.LearnString("hello friend", cleaner.VariantPlain)
	mt.LearnString("hello people", cleaner.VariantPlain)

	assert.Equal(t, []string{"", "hello", "world", "friend", "people"}, mt.RevWords)
	wID, ok := mt.FindWord("hello")

	assert.True(t, ok)
	assert.Equal(t, MarkovID(1), wID)

	cnt := mt.WordFragmentCount(wID)

	assert.Equal(t, 3, cnt)
}

func TestLearningSubTokens(t *testing.T) {
	mt := MakeMarkov("")
	mt.LearnString("Hello :jiggled:", cleaner.VariantSlack)

	assert.Equal(t, []string{"", "Hello", "\u0000Emoticon"}, mt.RevWords)
}

func zTestDonaldTrump(t *testing.T) {
	read, err := os.OpenFile("../../infra/war_of_the_worlds.txt", os.O_RDONLY, 0444)
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
		mt.LearnString(line, cleaner.VariantBook)
	}

	vars := mt.FindWordVariants("thousand")
	fmt.Printf("Found Variants: %+v\n", mt.WordIDsToWords(vars))

	random.UseTimeBasedSeed()

	var ids []MarkovID
	for k := 0; k < 3; k++ {
		ids = mt.GenerateGraph(vars[random.Int(len(vars))])
		//assert.True(t, len(ids) > 0)
		fmt.Println(strings.TrimSpace(strings.Join(mt.WordIDsToWords(ids), " ")))
		//spew.Dump(strings.Join(mt.WordIDsToWords(ids), " "))
		//ids = mt.GenerateGraph( vars )
		//spew.Dump(strings.Join(mt.WordIDsToWords(ids), " "))
	}
}
