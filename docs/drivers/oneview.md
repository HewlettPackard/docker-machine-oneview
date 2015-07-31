<!--[metadata]>
+++
title = "OneView"
description = "HP OneView driver for machine"
keywords = ["machine, OneView, driver"]
[menu.main]
parent="smn_machine_drivers"
+++
<![end-metadata]-->

# HP OneView
Create machines using HP OneView API 1.20


Options:

 - `--oneview-api-url`: **required** API URL for HP OneView server.

 TODO: these might not be needed and are just examples.

 - `--oneview-ssh-user`: SSH username used to connect.
 - `--oneview-ssh-key`: Path to the SSH user private key.
 - `--oneview-ssh-port`: Port to use for SSH.

> **Note**: You must use a base operating system supported by Machine.

Environment variables and default values:

| CLI option                 | Environment variable | Default             |
|----------------------------|----------------------|---------------------|
| **`--oneview-api-url`** | -                    | -                   |
| `--oneview-ssh-user`       | -                    | `root`              |
| `--oneview-ssh-key`        | -                    | `$HOME/.ssh/id_rsa` |
| `--oneview-ssh-port`       | -                    | `22`                |
