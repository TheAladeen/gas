package agents

import (
	"github.com/emicklei/go-restful"
	"github.com/featen/ags/service/config"
	db "github.com/featen/utils/db"
	log "github.com/featen/utils/log"
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

	allagents, ret := info.FetchInfoRows(" t.status=1 ")
	if ret == http.StatusOK {
		resp.WriteEntity(allagents)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func findAgentById(req *restful.Request, resp *restful.Response) {
	agent := new(Agent)
	id := req.PathParameter("agent-id")

	allgents, ret := info.FetchInfoRows("t.id=" + id)

	if ret == http.StatusOK && len(allgents) == 1 {
		agent.Id = allgents[0].Id
		agent.Status = allgents[0].Status
		agent.Info = allgents[0].Info

		resp.WriteEntity(agent)
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

