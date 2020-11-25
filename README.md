# go-vrop

<a href="https://github.com/greenpau/go-vrop/actions/" target="_blank"><img src="https://github.com/greenpau/go-vrop/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/go-vrop" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>

vRealize API Client Library

<!-- begin-markdown-toc -->
## Table of Contents

* [Getting Started](#getting-started)
* [References](#references)

<!-- end-markdown-toc -->

## Getting Started

First, install `vropcli`:

```bash
go get -u github.com/greenpau/go-vrop/cmd/vropcli
```

Next, set environment variables for vRealize API Token:

```bash
export VROP_HOST=vrop
export VROP_USERNAME=admin
export VROP_PASSWORD=My@Password
```

Alternatively, the settings could be passed in a configuration file. There are
two options:

1. The `vropcli.yaml` should be located in `$HOME/.config/vropcli` or current directory
2. Pass the location via `-config` flag

```yaml
---
host: vrop
username: admin
password: password
```

The following command fetches virtual machines data from vRealize API:

```bash
vropcli -get-virtual-machines
```
