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

var _ = suite("vpcs/list", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
		cmd    *exec.Cmd
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/vpcs":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(vpcsListResponse))
			default:
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					t.Fatal("failed to dump request")
				}

				t.Fatalf("received unknown request: %s", dump)
			}
		}))

		cmd = exec.Command(builtBinaryPath,
			"-t", "some-magic-token",
			"-u", server.URL,
			"vpcs",
		)

	})

	when("command is list", func() {
		it("lists all VPCs", func() {
			cmd.Args = append(cmd.Args, []string{"list"}...)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(vpcsListOutput), strings.TrimSpace(string(output)))
		})
	})

	when("command is ls", func() {
		it("lists all VPCs", func() {
			cmd.Args = append(cmd.Args, []string{"ls"}...)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(vpcsListOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	vpcsListOutput = `
ID      Name           Description    IP Range         Region    Created At                       Default
1001    my-new-vpc                    10.10.10.0/24    syd       2020-03-13 19:20:47 +0000 UTC    false
1002    default-syd                   10.102.0.0/20    syd       2020-03-13 19:29:20 +0000 UTC    true
1003    default-mel                   10.100.0.0/20    mel       2019-11-19 22:19:35 +0000 UTC    true
`
	vpcsListResponse = `
{
  "vpcs": [
    {
      "id": 1001,
      "urn": "bl:vpc:1001",
      "name": "my-new-vpc",
      "description": "",
      "region": "syd",
      "ip_range": "10.10.10.0/24",
      "created_at": "2020-03-13T19:20:47Z",
      "default": false
    },
    {
      "id": 1002,
      "urn": "bl:vpc:1002",
      "name": "default-syd",
      "description": "",
      "region": "syd",
      "ip_range": "10.102.0.0/20",
      "created_at": "2020-03-13T19:29:20Z",
      "default": true
    },
    {
      "id": 1003,
      "urn": "bl:vpc:1003",
      "name": "default-mel",
      "description": "",
      "region": "mel",
      "ip_range": "10.100.0.0/20",
      "created_at": "2019-11-19T22:19:35Z",
      "default": true
    }
  ],
  "links": {
  },
  "meta": {
    "total": 3
  }
}
`
)
