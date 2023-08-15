package testbench

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goto/salt/log"
	"github.com/julienschmidt/httprouter"
)

func startMockServer(ctx context.Context, logger *log.Zap, port int) {
	var (
		internalServerErrorWriter = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}
	)
	router := httprouter.New()
	router.GET("/api/ping", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("pong"))
	})
	router.POST("/api/resource", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			internalServerErrorWriter(w)
			return
		}

		var reqBody map[string]string
		if err := json.Unmarshal(b, &reqBody); err != nil {
			internalServerErrorWriter(w)
			return
		}

		var orgName = ""
		if hOrg, ok := r.Header["X-Shield-Org"]; ok {
			orgName = hOrg[0]
		}

		reqBody["org"] = orgName
		reqBody["urn"] = reqBody["name"]

		respBytes, err := json.Marshal(reqBody)
		if err != nil {
			internalServerErrorWriter(w)
			return
		}

		w.Write(respBytes)

		// jsonResponse :=
		// 	resource:
		// 	key: urn
		// 	type: json_payload
		// 	source: response
		//   project:
		// 	key: project
		// 	type: json_payload
		// 	source: request
		//   group:
		// 	key: group
		// 	type: json_payload
		// 	source: request
		//   resource_type:
		// 	value: "firehose"
		// 	type: constant
	})

	logger.Info("starting up mock server...", "port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err.Error())
	}
}
