package articles

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

var info db.InfoFetcher

func Init() {
	obj := db.InfoTable{Dbfile: config.GetValue("DbFile"), Tablename: "articles", Keyattrs: []string{"Title", "CreateTime"}}
	info = obj
}

func InitTable() {
	log.Info("create table articles")
	info.CreateTable()
}

func Register() {
	log.Info("articles registered")

	ws := new(restful.WebService)
	ws.Path("/service/articles").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	//	ws.Route(ws.GET("").To(getAllArticles))
	ws.Route(ws.GET("/{article-id}").To(getArticleById))
	ws.Route(ws.GET("/name/{nav}").To(getArticleByNav))
	ws.Route(ws.PUT("/{article-id}").To(updateArticle).Filter(auth.AuthEmployeeFilter))
	ws.Route(ws.POST("").To(createArticle).Filter(auth.AuthEmployeeFilter))
	ws.Route(ws.DELETE("/{article-id}").To(removeArticle).Filter(auth.AuthFilter))
	//extra apis for page rendering.
	ws.Route(ws.GET("/totalpage/number").To(getTotalPageNumber))
	ws.Route(ws.GET("/page/{pageNumber}").To(getPageArticles))

	restful.Add(ws)
}

func getAllArticles(req *restful.Request, resp *restful.Response) {
	log.Debug("get all articles")

	all, ret := info.SelectRows(" status=1 ")
	if ret == http.StatusOK {
		resp.WriteEntity(all)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func getArticleById(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("article-id")
	log.Debug("get article by id %s", id)

	all, ret := info.SelectRows(" id=" + id)
	if ret == http.StatusOK && len(all) == 1 {
		resp.WriteEntity(all[0])
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func getArticleByNav(req *restful.Request, resp *restful.Response) {
	nav := req.PathParameter("nav")
	log.Debug("get article by nav %s", nav)

	all, ret := info.SelectRows(" status=1 AND nav='" + nav + "'")
	if ret == http.StatusOK && len(all) == 1 {
		resp.WriteEntity(all[0])
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func updateArticle(req *restful.Request, resp *restful.Response) {
	obj := new(db.InfoObject)
	err := req.ReadEntity(&obj)
	if err == nil {
		if ret := info.UpdateRow(obj.Id, obj.Info); ret == http.StatusOK {
			resp.WriteHeader(http.StatusOK)
		} else {
			resp.WriteErrorString(ret, http.StatusText(ret))
		}
	} else {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}

func createArticle(req *restful.Request, resp *restful.Response) {
	obj := new(db.InfoObject)
	err := req.ReadEntity(&obj)
	if err == nil {
		_, ret := info.InsertRow(obj.Info)
		if ret == http.StatusOK {
			resp.WriteHeader(http.StatusCreated)
		} else {
			resp.WriteErrorString(ret, http.StatusText(ret))
		}
	} else {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}

func removeArticle(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("article-id")
	ret := info.DeleteRow(" id=" + id)
	if ret == http.StatusOK {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func getTotalPageNumber(req *restful.Request, resp *restful.Response) {
	pageArticlesLimit, err := strconv.ParseInt(config.GetValue("DealsPerPage"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
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
	if pagenumber <= 0 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	pageArticlesLimit, err := strconv.ParseInt(config.GetValue("DealsPerPage"), 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
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
