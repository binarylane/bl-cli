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

var _ = suite("compute/load-balancer/list", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
		cmd    *exec.Cmd
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

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(lbListResponse))
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

	})

	when("command is list", func() {
		it("lists all load balancers", func() {
			cmd.Args = append(cmd.Args, []string{"list"}...)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(lbListOutput), strings.TrimSpace(string(output)))
		})
	})

	when("command is ls", func() {
		it("lists all load balancers", func() {
			cmd.Args = append(cmd.Args, []string{"ls"}...)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(lbListOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	lbListOutput = `
ID      IP                 Name             Status    Created At              Algorithm      Region    Size         VPC ID    Tag    Server IDs    SSL      Sticky Sessions                                Health Check                                                                                                                   Forwarding Rules
2001    104.131.186.241    example-lb-01    new       2017-02-01T22:22:58Z    round_robin    venus3    lb-small     1001             3164444       false    type:none,cookie_name:,cookie_ttl_seconds:0    protocol:http,port:80,path:/,check_interval_seconds:10,response_timeout_seconds:5,healthy_threshold:5,unhealthy_threshold:3    entry_protocol:http,entry_port:80,target_protocol:http,target_port:80,certificate_id:,tls_passthrough:false
2002    104.131.188.204    example-lb-02    new       2017-02-01T20:44:58Z    round_robin    mars1     lb-medium    1001             3164445       false    type:none,cookie_name:,cookie_ttl_seconds:0    protocol:http,port:80,path:/,check_interval_seconds:10,response_timeout_seconds:5,healthy_threshold:5,unhealthy_threshold:3    entry_protocol:http,entry_port:80,target_protocol:http,target_port:80,certificate_id:,tls_passthrough:false
`
	lbListResponse = `
{
  "load_balancers": [
    {
      "id": 2001,
      "name": "example-lb-01",
      "ip": "104.131.186.241",
      "algorithm": "round_robin",
      "status": "new",
      "created_at": "2017-02-01T22:22:58Z",
      "forwarding_rules": [
        {
          "entry_protocol": "http",
          "entry_port": 80,
          "target_protocol": "http",
          "target_port": 80,
          "certificate_id": "",
          "tls_passthrough": false
        }
      ],
      "health_check": {
        "protocol": "http",
        "port": 80,
        "path": "/",
        "check_interval_seconds": 10,
        "response_timeout_seconds": 5,
        "healthy_threshold": 5,
        "unhealthy_threshold": 3
      },
      "sticky_sessions": {
        "type": "none"
      },
      "size": "lb-small",
      "region": {
        "name": "Venus",
        "slug": "venus3",
        "sizes": ["std-min"],
        "features": ["private_networking"],
        "available": true
      },
      "vpc_id": 1001,
      "tag": "",
      "server_ids": [3164444],
      "redirect_http_to_https": false,
      "enable_proxy_protocol": false
    },
    {
      "id": 2002,
      "name": "example-lb-02",
      "ip": "104.131.188.204",
      "algorithm": "round_robin",
      "status": "new",
      "created_at": "2017-02-01T20:44:58Z",
      "forwarding_rules": [
        {
          "entry_protocol": "http",
          "entry_port": 80,
          "target_protocol": "http",
          "target_port": 80,
          "certificate_id": "",
          "tls_passthrough": false
        }
      ],
      "health_check": {
        "protocol": "http",
        "port": 80,
        "path": "/",
        "check_interval_seconds": 10,
        "response_timeout_seconds": 5,
        "healthy_threshold": 5,
        "unhealthy_threshold": 3
      },
      "sticky_sessions": {
        "type": "none"
      },
      "size": "lb-medium",
      "region": {
        "name": "Mars",
        "slug": "mars1",
        "sizes": ["std-min"],
        "features": ["install_agent"],
        "available": true
      },
      "vpc_id": 1001,
      "tag": "",
      "server_ids": [3164445],
      "redirect_http_to_https": false,
      "enable_proxy_protocol": false
    }
  ],
  "links": {
  },
  "meta": {
    "total": 2
  }
}
`
)
