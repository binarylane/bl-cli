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

var _ = suite("compute/snapshot/get", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/snapshots/53344211":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(snapshotGetServerResponse))
			case "/v2/snapshots/0a343fac-eacf-11e9-b96b-0a58ac144633":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.Write([]byte(snapshotGetVolumeResponse))
			default:
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					t.Fatal("failed to dump request")
				}

				t.Fatalf("received unknown request: %s", dump)
			}
		}))
	})

	when("passed a server ID", func() {
		it("displays the server snapshot", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"snapshot",
				"get",
				"53344211",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(snapshotGetServerOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passed a volume ID", func() {
		it("displays the volume snapshot", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"snapshot",
				"get",
				"0a343fac-eacf-11e9-b96b-0a58ac144633",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(snapshotGetVolumeOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passing a format", func() {
		it("displays only those columns", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"snapshot",
				"get",
				"--format", "ID,ResourceType",
				"53344211",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(snapshotGetFormatOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passing no-header", func() {
		it("displays only values, no headers", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"snapshot",
				"get",
				"--no-header",
				"0a343fac-eacf-11e9-b96b-0a58ac144633",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(snapshotGetNoHeaderOutput), strings.TrimSpace(string(output)))
		})
	})

	when("passing a format and no header together", func() {
		it("displays only the value with no header", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"snapshot",
				"get",
				"--format", "ID",
				"--no-header",
				"53344211",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Equal(strings.TrimSpace(snapshotGetFormatNoHeaderOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	snapshotGetServerResponse = `
{
  "snapshot": {
    "id": "53344211",
    "name": "ubuntu-s-1vcpu-1gb-syd-01-1570651077842",
    "regions": [
      "syd"
    ],
    "created_at": "2019-10-09T19:57:59Z",
    "resource_id": "162347943",
    "resource_type": "server",
    "min_disk_size": 25,
    "size_gigabytes": 1.01,
    "tags": []
  }
}
`
	snapshotGetVolumeResponse = `
{
  "snapshot": {
    "id": "0a343fac-eacf-11e9-b96b-0a58ac144633",
    "name": "volume-syd-01-1570651053836",
    "regions": [
      "syd"
    ],
    "created_at": "2019-10-09T19:57:36Z",
    "resource_id": "e2068b37-eace-11e9-85ad-0a58ac14430f",
    "resource_type": "volume",
    "min_disk_size": 100,
    "size_gigabytes": 0,
    "tags": []
  }
}
`
	snapshotGetVolumeOutput = `
ID                                      Name                           Created at              Regions    Resource ID                             Resource Type    Min Disk Size    Size        Tags
0a343fac-eacf-11e9-b96b-0a58ac144633    volume-syd-01-1570651053836    2019-10-09T19:57:36Z    [syd]      e2068b37-eace-11e9-85ad-0a58ac14430f    volume           100              0.00 GiB
`
	snapshotGetServerOutput = `
ID          Name                                       Created at              Regions    Resource ID    Resource Type    Min Disk Size    Size        Tags
53344211    ubuntu-s-1vcpu-1gb-syd-01-1570651077842    2019-10-09T19:57:59Z    [syd]      162347943      server           25               1.01 GiB
`
	snapshotGetFormatOutput = `
ID          Resource Type
53344211    server
`
	snapshotGetNoHeaderOutput = `
0a343fac-eacf-11e9-b96b-0a58ac144633    volume-syd-01-1570651053836    2019-10-09T19:57:36Z    [syd]    e2068b37-eace-11e9-85ad-0a58ac14430f    volume    100    0.00 GiB
`
	snapshotGetFormatNoHeaderOutput = `
53344211
`
)
