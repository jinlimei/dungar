package main

import (
	"log"
	"os"

	"gitlab.int.magneato.site/dungar/prototype/internal/dcli"
)

func main() {
	if len(os.Args) <= 1 {
		printCommands()
		os.Exit(1)
	}

	switch os.Args[1] {

	case "store-text", "text-store", "store-txt", "txt-store":
		if len(os.Args) <= 2 {
			log.Println("Must provide <dir> parameter")
			log.Println("dungar store-text <dir>")
			os.Exit(4)
			return
		}

		dcli.StoreTextFiles(os.Args[2])

	case "reconstruct-m1":
		if len(os.Args) <= 2 {
			log.Println("Must provide <file> parameter")
			log.Println("dungar reconstruct-m1 <file>")
			os.Exit(3)
			return
		}

		dcli.ReconstructFragments(os.Args[2])

	case "run":
		dcli.ProtocolRunner()

	case "bot-info":
		dcli.PrintBotInfo()

	case "import-fortunes":
		if len(os.Args) <= 2 {
			log.Println("Must provide <file> parameter")
			log.Println("dungar import-fortunes <file>")
			os.Exit(2)
			return
		}

		dcli.ImportFortunes(os.Args[2])

	case "learn1-file":
		if len(os.Args) <= 2 {
			log.Println("Must provide <file> parameter")
			log.Println("dungar learn1-file <file>")
			os.Exit(2)
			return
		}

		dcli.M1LearnFile(os.Args[2])

	case "learn3-book":
		if len(os.Args) <= 3 {
			log.Println("Must provide <name> <file> parameter")
			log.Println("dungar learn3-book <file>")
			os.Exit(2)
			return
		}

		dcli.M3LearnBook(os.Args[2], os.Args[3], true)

	case "learn3-cleaned":
		if len(os.Args) <= 3 {
			log.Println("Must provide <name> <file> parameter")
			log.Println("dungar learn3-cleaned <file>")
			os.Exit(2)
			return
		}

		dcli.M3LearnCleaned(os.Args[2], os.Args[3], true)

	case "help", "-h", "-help", "--help":
		printCommands()

	default:
		log.Printf("unknown case: %s\n", os.Args[1])
		os.Exit(2)

	}
}

func printCommands() {
	log.Print(`
dungar (just another one of those markov bots)

usage: dungar <routine> [options]

routines:
  - dungar run
  - dungar bot-info
  - dungar import-fortunes <file>
  - dungar help
  - dungar learn1-file <file>
  - dungar learn3-book <id> <file>
  - dungar learn3-cleaned <id> <file>
  - dungar reconstruct-m1 <file>
  - dungar store-text <dir>

`)
}
