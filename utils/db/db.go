package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/featen/utils/log"
	_ "github.com/mattn/go-sqlite3"
)

type InfoTable struct {
	Dbfile    string
	Tablename string
	Keyattrs  []string
}

type InfoObject struct {
	Id     int64
	Status int64
	Nav    string
	Info   string //this contains all the info which will be decoded by js.
}

type InfoFetcher interface {
	CreateTable() int
	SelectRows(string) ([]InfoObject, int)
	SelectRowsCount(string) (int64, int)
	InsertRow(string) (int64, int)
	DeleteRow(string) int
	UpdateRow(int64, string) int
}

func VoidAttr(obj *InfoObject, attrs ...string) int {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(obj.Info), &m)
	if err != nil {
		log.Error("Unmarshal json failed %v:\n%s", err, obj.Info)
		return http.StatusNotAcceptable
	}

	for _, v := range attrs {
		m[v] = ""
	}

	b, err := json.Marshal(m)
	if err != nil {
		log.Error("Marshal json failed %v:\n%v", err, m)
		return http.StatusInternalServerError
	}
	obj.Info = string(b)

	return http.StatusOK
}

func (infotable InfoTable) CreateTable() int {
	fmt.Println(infotable.Dbfile)
	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		fmt.Println("dbHandler failed", err)
	}
	defer dbHandler.Close()

	var keys = ""
	for _, v := range infotable.Keyattrs {
		keys += fmt.Sprint(", ", v, " text")
	}
	s := "create table if not exists " + infotable.Tablename + " (id integer NOT NULL PRIMARY KEY, status int default 1, nav text unique, info text " + keys + ")"

	fmt.Println(s)
	_, err = dbHandler.Exec(s)
	if err != nil {
		fmt.Println("%q: %s\n", err, s)
	}
	return http.StatusOK
}

func (infotable InfoTable) AlterTable() int {
	return 0
}

func (infotable InfoTable) UpdateRow(id int64, infostr string) int {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(infostr), &m)
	if err != nil {
		log.Error("Unmarshal json failed %v:\n%s", err, infostr)
		return http.StatusInternalServerError
	}

	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		log.Fatal("open db failed: %v", err)
		return http.StatusInternalServerError
	}
	defer dbHandler.Close()

	ups := ""
	for i := 0; i < len(infotable.Keyattrs); i++ {
		v, ok := m[infotable.Keyattrs[i]]
		if ok {
			ups += fmt.Sprint(", ", infotable.Keyattrs[i], "='", v, "'")
		}
	}
	str := "UPDATE " + infotable.Tablename + " SET info=? " + ups + " WHERE id=" + strconv.FormatInt(id, 10)
	stmt, err := dbHandler.Prepare(str)
	if err != nil {
		log.Error("%s failed: %v", str, err)
		return http.StatusInternalServerError
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Error("%s failed: %v", str, err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (infotable InfoTable) DeleteRow(infostr string) int {
	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		log.Fatal("open db failed: %v", err)
		return http.StatusInternalServerError
	}
	defer dbHandler.Close()

	str := "UPDATE " + infotable.Tablename + " SET status=0 WHERE " + infostr
	stmt, err := dbHandler.Prepare(str)
	if err != nil {
		log.Error("%s failed: %v", str, err)
		return http.StatusInternalServerError
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Error("%s failed: %v", str, err)
		return http.StatusInternalServerError
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

	keys := ""
	values := ""
	for i := 0; i < len(infotable.Keyattrs); i++ {
		v, ok := m[infotable.Keyattrs[i]]
		if ok {
			keys += fmt.Sprint(", '", infotable.Keyattrs[i], "'")
			values += fmt.Sprint(", '", v, "'")
		}
	}
	str := "INSERT INTO " + infotable.Tablename + " (nav, info" + keys + " ) VALUES (?, ? " + values + ")"
	stmt, err := dbHandler.Prepare(str)
	if err != nil {
		log.Error("%v", err)
		return 0, http.StatusInternalServerError
	}
	defer stmt.Close()

	r, err := stmt.Exec(m["Nav"], infostr)
	if err != nil {
		log.Error("%v", err)
		return 0, http.StatusBadRequest
	}
	id, _ := r.LastInsertId()

	return id, http.StatusOK
}

func (infotable InfoTable) SelectRows(sqlstr string) ([]InfoObject, int) {
	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		log.Error("%v", err)
		return nil, http.StatusInternalServerError
	}
	defer dbHandler.Close()

	str := "SELECT id, status, nav, info FROM " + infotable.Tablename + " WHERE " + sqlstr
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

func (infotable InfoTable) SelectRowsCount(sqlstr string) (int64, int) {
	dbHandler, err := sql.Open("sqlite3", infotable.Dbfile)
	if err != nil {
		log.Error("%v", err)
		return 0, http.StatusInternalServerError
	}
	defer dbHandler.Close()

	str := "SELECT count(id) FROM " + infotable.Tablename + " WHERE " + sqlstr
	stmt, err := dbHandler.Prepare(str)
	if err != nil {
		log.Error("%v", err)
		return 0, http.StatusInternalServerError
	}
	defer stmt.Close()

	var count sql.NullInt64
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		log.Fatal("%v", err)
		return 0, http.StatusInternalServerError
	}

	return count.Int64, http.StatusOK
}
