# Ransidble

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

![Ransidble-logo](docs/images/logo_ransidble.png "Ransidble logo")

Ransidble is a utility that enables you to execute [Ansible](https://www.ansible.com/) commands on remote hosts. It functions as a wrapper over Ansible, allowing users to launch playbooks, roles, and tasks remotely. This is achieved by exposing a REST API in front of Ansible, facilitating the execution of Ansible commands on remote hosts.

- [Ransidble](#ransidble)
  - [Why Ransidble?](#why-ransidble)
  - [How about the name?](#how-about-the-name)
  - [Usage Reference](#usage-reference)
  - [Development Reference](#development-reference)
    - [Contributing](#contributing)
    - [Code of Conduct](#code-of-conduct)
    - [Roadmap](#roadmap)
  - [Acknowledgements](#acknowledgements)
  - [License](#license)

## Why Ransidble?

## How about the name?

Ransidble is a blend of 'Remote' and 'Ansible,' with a nod to the punk rock band [Rancid](https://rancidrancid.com/).

## Usage Reference

- Start the Ransidble server:

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

- Perform a request to execute an Ansible playbook:

```bash
❯ curl -i -s -H "Content-Type: application/json" -XPOST 0.0.0.0:8080/command/ansible-playbook -d '{"playbooks": ["test/site.yml"], "inventory": "127.0.0.1,", "connection": "local", "project": "project"}'
HTTP/1.1 202 Accepted
Content-Type: application/json
Vary: Accept-Encoding
Date: Sat, 18 May 2024 09:29:12 GMT
Content-Length: 46

{"id":"2aff4156-00e7-413a-a54c-a372d009a3f3"}
```

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
