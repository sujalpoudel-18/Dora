# DORA: Detection, Observation, Reconnaissance & Assessment

[![Build & Test](https://github.com/dora/dora/actions/workflows/test.yml/badge.svg)](https://github.com/dora/dora/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dora/dora)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

DORA is a serious, modular, high-performance cybersecurity vulnerability assessment scanner written in Go. Inspired by modern template-driven security engines, DORA provides fast, reproducible, and safe security observations across web applications, networks, API endpoints, and cloud services.

## Core Features
- **Template-Driven Architecture**: Easily extensible YAML security checks without recompiling binaries.
- **Safe Verification Principles**: Uses deterministic response matching, inert test markers, baseline request comparisons, and controlled callbacks rather than destructive payloads.
- **Multi-Protocol Engine**: HTTP/HTTPS, DNS analysis, TLS inspection, and TCP probing.
- **Technology Fingerprinting**: Multi-signal detection engine mapping headers, cookies, meta tags, and path signatures with confidence scoring.
- **Structured Output**: Support for JSON, JSON Lines (JSONL), CSV, SARIF, and interactive HTML reports.
- **SQLite & REST API**: Native persistent storage engine paired with a live REST API and web dashboard interface.

---

## Installation

### Prerequisites
- Go 1.22 or higher

### Build from Source
```bash
git clone https://github.com/dora/dora.git
cd dora
make build
```

### Install Binary
```bash
sudo make install
dora -version
```

---

## Quick Start

### Basic Target Assessment
```bash
dora -u https://example.com
```

### Bulk Scanning from File / Stdin
```bash
dora -list targets.txt -concurrency 50 -rate-limit 200
cat targets.txt | dora -stdin -json -output results.json
```

---

## CLI Reference

| Flag | Description |
|---|---|
| `-u` | Single target URL, hostname, or IP address |
| `-list` | Path to target list file |
| `-stdin` | Read targets directly from standard input |
| `-concurrency` | Global concurrency limit (default 25) |
| `-rate-limit` | Maximum HTTP requests per second (default 150) |
| `-timeout` | Timeout in seconds per HTTP request (default 10) |
| `-json` | Output findings in JSON format |
| `-version` | Print DORA version and exit |
| `-help` | Display CLI usage guide |

---

## Safe Verification Principles

DORA is designed strictly for **authorized security testing**, local lab environments, CTF challenges, and staging deployments. DORA enforces safe verification methods:
1. **Inert Markers**: Uses harmless, unique random strings (e.g. `DORA_XSS_MARKER`) to check reflection contexts without launching unmanaged script execution.
2. **Differential Comparison**: Verifies SQL injection or logic bugs via baseline vs modified response structures rather than invoking data dumping.
3. **Controlled Callbacks**: Out-of-band SSRF verification utilizes explicitly configured callback listener servers.

---

## License
MIT License. See [LICENSE](LICENSE) for details.
