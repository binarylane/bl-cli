package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os/exec"
	"strings"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
)

var _ = suite("compute/server/actions", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/servers/1111/actions":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(serverActionsResponse))
			default:
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					t.Fatal("failed to dump request")
				}

				t.Fatalf("received unknown request: %s", dump)
			}
		}))
	})

	when("all required flags are passed", func() {
		it("lists server actions", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"actions",
				"1111",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(serverActionsOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	serverActionsOutput = `
ID    Status       Type    Started At                       Completed At                     Resource ID    Resource Type    Region
2     completed            2014-11-14 16:37:39 +0000 UTC    2014-11-14 16:37:40 +0000 UTC    0              server
`
	serverActionsResponse = `
{
    "actions": [
      {
        "id": 2,
        "slug": "silly",
        "started_at": "2014-11-14T16:37:39Z",
        "completed_at":  "2014-11-14T16:37:40Z",
        "status": "completed",
        "resource_type": "server"
      }
    ]
}`
)
