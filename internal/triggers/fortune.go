package triggers

import (
	"database/sql"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var fortunes []*fortune

type fortune struct {
	ID            int
	Fortune       string
	LastUsed      time.Time
	LastUsedValid bool
}

func (f *fortune) use() string {
	f.LastUsed = time.Now()

	query := `
		UPDATE fortunes
		SET last_used = NOW()
		WHERE id = $1
	`

	db.ConMustExec(query, f.ID)

	return f.Fortune
}

func (f *fortune) usedRecently() bool {
	return (time.Now().Unix() - f.LastUsed.Unix()) <= 3600
}

func fortuneHandler(msg *core2.IncomingMessage) []*core2.Response {
	if !strings.HasPrefix(msg.Contents, "!f") && !strings.HasPrefix(msg.Contents, "!F") {
		return core2.EmptyRsp()
	}

	if len(fortunes) == 0 {
		loadFortunesFromDatabase()
	}

	fortuneLen := len(fortunes)
	picked := ""

	if fortuneLen == 0 {
		return core2.MakeSingleRsp("I'm not a fortune teller maybe go outside or something idk")
	}

	for limit := 0; limit < 10; limit++ {
		id := random.Int(fortuneLen)
		fortune := fortunes[id]

		if !fortune.usedRecently() {
			picked = fortune.use()
			break
		}
	}

	if picked == "" {
		return core2.MakeSingleRsp("I'm not a fortune teller maybe go outside or something")
	}

	return core2.MakeSingleRsp(picked)
}

func loadFortunesFromDatabase() {
	query := `
		SELECT id, fortune, last_used
		FROM fortunes
		WHERE active = 1
	`

	res := db.ConMustQuery(query)

	out := make([]*fortune, 0)

	defer res.Close()
	for res.Next() {
		var id int
		var fortuneStr string
		var lastUsedStr sql.NullString
		ts := time.Unix(1, 1)

		if err := res.Scan(&id, &fortuneStr, &lastUsedStr); err != nil {
			haltingErr("loadFortunesFromDatabase", err)
		}

		if lastUsedStr.Valid {
			t, err := time.Parse(time.RFC3339, lastUsedStr.String)
			haltingErr("time parse", err)
			ts = t
		}

		out = append(out, &fortune{
			ID:            id,
			Fortune:       fortuneStr,
			LastUsed:      ts,
			LastUsedValid: lastUsedStr.Valid,
		})
	}

	fortunes = out
}
