package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	cases := map[string]struct {
		port                   int
		changeSchema           bool
		testQuery              string
		expectedServerResponse map[string]interface{}
	}{
		"run server with port flag set": {
			port:         8081,
			changeSchema: false,
			expectedServerResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"__type": map[string]interface{}{
						"name": "Record",
						"fields": []interface{}{
							map[string]interface{}{
								"name": "id",
							},
							map[string]interface{}{
								"name": "myCustomField",
							},
						},
					},
				},
			},
		},
		"run server with the same schema": {
			changeSchema: false,
			expectedServerResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"__type": map[string]interface{}{
						"name": "Record",
						"fields": []interface{}{
							map[string]interface{}{
								"name": "id",
							},
							map[string]interface{}{
								"name": "myCustomField",
							},
						},
					},
				},
			},
		},
		"run server after changing the schema": {
			changeSchema: true,
			// expecting the query name to become find instead of search after go generate
			expectedServerResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"__type": map[string]interface{}{
						"name": "Record",
						"fields": []interface{}{
							map[string]interface{}{
								"name": "id",
							},
							map[string]interface{}{
								"name": "myCustomField",
							},
							map[string]interface{}{
								"name": "myNewField",
							},
						},
					},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			atomic.StoreUint64(&serverPID, 0)

			dir, err := createTempDir()
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			err = initializeProject([]string{})
			require.NoError(t, err)

			defer os.Unsetenv("TILORES_PORT")
			port = c.port

			if c.changeSchema {
				err = changeQuerySchema()
				require.NoError(t, err)
			}

			go func() {
				err := runGraphQLServer()
				require.EqualError(t, err, "an error occurred while waiting on server process: signal: killed", "expected test to kill the web server after the tests are done")
			}()
			defer killWebserver()
			jsonData := map[string]string{
				"query": `{__type(name: "Record"){name,fields{name}}}`,
			}
			jsonValue, _ := json.Marshal(jsonData)
			var url string
			if port == 0 {
				url = "http://localhost:8080/query"
			} else {
				url = "http://localhost:" + strconv.Itoa(port) + "/query"
			}
			request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
			request.Header.Add("Content-Type", "application/json")
			require.NoError(t, err)
			data := requestServerUntilTimeout(t, request)
			actual := map[string]interface{}{}
			err = json.Unmarshal(data, &actual)
			require.NoError(t, err)
			assert.Equal(t, c.expectedServerResponse, actual)
		})
	}
}

func changeQuerySchema() error {
	schemaFile, err := os.Create("schema/record.graphqls")
	if err != nil {
		return err
	}

	_, err = schemaFile.WriteString(`
input RecordInput {
	id: ID!
  myCustomField: String!
	myNewField: String!
}

type Record {
	id: ID!
  myCustomField: String!
	myNewField: String!
}
`)
	if err != nil {
		return err
	}

	err = schemaFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func killWebserver() {
	pid := atomic.LoadUint64(&serverPID)
	if pid != 0 {
		_ = syscall.Kill(-int(pid), syscall.SIGKILL)
	}
}

func requestServerUntilTimeout(t *testing.T, req *http.Request) []byte {
	futureTime := time.Now().Add(30 * time.Second)
	client := &http.Client{}
	var response *http.Response
	for {
		var err error
		time.Sleep(500 * time.Millisecond)
		response, err = client.Do(req)
		if time.Now().After(futureTime) {
			require.Fail(t, "no successful response from server within 30 seconds")
		}
		if err == nil {
			break
		}
	}
	data, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	_ = response.Body.Close()
	return data
}
