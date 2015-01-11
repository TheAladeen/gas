package service

import (
	"database/sql"
	"fmt"

	"github.com/featen/gas/service/article"
	"github.com/featen/gas/service/auth"
	"github.com/featen/gas/service/config"
	"github.com/featen/gas/service/product"
	"github.com/featen/gas/service/share"
	"github.com/featen/gas/service/user"
	log "github.com/featen/gas/utils/log"
	_ "github.com/mattn/go-sqlite3"
)

func createDb() {
	dbHandler, err := sql.Open("sqlite3", config.GetValue("DbFile"))
	if err != nil {
		log.Fatal("%v", err)
		fmt.Println("dbHandler failed", err)
	}
	defer dbHandler.Close()

	sqls := []string{
		//init dict table
		"create table if not exists dict (id integer NOT NULL PRIMARY KEY, q text, fanyi text)",
	}

	for _, s := range sqls {
		_, err := dbHandler.Exec(s)
		if err != nil {
			log.Fatal("%q: %s\n", err, s)
		}
	}

	article.InitTable()
	product.InitTable()
}

func RegService() {
	config.InitConfigs("data/config.js")
	article.Init()
	product.Init()

	auth.SetSysMagicNumber([]byte(config.GetValue("SysMagicNumber")))
	inited := config.IsConfigInited()
	if !inited {
		createDb()
		config.SetValue("dbInited", "Y")
	}

	user.Register()
	article.Register()
	share.Register()
	product.Register()
	//dict.Register()
}
