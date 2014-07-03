package customer

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/featen/ags/service/auth"
	"github.com/featen/ags/service/config"
	db "github.com/featen/ags/utils/db"
	log "github.com/featen/ags/utils/log"
)

type SearchCount struct {
	Total     int64
	PageLimit int
}

const pageItemLimit = 10

var info db.InfoFetcher

func Init() {
	obj := db.InfoTable{Dbfile: config.GetValue("DbFile"),
		Tablename: "customer",
		Keyattrs:  []string{"Name", "Email"}}
	info = obj
}

func InitTable() {
	log.Info("create table customer")
	info.CreateTable()
}

func Register() {
	log.Info("customer registered")
	ws := new(restful.WebService)
	ws.Path("/service/customer").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	//standard apis
	ws.Route(ws.GET("/{id}").To(getCustomer).Filter(auth.AuthFilter))
	ws.Route(ws.POST("").To(addCustomer).Filter(auth.AuthFilter))
	ws.Route(ws.PUT("").To(updateCustomer).Filter(auth.AuthFilter))
	ws.Route(ws.DELETE("/{id}").To(delCustomer).Filter(auth.AuthFilter))

	ws.Route(ws.GET("/search/{searchtext}/page/{pagenumber}").To(searchCustomers).Filter(auth.AuthFilter))
	ws.Route(ws.GET("/search/{searchtext}/count").To(searchCustomersCount).Filter(auth.AuthFilter))
	restful.Add(ws)
}

func searchCustomers(req *restful.Request, resp *restful.Response) {
	t := req.PathParameter("searchtext")
	p, err := strconv.Atoi(req.PathParameter("pagenumber"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	offset := pageItemLimit * (p - 1)
	all, ret := info.SelectRows(fmt.Sprintf(" status=1 AND (Name like '%%%s%%' or Email like '%%%s%%') order by id limit %d offset %d", t, t, pageItemLimit, offset))
	if ret == http.StatusOK && len(all) > 0 {
		for i := 0; i < len(all); i++ {
			db.VoidAttr(&all[i], "Log")
		}
		resp.WriteEntity(all)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func searchCustomersCount(req *restful.Request, resp *restful.Response) {
	t := req.PathParameter("searchtext")
	n, ret := info.SelectRowsCount(fmt.Sprintf(" status=1 AND (Name like '%%%s%%' or Email like '%%%s%%') ", t, t))
	if ret == http.StatusOK {
		resp.WriteEntity(SearchCount{n, pageItemLimit})
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func addCustomer(req *restful.Request, resp *restful.Response) {
	info.Add(req, resp)
}

func delCustomer(req *restful.Request, resp *restful.Response) {
	info.Del(req, resp)
}

func updateCustomer(req *restful.Request, resp *restful.Response) {
	info.Update(req, resp)
}

func getCustomer(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	all, ret := info.SelectRows(" status=1 AND id='" + id + "'")
	if ret == http.StatusOK && len(all) == 1 {
		resp.WriteEntity(all[0])
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}
