package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type TatSQLVocab struct {
	db *sql.DB
}

// New creates new local sqlite storage.
func New(path string) (*TatSQLVocab, error) {
	dbSql, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err = dbSql.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &TatSQLVocab{db: dbSql}, nil
}

func (tatSQLVocab TatSQLVocab) DoRequest(word string) (bool, error) {
	q := `SELECT COUNT(*) FROM tatar_dict WHERE word = ?`

	var count int

	if err := tatSQLVocab.db.QueryRow(q, strings.ToUpper(word)).Scan(&count); err != nil {
		log.Println("COUNT FROM ERROR", count)
		return false, fmt.Errorf("can't check if word exists: %w", err)
	}
	log.Println("COUNT", count)
	return count > 0, nil
}

func (tatSQLVocab TatSQLVocab) ChooseWord(day int) (string, error) {
	q := `SELECT WORD FROM tatar_dict WHERE dayNumb = ?`
	var word string

	err := tatSQLVocab.db.QueryRow(q, day).Scan(&word)
	if err == sql.ErrNoRows {
		return tatSQLVocab.chooseNewWord(day)

	}
	if err != nil {
		return "", fmt.Errorf("can't choose a new word: %w", err)
	}
	return strings.ToLower(word), nil

}

func (tatSQLVocab TatSQLVocab) chooseNewWord(day int) (string, error) {
	q := `SELECT WORD FROM tatar_dict WHERE dayNumb IS NULL 
	AND meaning NOT LIKE '%иск.%' 
	AND meaning NOT LIKE '%мед.%' 
	AND meaning NOT LIKE '%рус%' 
	AND meaning NOT LIKE '%нем%' 
	AND meaning NOT LIKE '%лат%' 
	AND meaning NOT LIKE '%фр%' 
	AND meaning NOT LIKE '%ингл%' 
	AND meaning NOT LIKE '%гр%' 
	AND meaning NOT LIKE '%гaр%' 
	AND meaning NOT LIKE '%яп%' 
	AND meaning NOT LIKE '%мыск%' 
	AND meaning NOT LIKE '%диал%' 
	ORDER BY random() LIMIT 1;`
	var word string

	err := tatSQLVocab.db.QueryRow(q).Scan(&word)
	if err != nil {
		return "", fmt.Errorf("can't choose a new word: %w", err)
	}

	q = `UPDATE tatar_dict SET dayNumb = ? WHERE word = ?;`
	tatSQLVocab.db.Exec(q, day, word)
	return strings.ToLower(word), nil
}
