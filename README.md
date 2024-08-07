# Spacetraders

My Spacetraders API client

## Usage

### Requirements 
- [Jet (query-builder)](https://github.com/go-jet/jet)
- [Taskfile](https://taskfile.dev/)


### Initial setup
1. Clone the repository
2. Run docker-compose.yaml
3. Initialize the database with your starting system `task init [system]`

### Updating the query builder from the database schema

```bash
jet -dsn=postgresql://spacetraders@localhost:5432/spacetraders?sslmode=disable -schema=public -path=./.gen
```

## Tasks

Tasks are managed using [Taskfile](https://taskfile.dev/)

### openapi

Generates the OpenAPI client, to be used when the game updates.
Note: Some manual works is required after this command due to conflict with some constant names but nothing taking more than 5 minutes.
Also, search for `data failed to match schemas in oneOf(ExtractResources201ResponseDataEventsInner)` and make so that it doesn't returns an error anymore. THIS NEEDS TO BE FIXED!!!

```bash
task openapi
```

### start

Run the project (cmd/spacetraders/main.go)

```bash
`task start` 
```

### Init

Initialize the database with the provided system (cmd/init/main.go)


```bash
task init -- [system]
```
```
```
