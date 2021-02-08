# Probe

[![Build Status](https://github.com/petaki/probe/workflows/tests/badge.svg)](https://github.com/petaki/probe/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

A small GO based agent for monitoring CPU, Memory and Disk usage.

## Getting Started

Before you start, you need to install the prerequisites.

### Prerequisites

- Redis: `Version >= 2.8` for data logging
- GO: `Version >= 1.11` for building

### Install from source

Currently, you can only install from source.

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

The configruation is stored in the `.env` file.

#### Redis connection:

```
PROBE_REDIS_URL=redis://127.0.0.1:6379/0
```

#### Redis key prefix:

```
PROBE_REDIS_KEY_PREFIX=probe:
```

#### Redis key timeout (in seconds):

```
PROBE_REDIS_KEY_TIMEOUT=2592000
```

## Running the tests

You can run the tests using the following command:

```
go test -v ./...
```

## Data visualization

You can display the collected data with the [Carrier](https://github.com/petaki/carrier).

## License

The Probe is open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT).
