package dcli

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

type metaInfo struct {
	Name    string `json:"name"`
	Variant string `json:"variant"`
}

// StoreTextFiles will store text files to do things haha yes
func StoreTextFiles(dir string) {
	utils.LoadSettingsAndSecrets()
	db.ConnectToDatabase()

	var (
		tx  *sql.Tx
		err error

		meta  []byte
		infos []metaInfo
		files []os.DirEntry
	)

	tx = db.BeginTransaction()

	meta, err = os.ReadFile(dir + "/meta.json")

	err = json.Unmarshal(meta, &infos)
	utils.HaltingError("StoreTextFiles.Unmarshal", err)

	files, err = os.ReadDir(dir)
	utils.HaltingError("StoreTextFiles.ReadDir", err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var (
			name string
			info metaInfo
			ok   bool
			data []byte
		)

		name = file.Name()

		fmt.Printf("Looking at name '%s' ", name)

		if len(name) < 4 || name[len(name)-3:] != "txt" {
			fmt.Print("Skipping\n")
			continue
		}

		info, ok = findMetaInfo(name[0:len(name)-4], infos)
		if !ok {
			fmt.Printf("TXT file '%s' with no meta entry in json\n", name)
			continue
		}

		data, err = os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			fmt.Printf("Could not read '%s': %v\n", file.Name(), err)
			continue
		}

		storeFile(tx, info.Name, string(data), info.Variant)
		fmt.Print("Inserted File\n")
	}

	db.CommitOrRollback(tx)
	fmt.Println("Committing...")
}

func findMetaInfo(name string, infos []metaInfo) (metaInfo, bool) {
	for _, info := range infos {
		if info.Name == name {
			return info, true
		}
	}

	return metaInfo{}, false
}

func findFile(tx *sql.Tx, title string) (int, bool, error) {
	qry := `
		SELECT id
		FROM tomes
		WHERE title = $1
	`

	row := db.QueryRow(tx, qry, title)

	var id int

	if row == nil {
		return id, false, errors.New("row returned nil")
	}

	if row.Err() != nil {
		fmt.Print(" row.Err() ")
		return id, false, row.Err()
	}

	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return id, false, nil
		}

		fmt.Print(" err from row.Scan() ")
		return id, false, err
	}

	return id, true, nil
}

func storeFile(tx *sql.Tx, title, contents, variant string) {
	id, exists, err := findFile(tx, title)

	if err != nil {
		fmt.Printf(" error searching (%v) ", err)
		return
	}

	if !exists {
		qry := `
			INSERT INTO tomes (title, contents, variant, active, created_at, updated_at)
			VALUES($1, $2, $3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`

		db.MustExec(tx, qry, title, contents, variant)
	} else {
		qry := `
			UPDATE tomes
			SET contents   = $2,
					variant    = $3,
					updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`

		db.MustExec(tx, qry, id, contents, variant)
	}
}
