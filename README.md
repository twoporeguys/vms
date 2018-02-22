# Version Management Service

## Overview

VMS is intended to be deployed as a server that accepts version information for various applications & environments and can answer queries regarding current version to deploy.

![VMS Diagram](https://github.com/twoporeguys/vms/blob/master/docs/VMS-Diagram.png)

## Requirements
Redis

## Configuration
The configuration relies on environment variables:
- `VMS_PORT`: The port number to bind the service to. Defaults to `8080`.
- `VMS_REDIS`: The address of the Redis instance used as a data storage. Defaults to `127.0.0.1:6379`.

## API
The API reads and returns JSON documents (`content-type: application/json`).
### Read (HTTP `GET`)
- `/`: Returns the whole state of all applications in all environments.
  - Request: `GET /`
  - Response:
```json
    {
      "vms": {
        "prod": {
          "backend": "0.9.5",
          "frontend": "1.0.6"
        },
        "stg": {
          "backend": "1.0.1",
          "frontend": "1.1.0"
        }
      }
    }
```
- `/{app}`: Returns the state of a single application in all environments.
  - Request: `GET /vms`
  - Response:
```json
    {
      "prod": {
        "backend": "0.9.5",
        "frontend": "1.0.6"
      },
      "stg": {
        "backend": "1.0.1",
        "frontend": "1.1.0"
      }
    }
```
- `/{app}/{environment}`: Returns the state of a single application in one environment.
  - Request: `GET /vms/prod`
  - Response:
```json
    {
      "backend": "0.9.5",
      "frontend": "1.0.6"
    }
```
- `/{app}/{environment}/{component}`: Returns the version of a component of a single application in one environment.
  - Request: `GET /vms/prod/backend`
  - Response:
```json
"0.9.5"
```

### Write (HTTP `POST`)
- `/`: Returns the whole state of all applications in all environments.
  - Request: `POST /`
  - Body:
```json
    {
      "vms": {
        "prod": {
          "backend": "0.9.5",
          "frontend": "1.0.6"
        },
        "stg": {
          "backend": "1.0.1",
          "frontend": "1.1.0"
        }
      }
    }
```
- `/{app}`: Returns the state of a single application in all environments.
  - Request: `POST /vms`
  - Body:
```json
    {
      "prod": {
        "backend": "0.9.5",
        "frontend": "1.0.6"
      },
      "stg": {
        "backend": "1.0.1",
        "frontend": "1.1.0"
      }
    }
```
- `/{app}/{environment}`: Returns the state of a single application in one environment.
  - Request: `POST /vms/prod`
  - Body:
```json
    {
      "backend": "0.9.5",
      "frontend": "1.0.6"
    }
```
- `/{app}/{environment}/{component}`: Returns the version of a component of a single application in one environment.
  - Request: `POST /vms/prod/backend`
  - Body:
```json
"0.9.5"
```
