package agents

import (
	"github.com/emicklei/go-restful"
	"github.com/featen/ags/service/config"
	db "github.com/featen/ags/utils/db"
	log "github.com/featen/ags/utils/log"
	"net/http"
)

type Agent struct {
	Id     int64
	Status int64
	Info   string
}

const timeLayout = "2006-01-02 3:04pm"

var info db.InfoFetcher

func Init() {
	obj := db.InfoTable{Dbfile: config.GetValue("DbFile"), Tablename: "agents", Keyattrs: []string{"Title"}}
	info = obj
}

func InitTable() {
	log.Info("create table agents")
	info.CreateTable()
}

func Register() {
	log.Info("agents registered")

	ws := new(restful.WebService)
	ws.Path("/service/agents").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("").To(getAllAgents))
	ws.Route(ws.GET("/{agent-id}").To(findAgentById))
	ws.Route(ws.POST("").To(createAgent))

	restful.Add(ws)
}

func getAllAgents(req *restful.Request, resp *restful.Response) {
	log.Debug("get all agents")

	all, ret := info.FetchInfoRows(" t.status=1 ")
	if ret == http.StatusOK {
		resp.WriteEntity(all)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func findAgentById(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("agent-id")
	log.Debug("get agent by id %s", id)
	all, ret := info.FetchInfoRows("t.id=" + id)
	if ret == http.StatusOK && len(all) == 1 {
		resp.WriteEntity(all[0])
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func createAgent(req *restful.Request, resp *restful.Response) {
	agent := new(Agent)
	err := req.ReadEntity(&agent)
	if err == nil {
		id, ret := info.InsertRow(agent.Info)
		if ret == http.StatusOK {
			agent.Id = id
			resp.WriteHeader(http.StatusCreated)
			resp.WriteEntity(agent)
		} else {
			resp.WriteErrorString(ret, http.StatusText(ret))
		}
	} else {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}
