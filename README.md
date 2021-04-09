<h1 align="center">bl-cli</h1>

<p align="center">
  <a href="https://travis-ci.org/binarylane/bl-cli">
    <img src="https://travis-ci.org/binarylane/bl-cli.svg?branch=master" alt="Build Status" />
  </a>
  <a href="https://godoc.org/github.com/binarylane/bl-cli">
    <img src="https://godoc.org/github.com/binarylane/bl-cli?status.svg" alt="GoDoc" />
  </a>
  <a href="https://goreportcard.com/report/github.com/binarylane/bl-cli">
    <img src="https://goreportcard.com/badge/github.com/binarylane/bl-cli" alt="Go Report Card" />
  </a>
</p>

```
bl is a command line interface (CLI) for the BinaryLane API.

Usage:
  bl [command]

Available Commands:
  account         Display commands that retrieve account details
  auth            Display commands for authenticating bl with an account
  balance         Display commands for retrieving your account balance
  billing-history Display commands for retrieving your billing history
  completion      Modify your shell so bl commands autocomplete with TAB
  compute         Display commands that manage infrastructure
  help            Help about any command
  invoice         Display commands for retrieving invoices for your account
  projects        Manage projects and assign resources to them
  version         Show the current version
  vpcs            Display commands that manage VPCs

Flags:
  -t, --access-token string   API V2 access token
  -u, --api-url string        Override default API endpoint
  -c, --config string         Specify a custom config file (default "$HOME/.config/bl/config.yaml")
      --context string        Specify a custom authentication context name
  -h, --help                  help for bl
  -o, --output string         Desired output format [text|json] (default "text")
      --trace                 Show a log of network activity while performing a command

Use "bl [command] --help" for more information about a command.
```

