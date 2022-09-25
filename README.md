# Probe

[![Build Status](https://github.com/petaki/probe/workflows/tests/badge.svg)](https://github.com/petaki/probe/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

A small GO based agent for monitoring CPU, Memory and Disk usage.

## Getting Started

Before you start, you need to install the prerequisites.

### Prerequisites

- Redis: `Version >= 5.0` for data logging
- GO: `Version >= 1.19` for building

### Install from binary

Downloads can be found at releases page on [GitHub](https://github.com/petaki/probe/releases).

### Install from source

1. Clone the repository:

```
git clone git@github.com:petaki/probe.git
```

2. Open the folder:

```
cd probe
```

3. Build the Probe:

```
go build
```

4. Copy the example configuration:

```
cp .env.example .env
```

## Configuration

The configruation is stored in the `.env` file.

#### Disk Ignored

- `PATTERN*` - Prefix
- `*PATTERN` - Suffix
- `*PATTERN*` - Contains
- `PATTERN` - Exact match

```
PROBE_DISK_IGNORED=/dev,/var/lib/docker/*
```

---

#### Redis URL:

```
PROBE_REDIS_URL=redis://127.0.0.1:6379/0
```

#### Redis Key Prefix:

```
PROBE_REDIS_KEY_PREFIX=probe:
```

---

#### Data Log Enabled:

```
PROBE_DATA_LOG_ENABLED=true
```

#### Data Log Timeout (in seconds):

```
PROBE_DATA_LOG_TIMEOUT=2592000
```

---

#### Alarm Enabled:

```
PROBE_ALARM_ENABLED=false
```

#### Alarm Timeout (in seconds between alarms):

- Required: `PROBE_DATA_LOG_ENABLED=true`

```
PROBE_ALARM_TIMEOUT=300
```

#### Alarm CPU Percent:

- `0` - Disabled

```
PROBE_ALARM_CPU_PERCENT=30
```

#### Alarm Memory Percent:

- `0` - Disabled

```
PROBE_ALARM_MEMORY_PERCENT=50
```

#### Alarm Disk Percent:

- `0` - Disabled

```
PROBE_ALARM_DISK_PERCENT=80
```

#### Alarm Webhook Method:

```
PROBE_ALARM_WEBHOOK_METHOD=POST
```

#### Alarm Webhook URL:

```
PROBE_ALARM_WEBHOOK_URL=http://127.0.0.1:4000/alarm
```

#### Alarm Webhook Header:

```
PROBE_ALARM_WEBHOOK_HEADER='{"Authorization": "Bearer TOKEN", "Accept": "application/json"}'
```

#### Alarm Webhook Body:

- `%p` - Probe
- `%n` - Name of the watcher
- `%a` - Alarm percent
- `%u` - Used percent
- `%t` - Timestamp in `RFC3339` format

```
PROBE_ALARM_WEBHOOK_BODY='{"probe": "%p", "name": "%n", "alarm": %a, "used": %u, "timestamp": "%t"}'
```

## Running the tests

You can run the tests using the following command:

```
go test -v ./...
```

## Data visualization

You can display the collected data with the [Satellite](https://github.com/petaki/satellite).

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/probe/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
