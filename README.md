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
❯ go run cmd/main.go serve
&{HTTPListenAddress::8080}
2024/05/03 19:47:47 Starting server on :8080
[NOT IMPLEMENTED] Running Ansible playbook
{"time":"2024-05-03T19:47:52.279116069+02:00","id":"","remote_ip":"127.0.0.1","host":"0.0.0.0:8080","method":"POST","uri":"/command/ansible-playbook","user_agent":"curl/7.81.0","status":200,"error":"","latency":427317,"latency_human":"427.317µs","bytes_in":0,"bytes_out":38}
```

- Perform a request to execute an Ansible playbook:

```bash
❯ curl --output /dev/stdout -s -XPOST 0.0.0.0:8080/command/ansible-playbook
Ansible playbook executed successfully
```

- Perform a request accepting gzip encoding to execute an Ansible playbook:

```bash
❯ curl -H 'Accept-Encoding: gzip' --output /dev/stdout -s -XPOST 0.0.0.0:8080/command/ansible-playbook
r+LIU(ILVHHM.-IMQ(.MNN-.N+'NyI&
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
