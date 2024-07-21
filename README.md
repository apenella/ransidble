# Ransidble

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

![Ransidble-logo](docs/images/logo_ransidble.png "Ransidble logo")

Ransidble is a utility that enables you to execute [Ansible](https://www.ansible.com/) commands on remote hosts. It functions as a wrapper over Ansible, allowing users to launch playbooks, roles, and tasks remotely. This is achieved by exposing a REST API in front of Ansible, facilitating the execution of Ansible commands on remote hosts.

- [Ransidble](#ransidble)
  - [Why Ransidble?](#why-ransidble)
  - [How about the name?](#how-about-the-name)
  - [Server Usage Reference](#server-usage-reference)
    - [Configuration](#configuration)
    - [Initiate the Ransidble server](#initiate-the-ransidble-server)
  - [User Reference](#user-reference)
    - [Perform a Request to Execute an Ansible playbook](#perform-a-request-to-execute-an-ansible-playbook)
    - [Perform a Request Accepting Gzip Encoding](#perform-a-request-accepting-gzip-encoding)
    - [Get the Status of an Execution](#get-the-status-of-an-execution)
  - [REST API Reference](#rest-api-reference)
    - [Command](#command)
      - [Execution Data Object](#execution-data-object)
      - [Dependencies Object](#dependencies-object)
        - [Roles Dependencies Object](#roles-dependencies-object)
        - [Collections Dependencies Object](#collections-dependencies-object)
    - [Task](#task)
  - [Development Reference](#development-reference)
    - [Contributing](#contributing)
    - [Code of Conduct](#code-of-conduct)
    - [Roadmap](#roadmap)
  - [Acknowledgements](#acknowledgements)
  - [License](#license)

## Why Ransidble?

## How about the name?

Ransidble is a blend of 'Remote' and 'Ansible,' with a nod to the punk rock band [Rancid](https://rancidrancid.com/).

## Server Usage Reference

### Configuration

The Ransidble server can be configured using environment variables. The following table lists the available environment variables:

| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| RANSIDBLE_SERVER_HTTP_LISTEN_ADDRESS | The port where the server listens for incoming requests | :8080 |
| RANSIDBLE_SERVER_LOG_LEVEL | The log level for the server | info |
| RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH | The path where the projects are stored | projects |
| RANSIDBLE_SERVER_PROJECT_STORAGE_TYPE | The type of storage used to store the projects | local |
| RANSIDBLE_SERVER_WORKER_POOL_SIZE | The number of workers to execute the commands | 1 |

Ransidble can be also configured using a configuration file. In this case, the file must be named `ransidble.yaml` and placed in the same directory as the binary. Environment variables take precedence over the configuration file.
The following is an example of a configuration file:

```yaml
server:
  http_listen_address: ":8080"
  log_level: info
  worker_pool_size: 5
  project:
    local_storage_path: projects
    storage_type: local
```

### Initiate the Ransidble server

```bash
RANSIDBLE_SERVER_LOG_LEVEL=info RANSIDBLE_SERVER_WORKER_POOL_SIZE=3 RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH=test/projects  go run cmd/main.go serve
{"data":null,"level":"info","msg":"Starting server on :8080","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Starting worker 4838bac7-3efe-4a9d-837d-3e61924c5f35","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Starting worker 50498ab1-72b0-4b8f-a197-a83b527ec874","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Starting worker 7f05137a-7701-4e88-8aee-1b8fe65d2e80","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Executing task ecfc92dc-6323-40d6-9bf8-71c4d4d98640","time":"2024-07-21T20:01:19+02:00"}
{"time":"2024-07-21T20:01:19.534469135+02:00","id":"","remote_ip":"127.0.0.1","host":"0.0.0.0:8080","method":"POST","uri":"/task/ansible-playbook/project-1","user_agent":"curl/7.81.0","status":202,"error":"","latency":518245,"latency_human":"518.245µs","bytes_in":77,"bytes_out":46}
{"data":null,"level":"info","msg":"Setup project project-1 to /tmp/ransidble846872083/4838bac7-3efe-4a9d-837d-3e61924c5f35/project-1/ecfc92dc-6323-40d6-9bf8-71c4d4d98640","time":"2024-07-21T20:01:19+02:00"}
[DEPRECATION WARNING]: ANSIBLE_COLLECTIONS_PATHS option, does not fit var
naming standard, use the singular form ANSIBLE_COLLECTIONS_PATH instead. This
feature will be removed from ansible-core in version 2.19. Deprecation warnings
 can be disabled by setting deprecation_warnings=False in ansible.cfg.

PLAY [all] *********************************************************************

TASK [wait for 5 seconds] ******************************************************
Pausing for 5 seconds
(ctrl+C then 'C' = continue early, ctrl+C then 'A' = abort)
ok: [127.0.0.1]

TASK [ansibleplaybook-simple] **************************************************
ok: [127.0.0.1] =>
  msg: Your are running 'ansibleplaybook-simple' example

PLAY RECAP *********************************************************************
127.0.0.1                  : ok=2    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

## User Reference

### Perform a Request to Execute an Ansible playbook

The following example demonstrates how to execute an Ansible playbook using the Ransidble server. Please refer to the [REST API Reference](#rest-api-reference) section for more information.

```bash
curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/task/ansible-playbook/project-1 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'
HTTP/1.1 202 Accepted
Content-Type: application/json
Vary: Accept-Encoding
Date: Sun, 21 Jul 2024 18:01:19 GMT
Content-Length: 46

{"id":"ecfc92dc-6323-40d6-9bf8-71c4d4d98640"}
```

### Perform a Request Accepting Gzip Encoding

```bash
❯ curl -H 'Accept-Encoding: gzip' --output /dev/stdout -s -XPOST 0.0.0.0:8080/command/ansible-playbook
r+LIU(ILVHHM.-IMQ(.MNN-.N+'NyI&
```

### Get the Status of an Execution

```bash
❯ curl -s -GET 0.0.0.0:8080/task/ecfc92dc-6323-40d6-9bf8-71c4d4d98640 | jq
{
  "command": "ansible-playbook",
  "completed_at": "2024-07-21T20:01:25+02:00",
  "created_at": "2024-07-21T20:01:19+02:00",
  "executed_at": "2024-07-21T20:01:19+02:00",
  "id": "ecfc92dc-6323-40d6-9bf8-71c4d4d98640",
  "parameters": {
    "playbooks": [
      "site.yml"
    ],
    "inventory": "127.0.0.1,",
    "connection": "local"
  },
  "project": {
    "name": "project-1",
    "reference": "test/projects/project-1",
    "type": "local"
  },
  "status": "SUCCESS"
}
```

## REST API Reference

### Command

#### Execution Data Object

The JSON schema payload to execute an Ansible playbook command is defined [here](api/schemas/input/command/ansible-playbook-parameters.json). However, the following table provides a summary of the available attributes:

| JSON Attribute | Type | Description |
|----------------|------|-------------|
| playbooks | []string | The ansible's playbooks list to be executed |
| check | bool | Don't make any changes; instead, try to predict some of the changes that may occur |
| diff | bool | When changing (small) files and templates, show the differences in those files; works great with --check |
| dependencies | object | A list of role and collection dependencies. The object is described in the section [Dependencies object](#dependencies-object) |
| extra_vars | map[string]interface{} | A map of extra variables used on ansible-playbook execution |
| extra_vars_file | []string | A list of files used to load extra-vars |
| flush_cache | bool | The flush cache flag for ansible-playbook |
| force_handlers | bool | Run handlers even if a task fails |
| forks | int | Specify number of parallel processes to use (default=50) |
| inventory | string | Specify inventory host path |
| limit | string | Selected hosts additional pattern |
| list_hosts | bool | Outputs a list of matching hosts |
| list_tags | bool | The list tags flag for ansible-playbook |
| list_tasks | bool | The list tasks flag for ansible-playbook |
| skip_tags | string | Only run plays and tasks whose tags do not match these values |
| start_at_task | string | Start the playbook at the task matching this name |
| syntax_check | bool | The syntax check flag for ansible-playbook |
| tags | string | The tags flag for ansible-playbook |
| vault_id | string | The vault identity to use |
| vault_password_file | string | Path to the file holding vault decryption key |
| verbose | bool | Verbose mode enabled |
| version | bool | Show program's version number, config file location, configured module search path, module location, executable location and exit |
| connection | string | The type of connection used by ansible-playbook |
| scp_extra_args | string | Specify extra arguments to pass to scp only |
| sftp_extra_args | string | Specify extra arguments to pass to sftp only |
| ssh_common_args | string | Specify common arguments to pass to sftp/scp/ssh |
| ssh_extra_args | string | Specify extra arguments to pass to ssh only |
| timeout | int | The connection timeout on ansible-playbook |
| user | string | The user to use to connect to a host |
| become | bool | Ansible-playbook's become flag |
| become_method | string | Ansible-playbook's become method |
| become_user | string | Ansible-playbook's become user |

#### Dependencies Object

| JSON Attribute | Type | Description |
|----------------|------|-------------|
| roles | object | Defines how to install roles dependencies. The object is described in the section [Roles dependencies object](#roles-dependencies-object) |
| collections | object | Defines how to install collections dependencies. The object is described in the section [Collections dependencies object](#collections-dependencies-object) |

##### Roles Dependencies Object

| JSON Attribute | Type | Description |
|----------------|------|-------------|
| roles | []string | A list of roles to install. |
| api_key | string | The API key to use to authenticate against the galaxy server. Same as --token. |
| ignore_errors | bool | Whether to continue processing even if a role fails to install. |
| no_deps | bool | Whether to install dependencies. |
| role_file | string | The path to a file containing a list of roles to install. |
| server | string | The flag to specify the galaxy server to use. |
| timeout | string | The time to wait for operations against the galaxy server, defaults to 60s. |
| token | string | The token to use to authenticate against the galaxy server. Same as --api-key. |
| verbose | bool | Verbose mode enabled. |

##### Collections Dependencies Object

| JSON Attribute | Type | Description |
|----------------|------|-------------|
| collections | []string | A list of collections to install. |
| api_key | string | The Ansible Galaxy API key. |
| force_with_deps | bool | Forces overwriting an existing collection and its dependencies. |
| pre | bool | Includes pre-release versions. Semantic versioning pre-releases are ignored by default. |
| timeout | string | The time to wait for operations against the galaxy server, defaults to 60s. |
| token | string | The Ansible Galaxy API key. |
| ignore_errors | bool | Ignores errors during installation and continue with the next specified collection. |
| requirements_file | string | A file containing a list of collections to be installed. |
| server | string | The Galaxy API server URL. |
| verbose | bool | Verbose mode enabled. |

### Task

## Development Reference

### Contributing

Thank you for your interest in contributing to Ransidble. All contributions are welcome, whether they are bug reports, feature requests, or code contributions!
Please, read the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information.

### Code of Conduct

The Ransidble project is committed to providing a friendly, safe and welcoming environment for all, regardless of gender, sexual orientation, disability, ethnicity, religion, or similar personal characteristic.

We expect all contributors, users, and community members to follow this code of conduct. This includes all interactions within the Ransidble community, whether online, in person, or otherwise.

Please, read the [CODE-OF-CONDUCT.md](CODE-OF-CONDUCT.md) file for more information.

### Roadmap

The roadmap is available in the [ROADMAP.md](ROADMAP.md) file.

## Acknowledgements

## License

Ransidble is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
