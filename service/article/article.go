package article

import (
	"math"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/featen/ags/service/auth"
	"github.com/featen/ags/service/config"
	db "github.com/featen/ags/utils/db"
	log "github.com/featen/ags/utils/log"
)

const timeLayout = "2006-01-02 3:04pm"
const pageArticlesLimit = 5

var info db.InfoFetcher

func Init() {
	obj := db.InfoTable{Dbfile: config.GetValue("DbFile"), Tablename: "article", Keyattrs: []string{"Title"}}
	info = obj
}

func InitTable() {
	log.Info("create table article")
	info.CreateTable()
}

func Register() {
	log.Info("article registered")

	ws := new(restful.WebService)
	ws.Path("/service/article").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{id}").To(getArticle))
	ws.Route(ws.PUT("").To(updateArticle).Filter(auth.AuthFilter))
	ws.Route(ws.POST("").To(addArticle).Filter(auth.AuthFilter))
	ws.Route(ws.DELETE("/{id}").To(delArticle).Filter(auth.AuthFilter))
	//extra apis for page rendering.
	ws.Route(ws.GET("/totalpagenumber").To(getTotalPageNumber))
	ws.Route(ws.GET("/page/{pageNumber}").To(getPageArticles))

	restful.Add(ws)
}

func getArticleById(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	log.Debug("get article by id %s", id)

	all, ret := info.SelectRows(" id=" + id)
	if ret == http.StatusOK && len(all) == 1 {
		resp.WriteEntity(all[0])
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func getArticle(req *restful.Request, resp *restful.Response) {
	info.Get(req, resp)
}

func updateArticle(req *restful.Request, resp *restful.Response) {
	info.Update(req, resp)
}

func addArticle(req *restful.Request, resp *restful.Response) {
	info.Add(req, resp)
}

func delArticle(req *restful.Request, resp *restful.Response) {
	info.Del(req, resp)
}

func getTotalPageNumber(req *restful.Request, resp *restful.Response) {
	total, ret := info.SelectRowsCount(" status=1 ")
	if ret == http.StatusOK {
		pageNumber := math.Ceil(float64(total) / float64(pageArticlesLimit))
		resp.WriteHeader(http.StatusOK)
		resp.WriteEntity(pageNumber)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func getPageArticles(req *restful.Request, resp *restful.Response) {
	pagenumber, err := strconv.ParseInt(req.PathParameter("pageNumber"), 10, 64)
	if err != nil || pagenumber <= 0 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	offset := (pagenumber - 1) * pageArticlesLimit

	all, ret := info.SelectRows(" status=1 order by id desc limit " + strconv.FormatInt(pageArticlesLimit, 10) + " offset " + strconv.FormatInt(offset, 10))
	if ret == http.StatusOK {
		for i := 0; i < len(all); i++ {
			db.VoidAttr(&all[i], "Content")
		}
		resp.WriteEntity(all)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}
