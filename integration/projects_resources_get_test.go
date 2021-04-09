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

var _ = suite("projects/resources/get", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/servers/1111":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(serverGetResponse))
			case "/v2/floating_ips/1111":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(projectsResourcesGetFloatingIPResponse))
			case "/v2/load_balancers/1111":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(projectsResourcesGetLoadbalancerResponse))
			case "/v2/domains/1111":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(projectsResourcesGetDomainResponse))
			default:
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					t.Fatal("failed to dump request")
				}

				t.Fatalf("received unknown request: %s", dump)
			}
		}))
	})

	when("passing a server urn", func() {
		it("gets that resource for the project", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"projects",
				"resources",
				"get",
				"bl:server:1111",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(projectsResourcesGetServerOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passing a floatingip urn", func() {
		it("gets that resource for the project", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"projects",
				"resources",
				"get",
				"bl:floatingip:1111",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(projectsResourcesGetFloatingIPOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passing a loadbalancer urn", func() {
		it("gets that resource for the project", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"projects",
				"resources",
				"get",
				"bl:loadbalancer:1111",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(projectsResourcesGetLoadbalancerOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passing a domain urn", func() {
		it("gets that resource for the project", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"projects",
				"resources",
				"get",
				"bl:domain:1111",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(projectsResourcesGetDomainOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	projectsResourcesGetServerOutput = `
ID      Name                Public IPv4    Private IPv4    Public IPv6    Memory    VCPUs    Disk    Region              Image                          VPC ID    Status    Tags    Features    Volumes
5555    some-server-name                                                  0         0        0       some-region-slug    some-distro some-image-name              active    yes     remotes     some-volume-id
`
	projectsResourcesGetFloatingIPOutput = `
IP             Region    Server ID    Server Name
45.55.96.47    syd
`
	projectsResourcesGetFloatingIPResponse = `
{
  "floating_ip": {
    "ip": "45.55.96.47",
    "server": null,
    "region": {
      "name": "Sydney",
      "slug": "syd",
      "sizes": [ "std-min" ],
      "features": [ "metadata" ],
      "available": true
    },
    "locked": false
  }
}
`
	projectsResourcesGetLoadbalancerOutput = `
ID      IP                 Name             Status    Created At              Algorithm      Region    Size        VPC ID    Tag    Server IDs    SSL      Sticky Sessions                                Health Check                                                                                                            Forwarding Rules
1234    104.131.186.241    example-lb-01    new       2017-02-01T22:22:58Z    round_robin    syd       lb-small    1001             3164445       false    type:none,cookie_name:,cookie_ttl_seconds:0    protocol:,port:0,path:,check_interval_seconds:0,response_timeout_seconds:0,healthy_threshold:0,unhealthy_threshold:0    entry_protocol:https,entry_port:444,target_protocol:https,target_port:443,certificate_id:,tls_passthrough:true
`
	projectsResourcesGetLoadbalancerResponse = `
{
  "load_balancer": {
    "id": 1234,
    "name": "example-lb-01",
    "ip": "104.131.186.241",
    "algorithm": "round_robin",
    "status": "new",
    "created_at": "2017-02-01T22:22:58Z",
    "forwarding_rules": [
      {
        "entry_protocol": "https",
        "entry_port": 444,
        "target_protocol": "https",
        "target_port": 443,
        "certificate_id": "",
        "tls_passthrough": true
      }
    ],
    "health_check": {},
    "sticky_sessions": {
      "type": "none"
	},
	"size": "lb-small",
    "region": {
      "name": "Sydney",
      "slug": "syd",
      "sizes": [
        "std-min"
      ],
      "features": [ "install_agent" ],
      "available": true
    },
    "vpc_id": 1001,
    "server_ids": [ 3164445 ],
    "redirect_http_to_https": false,
    "enable_proxy_protocol": false
  }
}
`
	projectsResourcesGetDomainOutput = `
Domain         TTL
example.com    1800
`
	projectsResourcesGetDomainResponse = `
{
  "domain": {
    "name": "example.com",
    "ttl": 1800,
    "zone_file": "some zone file with crazy data"
  }
}
`
)
