package dict

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/featen/gas/service/config"
	log "github.com/featen/gas/utils/log"
	"html/template"
	"net/http"
)

type Client struct {
	BaseURL string
	Keyfrom string
	Key     string
}

var (
	YoudaoBaseURL = "http://fanyi.youdao.com/"
	qs            = make(map[string]string)
	client        *Client
)

func Register() {
	log.Info("youdao dict registered")
	ws := new(restful.WebService)
	ws.Path("/service/dict").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.GET("/{cond}").To(genDataByCond))
	restful.Add(ws)

	client = NewClient()
}

func NewClient() *Client {
	return &Client{
		BaseURL: YoudaoBaseURL,
		Keyfrom: config.GetValue("youdaoKeyfrom"),
		Key:     config.GetValue("youdaoKey"),
	}
}

type Result struct {
	ErrorCode   int      `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
	Basic       *struct {
		Phonetic string   `json:"phonetic"`
		Explains []string `json:"explains"`
	} `json:"basic"`
	Web []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"web"`
}

func (c *Client) Query(q string) (*Result, error) {
	return c.QueryHttp(http.DefaultClient, q)
}

func (c *Client) QueryHttp(httpClient *http.Client, q string) (*Result, error) {
	requestURL := fmt.Sprintf(
		"%sopenapi.do?keyfrom=%s&key=%s&type=data&doctype=json&version=1.1&q=%s",
		c.BaseURL, template.URLQueryEscaper(c.Keyfrom),
		template.URLQueryEscaper(c.Key), template.URLQueryEscaper(q))
	fmt.Println(requestURL)

	resp, err := httpClient.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var res Result
	err = dec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func insertYoudaodict(q string, fanyi string) {
	dbHandler, err := sql.Open("sqlite3", config.GetValue("DbFile"))
	if err != nil {
		log.Fatal("%v", err)
		fmt.Println("dbHandler failed", err)
	}
	defer dbHandler.Close()

	s := "insert or replace into dict (q, fanyi) values (?, ?)"
	_, err = dbHandler.Exec(s, q, fanyi)
	if err != nil {
		log.Fatal("%q: %s\n", err, s)
	}

	qs[q] = fanyi
}

func genDataByCond(req *restful.Request, resp *restful.Response) {
	cond := req.PathParameter("cond")

	r, ret := dbGenDataByCond(cond)
	if ret == http.StatusOK {
		resp.WriteEntity(r)
	} else {
		resp.WriteErrorString(ret, http.StatusText(ret))
	}
}

func getResult(s string) (*Result, int) {
	var r Result

	err := json.Unmarshal([]byte(s), &r)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	return &r, http.StatusOK
}

func dbGenDataByCond(cond string) (*Result, int) {
	f, ok := qs[cond]
	if ok {
		return getResult(f)
	}

	dbHandler, err := sql.Open("sqlite3", config.GetValue("DbFile"))
	if err != nil {
		log.Fatal("%v", err)
		return nil, http.StatusInternalServerError
	}
	defer dbHandler.Close()

	querySql := "select fanyi from dict where q=? limit 1"
	var fanyi sql.NullString
	err = dbHandler.QueryRow(querySql, cond).Scan(&fanyi)
	if err == nil {
		return getResult(fanyi.String)
	}

	res, err := client.Query(cond)
	if err != nil {
		log.Error("%v", err)
		return nil, http.StatusInternalServerError
	}
	r, err := json.Marshal(res)
	if err != nil {
		log.Error("json marshal failed: %v", err)
	} else {
		if res.ErrorCode == 0 {
			insertYoudaodict(cond, string(r))
		}
	}

	return res, http.StatusOK
}