See the [full reference documentation](https://api.binarylane.com.au/reference/) for information about each available command.

- [bl-cli](#bl-cli---)
    - [Installing `bl-cli`](#installing-bl-cli)
        - [Downloading a Release from GitHub](#downloading-a-release-from-github)
        - [Building with Docker](#building-with-docker)
        - [Building the Development Version from Source](#building-the-development-version-from-source)
            - [Dependencies](#dependencies)
    - [Authenticating with BinaryLane](#authenticating-with-binarylane)
        - [Logging in to multiple BinaryLane accounts](#logging-in-to-multiple-binarylane-accounts)
    - [Configuring Default Values](#configuring-default-values)
    - [Enabling Shell Auto-Completion](#enabling-shell-auto-completion)
        - [Linux](#linux-auto-completion)
        - [macOS](#macos-auto-completion)
    - [Examples](#examples)
    - [bl-cli Releases](https://github.com/binarylane/bl-cli/releases)


## Installing `bl-cli`

### Downloading a Release from GitHub

Visit the [Releases
page](https://github.com/binarylane/bl-cli/releases) for the
[`bl-cli` GitHub project](https://github.com/binarylane/bl-cli), and find the
appropriate archive for your operating system and architecture.
Download the archive from your browser or copy its URL and
retrieve it to your home directory with `wget` or `curl`.

For example, with `wget`:

```
cd ~
wget https://github.com/binarylane/bl-cli/releases/download/v<version>/bl-cli-<version>-linux-amd64.tar.gz
```

Or with `curl`:

```
cd ~
curl -OL https://github.com/binarylane/bl-cli/releases/download/v<version>/bl-cli-<version>-linux-amd64.tar.gz
```

Extract the binary:

```
tar xf ~/bl-cli-<version>-linux-amd64.tar.gz
```

Or download and extract with this oneliner:
```
curl -sL https://github.com/binarylane/bl-cli/releases/download/v<version>/bl-cli-<version>-linux-amd64.tar.gz | tar -xzv
```

where `<version>` is the full semantic version, e.g., `0.1.0`.

On Windows systems, you should be able to double-click the zip archive to extract the `bl-cli` executable.

Move the `bl` binary to somewhere in your path. For example, on GNU/Linux and OS X systems:

```
sudo mv ~/bl /usr/local/bin
```

Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `bl` to their `PATH`.

### Building with Docker

If you have Docker configured, you can build a local Docker image using `bl-cli`'s
[Dockerfile](https://github.com/binarylane/bl-cli/blob/master/Dockerfile)
and run `bl-cli` within a container.

```
docker build --tag=bl-cli .
```

Then you can run it within a container.

```
docker run --rm --interactive --tty --env=BINARYLANE_ACCESS_TOKEN="your_BL_token" bl any_command
```

### Building the Development Version from Source

If you have a Go environment configured, you can install the development version of `bl-cli` from
the command line.

```
go get github.com/binarylane/bl-cli/cmd/bl
```

While the development version is a good way to take a peek at
`bl-cli`'s latest features before they get released, be aware that it
may have bugs. Officially released versions will generally be more
stable.

### Dependencies

`bl-cli` uses Go modules with vendoring.

## Authenticating with BinaryLane

To use `bl-cli`, you need to authenticate with BinaryLane by providing an access token, which can be created from the [Developer API](https://home.binarylane.com.au/api-info) section of the Control Panel.

Docker users will have to use the `BINARYLANE_ACCESS_TOKEN` environmental variable to authenticate, as explained in the Installation section of this document.

If you're not using Docker to run `bl-cli`, authenticate with the `auth init` command.

```
bl auth init
```

You will be prompted to enter the BinaryLane access token that you generated in the BinaryLane control panel.

```
BinaryLane access token: your_access_token
```

After entering your token, you will receive confirmation that the credentials were accepted. If the token doesn't validate, make sure you copied and pasted it correctly.

```
Validating token: OK
```

This will create the necessary directory structure and configuration file to store your credentials.

### Logging into multiple BinaryLane accounts

`bl-cli` allows you to log in to multiple BinaryLane accounts at the same time and easily switch between them with the use of authentication contexts.

By default, a context named `default` is used. To create a new context, run `bl auth init --context <new-context-name>`. You may also pass the new context's name using the `BINARYLANE_CONTEXT` environment variable. You will be prompted for your API access token which will be associated with the new context.

To use a non-default context, pass the context name to any `bl` command. For example:

```
bl compute server list --context <new-context-name>
```

To set a new default context, run `bl auth switch --context <new-context-name>`. This command will save the current context to the config file and use it for all commands by default if a context is not specified.

The `--access-token` flag or `BINARYLANE_ACCESS_TOKEN` variable are acknowledged only if the `default` context is used. Otherwise, they will have no effect on what API access token is used. To temporarily override the access token if a different context is set as default, use `bl --context default --access-token your_access_token ...`.

## Configuring Default Values

The `bl-cli` configuration file is used to store your API Access Token as well as the defaults for command flags. If you find yourself using certain flags frequently, you can change their default values to avoid typing them every time. This can be useful when, for example, you want to change the username or port used for SSH.

On OS X, `bl-cli` saves its configuration as `${HOME}/Library/Application Support/bl/config.yaml`. The `${HOME}/Library/Application Support/bl/` directory will be created once you run `bl auth init`.

On Linux, `bl-cli` saves its configuration as `${XDG_CONFIG_HOME}/bl/config.yaml` if the `${XDG_CONFIG_HOME}` environmental variable is set, or `~/.config/bl/config.yaml` if it is not. On Windows, the config file location is `%APPDATA%\bl\config.yaml`.

The configuration file is automatically created and populated with default properties when you authenticate with `bl-cli` for the first time. The typical format for a property is `category.command.sub-command.flag: value`. For example, the property for the `force` flag with tag deletion is `tag.delete.force`.

To change the default SSH user used when connecting to a server with `bl-cli`, look for the `compute.ssh.ssh-user` property and change the value after the colon. In this example, we changed it to the username **admin**.

```
. . .
compute.ssh.ssh-user: admin
. . .
```

Save and close the file. The next time you use `bl-cli`, the new default values you set will be in effect. In this example, that means that it will SSH as the **admin** user (instead of the default **root** user) next time you log into a server.

## Enabling Shell Auto-Completion

`bl-cli` also has auto-completion support. It can be set up so that if you partially type a command and then press `TAB`, the rest of the command is automatically filled in. For example, if you type `bl comp<TAB><TAB> ser<TAB><TAB>` with auto-completion enabled, you'll see `bl-cli compute server` appear on your command prompt.

**Note:** Shell auto-completion is not available for Windows users.

`bl-cli` can generate an auto-completion script with the `bl completion your_shell_here` command. Valid arguments for the shell are Bash (`bash`) and ZSH (`zsh`). By default, the script will be printed to the command line output.  For more usage examples for the `completion` command, use `bl completion --help`.

### Linux Auto Completion

The most common way to use the `completion` command is by adding a line to your local profile configuration. At the end of your `~/.profile` file, add this line:

```
source <(bl completion your_shell_here)
```

Then refresh your profile.

```
source ~/.profile
```

### MacOS

macOS users will have to install the `bash-completion` framework to use the auto-completion feature.

```
brew install bash-completion
```

After it's installed, load `bash_completion` by adding the following line to your `.profile` or `.bashrc`/`.zshrc` file.

```
source $(brew --prefix)/etc/bash_completion
```

Then refresh your profile using the appropriate command for the bash configurations file.

```
source ~/.profile
source ~/.bashrc
source ~/.zshrc
```

## Examples

`bl-cli` is able to interact with all of your BinaryLane resources. Below are a few common usage examples. 

* List all servers on your account:
```
bl compute server list
```
* Create a server:
```
bl compute server create <name> --region <region-slug> --image <image-slug> --size <size-slug>
```
* Create a new A record for an existing domain:
```
bl compute domain records create --record-type A --record-name www --record-data <ip-addr> <domain-name>
```

`bl-cli` also simplifies actions without an API endpoint. For instance, it allows you to SSH to your server by name:
```
bl compute ssh <server-name>
```

By default, it assumes you are using the `root` user. If you want to SSH as a specific user, you can do that as well:
```
bl compute ssh <user>@<server-name>
```
