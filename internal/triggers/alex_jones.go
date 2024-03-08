package triggers

import (
	"fmt"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var alexJonesInitialSubject = []string{
	"The new world order",        // A
	"Barack Obama",               // B
	"Hillary Clinton",            // C
	"George Soros",               // D
	"Black Lives Matter",         // E
	"Inter-dimensional Vampires", // F
	"The Elites",                 // G
	"Reptilians",                 // H
	"China",                      // I
	"Mexicans",                   // J
	"Antifa",                     // K
	"The liberal media",          // L
	"Demons",                     // M
	"Satanists",                  // N
	"The gays",                   // O
	"Baby-eating atheists",       // P
	"Muslims",                    // Q
	"The Jews",                   // R
	"Globalists",                 // S
	"Commie Sanders",             // T
	"The CIA",                    // U
	"SJWs",                       // V
	"Jack Dorsey",                // W
	"Mark Zuckerberg",            // X
	"The NSA",                    // Y
	"ISIS",                       // Z
	// Added
	"Communists",
	"The Devil",
	"My Neighbors",
	"Paid Actors",
	"RINOs",
	"The WHO",
	"Elite Chinese Medical Researchers",
	"FEMA",
}

var alexJonesVerbs = []string{
	"is spying on",                               // A
	"is collection data on",                      // B
	"is profiting from",                          // C
	"is putting chemicals into",                  // D
	"is corrupting",                              // E
	"is waging war on",                           // F
	"running a sci-op on",                        // G
	"is lying about",                             // H
	"is mind controlling",                        // I
	"is inter-dimensionally crossbreeding with",  // J
	"is Islamising",                              // K
	"is unleashing Satanic powers against",       // L
	"is sacrificing",                             // M
	"is brainwashing",                            // N
	"is plotting a false-flag operation against", // O
	"is fabricating evidence against",            // P
	"is censoring",                               // Q
	"is conspiring against",                      // R
	"is micro-chipping",                          // S
	"is tracking",                                // T
	"is plotting an assassination on",            // U
	"is experimenting on",                        // V
	"is hoarding wealth from",                    // W
	"is using concentration camps to round up",   // X
	"is kidnapping",                              // Y
	"is molesting",                               // Z
	// Additional
	"is using FEMA death camps to round up",
	"is paying",
	"is seducing",
	"is smuggling",
	"is waterboarding",
	"is deleting Hillary Clinton's emails to hide the truth from",
	"is protesting against",
	"is deleting my tweets about",
	"is banning my YouTube channel, harming",
}

var alexJonesPrepositionVictims = []string{
	"honest",
	"brave",
	"patriotic",
	"faithful",
	"studious",
	"enthusiastic",
	"proud",
	"true",
}

var alexJonesVictimSubjects = []string{
	"christian children",          // 1
	"white men",                   // 2
	"Donald Trump",                // 3
	"conservatives",               // 4
	"gun owners",                  // 5
	"frogs",                       // 6
	"God-fearing Americans",       // 7
	"natural-born Americans",      // 8
	"teenage boys",                // 9
	"straight married couples",    // 10
	"unborn babies",               // 11
	"members of right-wing media", // 12
	// Additional
	"white people",
	"daughters of the confederacy",
	"free-thinkers",
	"free-speech",
	"Alex Jones",
	"Sean Hannity",
	"Ben Shapiro",
	"Jordan Peterson",
	"Bannon",
	"Miller",
	"Epstein",
}

var alexJonesConclusions = []string{
	"to establish a New World Order",    // Red
	"to turn them gay",                  // Orange
	"to cover up the truth behind 9/11", // Yellow
	"to destroy the nuclear family",     // Green
	"to achieve eternal life",           // Blue
	"to spread a Satanist agenda",       // Purple
	"to start World War III",            // White
	"to suppress the population",        // Black
	"to control the media narrative",    // Brown
	"to take away our guns",             // Gray
	"to overthrow the US government",    // None
	// Additional
	"to change us to communism",
	"to change us to socialism",
	"to destroy my platform of free speech",
	"to get Donald Trump to love Tiffany Trump",
	"to destroy the universe",
	"to summon the devil",
	"to eat my neighbors before I can",
	"to hide Hillary Clinton's emails",
}

func alexJonesHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	channel, _ := svc.GetChannel(msg.ChannelID, msg.ServerID)

	if !strings.Contains(strings.ToLower(channel.Name), "politics") {
		return core2.EmptyRsp()
	}

	if strings.HasPrefix(strings.ToLower(msg.Contents), "!hannity") {
		timesWaterboarded := []weightedChoice{
			{0.45, fmt.Sprintf("%d", random.Int32Range(2, 20))},
			{0.45, fmt.Sprintf("%d", random.Int64Range(21, 200))},
			{0.10, fmt.Sprintf("%d", random.Int64Range(201, 2300000000))},
		}

		return core2.MakeSingleRsp(fmt.Sprintf(
			"In parallel universe #%d Sean Hannity has been waterboarded ~%s times",
			random.Int64Range(10, 2_300_000_000),
			pickWeightedChoice(timesWaterboarded),
		))
	}

	if !strings.HasPrefix(strings.ToLower(msg.Contents), "!alexjones") {
		return core2.EmptyRsp()
	}

	output := random.PickString(alexJonesInitialSubject) +
		" " + random.PickString(alexJonesVerbs)

	if fromBasicChance("alexJonesHandler--prepositions") {
		output += " " + random.PickString(alexJonesPrepositionVictims)
	}

	output += " " + random.PickString(alexJonesVictimSubjects) +
		" in order " + random.PickString(alexJonesConclusions)

	return core2.MakeSingleRsp(output)
}
