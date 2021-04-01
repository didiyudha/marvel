![](marvel.png?raw=true)

# Marvel API Implementation
This is a simple implementation to consume the API of character Marvel. 
For more information you can take look [here](https://developer.marvel.com/). 
You need to register and generate public and private key to access the API.

## Requirements
* Golang version go1.16+ (Since we use embed functionality for migration)
* Postgres 9.6+
* Redis
* Marvel API credentials
* Docker

## Note
At this moment. I assume you already installed Postgres and Redis in your local machine. In addition, please 
create database with name `marvel`

## Run Database Migration
```shell
make migration
```

## Install Dependencies
```shell
make deps
```

## Build Binary
There are two command that is provided to build binary. First, to build binary that runs on Linux machine. The second
one to build binary that runs on macOS machine.

```shell
make marvel-linux
```
```shell
make marvel-osx
```

```shell
make marvel-worker-linux
```
```shell
make marvel-worker-osx
```

## Generate Mock
```shell
make mock
```

## Run Unit Test
```shell
make test
```

## Initialize Data
The following command will delete all data from `characters` table in the `marvel` database, calling Marvel API 
to get the character data and insert them into characters table. Here's the command:
```shell
make characters
```

## Serve Swagger UI
```shell
cd openapi 
```
```shell
docker-compose up 
```

![](swagger.png?raw=true)

## System Design
![](system-design.png?raw=true)

## Configuration
```yaml
port: 8080
publicKey: <public key credential>
privateKey: <private key credential>
marvelHost: <marvel host>
db:
  user: <postgres username>
  password: <postgress password>
  host: <postgres host>
  name: <database name>
  maxIdleConns: <max idle connections>
  maxOpenConns: <max open connections>
  disableTLS: <true or false>
caching:
  addr: <redis address>
  password: <redis password>
  db: <redis database>
```