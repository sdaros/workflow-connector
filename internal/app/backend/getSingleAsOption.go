package backend

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/signavio/workflow-connector/internal/pkg/config"
	"github.com/signavio/workflow-connector/internal/pkg/formatting"
	"github.com/signavio/workflow-connector/internal/pkg/log"
	"github.com/signavio/workflow-connector/internal/pkg/query"
	"github.com/signavio/workflow-connector/internal/pkg/util"
)

func (b *Backend) GetSingleAsOption(rw http.ResponseWriter, req *http.Request) {
	log.When(config.Options.Logging).Infoln("[handler] GetSingleAsOption")
	routeName := mux.CurrentRoute(req).GetName()
	id := mux.Vars(req)["id"]
	table := req.Context().Value(util.ContextKey("table")).(string)
	uniqueIDColumn := req.Context().Value(util.ContextKey("uniqueIDColumn")).(string)
	columnAsOptionName := req.Context().Value(util.ContextKey("columnAsOptionName")).(string)
	queryUninterpolated := b.GetQueryTemplate(routeName)
	queryTemplate := &query.QueryTemplate{
		Vars: []string{queryUninterpolated},
		TemplateData: struct {
			TableName          string
			UniqueIDColumn     string
			ColumnAsOptionName string
		}{
			TableName:          table,
			UniqueIDColumn:     uniqueIDColumn,
			ColumnAsOptionName: columnAsOptionName,
		},
	}
	log.When(config.Options.Logging).Infof("[handler] %s", routeName)

	log.When(config.Options.Logging).Infoln("[handler -> backend] interpolate query string")
	queryString, err := queryTemplate.Interpolate()
	if err != nil {
		msg := &util.ResponseMessage{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}
		http.Error(rw, msg.Error(), http.StatusInternalServerError)
		return
	}
	log.When(config.Options.Logging).Infof("[handler <- backend]\n%s\n", queryString)

	log.When(config.Options.Logging).Infoln("[handler -> db] get query results")
	results, err := b.QueryContext(req.Context(), queryString, id)
	if err != nil {
		msg := &util.ResponseMessage{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}
		http.Error(rw, msg.Error(), http.StatusInternalServerError)
		return
	}
	log.When(config.Options.Logging).Infof("[handler <- db] query results: \n%s\n",
		results,
	)
	if len(results) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	log.When(config.Options.Logging).Infoln("[handler -> formatter] format results as json")
	formattedResults, err := formatting.WorkflowAccelerator.Format(req, results)
	if err != nil {
		msg := &util.ResponseMessage{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}
		http.Error(rw, msg.Error(), http.StatusInternalServerError)
		return
	}
	log.When(config.Options.Logging).Infof("[handler <- formatter] formatted results: \n%s\n",
		formattedResults,
	)

	rw.Write(formattedResults)
	return
}