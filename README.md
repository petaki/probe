# Probe

A lightweight Go agent for monitoring system resources and log files.

Probe collects every metric once a minute, on the minute.

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
- **Process** - the top 3 processes by CPU and the top 3 by memory. Process CPU is measured per core, so a busy process can exceed 100%.
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

- Go: `Version >= 1.27`

#### Steps

1. Clone the repository:

```bash
git clone https://github.com/petaki/probe.git
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

4. Run Probe from the directory that holds `.env`:

```bash
./probe
```

## Configuration

Probe reads its configuration from environment variables. On startup, it loads `.env` from the current working directory, not from the directory of the binary. A variable that is already set in the environment overrides the value in `.env`.

A ✅ marks where the variable is required.

| Option | Always | Data Log | Alarm | Alarm Filter | Log Tail |
|--------|:------:|:--------:|:-----:|:------------:|:--------:|
| `PROBE_NAME` | ✅ | | | | |
| `PROBE_DISK_IGNORED` | ✅ | | | | |
| `PROBE_DATA_LOG_ENABLED` | ✅ | | | | |
| `PROBE_ALARM_ENABLED` | ✅ | | | | |
| `PROBE_ALARM_FILTER_ENABLED` | ✅ | | | | |
| `PROBE_LOG_TAIL_ENABLED` | ✅ | | | | |
| `PROBE_REDIS_URL` | | ✅ | | ✅ | |
| `PROBE_DATA_LOG_TIMEOUT` | | ✅ | | | |
| `PROBE_ALARM_CPU_PERCENT` | | | ✅ | | |
| `PROBE_ALARM_MEMORY_PERCENT` | | | ✅ | | |
| `PROBE_ALARM_DISK_PERCENT` | | | ✅ | | |
| `PROBE_ALARM_LOAD_VALUE` | | | ✅ | | |
| `PROBE_ALARM_WEBHOOK_METHOD` | | | ✅ | | |
| `PROBE_ALARM_WEBHOOK_URL` | | | ✅ | | |
| `PROBE_ALARM_WEBHOOK_HEADER` | | | ✅ | | |
| `PROBE_ALARM_WEBHOOK_BODY` | | | ✅ | | |
| `PROBE_ALARM_FILTER_WAIT` | | | | ✅ | |
| `PROBE_ALARM_FILTER_SLEEP` | | | | ✅ | |
| `PROBE_LOG_TAIL_FILES` | | | | | ✅ |
| `PROBE_LOG_TAIL_LINES` | | | | | ✅ |
| `PROBE_LOG_TAIL_BUFFER_SIZE` | | | | | ✅ |
| `PROBE_LOG_TAIL_LIMIT` | | | | | ✅ |
| `PROBE_LOG_TAIL_TIMEOUT` | | | | | ✅ |

### General

#### Probe Name

Identifier of this probe instance. Used as the `%p` webhook placeholder and the Redis key prefix.

```
PROBE_NAME=probe
```

#### Disk Ignored

Comma-separated patterns to exclude disk partitions from monitoring. Spaces around entries are ignored.

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

Required when data logging or alarm filtering is enabled. Keys are prefixed with `<PROBE_NAME>:`.

#### Redis URL

```
PROBE_REDIS_URL=redis://127.0.0.1:6379/0
```

---

### Data Log

Persists collected metrics to Redis. Requires Redis to be configured.

#### Data Log Enabled

```
PROBE_DATA_LOG_ENABLED=true
```

#### Data Log Timeout (in seconds)

Probe stores each day's data under its own keys. The keys expire this many seconds after the first sample of the day. The value must be at least `1`.

```
PROBE_DATA_LOG_TIMEOUT=2592000
```

---

### Log Tail

Tails the last N lines of the specified log files every minute.

At most 1 MiB is read from the end of each file.

#### Log Tail Enabled

```
PROBE_LOG_TAIL_ENABLED=false
```

#### Log Tail Files (comma-separated)

Spaces around entries are ignored and empty entries are skipped. At least one path is required.

```
PROBE_LOG_TAIL_FILES=/var/log/syslog,/var/log/auth.log
```

#### Log Tail Lines

The value must be at least `1`.

```
PROBE_LOG_TAIL_LINES=10
```

#### Log Tail Buffer Size (in bytes)

The number of bytes that Probe reads per step while it searches backwards for line breaks. The value must be at least `1`.

```
PROBE_LOG_TAIL_BUFFER_SIZE=4096
```

#### Log Tail Limit (max entries per log file per day in Redis)

Probe stores each day's tails under their own keys, so the limit applies to each log file for each day. The value must be at least `1`.

```
PROBE_LOG_TAIL_LIMIT=60
```

#### Log Tail Timeout (in seconds)

The keys for each day's tails expire this many seconds after the first tail of the day. The value must be at least `1`.

```
PROBE_LOG_TAIL_TIMEOUT=172800
```

---

### Alarm

Sends a webhook notification when a metric reaches its threshold. Probe checks every minute, so it sends a notification every minute while the metric stays at or above the threshold. Use the [alarm filter](#alarm-filter) to reduce repeats.

Set a threshold to `0` to disable it. The percent thresholds must be between `0` and `100`.

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

A load alarm fires when any of the 1, 5, or 15 minute load averages reaches the value. The value must be at least `0`.

```
PROBE_ALARM_LOAD_VALUE=1.0
```

---

### Alarm Webhook

#### Alarm Webhook Method

One of `GET`, `HEAD`, `POST`, `PUT`, `PATCH`, `DELETE` or `OPTIONS`, in uppercase.

```
PROBE_ALARM_WEBHOOK_METHOD=POST
```

#### Alarm Webhook URL

An absolute `http` or `https` URL.

```
PROBE_ALARM_WEBHOOK_URL=http://127.0.0.1:4000/alarm
```

#### Alarm Webhook Header

A JSON object with string values. Probe does not set `Content-Type`, so set it here if your receiver needs it.

```
PROBE_ALARM_WEBHOOK_HEADER='{"Authorization": "Bearer TOKEN", "Accept": "application/json", "Content-Type": "application/json"}'
```

#### Alarm Webhook Body

The body supports the following placeholders:

| Placeholder | Description |
|-------------|-------------|
| `%p` | Probe name |
| `%n` | Watcher name: `CPU`, `Memory`, `Load`, or `Disk:<path>` |
| `%a` | Alarm threshold |
| `%u` | Current value as a number. For `Load`, a quoted string of the 1, 5, and 15 minute averages, such as `"1.50,0.80,0.40"` |
| `%t` | Timestamp (`RFC3339`) |
| `%x` | Timestamp (`Unix`) |
| `%l` | Satellite link (relative) |

`%p` and `%n` are escaped for use inside a JSON string. Put `%u` in the body without quotes, because the `Load` value already has them.

Probe waits up to 10 seconds for a response. A status code of `400` or higher counts as a failure.

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

Set to `0` or `1` to disable. With a value of N, an alarm fires on the Nth consecutive minute at or above the threshold. A minute without data, such as a minute after a restart, starts the count again.

Values above `1` require `PROBE_DATA_LOG_ENABLED=true`.

```
PROBE_ALARM_FILTER_WAIT=5
```

#### Alarm Filter Sleep (in seconds between alarms)

Set to `0` to disable. The cooldown applies to each metric separately, and to each disk path. It starts after a webhook succeeds. Probe retries a failed webhook the next minute.

```
PROBE_ALARM_FILTER_SLEEP=300
```

## Testing

```bash
go test ./...
```

Tests that need Redis are skipped unless `PROBE_TEST_REDIS_URL` is set. These tests write short-lived `probe-test-*` keys, so point the variable at a database that you can spare:

```bash
PROBE_TEST_REDIS_URL=redis://127.0.0.1:6379/15 go test ./...
```

## Data Visualization

Collected data can be visualized with [Satellite](https://github.com/petaki/satellite).

## Contributors

- [@dyipon](https://github.com/dyipon) - development ideas, bug reports and testing

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/probe/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
