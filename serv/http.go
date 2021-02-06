package serv

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxReadBytes       = 100000 // 100Kb
	introspectionQuery = "IntrospectionQuery"
	openVar            = "{{"
	closeVar           = "}}"
)

var (
	upgrader        = websocket.Upgrader{}
	errNoUserID     = errors.New("no user_id available")
	errUnauthorized = errors.New("not authorized")
)

type gqlReq struct {
	OpName string    `json:"operationName"`
	Query  string    `json:"query"`
	Vars   variables `json:"variables"`
}

type variables map[string]interface{}

type gqlResp struct {
	Error      string          `json:"error,omitempty"`
	Data       json.RawMessage `json:"data"`
	Extensions *extensions     `json:"extensions,omitempty"`
}

type extensions struct {
	Tracing *trace `json:"tracing,omitempty"`
}

type trace struct {
	Version   int           `json:"version"`
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Duration  time.Duration `json:"duration"`
	Execution execution     `json:"execution"`
}

type execution struct {
	Resolvers []resolver `json:"resolvers"`
}

type resolver struct {
	Path        []string      `json:"path"`
	ParentType  string        `json:"parentType"`
	FieldName   string        `json:"fieldName"`
	ReturnType  string        `json:"returnType"`
	StartOffset int           `json:"startOffset"`
	Duration    time.Duration `json:"duration"`
}

func apiv1Http(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if authFailBlock == authFailBlockAlways && authCheck(ctx) == false {
		http.Error(w, "Not authorized", 401)
		return
	}

	b, err := ioutil.ReadAll(io.LimitReader(r.Body, maxReadBytes))
	defer r.Body.Close()
	if err != nil {
		errorResp(w, err)
		return
	}

	req := &gqlReq{}
	if err := json.Unmarshal(b, req); err != nil {
		errorResp(w, err)
		return
	}

	if strings.EqualFold(req.OpName, introspectionQuery) {
		dat, err := ioutil.ReadFile("test.schema")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(dat)
		return
	}

	err = handleReq(ctx, w, req)

	if err == errUnauthorized {
		http.Error(w, "Not authorized", 401)
	}

	if err != nil {
		errorResp(w, err)
	}
}
