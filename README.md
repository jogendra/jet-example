## Salesforce Marketing Cloud Content Fetcher

Application fetch content blocks from Salesforce Marketing Cloud (SFMC) and store them in a chosen storage location (currently supporting local or S3 storage) on schedule.

### Features

* Securely fetches content blocks from SFMC.
* Caches access tokens for efficient API calls.
* Supports concurrent fetching of content blocks for improved performance.
* Configurable/extendable storage options:
    * Local file storage (not implemented yet)
    * Amazon S3 bucket
* Schedulable execution (e.g., once per day) using cron.
* Configuration via environment variables.

### Project Structure
```
├── cmd -> entry points for application
│   ├── lambda -> lambda hanlder (empty)
│   └── cli -> cron job scheduler
├── internal -> core domain logic and application components
│   ├── scheduler -> scheduling and orchestration of tasks
│   ├── config -> application configurations (env variables)
│   ├── uploader -> adapters for uploading content blocks
│   ├── fetcher -> adapters for fetching content blocks
│   └── domain -> core domain logic and interfaces
└── pkg -> reusable packages 
```
This arrangement emphasizes the entry point of the application (`cmd`) and then shows the core domain logic (`internal`) before the supporting packages (`pkg`).

### Architecture
This application is using the principles of [hexagonal architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) (also known as ports and adapters) in a good extent.

#### Ports (Interfaces)
The `Uploader` and `Fetcher` interfaces in `internal/domain/ports.go` act as the "ports" in hexagonal architecture.

These interfaces define how the core domain interacts with external concerns (like fetching from Salesforce or storing in S3) without depending on concrete implementations.

#### Adapters (Implementations)
It has separate packages for different implementations ("adapters") of the `Uploader` and `Fetcher` ports:

* `internal/fetcher/salesforce`: Implements the `Fetcher` interface to fetch content blocks from Salesforce.
* `internal/uploader/s3`: Implements the `Uploader` interface to store content blocks in an S3 bucket.
* `internal/uploader/local`: Provides a local implementation of the `Uploader` (although currently unimplemented).


The core domain (`internal/domain`) depends on abstractions (interfaces), not on concrete implementations. This allows you to easily switch between different implementations (e.g., using a different cloud storage provider) without modifying the core domain logic.

The `internal/scheduler/scheduler.go` acts as the application logic that orchestrates the interaction between the `Fetcher` and `Uploader` implementations.

#### Benefits of this approach

* **Testability**: We can easily test the core domain logic independently of external dependencies by mocking the `Uploader` and `Fetcher` interfaces.
* **Maintainability**: The code is more modular and easier to understand and maintain due to the clear separation of concerns.
* **Flexibility**: We can easily add new features or change implementations without affecting other parts of the application.
* **Extensibility**: We can easily add support for new data sources or storage backends by implementing the corresponding interfaces.

### Requirements

* Go 1.23.3
* Docker

### Build the application
```bash
make build
```

### Run the application
Please make sure to have updated environment variables in `.env.*` file.
```bash
make run
```

### Testing
Only few unit tests are added so far.
```bash
make test 
```

### Makefile Commands
This project uses a Makefile to streamline common development tasks. Here are some of the available commands (including few mentioned above):

* `make deps`: Downloads project dependencies.
* `make fmt`: Formats the code.
* `make tidy`: Tidies the Go module dependencies.
* `make build`: Builds the application.
* `make run`: Runs the application.
* `make test`: Runs the tests.
* `make clean`: Cleans the build artifacts.
* `make docker-build`: Builds the Docker image.
* `make docker-run`: Runs the application in a Docker container.
* `make help`: Displays the help message with available commands.

## Further Improvements / Future Work
* [ ] Logging and monitoring
* [ ] Configuration schedule time - currently 24h. This could be achieved by taking time as ARG.
* [ ] Include infra-as-code - deploy on action
* [ ] Run scheduled job as scheduled lambda
* [ ] Implement local uploader (currently it is unimplemented)
* [ ] Implement more uploader (e.g. Google Cloud Storage and Azure Blob Storage)
* [ ] Flexibility to choose uploader. This could be achieved by taking uploader type as ARG.
* [ ] Add more unit tests
* [ ] Add integrations tests
* [ ] Add more pre-commit checks
* [ ] Add GitHub workflows for on-pr and on-merge to lint code, build images, deploy etc.
* [ ] Use tools like [Mockery](https://github.com/vektra/mockery) to mock interfaces for testing.
* [ ] Delta updates: Instead of fetching all content blocks every time, implement a mechanism to identify and fetch only the content blocks that have been updated or added since the last run. This can be achieved by using a timestamp or versioning system.
* [ ] Implement more robust error handling and retry mechanisms for failed API calls or storage operations. This will improve the resilience of the application.
* [ ] A command-line interface for the application. This would allow users to easily interact with the application and perform operations like scheduling, configuration, and manual triggering of content fetching.
* [ ] Expose API endpoints for integrating the application with other systems or services. (could be easily achievable with a new `cmd/api` since building blocks are already there)

Thank you :)