package triggers

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	// necessary for doing sqlite things
	_ "modernc.org/sqlite"
	"modernc.org/sqlite/vfs"
)

//go:embed pins.sqlite
var pinsSqliteFile embed.FS

var (
	pinService *pinDBService

	pinServiceHasError bool
)

type pinDBService struct {
	db *sql.DB
}

type simplePin struct {
	UserName string
	Text     string
}

func (pds *pinDBService) requestAllPins() ([]simplePin, error) {
	qry := `
		SELECT userDisplay, messageText
		FROM messages
`
	rows, err := pds.db.Query(qry)

	if err != nil {
		return nil, err
	}

	return pds.completeRows(rows, err)
}

func (pds *pinDBService) requestPinsFromText(text string) ([]simplePin, error) {
	qry := `
		SELECT userDisplay, messageText
		FROM messages
		WHERE messageText LIKE $1
`

	return pds.completeQuerySearch(qry, text)
}

func (pds *pinDBService) requestPinsFromUser(name string) ([]simplePin, error) {
	qry := `
		SELECT userDisplay, messageText
		FROM messages
		WHERE userDisplay LIKE $1
`

	return pds.completeQuerySearch(qry, name)
}

func (pds *pinDBService) completeQuerySearch(qry string, text string) ([]simplePin, error) {
	rows, err := pds.db.Query(
		qry,
		"%"+strings.ReplaceAll(text, "%", "")+"%",
	)

	if err != nil {
		return nil, err
	}

	return pds.completeRows(rows, err)
}

func (pds *pinDBService) completeRows(rows *sql.Rows, err error) ([]simplePin, error) {
	var (
		pins = make([]simplePin, 0)

		userDisplay string
		messageText string
	)

	for rows.Next() {
		err = rows.Scan(&userDisplay, &messageText)

		if err != nil {
			return nil, err
		}

		pins = append(pins, simplePin{
			UserName: userDisplay,
			Text:     messageText,
		})
	}

	return pins, nil
}

func startPinService() (*pinDBService, error) {
	fn, _, err := vfs.New(pinsSqliteFile)

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", "file:pins.sqlite?vfs="+fn)

	if err != nil {
		return nil, err
	}

	return &pinDBService{db: db}, nil
}

func pinDBHandler(msg *core2.IncomingMessage) []*core2.Response {
	if !strings.HasPrefix(msg.Contents, "!pins") {
		return core2.EmptyRsp()
	}

	if pinServiceHasError {
		return core2.EmptyRsp()
	}

	var (
		parts = strings.Split(msg.Contents, " ")

		err  error
		pins []simplePin
	)

	pinService, err = startPinService()

	if err != nil {
		log.Printf("ERROR: Failed to start pin service: %v", err)

		pinServiceHasError = true
		return core2.MakeSingleRsp("Oh god pins caught on fire")
	}

	if len(parts) == 1 {
		pins, err = pinService.requestAllPins()

		if err != nil {
			pinServiceHasError = true
			return core2.MakeSingleRsp(fmt.Sprintf("oh god requesting all pins caught fire: %v", err))
		}

		// We'll just pick one at random fuck it yolo
	} else if len(parts) >= 2 {
		if strings.HasPrefix(parts[1], "@") {
			pins, err = pinService.requestPinsFromUser(parts[1][1:])
		} else {
			pins, err = pinService.requestPinsFromText(parts[1])
		}

		if err != nil {
			pinServiceHasError = true
			return core2.MakeSingleRsp(fmt.Sprintf("oh god requesting some pins caught fire: %v", err))
		}
	}

	if len(pins) == 0 {
		return core2.MakeSingleRsp(markovGenerate(parts[1]))
	}

	pinNum := random.Int(len(pins))
	picked := pins[pinNum]

	return core2.MakeSingleRsp(fmt.Sprintf(
		"%s\n> pin from %s",
		picked.Text,
		picked.UserName,
	))
}
