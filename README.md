## Test task

Go version used: `1.13.4`

This is a test solution that provides API to query amount of trips specific cab has done on a specific date.

### Solution peculiarities
* I assume that API receives one date and N cabs' medallions, where N >= 1.
* To cache data I use Redis, data has a TTL, that can be specified in `config.json`
* API responds with JSON, that has 3 blocks:
    * `trace_id` - that's the request ID assigned to a request.
    * `data` - block with key-value pairs, where key is cab's medallion and value is amount of trips performed by that cab.
    * `error` - contains error message in case if something went wrong.
* `data` provides only list of cabs with available data. It doesn't fill in missing data with 0. If this is required it can be easily done as the code is designed to be flexible to such changes.
* response can contain mixture of data from cache and DB, depending on what was cached.
* If you have a look at the code, you might notice that sometimes I use `interface`s to make abstraction over something and sometimes I define a type to `func`. I use interfaces when methods are related to one another and use common data/objects, and I use functions, when there will be an interface/object with a single method.
* Mocks for unit-test where generated mostly by mockgen, unfortunately it can't mock functions, so function's mocks I did manually using the same approach.

### Contents
This project contains packages that can be treated as utility packages:
* [logs](logs/README.md) - logging interface and implementations
* [utils](utils/README.md) - utility functions
* [mocks](mocks/README.md) - mocks for unit tests
* mysql - some files related to containers created by `docker-compose`.

and solution package:
* [logic](logic/README.md) - core interfaces and implementations to solve the task

### Endpoints:
* GET `http://<binded-host:binded-port>/v1/cab?id=<ids>&date=<date>&nocache=<skip_cache>` - gets information about cabs
    * ids - required, comma separated list of cabs' medallions, at least 1.
    * date - required, date on which we should count trips
    * nocache - optional, set to `true` to bypass cache
* DELETE `http://<binded-host:binded-port>/v1/caches` - wipes cache completely
* DELETE `http://<binded-host:binded-port>/v1/cache?id=<ids>&date=<date>` - wipes data for specified cabs and on specified date

### Unit tests
Package `logic` is mostly covered with tests, as it contains a core business logic, but several unit-tests still have to written.

### Things to improve:
* add more unit-tests and reduce code duplication in existing tests. It is possible to add few more tests in `logic` package
and to cover code with tests in `utils` and `logs` packages.
* add more logging. To keep code simple, I did less logging.
* add more options to redis config.
* clean up bash-scripts and make `docker-compose` a bit more flexible - right now it has password hardcoded, for example.

### Things to keep in mind:
* if it is a production release, one should consider using https and providing `OPTIONS`-method for existing handlers.
 This can be done either using a proxy/load-balancer before hitting this service or (worst case scenario) via using
 ambassador template with `nginx` inside container, that will run in the same network-namespace as a container with this service.


## How to run the service:
Modify file `./config.json`:
* `mysql` - mysql configureation. **MySQL password is provided via command line**.
* `cache_ttl_seconds` - for how long we should keep cached value (in seconds).
* `redis` - block of redis configuration. Supports only address (`host:port`) and DB. **Redis password is provided via
command line**.
* `bind` - which IP and port should be used by the service.
* `app_timeout_seconds` - when `SIGINT` or `SIGTERM` is caught, application is informed and should stop withing this
time interval, otherwise it will be killed.

### Source code:
`go build && ./testdb -config=./config.json -redispwd='securepassword' -mysqlpwd='securepassword'`
### Docker-compose-way:
In root of repo just execute `docker-compose up` and wait while data is being restored. Once it's done - service will start listening and will be ready to use.
This approach assumes, that data to be restored is contained in `mysql/data/*.zip` file inside this repo.
### Docker-way:
We will pull mysql container, run it with `root` user, will pull redis container and run it without any authentication just to keep things simple.
1. Create a mysql container:
    * `docker run -d --rm --name mysql_db -e MYSQL_ROOT_PASSWORD=securepassword mysql:5.7.28`
2. Don't forget to restore data from dump.
3. Create a redis container:
    * `docker run -d --rm --name redistest --network=container:mysql_db redis:latest`
4. Build container with service inside (from root of this repo):
    * `docker build . -t testdb`
5. Run container with service:
    * `docker run -d --rm--network=container:mysql_db --name=testdb testdb -redispwd="" -mysqlpwd="securepassword"`

Those commands will spin a redis container without authentication, mysql container and will create a container with the service,
building it from source code and will spin that container with the service inside, providing same network namespace as
mysql container will have.

To remove container just stop them.

If you have redis/mysql with authentication running in another container or on host, you can execute `docker run` for `testdb`,
providing redis/mysql password in `-redispwd`/`-mysqlpwd` parameters:

`docker run -d  --name cachetest cachetest -redispwd="securepassword" -mysql="securepassword"`

Also you can use precompiled container ([this one](https://hub.docker.com/repository/docker/coldze/testdb)), substituting `config.json` file:

`docker run -d -v $(pwd)/build/config.json:/go/src/app/config.json --network=container:redistest --name testdb coldze/testdb -redispwd="" -mysqlpwd="securepassword""`

### Sample tests:
* With `curl` from host:
   * Get date: `curl http://<service-container-ip:port>/v1/cab?id=000318C2E3E6381580E5C99910A60668,00377E15077848677B32CE184CE7E871&date=2013-12-03`
   * Wipe cache: `curl -X "DELETE" "http://<service-container-ip:port>/v1/caches"`
   * Delete entry from cache: `curl -X "DELETE" "http://<service-container-ip:port>/v1/cache?id=000318C2E3E6381580E5C99910A60668&date=2013-12-03"`
* With `curl` inside separate container, if you were following `docker-compose-way` or `docker-way` (in this case you don't need IP of container):
   * Get date: `docker run --rm --network=container:mysql_db coldze/alpine:curl "http://<service-container-ip:port>/v1/cab?id=000318C2E3E6381580E5C99910A60668,00377E15077848677B32CE184CE7E871&date=2013-12-03"`
   * Wipe cache: `docker run --rm --network=container:mysql_db coldze/alpine:curl -X "DELETE" "http://<service-container-ip:port>/v1/caches"`
   * Delete entry from cache: `docker run --rm --network=container:mysql_db coldze/alpine:curl -X "DELETE" "http://<service-container-ip:port>/v1/cache?id=000318C2E3E6381580E5C99910A60668&date=2013-12-03"`

## How to run unit-tests:
From root of repo run following command:

`go test ./...`

With coverage:
`go test -cover -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`
