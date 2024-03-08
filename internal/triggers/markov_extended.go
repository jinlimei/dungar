package triggers

import (
	"log"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/learning"
	"gitlab.int.magneato.site/dungar/prototype/internal/markov3"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var (
	altMarkovs = make(map[string]*markov3.Markov)

	lastSentOfferings int64
)

func buildAlternativeMarkovs() {
	if utils.InTestEnv() {
		return
	}

	qry := `
		SELECT id, title, contents, variant
		FROM tomes
		WHERE active = 1
	`

	var (
		id       int
		title    string
		contents string
		variant  string
	)

	rows, err := db.ConQuery(qry)
	if err != nil {
		log.Printf("Failed to query database: %v\n", err)
		return
	}

	for rows.Next() {
		err = rows.Scan(&id, &title, &contents, &variant)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return
		}

		_, ok := altMarkovs[title]
		if ok {
			continue
		}

		log.Printf("learning %s...\n", title)
		start := time.Now()

		m := markov3.MakeMarkov(markov3.MarkovSpaceID(title))
		m.MaxChainDistance = 200
		m.MinChainDistance = 20

		switch variant {
		case "book":
			lines := learning.ReadABook(contents)
			for _, line := range lines {
				m.LearnString(line, cleaner.VariantBook)
			}
		case "script", "sonnet":
			lines := strings.Split(contents, "\n")
			for _, line := range lines {
				m.LearnString(line, cleaner.VariantSonnet)
			}
		case "bible":
			lines := learning.ReadABible(contents)
			for _, line := range lines {
				m.LearnString(line, cleaner.VariantBible)
			}
		case "twitter":
			lines := strings.Split(contents, "\n")
			for _, line := range lines {
				m.LearnString(line, cleaner.VariantTwitter)
			}
		default:
			log.Printf("Unknown Variant '%s'\n", variant)
		}

		log.Printf("finished '%s', took %s\n",
			title, time.Now().Sub(start))

		altMarkovs[title] = m
	}
}

func markovExtendedHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if !isDirectedAtDungar(svc, msg) {
		return core2.EmptyRsp()
	}

	if !strings.Contains(msg.Lowered(), "tell me about") {
		return core2.EmptyRsp()
	}

	if len(altMarkovs) == 0 {
		buildAlternativeMarkovs()
	}

	offerings := ""
	for title, mkv := range altMarkovs {
		title = strings.ReplaceAll(title, "_", " ")
		offerings += title + "\n"
		log.Printf("checking input '%s' against '%s'\n",
			msg.Lowered(), title)

		if strings.Contains(msg.Lowered(), title) {
			return core2.MakeSingleRsp(mkv.Generate(mkv.GetRandomWordStr()))
		}
	}

	if time.Now().Unix()-lastSentOfferings > 1800 {
		return core2.MakeSingleRsp("I can tell you about:\n" + offerings)
	}

	return core2.EmptyRsp()
}
