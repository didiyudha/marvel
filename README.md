![](marvel.png?raw=true)

This service is built with :heart: and Go.

# Marvel API Implementation
This is a simple implementation to consume the API of character Marvel. 
For more information you can take look [here](https://developer.marvel.com/). 
You need to register and generate public and private key to access the API.

## Requirements
* [Golang](https://golang.org/dl/) version go1.16+ (Since we use embed functionality for migration)
* [Postgres](https://www.postgresql.org/download/) 9.6+
* [Redis](https://redis.io/)
* [Marvel API](https://developer.marvel.com/) credentials
* [Docker](https://www.docker.com/)

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

## Run Source Code
Before running the source code, you need to set the environment variable `MARVEL_CONFIG` with the value is the path 
of your configuration file. `export MARVEL_CONFIG=$(pwd)/config.yaml`. I suggest you to run `make characters` to get the data
from the API and store it into the database before running the application.

* ### Run Service
    ```shell
    make run-service
    ```

* ### Run Worker
    At this moment. Worker will be running `@hourly`
    ```shell
    make run-worker
    ```

## Run Binary
Don't forget to set `MARVEL_CONFIG` environment variable like the previous step.

* ### Run service binary
    ```shell
    ./marvel-linux
    ```

* ### Run worker binary
    ```shell
    ./marvel-worker-linux
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

## System Design and Caching Strategy
![](system-design.png?raw=true)

The main idea of the design is, we initialize the data for the first time using `marvel command`. The command
will delete all the data from `characters` table, call the Marvel API to get character data and store it to
the database (PostgreSQL). The worker will play crucial part since it will call Marvel API periodically, and compare
data from database to Marvel API data. The comparison purpose is to get the new data from the API. Only new data that
will be stored to database. At this time the worker will be running `@hourly`. The Marvel service provides API to get 
characters data. Redis will be the first option to get the data. If data that we are looking for is not exists
in the Redis, Marvel service will find it to the database. 

* ### Pros of the design
    * If the 3rd party API is broken, we still be able to serve the data since it is stored in our internal database.
* ### Cons of the design
  * Complexity of the architecture is higher than we just have one service to do it all.

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