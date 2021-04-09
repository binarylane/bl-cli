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

var _ = suite("vpcs/get", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/vpcs/1234":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(vpcsGetResponse))
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
		it("gets the specified VPC", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"vpcs",
				"get",
				"1234",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(vpcsGetOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passing the format flag", func() {
		it("changes the output", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"vpcs",
				"get",
				"1234",
				"--format", "Description",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(vpcsGetFormattedOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	vpcsGetOutput = `
ID      Name          Description        IP Range         Region    Created At                       Default
1234    my-new-vpc    vpc description    10.10.10.0/24    syd       2020-03-13 19:20:47 +0000 UTC    false
`
	vpcsGetFormattedOutput = `
Description
vpc description
`
	vpcsGetResponse = `
{
  "vpc": {
    "id": 1234,
    "urn": "bl:vpc:1234",
    "name": "my-new-vpc",
    "description": "vpc description",
    "region": "syd",
    "ip_range": "10.10.10.0/24",
    "created_at": "2020-03-13T19:20:47Z",
    "default": false
  }
}
`
)
