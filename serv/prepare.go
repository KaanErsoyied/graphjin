package serv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dosco/super-graph/psql"
	"github.com/dosco/super-graph/qcode"
	"github.com/go-pg/pg"
	"github.com/valyala/fasttemplate"
)

type preparedItem struct {
	stmt    *pg.Stmt
	args    [][]byte
	skipped uint32
	qc      *qcode.QCode
}

var (
	_preparedList map[string]*preparedItem
)

func initPreparedList() {
	_preparedList = make(map[string]*preparedItem)

	for k, v := range _allowList.list {
		err := prepareStmt(k, v.gql, v.vars)
		if err != nil {
			panic(err)
		}
	}
}

func prepareStmt(key, gql string, varBytes json.RawMessage) error {
	if len(gql) == 0 || len(key) == 0 {
		return nil
	}

	qc, err := qcompile.Compile([]byte(gql))
	if err != nil {
		return err
	}

	var vars map[string]json.RawMessage

	if len(varBytes) != 0 {
		vars = make(map[string]json.RawMessage)

		if err := json.Unmarshal(varBytes, &vars); err != nil {
			return err
		}
	}

	buf := &bytes.Buffer{}

	skipped, err := pcompile.Compile(qc, buf, psql.Variables(vars))
	if err != nil {
		return err
	}

	t := fasttemplate.New(buf.String(), `{{`, `}}`)
	am := make([][]byte, 0, 5)
	i := 0

	finalSQL := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		am = append(am, []byte(tag))
		i++
		return w.Write([]byte(fmt.Sprintf("$%d", i)))
	})

	if err != nil {
		return err
	}

	pstmt, err := db.Prepare(finalSQL)
	if err != nil {
		return err
	}

	_preparedList[key] = &preparedItem{
		stmt:    pstmt,
		args:    am,
		skipped: skipped,
		qc:      qc,
	}

	return nil
}
