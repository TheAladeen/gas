package db

import (
	"database/sql"
	"fmt"
	log "github.com/featen/utils/log"
	"net/http"
    _ "github.com/mattn/go-sqlite3"
    "encoding/json"
)

type InfoTable struct {
	Dbfile    string
	Tablename string
    Keyattrs []string
}

type InfoObject struct {
	Id     int64
	Status int64
	Nav    string
	Info   string
}

type InfoFetcher interface {
	CreateTable() int
	FetchInfoRows(string) ([]InfoObject, int)
	InsertRow(string) (int64, int)
}

func (infotable InfoTable) CreateTable() int {
    fmt.Println(infotable.Dbfile)
	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		fmt.Println("dbHandler failed", err)
	}
	defer dbHandler.Close()

	sqls := []string{
		"create table if not exists " + infotable.Tablename + " (id integer NOT NULL PRIMARY KEY, status int default 1, nav text unique, info text)",
		"create table if not exists " + infotable.Tablename + "_attr (obj_id int, key text, value text)",
	}

	for _, s := range sqls {
        fmt.Println(s)
		_, err := dbHandler.Exec(s)
		if err != nil {
			fmt.Println("%q: %s\n", err, s)
		}
	}
	return http.StatusOK
}

func (infotable InfoTable) InsertRow(infostr string) (int64, int) {
    var m map[string]interface{}

    err := json.Unmarshal([]byte(infostr), &m)
    if err != nil {
        fmt.Println("error:", err)
        return 0, http.StatusBadRequest
    }


	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		log.Fatal("%v", err)
        return 0, http.StatusInternalServerError
	}
	defer dbHandler.Close()

    str := "INSERT INTO " + infotable.Tablename + " (nav, info) VALUES (?, ?)"
	stmt, err := dbHandler.Prepare(str)
	if err != nil {
		log.Error("%v", err)
		return 0, http.StatusInternalServerError
	}
	defer stmt.Close()

	r, err := stmt.Exec(m["NavName"], infostr)
	if err != nil {
		log.Error("%v", err)
		return 0, http.StatusBadRequest
	}
	id, _ := r.LastInsertId()


    tx, err := dbHandler.Begin()
	if err != nil {
		log.Fatal("%v", err)
        return 0, http.StatusInternalServerError
	}
	stmt, err = tx.Prepare("insert into " + infotable.Tablename + "_attr  (obj_id, key, value) values(?, ?, ?)")
	if err != nil {
		log.Fatal("%v", err)
        tx.Rollback()
        return 0, http.StatusInternalServerError
	}
	defer stmt.Close()

	for i := 0; i < len(infotable.Keyattrs); i++ {
        v, ok := m[infotable.Keyattrs[i]]
        if ok {
            _, err = stmt.Exec(id, infotable.Keyattrs[i], v)
    		if err != nil {
    			log.Fatal("%v", err)
                tx.Rollback()
                return 0, http.StatusInternalServerError
    		}
        }

	}
	tx.Commit()


	return id, http.StatusOK
}

func (infotable InfoTable) FetchInfoRows(sqlstr string) ([]InfoObject, int) {
	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		log.Error("%v", err)
		return nil, http.StatusInternalServerError
	}
	defer dbHandler.Close()

    str := "SELECT t.id, t.status, t.nav, t.info FROM " + infotable.Tablename + " t, " + infotable.Tablename + "_attr attr WHERE t.id=attr.obj_id AND " + sqlstr
	stmt, err := dbHandler.Prepare(str)
	if err != nil {
		log.Error("%v", err)
		return nil, http.StatusInternalServerError
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("%v", err)
		return nil, http.StatusInternalServerError
	}
	defer rows.Close()

	all := make([]InfoObject, 0)
	for rows.Next() {
		var nav, info sql.NullString
		var id, status sql.NullInt64
		rows.Scan(&id, &status, &nav, &info)

        fmt.Println("one row found")
		all = append(all, InfoObject{id.Int64, status.Int64, nav.String, info.String})
	}
    if len(all) == 0 {
        return nil, http.StatusNotFound
    }
    fmt.Println("total rows", len(all))

	return all, http.StatusOK
}
