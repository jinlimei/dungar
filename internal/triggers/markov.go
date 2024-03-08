package triggers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"gitlab.int.magneato.site/dungar/prototype/internal/markov"
	"gitlab.int.magneato.site/dungar/prototype/internal/markov3"
)

type markovVer uint8

const (
	mV1 markovVer = 1
	mV3 markovVer = 3
)

var (
	learningMutex       = sync.Mutex{}
	activeMarkovVersion = mV3
	lastSpokeM3         = time.Now()

	m3           *markov3.Markov
	lastLoadedM3 time.Time
)

func setMarkovVersion(ver markovVer) {
	activeMarkovVersion = ver
}

func markovUsingV3() *markov3.Markov {
	if m3 == nil {
		m3 = markov3.MakeMarkov("version3")
	}

	return m3
}

func markovV3Learn() {
	if m3 != nil {
		log.Println("[markovV3Learn] Learning.......")
		learningMutex.Lock()

		lastLoadedM3 = time.Now()
		modern := m3.LearnFromRawMessages()
		legacy := 0 //m3.LearnFromRawMessagesM1()

		learningMutex.Unlock()
		log.Printf("[markovV3Learn] Learned (legacy=%d, modern=%d) messages!\n",
			legacy, modern)
	}
}

func markovPickWord() string {
	switch activeMarkovVersion {
	case mV1:
		return markov.PickWord()
	case mV3:
		return markovUsingV3().GetRandomWordStr()
	}

	return ""
}

func markovGenerate(word string) string {
	switch activeMarkovVersion {
	case mV1:
		return markov.Generate(word)
	case mV3:
		now := time.Now()

		if now.Sub(lastLoadedM3).Hours() > 1 {
			log.Println("[markovGenerate] V3 learning expired... learning more.")

			if len(markovUsingV3().RevWords) <= 2 {
				lim := 0
				for len(markovUsingV3().RevWords) <= 2 {
					markovV3Learn()
					lastLoadedM3 = time.Now()

					time.Sleep(1 * time.Second)

					lim++
					if lim > 10 {
						break
					}
				}
			} else {
				markovV3Learn()
				lastLoadedM3 = now
			}

			if lastSpokeM3.Unix() > 0 && now.Sub(lastSpokeM3).Hours() >= 24 {
				lastSpokeM3 = time.Now()
				log.Printf("[markovGenerate] last spoke was >=12 hrs, so updating channel")

				usingOutgoingQueue().Push(&core2.ScheduledMessage{
					ChannelID: "dungar-test",
					Cancelled: false,
					Contents:  generateM3Stats(),
					SentAt:    now,
				})
			}
		}

		return markovUsingV3().Generate(word)
	}

	return ""
}

func markovSettingHandler(msg *core2.IncomingMessage) []*core2.Response {
	var (
		txt    = msg.Contents
		pieces = strings.Split(txt, " ")
	)

	if !strings.HasPrefix(txt, "!m") {
		return core2.EmptyRsp()
	}

	if strings.HasPrefix(txt, "!markov ") && len(pieces) >= 2 {
		if !isAdmin(msg.UserID) {
			if fromBasicChance("markovHandler--badTouch") {
				return core2.PrefixedSingleRsp("NO NOT YOU BAD TOUCH")
			}

			return core2.EmptyRsp()
		}

		switch pieces[1] {
		case "set-version":
			if len(pieces) < 3 {
				return core2.PrefixedSingleRsp("Set version TO WHAT?")
			}

			ver, err := strconv.Atoi(pieces[2])
			if err != nil {
				return core2.PrefixedSingleRsp("You fucked up")
			}

			if ver == 1 || ver == 3 {
				setMarkovVersion(markovVer(ver))
				return core2.MakeSingleRsp(fmt.Sprintf("Set version to '%d'", ver))
			}
		case "get-version":
			return core2.MakeSingleRsp(fmt.Sprintf("Actively using version '%d'", activeMarkovVersion))
		case "load-m3", "m3-load":
			if utils.MustUseEnvVars() {
				go func() {
					markovV3Learn()
					usingOutgoingQueue().Push(&core2.ScheduledMessage{
						ChannelID: "dungar-test",
						Cancelled: false,
						Contents:  "Updated Raw Messages!",
						SentAt:    time.Now(),
					})
				}()
			} else {
				lastLoadedM3 = time.Unix(0, 0)
			}

			return core2.MakeSingleRsp("Yeah sure whatever")
		case "m3-stats", "stats-m3":
			return core2.MakeSingleRsp(generateM3Stats())
		}
	}

	if strings.HasPrefix(txt, "!m3") {
		mm := markovUsingV3()
		word := ""
		out := ""

		if len(pieces) < 2 {
			word = mm.GetRandomWordStr()
			out = mm.Generate(word)
		} else {
			word = strings.ToLower(pieces[1])
			out = mm.Generate(word)
		}

		if out == "" {
			markovV3Learn()
			out = mm.Generate(word)
		}

		return core2.MakeSingleRsp(out)
	}

	if strings.HasPrefix(txt, "!m1") {
		if len(pieces) < 2 {
			return core2.MakeSingleRsp(markov.Generate(markov.PickWord()))
		}

		return core2.MakeSingleRsp(markov.Generate(pieces[1]))
	}

	return core2.EmptyRsp()
}

func generateM3Stats() string {
	m3s := markovUsingV3()

	return fmt.Sprintf(
		"Words: %d\nFragments: %d\nFragmentWords: %d\nLast Loaded: %s\n",
		len(m3s.RevWords),
		len(m3s.Fragments),
		len(m3s.FragmentWords),
		lastLoadedM3.Format(time.RFC822Z),
	)
}
