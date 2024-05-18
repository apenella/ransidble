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
    - [Perform a request to execute an Ansible playbook](#perform-a-request-to-execute-an-ansible-playbook)
    - [REST API Reference](#rest-api-reference)
      - [Command](#command)
        - [Execution Data Object](#execution-data-object)
        - [Dependencies Object](#dependencies-object)
          - [Roles Dependencies Object](#roles-dependencies-object)
          - [Collections Dependencies Object](#collections-dependencies-object)
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
| RANSIDBLE_WORKER_POOL_SIZE | The number of workers to execute the commands | 1 |
| RANSIDBLE_HTTP_LISTEN_ADDRESS | The port where the server listens for incoming requests | :8080 |
| RANSIDBLE_LOG_LEVEL | The log level for the server | info |

Ransidble can be also configured using a configuration file. In this case, the file must be named `ransidble.yaml` and placed in the same directory as the binary. Environment variables take precedence over the configuration file.
The following is an example of a configuration file:

```yaml
http_listen_address: ":8080"
log_level: "info"
worker_pool_size: 5
```

### Initiate the Ransidble server

```bash
❯ RANSIDBLE_WORKER_POOL_SIZE=5 go run cmd/main.go serve
{"data":null,"level":"info","msg":"Starting server on :8080","time":"2024-05-15T17:57:50+02:00"}
{"data":null,"level":"info","msg":"Starting worker ecb8cc75-322d-49af-b7dd-391556ef2fb4","time":"2024-05-15T17:57:50+02:00"}
{"data":null,"level":"info","msg":"Starting worker 5b41bc79-aefb-44f6-9980-52e11c27a5da","time":"2024-05-15T17:57:50+02:00"}
{"data":null,"level":"info","msg":"Starting worker 66a3b4f7-b875-44eb-8f70-0bb79ca86a09","time":"2024-05-15T17:57:50+02:00"}
{"data":null,"level":"info","msg":"Starting worker 23ee18fc-d70e-4039-85d6-30c7c81933b9","time":"2024-05-15T17:57:50+02:00"}
{"data":null,"level":"info","msg":"Starting worker e53efbc4-8492-483a-97f2-0fff0556e358","time":"2024-05-15T17:57:50+02:00"}


{"time":"2024-05-15T17:58:21.603022413+02:00","id":"","remote_ip":"127.0.0.1","host":"0.0.0.0:8080","method":"POST","uri":"/command/ansible-playbook","user_agent":"curl/7.81.0","status":202,"error":"","latency":277778,"latency_human":"277.778µs","bytes_in":82,"bytes_out":274}
&{ [test/site.yml] false false <nil> map[] [] false false  127.0.0.1,  false false false    false    false false local     0  false  }

PLAY [all] *********************************************************************

TASK [Gathering Facts] *********************************************************
ok: [127.0.0.1]

TASK [ansibleplaybook-simple] **************************************************
ok: [127.0.0.1] =>
  msg: Your are running 'ansibleplaybook-simple' example

PLAY RECAP *********************************************************************
127.0.0.1                  : ok=2    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

## User Reference

### Perform a request to execute an Ansible playbook

The following example demonstrates how to execute an Ansible playbook using the Ransidble server. Please refer to the [REST API Reference](#rest-api-reference) section for more information.

```bash
❯ curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/command/ansible-playbook -d '{"playbooks": ["test/site.yml"], "inventory": "127.0.0.1,", "connection": "local", "project": "project"}'
HTTP/1.1 202 Accepted
Content-Type: application/json
Vary: Accept-Encoding
Date: Sat, 18 May 2024 09:29:12 GMT
Content-Length: 46

{"id":"2aff4156-00e7-413a-a54c-a372d009a3f3"}
```

### REST API Reference

#### Command

##### Execution Data Object

The JSON schema payload to execute an Ansible playbook command is defined [here](api/schemas/input/command/ansible-playbook-parameters.json). However, the following table provides a summary of the available attributes:

| JSON Attribute | Type | Description |
|----------------|------|-------------|
| project | string | The project name |
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

##### Dependencies Object

| JSON Attribute | Type | Description |
|----------------|------|-------------|
| roles | object | Defines how to install roles dependencies. The object is described in the section [Roles dependencies object](#roles-dependencies-object) |
| collections | object | Defines how to install collections dependencies. The object is described in the section [Collections dependencies object](#collections-dependencies-object) |

###### Roles Dependencies Object

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

###### Collections Dependencies Object

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

- Perform a request accepting gzip encoding to execute an Ansible playbook:

```bash
❯ curl -H 'Accept-Encoding: gzip' --output /dev/stdout -s -XPOST 0.0.0.0:8080/command/ansible-playbook
r+LIU(ILVHHM.-IMQ(.MNN-.N+'NyI&
```

- Get the status of the execution:

```bash
❯ curl -GET 0.0.0.0:8080/task/2aff4156-00e7-413a-a54c-a372d009a3f3 |jq
{
  "command": "ansible-playbook",
  "completed_at": "",
  "created_at": "2024-05-18T11:29:12+02:00",
  "executed_at": "2024-05-18T11:29:12+02:00",
  "id": "2aff4156-00e7-413a-a54c-a372d009a3f3",
  "parameters": {
    "project": "project",
    "playbooks": [
      "test/site.yml"
    ],
    "inventory": "127.0.0.1,",
    "connection": "local"
  },
  "status": "RUNNING"
}
```

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
