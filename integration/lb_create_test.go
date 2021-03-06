package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os/exec"
	"strings"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
)

var _ = suite("compute/load-balancer/create", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect   *require.Assertions
		server   *httptest.Server
		cmd      *exec.Cmd
		baseArgs []string
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/load_balancers":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodPost {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				reqBody, err := ioutil.ReadAll(req.Body)
				expect.NoError(err)

				expect.JSONEq(lbCreateRequest, string(reqBody))

				w.Write([]byte(lbCreateResponse))
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
			"compute",
			"load-balancer",
		)

		baseArgs = []string{
			"--server-ids", "22,66",
			"--name", "my-lb-name",
			"--region", "venus",
			"--size", "lb-small",
			"--redirect-http-to-https",
			"--enable-proxy-protocol",
			"--enable-backend-keepalive",
			"--tag-name", "magic-lb",
			"--vpc-id", "1001",
		}
	})

	when("command is create", func() {
		it("creates a load balancer", func() {
			args := append([]string{"create"}, baseArgs...)
			cmd.Args = append(cmd.Args, args...)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(lbCreateOutput), strings.TrimSpace(string(output)))
		})
	})

	when("command is c", func() {
		it("creates a load balancer", func() {
			args := append([]string{"c"}, baseArgs...)
			cmd.Args = append(cmd.Args, args...)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(lbCreateOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	lbCreateOutput = `
ID      IP    Name             Status    Created At              Algorithm      Region    Size        VPC ID    Tag    Server IDs         SSL     Sticky Sessions                                Health Check                                                                                                            Forwarding Rules
1234          example-lb-01    new       2017-02-01T22:22:58Z    round_robin    syd       lb-small    1001             3164444,3164445    true    type:none,cookie_name:,cookie_ttl_seconds:0    protocol:,port:0,path:,check_interval_seconds:0,response_timeout_seconds:0,healthy_threshold:0,unhealthy_threshold:0
`
	lbCreateResponse = `
{
  "load_balancer": {
    "id": 1234,
    "name": "example-lb-01",
    "ip": "",
    "algorithm": "round_robin",
    "status": "new",
    "created_at": "2017-02-01T22:22:58Z",
    "forwarding_rules": [],
    "health_check": {},
    "sticky_sessions": {
      "type": "none"
    },
    "region": {
      "name": "Sydney",
      "slug": "syd",
      "sizes": [
        "std-min"
      ],
      "features": [
        "install_agent"
      ],
      "available": true
	},
	"size": "lb-small",
    "vpc_id": 1001,
    "tag": "",
    "server_ids": [
      3164444,
      3164445
    ],
    "redirect_http_to_https": true,
    "enable_proxy_protocol": true,
    "enable_backend_keepalive": true
  }
}`
	lbCreateRequest = `
{
  "name":"my-lb-name",
  "algorithm":"round_robin",
  "region":"venus",
  "size": "lb-small",
  "health_check":{},
  "sticky_sessions":{},
  "server_ids":[22,66],
  "tag":"magic-lb",
  "redirect_http_to_https":true,
  "enable_proxy_protocol":true,
  "enable_backend_keepalive":true,
  "vpc_id": 1001
}`
)
