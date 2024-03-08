package triggers

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"regexp"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

var garbagePrefix = regexp.MustCompile("^(<[^>]+>|\\(?[^)]+\\))")

var foodChoices = []string{
	"chipotle",
	"wings",
	"pizza",
	"Panda Express",
	"chinese",
	"soylent",
	"burger",
	"chinese",
	"vietnamese",
	"water",
	"In 'n Out",
	"Five Guys",
}

func randomFood(_ string) string {
	return random.PickString(foodChoices)
}

var nouns = []string{
	"fruit",
	"samara",
	"samurai",
	"google",
	"drupe",
	"legume",
	"pome",
	"pear",
	"apple",
	"grape",
	"polonium",
	"salt",
	"juice",
	"jewelry",
	"butt",
	"ccp",
	"blizzard",
	"twitch",
	"destiny",
	"overwatch",
	"romeo",
	"mysql",
	"js",
	"raw",
	"bitcoin",
	"ethereum",
	"fresh",
	"top",
	"bottom",
	"trump",
	"buzzword",
	"battlestar",
	"fork",
	"battleship",
	"pirate",
	"news",
	"russian",
	"canadian",
	"weeb",
	"keeb",
	"kpop",
	"kappa",
	"panzi",
	"roulette",
	"obama",
	"trump",
}

func randomNoun(_ string) string {
	return random.PickString(nouns)
}

var platforms = []string{
	"MongoDB",
	"Laravel",
	"Symfony",
	"Drupal",
	"PHP",
	"Python",
	"Azure",
	"AWS",
	"Mobile",
	"iOS",
	"macOS",
	"Windows",
	"Linux",
	"MySQL",
	"Postgres",
	"Oracle DB",
	"MSSQL",
	"ODBC",
	"Cloud",
	"GCloud",
	"Kubernetes",
	"Docker",
	"Django",
	"jQuery",
	"memcached",
	"redis",
	"Cassandra",
	"Blizzard",
	"Activision",
	"America",
	"USA",
	"Canada",
	"Europe",
	"UK",
}

func randomPlatform(_ string) string {
	return random.PickString(platforms)
}

func randomNick(serverID string) string {
	if utils.InTestEnv() {
		return ""
	}

	posters := db.GetTopPosters(serverID)
	nicks := make([]string, len(posters))

	pos := 0
	for _, poster := range posters {
		nicks[pos] = poster.Nick
		pos++
	}

	return random.PickString(nicks)
}

func lotsOfSpaceHandler(svc *core2.Service, rsps []*core2.Response) []*core2.Response {
	for _, rsp := range rsps {
		str := rsp.Contents
		str = strings.TrimSpace(str)

		for strings.Contains(str, "  ") {
			str = strings.Replace(str, "  ", " ", -1)
		}

		rsp.Contents = str
	}

	return rsps
}

func removeGarbagePrefixedHandler(svc *core2.Service, rsps []*core2.Response) []*core2.Response {
	for _, rsp := range rsps {
		if garbagePrefix.MatchString(rsp.Contents) {
			rsp.Contents = strings.TrimSpace(garbagePrefix.ReplaceAllString(rsp.Contents, ""))
		}
	}

	return rsps
}
