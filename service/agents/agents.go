package agents

import (
	"database/sql"
	"github.com/emicklei/go-restful"
	"github.com/featen/ags/service/config"
	log "github.com/featen/utils/log"
	"net/http"
	"strconv"
)

type Agent struct {
	Id string
    Status int64
    Info string
}

const timeLayout = "2006-01-02 3:04pm"

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
	log.Debug("get all articles")
	dbHandler, err := sql.Open("sqlite3", config.GetValue("DbFile"))
	if err != nil {
		log.Fatal("%v", err)
        resp.WriteErrorString(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
        return
	}
	defer dbHandler.Close()

	stmt, err := dbHandler.Prepare("SELECT id, info FROM agents ORDER BY a.id DESC")
	if err != nil {
		log.Error("%v", err)
        resp.WriteErrorString(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
        return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("%v", err)
        resp.WriteErrorString(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer rows.Close()

	allagents := make([]Agent, 0)
	for rows.Next() {
		var info sql.NullString
		var id sql.NullInt64
		rows.Scan(&id, &info)

		allagents = append(allagents, Agent{strconv.FormatInt(id.Int64, 10), 0, info.String})
	}
    resp.WriteEntity(allagents)
}

func findAgentById(req *restful.Request, resp *restful.Response) {
	agent := new(Agent)
	agent.Id = req.PathParameter("agent-id")

    log.Debug("try to find agent with id : %v", agent.Id)
	dbHandler, err := sql.Open("sqlite3", config.GetValue("DbFile"))
	if err != nil {
		log.Fatal("%v", err)
        return
	}
	defer dbHandler.Close()

	stmt, err := dbHandler.Prepare("SELECT info FROM agents WHERE id = ? ")
	if err != nil {
		log.Error("%v", err)
        resp.WriteErrorString(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer stmt.Close()

	var info sql.NullString
	err = stmt.QueryRow(agent.Id).Scan(&info)
	if err != nil {
		log.Error("%v", err)
		if err == sql.ErrNoRows {
           resp.WriteErrorString(http.StatusInternalServerError, http.StatusText(http.StatusNotFound))
			return
		} else {
           resp.WriteErrorString(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	if !info.Valid {
        resp.WriteErrorString(http.StatusInternalServerError, http.StatusText(http.StatusNotFound))
		return
	} else {
		agent.Info = info.String
	}

	resp.WriteEntity(agent)
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
	agent.Id = strconv.FormatInt(id, 10)

	return http.StatusOK
}


