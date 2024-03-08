package dcli

import (
	"bufio"
	"fmt"
	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/learning"
	"gitlab.int.magneato.site/dungar/prototype/internal/markov3"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// M3LearnCleaned is a methodology for learning a "cleaned" file where each
// line is a singular sentence.
func M3LearnCleaned(id, file string, makeFile bool) {
	var (
		err     error
		line    string
		sock    *os.File
		step    time.Time
		stop    time.Time
		scanner *bufio.Scanner

		pos = 0

		markov *markov3.Markov

		start = time.Now()
	)

	sock, err = os.OpenFile(file, os.O_RDONLY, 0444)
	utils.HaltingError("M3LearnCleaned", err)

	start = time.Now()
	stop = time.Now()
	step = time.Now()

	markov = markov3.MakeMarkov(markov3.MarkovSpaceID(id))
	scanner = bufio.NewScanner(sock)

	for scanner.Scan() {
		if (pos % 500) == 0 {
			stop = time.Now()

			fmt.Printf("Learned %d lines (%v)\n", pos,
				stop.Sub(step).String())

			step = stop
		}

		line = scanner.Text()

		markov.LearnString(line, cleaner.VariantPlain)
		pos++
	}

	if makeFile {
		out := markov.Serialize()
		err = ioutil.WriteFile(fmt.Sprintf("%s.markov.gob", id), out, 0664)

		if err != nil {
			log.Fatalf("Failed to write gob: %v\n", err)
		}
	}

	stop = time.Now()

	fmt.Printf("Learned %d lines (%v, total %v)\n", pos,
		stop.Sub(step).String(), stop.Sub(start).String())
}

// M3LearnBook is our version which allows us to learn text formatted like
// the Alice-in-Wonderland text is formatted.
func M3LearnBook(id, file string, makeFile bool) {
	var (
		bytes []byte
		err   error
		lines []string
		step  time.Time
		stop  time.Time
		max   int

		markov *markov3.Markov

		start = time.Now()
	)

	bytes, err = ioutil.ReadFile(file)
	utils.HaltingError("M3LearnBook", err)

	file = strings.ReplaceAll(file, "\\", "/")

	lines = learning.CleanToLines(string(bytes))

	max = len(lines) - 1

	start = time.Now()
	stop = time.Now()
	step = time.Now()

	markov = markov3.MakeMarkov(markov3.MarkovSpaceID(id))

	for pos, line := range lines {
		if (pos % 500) == 0 {
			stop = time.Now()

			fmt.Printf("Learning line %d of %d lines (%v)\n", pos, max,
				stop.Sub(step).String())

			step = stop
		}

		markov.LearnString(line, cleaner.VariantBook)
	}

	if makeFile {
		out := markov.Serialize()
		err = ioutil.WriteFile(fmt.Sprintf("%s.markov.gob", id), out, 0664)

		if err != nil {
			log.Fatalf("Failed to write gob: %v\n", err)
		}
	}

	stop = time.Now()

	fmt.Printf("Learned %d of %d lines (%v, total %v)\n", max, max,
		stop.Sub(step).String(), stop.Sub(start).String())
}
