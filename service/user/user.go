package user

import (
	"net/http"

	"github.com/featen/gas/service/config"

	"github.com/emicklei/go-restful"
	"github.com/featen/gas/service/auth"
	log "github.com/featen/gas/utils/log"
)

const timeLayout = "2006-01-02 3:04pm"

type VerifyInfo struct {
	NameOrEmail, Passwd string
}

func Register() {
	ws := new(restful.WebService)
	ws.
		Path("/service/user").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/").To(currUser).Filter(auth.AuthFilter))
	ws.Route(ws.POST("/signin").To(signin))
	ws.Route(ws.GET("/signout").To(signout).Filter(auth.AuthFilter))

	restful.Add(ws)

	log.Info("user registered! ")
}

func currUser(req *restful.Request, resp *restful.Response) {
	resp.WriteEntity(config.GetValue("AdminEmail"))
}

func checkAdminUser(vi *VerifyInfo) int {
	if (vi.NameOrEmail == config.GetValue("AdminName") ||
		vi.NameOrEmail == config.GetValue("AdminEmail")) &&
		vi.Passwd == config.GetValue("AdminPassword") {
		return http.StatusOK
	} else {
		return http.StatusForbidden
	}
}

func signin(req *restful.Request, resp *restful.Response) {
	vi := new(VerifyInfo)
	err := req.ReadEntity(&vi)
	if err == nil {
		if ret := checkAdminUser(vi); ret == http.StatusOK {
			auth.AddCookie(req.Request, resp.ResponseWriter)
			resp.WriteEntity(config.GetValue("AdminEmail"))
		} else {
			resp.WriteErrorString(ret, http.StatusText(ret))
		}
	} else {
		resp.WriteError(http.StatusBadRequest, err)
	}
}

func signout(req *restful.Request, resp *restful.Response) {
	auth.DelCookie(req, resp)
	resp.WriteHeader(http.StatusOK)
}
