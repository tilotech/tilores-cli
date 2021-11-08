package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	atomic.StoreUint64(&serverPID, 0)

	dir, err := createTempDir()
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = initializeProject([]string{})
	require.NoError(t, err)

	go func() {
		err := startWebserver()
		require.EqualError(t, err, "an error occurred while waiting on server process: signal: killed", "expected test to kill the web server after the tests are done")
	}()
	time.Sleep(1 * time.Second)
	defer syscall.Kill(-int(atomic.LoadUint64(&serverPID)), syscall.SIGKILL)
	jsonData := map[string]string{
		"query": `
            {__typename}
        `,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(jsonValue))
	request.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)
	client := &http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	actual := map[string]interface{}{}
	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)
	expected := map[string]interface{}{
		"data": map[string]interface{}{"__typename": "Query"},
	}
	assert.Equal(t, expected, actual)
}
