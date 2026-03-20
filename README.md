# Probe

A lightweight Go agent for monitoring system resources and log files.

```

  🔍 Starting Probe...

  📡 Data logging is enabled.

  🚨 Alarm is armed.

  🤖 Probe is watching.

```

## Badges

[![Build Status](https://github.com/petaki/probe/workflows/tests/badge.svg)](https://github.com/petaki/probe/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Features

- **CPU** - overall CPU usage percentage
- **Memory** - overall memory usage percentage
- **Process** - per-process CPU and memory usage
- **Load** - system load averages (1, 5, 15 min)
- **Disk** - per-partition disk usage percentage (with ignore patterns)
- **Log Tail** - tail the last N lines of specified log files
- **Data Logging** - persist collected metrics to Redis with configurable TTL
- **Alarms** - webhook notifications when thresholds are exceeded
- **Alarm Filtering** - wait and sleep intervals to reduce alarm noise
- **Console Mode** - print metrics to stdout when data logging is disabled

## Getting Started

Follow the steps below to install and configure Probe.

### Prerequisites

- Redis: `Version >= 5.0` (optional, required for data logging and alarm filtering)

### Install from Binary

Download the latest release for your platform from the [GitHub Releases](https://github.com/petaki/probe/releases) page.

---

### Install from Source

#### Prerequisites

- Go: `Version >= 1.26`

#### Steps

1. Clone the repository:

```bash
git clone git@github.com:petaki/probe.git
```

2. Build the binary:

```bash
cd probe
go build
```

3. Copy and edit the configuration:

```bash
cp .env.example .env
```

## Configuration

All configuration is done through environment variables in the `.env` file.

### General

#### Disk Ignored

Comma-separated patterns to exclude disk partitions from monitoring:

| Pattern | Match Type |
|---------|------------|
| `PATTERN*` | Prefix |
| `*PATTERN` | Suffix |
| `*PATTERN*` | Contains |
| `PATTERN` | Exact |

```
PROBE_DISK_IGNORED=/dev,/var/lib/docker/*
```

---

### Redis

#### Redis URL

```
PROBE_REDIS_URL=redis://127.0.0.1:6379/0
```

#### Redis Key Prefix

```
PROBE_REDIS_KEY_PREFIX=probe:
```

---

### Data Log

Persists collected metrics to Redis. Requires Redis to be configured.

#### Data Log Enabled

```
PROBE_DATA_LOG_ENABLED=true
```

#### Data Log Timeout (in seconds)

```
PROBE_DATA_LOG_TIMEOUT=2592000
```

---

### Log Tail

Tails the last N lines of specified log files on each monitoring cycle.

#### Log Tail Enabled

```
PROBE_LOG_TAIL_ENABLED=false
```

#### Log Tail Files (comma-separated)

```
PROBE_LOG_TAIL_FILES=/var/log/syslog,/var/log/auth.log
```

#### Log Tail Lines

```
PROBE_LOG_TAIL_LINES=10
```

#### Log Tail Buffer Size (in bytes)

```
PROBE_LOG_TAIL_BUFFER_SIZE=4096
```

#### Log Tail Limit (max entries per log file in Redis)

```
PROBE_LOG_TAIL_LIMIT=60
```

#### Log Tail Timeout (in seconds)

```
PROBE_LOG_TAIL_TIMEOUT=259200
```

---

### Alarm

Sends webhook notifications when a metric exceeds its threshold. Set a threshold to `0` to disable it.

#### Alarm Enabled

```
PROBE_ALARM_ENABLED=false
```

#### Alarm CPU Percent

```
PROBE_ALARM_CPU_PERCENT=30
```

#### Alarm Memory Percent

```
PROBE_ALARM_MEMORY_PERCENT=50
```

#### Alarm Disk Percent

```
PROBE_ALARM_DISK_PERCENT=80
```

#### Alarm Load Value

```
PROBE_ALARM_LOAD_VALUE=1.0
```

---

### Alarm Webhook

#### Alarm Webhook Method

```
PROBE_ALARM_WEBHOOK_METHOD=POST
```

#### Alarm Webhook URL

```
PROBE_ALARM_WEBHOOK_URL=http://127.0.0.1:4000/alarm
```

#### Alarm Webhook Header

```
PROBE_ALARM_WEBHOOK_HEADER='{"Authorization": "Bearer TOKEN", "Accept": "application/json"}'
```

#### Alarm Webhook Body

The body supports the following placeholders:

| Placeholder | Description |
|-------------|-------------|
| `%p` | Probe name |
| `%n` | Watcher name |
| `%a` | Alarm threshold |
| `%u` | Current value |
| `%t` | Timestamp (`RFC3339`) |
| `%x` | Timestamp (`Unix`) |
| `%l` | Satellite link (relative) |

```
PROBE_ALARM_WEBHOOK_BODY='{"probe": "%p", "name": "%n", "alarm": %a, "used": %u, "timestamp_rfc3339": "%t", "timestamp_unix": %x, "link": "%l"}'
```

---

### Alarm Filter

Reduces alarm noise by requiring sustained threshold violations before firing and enforcing a cooldown between alarms. Requires Redis to be configured.

#### Alarm Filter Enabled

```
PROBE_ALARM_FILTER_ENABLED=false
```

#### Alarm Filter Wait (in minutes before first alarm)

Set to `0` to disable.

```
PROBE_ALARM_FILTER_WAIT=5
```

#### Alarm Filter Sleep (in seconds between alarms)

Set to `0` to disable.

```
PROBE_ALARM_FILTER_SLEEP=300
```

## Testing

```bash
go test -v ./...
```

## Data Visualization

Collected data can be visualized with [Satellite](https://github.com/petaki/satellite).

## Contributors

- [@dyipon](https://github.com/dyipon) - development ideas, bug reports and testing

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/probe/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
