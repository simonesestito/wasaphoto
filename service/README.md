# Folder structure (in a nutshell)

TL;DR: the core is inside **features/** and it's packaged by features. Each with all the components it needs.

---

## features/

Each feature inside the `features` package follows the MVC architectural pattern.

More specifically, the code application is divided into layers:

* **Presentation Layer (Controller)**: it includes the Controllers, responsible for interfacing
  the platform independent Business Logic with the external world. In case one day the backend will move from REST to
  something else we still don't know, or we'll need to use XML instead of JSON,
  the only piece of code that must be changed is the Controllers code, leaving other layers untouched.
* **Service Layer**: it has all the business logic, without any dependency on the actual implementation.
  It MUST be written using only standard Go code, with no libraries of any sort.
* **Data Layer (DAO)**: it interfaces with the underlying data sources, like a SQL database or a file storage.
* **DTOs**: Data Transfer Object are used in data transmission between the external world and the application

A lot of effort is spent to have everything as abstract as possible, trying to implement the **Ports and Adapters
pattern**,
This approach helps us to achieve to make the business (domain) layer independent
of framework, UI, database or any other external components.

## Other folders

* `api` is the package with the **common** functionalities and types necessary to every other real controller or REST
  API endpoint.
* `utils` has all necessary utility functions, logically divided by type
* `ioc` since this app heavily uses **Dependency Injection**, the code here is responsible for
  creating instances of all interfaces providing real implementations (Inversion of Control).
