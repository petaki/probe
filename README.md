# Probe

A small GO based agent for monitoring CPU, Memory and Disk usage.

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

- CPU usage monitoring
- Memory usage monitoring
- Disk usage monitoring
- Process CPU and memory monitoring
- System load monitoring
- Log file tail watching
- Data logging to Redis
- Alarm notifications via webhooks
- Alarm filtering with wait and sleep intervals
- Console output mode

## Getting Started

Before you start, you need to install the prerequisites.

### Prerequisites

- Redis: `Version >= 5.0` for data logging

### Install from binary

Downloads can be found at releases page on [GitHub](https://github.com/petaki/probe/releases).

---

### Install from source

#### Prerequisites for building

- GO: `Version >= 1.21`

#### 1. Clone the repository:

```
git clone git@github.com:petaki/probe.git
```

#### 2. Open the folder:

```
cd probe
```

#### 3. Build the Probe:

```
go build
```

#### 4. Copy the example configuration:

```
cp .env.example .env
```

## Configuration

The configuration is stored in the `.env` file.

### General

#### Disk Ignored

Patterns to exclude disk partitions from monitoring:

- `PATTERN*` - Prefix
- `*PATTERN` - Suffix
- `*PATTERN*` - Contains
- `PATTERN` - Exact match

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

Requires Redis to be configured.

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

---

### Alarm

#### Alarm Enabled

```
PROBE_ALARM_ENABLED=false
```

#### Alarm CPU Percent

- `0` - Disabled

```
PROBE_ALARM_CPU_PERCENT=30
```

#### Alarm Memory Percent

- `0` - Disabled

```
PROBE_ALARM_MEMORY_PERCENT=50
```

#### Alarm Disk Percent

- `0` - Disabled

```
PROBE_ALARM_DISK_PERCENT=80
```

#### Alarm Load Value

- `0` - Disabled

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

- `%p` - Probe
- `%n` - Name of the watcher
- `%a` - Alarm percent or value
- `%u` - Used percent or value
- `%t` - Timestamp in `RFC3339` format
- `%x` - Timestamp in `Unix` format
- `%l` - Satellite link (relative)

```
PROBE_ALARM_WEBHOOK_BODY='{"probe": "%p", "name": "%n", "alarm": %a, "used": %u, "timestamp_rfc3339": "%t", "timestamp_unix": %x, "link": "%l"}'
```

---

### Alarm Filter

Requires Redis to be configured.

#### Alarm Filter Enabled

```
PROBE_ALARM_FILTER_ENABLED=false
```

#### Alarm Filter Wait (in minutes before first alarm)

- `0` - Disabled

```
PROBE_ALARM_FILTER_WAIT=5
```

#### Alarm Filter Sleep (in seconds between alarms)

- `0` - Disabled

```
PROBE_ALARM_FILTER_SLEEP=300
```

## Running the tests

You can run the tests using the following command:

```
go test -v ./...
```

## Data visualization

You can display the collected data with the [Satellite](https://github.com/petaki/satellite).

## Contributors

- [@dyipon](https://github.com/dyipon) for development ideas, bug reports and testing

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/probe/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
