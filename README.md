# Ransidble

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![swagger-editor](https://img.shields.io/badge/open--API-in--editor-brightgreen.svg?style=flat&label=open-api-v3)](https://editor.swagger.io/?url=https://raw.githubusercontent.com/apenella/ransidble/main/api/openapi.yaml)

![Ransidble-logo](docs/images/logo_ransidble.png "Ransidble logo")

Ransidble is a utility that enables you to execute [Ansible](https://www.ansible.com/) commands on remote hosts. It functions as a wrapper over Ansible, allowing users to launch playbooks, roles, and tasks remotely. This is achieved by exposing a REST API in front of Ansible, facilitating the execution of Ansible commands on remote hosts.

- [Ransidble](#ransidble)
  - [Why Ransidble?](#why-ransidble)
  - [How about the name?](#how-about-the-name)
  - [Server Usage Reference](#server-usage-reference)
    - [Configurating the Ransidble server](#configurating-the-ransidble-server)
    - [Initiating the Ransidble server](#initiating-the-ransidble-server)
  - [REST API Reference](#rest-api-reference)
  - [User Reference](#user-reference)
    - [Project Definition](#project-definition)
      - [Project Storage Types](#project-storage-types)
        - [Local Filesystem](#local-filesystem)
      - [Project Format Types](#project-format-types)
        - [Plain](#plain)
        - [Tar Gz](#tar-gz)
    - [Examples of Requests](#examples-of-requests)
      - [Performing a Request to Execute an Ansible playbook](#performing-a-request-to-execute-an-ansible-playbook)
      - [Perform a Request Accepting Gzip Encoding](#perform-a-request-accepting-gzip-encoding)
      - [Getting the Status of an Execution](#getting-the-status-of-an-execution)
      - [Getting the project details](#getting-the-project-details)
      - [Getting the list of projects](#getting-the-list-of-projects)
  - [Development Reference](#development-reference)
    - [Contributing](#contributing)
    - [Code of Conduct](#code-of-conduct)
    - [Roadmap](#roadmap)
  - [Acknowledgments](#acknowledgments)
  - [License](#license)

## Why Ransidble?

## How about the name?

Ransidble is a blend of 'Remote' and 'Ansible,' with a nod to the punk rock band [Rancid](https://rancidrancid.com/).

## Server Usage Reference

### Configurating the Ransidble server

The Ransidble server can be configured using environment variables. The following table lists the available environment variables:

| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| RANSIDBLE_SERVER_HTTP_LISTEN_ADDRESS | The port where the server listens for incoming requests | :8080 |
| RANSIDBLE_SERVER_LOG_LEVEL | The log level for the server | info |
| RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH | The path where the projects are stored | projects |
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

### Initiating the Ransidble server

```bash
RANSIDBLE_SERVER_LOG_LEVEL=info RANSIDBLE_SERVER_WORKER_POOL_SIZE=3 RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH=test/projects  go run cmd/main.go serve
{"data":null,"level":"info","msg":"Starting server on :8080","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Starting worker 4838bac7-3efe-4a9d-837d-3e61924c5f35","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Starting worker 50498ab1-72b0-4b8f-a197-a83b527ec874","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Starting worker 7f05137a-7701-4e88-8aee-1b8fe65d2e80","time":"2024-07-21T20:01:04+02:00"}
{"data":null,"level":"info","msg":"Executing task ecfc92dc-6323-40d6-9bf8-71c4d4d98640","time":"2024-07-21T20:01:19+02:00"}
{"time":"2024-07-21T20:01:19.534469135+02:00","id":"","remote_ip":"127.0.0.1","host":"0.0.0.0:8080","method":"POST","uri":"/tasks/ansible-playbook/project-1","user_agent":"curl/7.81.0","status":202,"error":"","latency":518245,"latency_human":"518.245µs","bytes_in":77,"bytes_out":46}
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

## REST API Reference

The Ransidble provides you with a Open API specification that you can use to interact with the server. The Open API specification is available in the [api/openapi.yaml](api/openapi.yaml) file or through the [Swagger Editor](https://editor.swagger.io/?url=https://raw.githubusercontent.com/apenella/ransidble/main/api/openapi.yaml).

## User Reference

### Project Definition

To execute an Ansible playbook, you must first define a project. A project serves as the source code for the Ansible playbook.

Each project has the following attributes:

- **Name**: The name of the project. Which is the unique identifier for the project.
- **Reference**: The reference where the project is located in the storage.
- **Storage Type**: The type of storage used to store the project.
- **Format**: The format of the project.

#### Project Storage Types

##### Local Filesystem

The local storage type stores the project in the local filesystem. You can define the path where projects are stored by using the `RANSIDBLE_SERVER_PROJECT_LOCAL_STORAGE_PATH` environment variable.

#### Project Format Types

##### Plain

The plain format is a directory containing the Ansible playbook files.
You can use the `plain` format for a project stored in the local filesystem.

##### Tar Gz

The `targz` format is a tarball compressed with gzip that contains the Ansible playbook files. Ransidble identifies a `targz` project by its `.tar.gz` extension.
You can use the `targz` format for a project stored in the local filesystem.

To prepare a project in the `targz` format, create a tarball compressed with gzip that contains the Ansible playbook files. The following example demonstrates how to create such a tarball from the files stored in the `my-project` directory:

```bash
tar -czvf my-project.tar.gz -C my-project .
```

### Examples of Requests

#### Performing a Request to Execute an Ansible playbook

The following example demonstrates how to execute an Ansible playbook using the Ransidble server. Please refer to the [REST API Reference](#rest-api-reference) section for more information.

```bash
curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/tasks/ansible-playbook/project-1 -d '{"playbooks": ["site.yml"], "inventory": "127.0.0.1,", "connection": "local"}'
HTTP/1.1 202 Accepted
Content-Type: application/json
Vary: Accept-Encoding
Date: Sun, 21 Jul 2024 18:01:19 GMT
Content-Length: 46

{"id":"ecfc92dc-6323-40d6-9bf8-71c4d4d98640"}
```

#### Perform a Request Accepting Gzip Encoding

```bash
❯ curl -H 'Accept-Encoding: gzip' --output /dev/stdout -s -XPOST 0.0.0.0:8080/command/ansible-playbook
r+LIU(ILVHHM.-IMQ(.MNN-.N+'NyI&
```

#### Getting the Status of an Execution

```bash
❯ curl -s -GET 0.0.0.0:8080/tasks/ecfc92dc-6323-40d6-9bf8-71c4d4d98640 | jq
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

#### Getting the project details

```bash
❯ curl -s 0.0.0.0:8080/projects/project-1 | jq
{
  "format": "plain",
  "name": "project-1",
  "reference": "test/projects/project-1",
  "type": "local"
}
```

#### Getting the list of projects

```bash
❯ curl -s 0.0.0.0:8080/projects | jq
[
  {
    "format": "plain",
    "name": "project-1",
    "reference": "test/projects/project-1",
    "type": "local"
  },
  {
    "format": "plain",
    "name": "project-2",
    "reference": "test/projects/project-2",
    "type": "local"
  }
]
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

## Acknowledgments

## License

Ransidble is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
