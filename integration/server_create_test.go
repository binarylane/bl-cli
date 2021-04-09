package integration

import (
	"encoding/json"
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

type serverRequest struct {
	Name string `json:"name"`
}

var _ = suite("compute/server/create", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect  *require.Assertions
		server  *httptest.Server
		reqBody []byte
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/servers":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodPost {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				var err error
				reqBody, err = ioutil.ReadAll(req.Body)
				expect.NoError(err)

				var dr serverRequest
				err = json.Unmarshal(reqBody, &dr)
				expect.NoError(err)

				if dr.Name == "waiting-on-name" {
					w.Write([]byte(serverCreateWaitResponse))
					return
				}

				w.Write([]byte(serverCreateResponse))
			case "/poll-for-server":
				w.Write([]byte(actionCompletedResponse))
			case "/v2/servers/777":
				// we don't really need another fake server here
				// since we've successfully tested all the behavior
				// at this point
				w.Write([]byte(serverCreateResponse))
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
		it("creates a server", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"create",
				"some-server-name",
				"--image", "a-test-image",
				"--region", "a-test-region",
				"--size", "a-test-size",
				"--vpc-id", "1001",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(serverCreateOutput), strings.TrimSpace(string(output)))

			request := &struct {
				Name   string
				Image  string
				Region string
				Size   string
				VPCID  int `json:"vpc_id,float64"`
			}{}

			err = json.Unmarshal(reqBody, request)
			expect.NoError(err)

			expect.Equal("some-server-name", request.Name)
			expect.Equal("a-test-image", request.Image)
			expect.Equal("a-test-region", request.Region)
			expect.Equal("a-test-size", request.Size)
			expect.Equal(1001, request.VPCID)
		})
	})

	when("the wait flag is passed", func() {
		it("polls until the server is created", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"create",
				"waiting-on-name",
				"--wait",
				"--image", "a-test-image",
				"--region", "a-test-region",
				"--size", "a-test-size",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
		})
	})

	when("missing required arguments", func() {
		base := []string{
			"-t", "some-magic-token",
			"-u", "https://www.example.com",
			"compute",
			"server",
			"create",
		}

		baseErr := `Error: (server.create%s) command is missing required arguments`

		cases := []struct {
			desc string
			err  string
			args []string
		}{
			{desc: "missing all", err: fmt.Sprintf(baseErr, ""), args: base},
			{desc: "missing only name", err: fmt.Sprintf(baseErr, ""), args: append(base, []string{"--size", "test", "--region", "test", "--image", "test"}...)},
			{desc: "missing only region", err: fmt.Sprintf(baseErr, ".region"), args: append(base, []string{"some-name", "--size", "test", "--image", "test"}...)},
			{desc: "missing only size", err: fmt.Sprintf(baseErr, ".size"), args: append(base, []string{"some-name", "--image", "test", "--region", "test"}...)},
			{desc: "missing only image", err: fmt.Sprintf(baseErr, ".image"), args: append(base, []string{"some-name", "--size", "test", "--region", "test"}...)},
		}

		for _, c := range cases {
			commandArgs := c.args
			expectedErr := c.err

			when(c.desc, func() {
				it("returns an error", func() {
					cmd := exec.Command(builtBinaryPath, commandArgs...)

					output, err := cmd.CombinedOutput()
					expect.Error(err)
					expect.Contains(string(output), expectedErr)
				})
			})
		}
	})
})

const (
	serverCreateResponse = `
{
  "server": {
    "id": 1111,
    "memory": 12,
    "vcpus": 13,
    "disk": 15,
    "name": "some-server-name",
    "networks": {
      "v4": [
        {"type": "public", "ip_address": "1.2.3.4"},
        {"type": "private", "ip_address": "7.7.7.7"}
      ]
    },
    "image": {
      "distribution": "some-distro",
      "name": "some-image-name"
    },
    "region": {
      "slug": "some-region-slug"
    },
	"status": "active",
	"vpc_id": 1001,
    "tags": ["yes"],
    "features": ["remotes"],
    "volume_ids": ["some-volume-id"]

  }
}`
	serverCreateWaitResponse = `
{"server": {"id": 777}, "links": {"actions": [{"id":1, "rel":"create", "href":"poll-for-server"}]}}
`
	actionCompletedResponse = `
{"action": "id": 1, "status": "completed"}
`
	serverCreateOutput = `
ID      Name                Public IPv4    Private IPv4    Public IPv6    Memory    VCPUs    Disk    Region              Image                          VPC ID    Status    Tags    Features    Volumes
1111    some-server-name    1.2.3.4        7.7.7.7                        12        13       15      some-region-slug    some-distro some-image-name    1001      active    yes     remotes     some-volume-id
`
)
