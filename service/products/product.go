package products

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

const productPageLimit = 10
const randProducts = "1"

var info db.InfoFetcher

func Init() {
	obj := db.InfoTable{Dbfile: config.GetValue("DbFile"),
		Tablename: "products",
		Keyattrs:  []string{"Name", "Title"}}
	info = obj
}

func InitTable() {
	log.Info("create table products")
	info.CreateTable()
}

func Register() {
	log.Info("product registered")
	ws := new(restful.WebService)
	ws.Path("/service/product").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("").To(getRandomProducts))

	//standard apis
	ws.Route(ws.GET("/{name}").To(getProduct))
	ws.Route(ws.POST("").To(addProduct).Filter(auth.AuthEmployeeFilter))
	ws.Route(ws.PUT("").To(updateProduct).Filter(auth.AuthEmployeeFilter))
	ws.Route(ws.DELETE("/{id}").To(delProduct).Filter(auth.AuthEmployeeFilter))

	ws.Route(ws.GET("/search/{searchtext}/page/{pagenumber}").To(searchProducts).Filter(auth.AuthEmployeeFilter))
	ws.Route(ws.GET("/search/{searchtext}/count").To(searchProductsCount).Filter(auth.AuthEmployeeFilter))
	restful.Add(ws)
}

func searchProducts(req *restful.Request, resp *restful.Response) {
	t := req.PathParameter("searchtext")
	p, err := strconv.Atoi(req.PathParameter("pagenumber"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	offset := productPageLimit * (p - 1)
	all, ret := info.SelectRows(fmt.Sprintf(" status=1 AND Name like '%%%s%%' order by id limit %d offset %d", t, productPageLimit, offset))
	if ret == http.StatusOK && len(all) > 0 {
		for i := 0; i < len(all); i++ {
			db.VoidAttr(&all[i], "Spec", "Introduction")
		}
		resp.WriteEntity(all)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func searchProductsCount(req *restful.Request, resp *restful.Response) {
	t := req.PathParameter("searchtext")
	n, ret := info.SelectRowsCount(fmt.Sprintf(" status=1 AND Name like '%%%s%%' ", t))
	if ret == http.StatusOK {
		resp.WriteEntity(SearchCount{n, productPageLimit})
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func getRandomProducts(req *restful.Request, resp *restful.Response) {
	total, ret := info.SelectRowsCount(" status=1 ")
	if ret == http.StatusOK && total > 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		offset := rand.Int63n(total)
		all, ret := info.SelectRows(" status=1 order by id desc limit " + randProducts + " offset " + strconv.FormatInt(offset, 10))
		if ret == http.StatusOK {
			for i := 0; i < len(all); i++ {
				db.VoidAttr(&all[i], "Spec", "Introduction")
			}
			resp.WriteEntity(all)
		} else {
			resp.WriteErrorString(ret, http.StatusText(ret))
		}
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func addProduct(req *restful.Request, resp *restful.Response) {
	info.Add(req, resp)
}

func delProduct(req *restful.Request, resp *restful.Response) {
	info.Del(req, resp)
}

func updateProduct(req *restful.Request, resp *restful.Response) {
	info.Update(req, resp)
}

func getProduct(req *restful.Request, resp *restful.Response) {
	info.Get(req, resp)
}
