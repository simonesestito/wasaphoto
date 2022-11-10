# WASAPhoto

**A structured full stack project in Go and Vue.js**

This repository contains the source code for the university project for the
[Web and Software Application](http://gamificationlab.uniroma1.it/en/wasa/) course.

WASAPhoto is a social network where users can post photos, leave likes, comments and
also ban other users, with all the implications about information hiding.

It consists of:

* Documented REST API (OpenAPI 3.0) with all the endpoints described.
  You can find the specification [here](doc/api.yaml)
* Golang backend which implements the REST API. According to the given project
  specification, an authentication mechanism is not provided. Instead,
  the User ID is sent as an Authorization Bearer header, as it was a token in some way.
* Vue.js frontend app, which of course interfaces with the implemented REST API.
* All distributed using a Docker image

## Project structure and architecture

* `service/` contains all the private application code (project-specific functionalities).
	* `service/features` contains all the application code, **packaged by features**, which is a more robust way
	  of packaging source code compared to the more naive approach of packaging by type (e.g.: controllers, services, ...)
	  Each feature follows an **architectural pattern**, [discussed below](#architectural-pattern).
	* `service/api` is the package with the **common** functionalities and types necessary to every other real
	  controller or REST API endpoint
	* `service/utils` has all necessary utility functions, logically divided by type
	* `service/ioc` since this app heavily uses **Dependency Injection**,
	  the code here is responsible for creating instances of all interfaces providing real implementations.
* `cmd/` contains all executables; Go programs here only do "executable-stuff",
  like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is a daemon for checking the health of servers daemons;
	  useful when the hypervisor is not providing HTTP readiness/liveliness probes (e.g., Docker engine)
	* `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the OpenAPI specification
* `vendor/` is [managed by Go](https://go.dev/ref/mod#vendoring), and contains a copy of all dependencies
* `webui/` is the frontend code, developed in Vue.js; it includes Go code for release embedding

## Architectural Pattern

Each feature inside the `service/features` package follows the MVC architectural pattern.

More specifically, the code application is divided into layers:

* **Presentation Layer**: it includes the Controllers, responsible for interfacing
  the platform independent Business Logic with the external world. In case one day the backend will move from REST to
  something else we still don't know, or we'll need to use XML instead of JSON,
  the only piece of code that must be changed is the Controllers code, leaving other layers untouched.
* **Service Layer**: it has all the business logic, without any dependency on the actual implementation.
  It MUST be written using only standard Go code, with no libraries of any sort.
* **Data Layer**: it interfaces with the underlying data sources, like a SQL database or a file storage.
* **DTOs**: Data Transfer Object are used in data transmission between the external world and the application

A lot of effort is spent to have everything as abstract as possible, trying to implement the **Ports and Adapters
pattern**,
This approach helps us to achieve to make the business (domain) layer independent
of framework, UI, database or any other external components.

Also, components MUST NOT depend on the actual implementation, but they have to use the interfaces.
By doing that, unit testing or using multiple implementations it's guaranteed to be quite easy,
and the software is much more robust.
We could even not know the implementation yet when developing another components which depends on it.

## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```shell
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-npm.sh
  # (inside the NPM container)
  npm run build-embed
  exit

# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

## License

    Copyright (C) 2022 Simone Sestito

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
