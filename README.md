# Project ISO_Auditing_Tool

The ISO Auditing Tool helps ISO Auditors fill an audit questionnaire that automates loading data into Confluence and Jira.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```
### Create Database Container

Create DB container
```bash
make docker-run
```

Run the docker container in your computer. Remember to log in with the following commands. Remember to use a .env file to safely manage your credentials.
```bash
mysql -u USER - p
```

Once you're connected to the database you can use MySQL commands to select a database and manage your data:

```sql
SHOW DATABASES;
USE database;
SHOW TABLES;
SELECT * FROM table1;
```

Shutdown DB container
```bash
make docker-down
```

### Database migrations
The project has database migrations where you can create tables and their relationships.
You can execute the following make command to run the migrations:

```bash
make migrate
```

You can also seed the database which inserts data based on csv files that match the database tables.
To seed the database you can execute the following makek command:

```bash
make migrate
```

Finally, you can delete all the data in the database tables.
To wipe out all the data of the database tables execute the following make command:

```bash
make truncate
```

### Run the application

live reload the application
```bash
make watch
```

### Run Tests
run the test suite
```bash
make test
```

run the test suite and return a short results
```bash
make test-short
```

clean up binary from the last build
```bash
make clean
```
