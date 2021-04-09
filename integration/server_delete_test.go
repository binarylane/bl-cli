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

var _ = suite("compute/server/delete", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
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
				if req.URL.RawQuery == "page=1&per_page=200&tag_name=one" {
					if req.Method != http.MethodGet {
						w.WriteHeader(http.StatusMethodNotAllowed)
						return
					}

					w.Write([]byte(`{"servers":[{"name":"some-server-name", "id": 1337}]}`))
				} else if req.URL.RawQuery == "tag_name=one" {
					if req.Method == http.MethodDelete {
						w.WriteHeader(http.StatusNoContent)
					} else {
						w.WriteHeader(http.StatusMethodNotAllowed)
						return
					}
				} else if req.URL.RawQuery == "page=1&per_page=200&tag_name=two" {
					if req.Method != http.MethodGet {
						w.WriteHeader(http.StatusMethodNotAllowed)
						return
					}

					w.Write([]byte(`{"servers":[{"name":"some-server-name", "id": 1337}, {"name":"another-server-name", "id": 7331}]}`))
				} else if req.URL.RawQuery == "tag_name=two" {
					if req.Method == http.MethodDelete {
						w.WriteHeader(http.StatusNoContent)
					} else {
						w.WriteHeader(http.StatusMethodNotAllowed)
						return
					}
				} else {
					if req.Method != http.MethodGet {
						w.WriteHeader(http.StatusMethodNotAllowed)
						return
					}

					w.Write([]byte(`{"servers":[{"name":"some-server-name", "id": 1337}]}`))
				}
			case "/v2/servers/1337":
				if req.Method != http.MethodDelete {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.WriteHeader(http.StatusNoContent)
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
		base := []string{
			"-t", "some-magic-token",
			"compute",
			"server",
		}

		cases := []struct {
			desc string
			args []string
		}{
			{desc: "command is delete", args: append(base, []string{"delete", "some-server-name", "--force"}...)},
			{desc: "command is rm", args: append(base, []string{"rm", "some-server-name", "--force"}...)},
			{desc: "command is d", args: append(base, []string{"d", "some-server-name", "--force"}...)},
			{desc: "command is del", args: append(base, []string{"del", "some-server-name", "--force"}...)},
		}

		for _, c := range cases {
			commandArgs := c.args

			when(c.desc, func() {
				it("deletes a server", func() {
					finalArgs := append([]string{"-u", server.URL}, commandArgs...)
					cmd := exec.Command(builtBinaryPath, finalArgs...)

					output, err := cmd.CombinedOutput()
					expect.NoError(err, fmt.Sprintf("received error output: %s", output))
					expect.Empty(output)
				})
			})
		}
	})

	when("deleting by tag name", func() {
		it("deletes the right Server", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"delete",
				"--tag-name", "one",
				"--force",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
			expect.Empty(output)
		})
	})

	when("deleting one Server without force flag", func() {
		it("correctly promts for confirmation", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"delete",
				"some-server-name",
			)

			output, err := cmd.CombinedOutput()
			expect.Error(err)
			expect.Equal(strings.TrimSpace(serverDelOutput), strings.TrimSpace(string(output)))
		})
	})

	when("deleting two Server without force flag", func() {
		it("correctly promts for confirmation", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"delete",
				"some-server-name",
				"another-server-name",
			)

			output, err := cmd.CombinedOutput()
			expect.Error(err)
			expect.Equal(strings.TrimSpace(multiServerDelOutput), strings.TrimSpace(string(output)))
		})
	})

	when("deleting one Server by tag without force flag", func() {
		it("correctly promts for confirmation", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"delete",
				"--tag-name", "one",
			)

			output, err := cmd.CombinedOutput()
			expect.Error(err)
			expect.Equal(strings.TrimSpace(tagServerDelOutput), strings.TrimSpace(string(output)))
		})
	})

	when("deleting two Server by tag without force flag", func() {
		it("correctly promts for confirmation", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"server",
				"delete",
				"--tag-name", "two",
			)

			output, err := cmd.CombinedOutput()
			expect.Error(err)
			expect.Equal(strings.TrimSpace(tagMultiServerDelOutput), strings.TrimSpace(string(output)))
		})
	})
})

const (
	serverDelOutput         = "Warning: Are you sure you want to delete this Server? (y/N) ? Error: Operation aborted."
	multiServerDelOutput    = "Warning: Are you sure you want to delete 2 Servers? (y/N) ? Error: Operation aborted."
	tagServerDelOutput      = `Warning: Are you sure you want to delete 1 Server tagged "one"? [affected Server: 1337] (y/N) ? Error: Operation aborted.`
	tagMultiServerDelOutput = `Warning: Are you sure you want to delete 2 Servers tagged "two"? [affected Servers: 1337 7331] (y/N) ? Error: Operation aborted.`
)
