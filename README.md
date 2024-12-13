# GPT Assistant (GPTA)

GPT Assistant (`gpta`) is a command-line app that uses OpenAI to assist with tasks and automation.

## Usage

```bash
$ gpta -h
Usage:
  gpta [OPTIONS]

Application Options:
  -k, --key=         OpenAI API key [$OPENAI_API_KEY]
  -t, --task=        Task to perform [$GPTA_TASK]
  -i, --interactive  Interactive mode (ask for confirmation before executing commands)
  -s, --shell=       Shell to use for executing commands (default: /bin/sh) [$GPTA_SHELL]
  -m, --model=       GPT model to use (default: gpt-4o) [$GPTA_MODEL]
  -v, --verbose      Verbose output
  -c, --config=      Configuration file (default: ~/gpta.system) [$GPTA_CONFIG]

Help Options:
  -h, --help         Show this help message
```

## Example Usage

Executing a task with `gpta` is as simple as providing a task string to the `-t` flag:

```bash
% gpta -t "which ports are opened"
tcp        0      0 127.0.0.53:53           0.0.0.0:*               LISTEN

% gpta -t "what is my IP?" 
x.x.x.x

% gpta -t "is this makefile correct? do not execute it" makefile 
The structure of the Makefile looks generally correct for a Go project. It defines 'compile' and 'clean' targets, with 'go build -o' used to build the binary and 'rm' to remove it. However, ensure paths are correct and handle potential errors during removal.

% gpta -t "can you spot any spelling mistakes here?" < README.md 
Spelling mistakes spotted: 'opened' should be 'open', 'shebang' should be capitalized as 'Shebang', and 'Print' should be 'Prints'.

% gpta -vt "install netstat tool"
2024/12/13 05:23:59 Executing: cat ~/gpta.conf
No output
2024/12/13 05:23:59 Executing: sudo apt-get update -y
Hit:1 http://security.ubuntu.com/ubuntu jammy-security InRelease
Hit:2 http://de.archive.ubuntu.com/ubuntu jammy InRelease
Hit:3 http://de.archive.ubuntu.com/ubuntu jammy-updates InRelease
Hit:4 https://packages.microsoft.com/repos/edge stable InRelease
Hit:5 http://de.archive.ubuntu.com/ubuntu jammy-backports InRelease
Hit:6 https://packages.microsoft.com/repos/code stable InRelease
Reading package lists...
2024/12/13 05:24:02 Executing: sudo apt-get install -y net-tools
Reading package lists...
Building dependency tree...
Reading state information...
net-tools is already the newest version (1.60+git20181103.0eebece-1ubuntu5).
The following packages were automatically installed and are no longer required:
  libirs-export161 libisccfg-export163 libmujs1 libsixel1 mpv
  python3-pycryptodome python3-pyxattr rtmpdump
Use 'sudo apt autoremove' to remove them.
0 upgraded, 0 newly installed, 0 to remove and 111 not upgraded.
2024/12/13 05:24:03 Executing: echo "Netstat is already installed as part of the net-tools package."
Netstat is already installed as part of the net-tools package.
2024/12/13 05:24:03 Exiting with code: 0
```

`gpta` can also read tasks from a file or standard input:

```bash
$ echo "What is my user name?" | gpta -v
2024/12/12 21:26:14 Executing: echo $USER
user0
2024/12/12 21:26:14 Exiting with code: 0
```

`gpta` can also be a shebang interpreter in scripts:

```bash
#!/usr/local/bin/gpta -i

Close port 3000 on my system. Use sudo to elevate the privileges.
```

The [examples/](examples/) directory contains various .gpta and .sh files demonstrating how `gpta` can be used:

- **[check_updates.gpta](examples/check_updates.gpta)**: Checks for available updates on the system.
- **[disk_usage.gpta](examples/disk_usage.gpta)**: Print top 5 largest files in /var/log.
- **[explain_errors.sh](examples/explain_errors.sh)**: Pipe an error message into `gpta` and ask it to explain.
- **[close_port.gpta](examples/close_port.gpta)**: Close port 3000 with elevated privileges.
- **[http.gpta](examples/http.gpta)**: Download and display the first 10 lines of a remote file.

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
EOF
```

## Guarantees

No guarantees are provided with this software. Use at your own risk. I mean it. All responsibility is yours. Have fun!
