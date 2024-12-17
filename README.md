# Ansible Task Code Generator (atcg)

### Build and Test Status
[![Build Status](https://github.com/kbcz1989/atcg/actions/workflows/release.yml/badge.svg)](https://github.com/kbcz1989/atcg/actions/workflows/release.yml)
[![Tests](https://img.shields.io/github/actions/workflow/status/kbcz1989/atcg/tests.yml?label=tests)](https://github.com/kbcz1989/atcg/actions/workflows/tests.yml)

### Versioning and Downloads
![Go Version](https://img.shields.io/github/go-mod/go-version/kbcz1989/atcg)
![Latest Release](https://img.shields.io/github/v/release/kbcz1989/atcg)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/kbcz1989/atcg)
![Downloads](https://img.shields.io/github/downloads/kbcz1989/atcg/latest/total)

### Code Quality and Metrics
[![Go Report Card](https://goreportcard.com/badge/github.com/kbcz1989/atcg)](https://goreportcard.com/report/github.com/kbcz1989/atcg)
[![Codecov](https://codecov.io/gh/kbcz1989/atcg/branch/main/graph/badge.svg)](https://codecov.io/gh/kbcz1989/atcg)

![Build Size](https://img.shields.io/github/languages/code-size/kbcz1989/atcg)
![Languages](https://img.shields.io/github/languages/top/kbcz1989/atcg)

### Community and Contributions
![Open Issues](https://img.shields.io/github/issues/kbcz1989/atcg)
![Contributors](https://img.shields.io/github/contributors/kbcz1989/atcg)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)

### License and Miscellaneous
![License](https://img.shields.io/github/license/kbcz1989/atcg)
![Markdown Style](https://img.shields.io/badge/markdown-friendly-yellow)
![Powered by Go](https://img.shields.io/badge/powered%20by-Go-blue?logo=go)
![Last Commit](https://img.shields.io/github/last-commit/kbcz1989/atcg)

A Go-based tool for dynamically generating Ansible tasks and playbooks by leveraging `ansible-doc` and user-defined module configurations.

## Features

- Parses `ansible-doc` JSON outputs for Ansible modules.
- Dynamically generates:
  - Task files for specified modules.
  - A `main.yml` playbook to include and loop over tasks.
- Supports flexible configuration of module parameters using defaults and overrides.
- Ensures clean, reusable Ansible playbooks with proper structure and formatting.

## Requirements

- **Ansible CLI**: Required for running `ansible-doc`.
- **Go 1.23+**: For building from source.

## Installation

### Linux:

```bash
ARCH=$(uname -m | grep -q 'aarch64' && echo 'arm64' || echo 'amd64')
sudo wget "https://github.com/kbcz1989/atcg/releases/latest/download/atcg-linux-$ARCH" -O /usr/local/bin/atcg
sudo chmod +x /usr/local/bin/atcg
```

### macOS:

```shell
ARCH=$(uname -m | grep -q 'arm64' && echo 'arm64' || echo 'amd64')
curl -L "https://github.com/kbcz1989/atcg/releases/latest/download/atcg-darwin-$ARCH" -o /usr/local/bin/atcg
chmod +x /usr/local/bin/atcg
```

### Windows:

```shell
$ARCH = if ($ENV:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
Invoke-WebRequest -Uri "https://github.com/kbcz1989/atcg/releases/latest/download/atcg-windows-$ARCH.exe" -OutFile "$Env:LOCALAPPDATA\atcg.exe" -UseBasicParsing
```

## Installation from Source

Clone the repository:

```bash
git clone https://github.com/kbcz1989/atcg.git
cd atcg
```

Build the project:

```bash
go build -o atcg ./cmd/atcg
```

Run the binary:

```bash
./atcg --help
```

## Usage

### Command-Line Options

| Flag            | Description                                              | Example                             |
| --------------- | -------------------------------------------------------- | ----------------------------------- |
| `--module, -m`  | Specify Ansible modules to generate tasks for.           | `-m ansible.windows.win_user_right` |
| `--output, -o`  | The output directory for generated tasks and `main.yml`. | `-o ./tasks`                        |
| `--help, -h`    | Show usage information.                                  |                                     |
| `--version, -v` | Show app version.                                        |                                     |

### Output Files

- **Task Files**: One task file per module (e.g., `win_user_right.yml`).
- **`main.yml`**: Includes and loops over the generated tasks.

## Tests

Run all tests:

```bash
make test
```

### Key Test Features

- Verifies proper parsing of `ansible-doc` outputs.
- Ensures tasks and `main.yml` generation match expected structure.
- Validates CLI flag parsing and output consistency.

## How `atcg` Works

1. **Parsing Module Documentation**
   - The tool takes Ansible modules specified via CLI flags.
   - It retrieves the module documentation by running `ansible-doc` in JSON mode.
   - The documentation includes the module’s attributes, descriptions, defaults, and requirements.

2. **Generating Individual Task Files**
   - For each specified module, `atcg` creates a dedicated task file (e.g., `win_user_right.yml`).
   - The task file includes:
     - All module attributes as task parameters.
     - Default values where applicable.
     - Conditional `omit` for optional parameters without defaults.

3. **Building the `main.yml` Playbook**
   - After generating individual task files, `atcg` creates a `main.yml` playbook.
   - The playbook includes:
     - A task for each module file using `ansible.builtin.include_tasks`.
     - Loops over variables corresponding to each module.
     - Tags to organize and filter tasks during execution.

4. **Flexible Output Directory**
   - All generated files are stored in the specified output directory (default: `./tasks`).
   - The directory contains:
     - One task file per module.
     - A `main.yml` playbook to include and manage the tasks.

5. **Ensuring Clean and Reusable Playbooks**
   - The generated task files use structured and reusable patterns.
   - Tasks dynamically handle defaults, optional parameters, and required fields.

---

## Real-World Scenario

Imagine you are tasked with writing an Ansible playbook to configure multiple resources, such as user rights in Windows or SSL certificates. Crafting such playbooks manually can be tedious because:

1. **You need to understand the module documentation in detail**: Modules often support dozens of parameters, with varying requirements and defaults.
2. **Ensuring flexibility**: To make your playbooks reusable, you must account for configurable parameters and sensible defaults.
3. **Maintaining consistency**: Large playbooks can become inconsistent if not structured properly.

This is where `atcg` can help.

### How `atcg` Solves This Problem

`atcg` generates Ansible task files and a playbook by leveraging `ansible-doc` outputs. Its approach is guided by best practices for creating modular, reusable, and dynamic playbooks:

- **Comprehensive parameter coverage**: All attributes supported by the module are included, ensuring you don’t miss any critical options.
- **Baseline for further customization**: The generated tasks serve as a starting point, which you can refine and extend to meet specific requirements.
- **Dynamic playbook generation**: The `main.yml` playbook automatically includes and loops over the generated task files, organizing tasks with tags.

For example, if you need to configure user rights on a Windows system, instead of manually writing tasks for the `ansible.windows.win_user_right` module, you can run `atcg` and instantly get:

1. A task file (`win_user_right.yml`) with all supported parameters.
2. A `main.yml` playbook to include the task and loop over your input variables.
3. Flexibility to add more modules and regenerate tasks and playbooks as needed.

This approach saves time, ensures consistency, and allows you to focus on high-value tasks like customizing logic and deploying infrastructure.

## Example Workflow

#### 1. Specify Modules

Use `atcg` to specify Ansible modules you want tasks for.

#### 2. Review Generated Tasks and Playbook

Navigate to the output directory (`./tasks`):

- Review individual task files (e.g., `win_user_right.yml`).
- Check the generated `main.yml` for proper inclusion and looping.

#### 3. Integrate into Your Ansible Project

Copy the files into your Ansible project directory and use them in your playbooks.

## Contributing

Contributions are welcome! Feel free to fork the repository and submit pull requests.

## Acknowledgments

Special thanks to:

- [Ansible](https://www.ansible.com/) for its amazing automation platform.
- [Go](https://golang.org) community for the fantastic libraries used in this project.
- **ChatGPT** for guiding and assisting in creating this tool and documentation.