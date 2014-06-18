package agents

import (
	"database/sql"
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

var obj = db.InfoTable{Dbfile: config.GetValue("DbFile"), Tablename: "agents"}
var info db.InfoFetcher = obj

func InitTable() {
	obj.CreateTable()
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
	info = obj

	allagents, ret := info.FetchInfoRows("t.Status=1")
	if ret == http.StatusOK {
		resp.WriteEntity(allagents)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func findAgentById(req *restful.Request, resp *restful.Response) {
	agent := new(Agent)
	id := req.PathParameter("agent-id")

	info = obj
	allgents, ret := info.FetchInfoRows("t.Id=" + id)

	if ret == http.StatusOK {
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
		ret := dbCreateAgent(agent)
		if ret == http.StatusOK {
			resp.WriteHeader(http.StatusCreated)
			resp.WriteEntity(agent)
		} else {
			resp.WriteErrorString(ret, http.StatusText(ret))
		}
	} else {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}

func dbCreateAgent(agent *Agent) int {
	log.Debug("try to create agent %v", agent)

	dbHandler, err := sql.Open("sqlite3", config.GetValue("DbFile"))
	if err != nil {
		log.Fatal("%v", err)
	}
	defer dbHandler.Close()

	stmt, err := dbHandler.Prepare("INSERT INTO agents (info) VALUES (?)")
	if err != nil {
		log.Error("%v", err)
		return http.StatusInternalServerError
	}
	defer stmt.Close()

	r, err := stmt.Exec(agent.Info)
	if err != nil {
		log.Error("%v", err)
		return http.StatusBadRequest
	}
	id, _ := r.LastInsertId()
	agent.Id = id

	return http.StatusOK
}
