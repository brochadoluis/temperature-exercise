# Temperature Exercise

This project implements a temperature tracking system consisting of an API server and a data scrapper.

## Features

- The API server accepts temperature data and saves it to a MongoDB database.
- The data scrapper fetches temperature data from an external source and sends it to the API server.
- The system uses gRPC for communication between the scrapper and API server.
- Temperature data is categorized and stored in different collections based on its status (success, error, alert).
- The project uses Docker containers to run the API server, data scrapper, and MongoDB database.

## Installation

1. Make sure you have Go and Docker installed on your machine.
2. Install the Go plugins for the protocol compiler by running the following commands:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

3. Update your PATH environment variable to include the Go installation's bin directory. Add the following line to your
   shell profile (e.g., ~/.bashrc or ~/.bash_profile):

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

4. Clone this repository to your local machine.
5. Navigate to the project directory.

## Usage

- Run `make help` to list all available commands and their descriptions.
- Start the containers using Docker Compose:

```
make start
```

- Wait for the containers to start up. You can check the logs for each container to ensure they are running
  successfully:

```
docker logs -f <container-name>
```

- Once the containers are running, you can interact with the API server using the provided endpoints (
  e.g., http://localhost:8080).

## Configuration

- The API server is exposed on port 8080.
- The scrapper server is exposed on port 50051.
- The database server is exposed on port 50053.
- The MongoDB database is exposed on port 27017.
- Configuration settings can be modified in the `docker-compose.yml` file.

## Dependencies

- Golang: The project is written in Go and requires Go to be installed.
- Docker: The project uses Docker containers for running the services.
- MongoDB: The project uses MongoDB as the database.
- gRPC: The project relies on gRPC for communication between services. Install the Go plugins for the protocol compiler
  as mentioned in the Installation section.

## Usage

The service can be used by calling `http://localhost:8080/getTemperature?latitude={value}&longitude={value}`

## Contributing

Contributions are welcome! If you have any suggestions or find any issues, please create a new issue or submit a pull
request.

## License

This project is licensed under the [MIT License](LICENSE).
