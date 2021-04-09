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

var _ = suite("compute/floating-ip/list", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		cmd    *exec.Cmd
		server *httptest.Server
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/floating_ips":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(floatingIPListResponse))
			default:
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					t.Fatal("failed to dump request")
				}

				t.Fatalf("received unknown request: %s", dump)
			}
		}))

	})

	when("required flags are passed", func() {
		it("lists all floating-ips", func() {
			aliases := []string{"list", "ls"}

			for _, alias := range aliases {
				cmd = exec.Command(builtBinaryPath,
					"-t", "some-magic-token",
					"-u", server.URL,
					"compute",
					"floating-ip",
					alias,
				)

				output, err := cmd.CombinedOutput()
				expect.NoError(err, fmt.Sprintf("received error output: %s", output))
				expect.Equal(strings.TrimSpace(floatingIPListOutput), strings.TrimSpace(string(output)))
			}
		})
	})
})

const (
	floatingIPListOutput = `
IP         Region    Server ID    Server Name
8.8.8.8    syd       8888         hello
1.1.1.1    syd       1111
`
	floatingIPListResponse = `
{
  "floating_ips": [
    {
      "ip": "8.8.8.8",
      "server": {"id": 8888, "name": "hello"},
      "region": {
        "name": "Sydney",
        "slug": "syd",
        "sizes": [ "std-min" ],
        "features": [ "metadata" ],
        "available": true
      },
      "locked": false
    },
    {
      "ip": "1.1.1.1",
      "server": {"id": 1111},
      "region": {
        "name": "Sydney",
        "slug": "syd",
        "sizes": [ "std-min" ],
        "features": [ "metadata" ],
        "available": true
      },
      "locked": false
    }
  ],
  "links": {},
  "meta": {
    "total": 2
  }
}
`
)
