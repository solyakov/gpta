# GPT Assistant (GPTA)

GPT Assistant (`gpta`) is a command-line app that uses OpenAI to assist with tasks and automation.

## Usage

```bash
$ ./gpta -h
Usage:
  gpta [OPTIONS]

Application Options:
  -k, --key=         OpenAI API key [$OPENAI_API_KEY]
  -t, --task=        Task to perform [$GPTA_TASK]
  -i, --interactive  Interactive mode (ask for confirmation before executing commands)
  -s, --shell=       Shell to use for executing commands (default: /bin/sh) [$GPTA_SHELL]
  -m, --model=       GPT model to use (default: gpt-4o) [$GPTA_MODEL]
  -v, --verbose      Verbose output
  -c, --config=      Configuration file (default: ~/gpta.conf) [$GPTA_CONFIG]

Help Options:
  -h, --help         Show this help message
```

## Example Use Cases

The [examples/](examples/) directory contains various .gpta and .sh files demonstrating how `gpta` can be used:

- **[check_updates.gpta](examples/check_updates.gpta)**: Checks for available updates on the system.
- **[disk_usage.gpta](examples/disk_usage.gpta)**: Print top 5 largest files in /var/log.
- **[explain_errors.sh](examples/explain_errors.sh)**: Pipe an error message into `gpta` and ask it to explain.
- **[close_port.gpta](examples/close_port.gpta)**: Close port 3000 with elevated privileges.
- **[http.gpta](examples/http.gpta)**: Download and display the first 10 lines of a remote file.

and more...

## Example Usage

Executing a task with `gpta` is as simple as providing a task string to the `-t` flag:

```bash
% ./gpta -t "which ports are opened"
tcp        0      0 127.0.0.53:53           0.0.0.0:*               LISTEN

% gpta -t "what is my IP?" 
x.x.x.x
```

`gpta` can also read tasks from a file or standard input:

```bash
$ echo "What is my user name?" | ./gpta -v
2024/12/12 21:26:14 Executing: echo $USER
user0
2024/12/12 21:26:14 Exiting with code 0
```

`gpta` can also be a shebang interpreter in scripts:

```bash
#!../gpta -i

Close port 3000 on my system. Use sudo to elevate the privileges.
```

and more...

## Installation

Build the `gpta` binary by running `make`:

```bash
$ make
```

Set the `OPENAI_API_KEY` environment variable to your OpenAI API key. And ask `gpta` to install itself:

```bash
$ export OPENAI_API_KEY=...
$ cat <<EOF | ./gpta --interactive --verbose
Install the 'gpta' binary from the current directory into a directory included in my PATH (using sudo if necessary). If 'gpta' is already installed, update it by replacing the existing binary.

Update my shell profile (e.g., ~/.bashrc) to export the OPENAI_API_KEY environment variable:
    export OPENAI_API_KEY=$OPENAI_API_KEY

Create an empty configuration file at ~/gpta.conf if it does not exist. Do not overwrite an existing file.
EOF
```

## Guarantees

No guarantees are provided with this software. Use at your own risk. I mean it. All responsibility is yours. Have fun!
