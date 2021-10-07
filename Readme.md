# Golang auth try
This is a simple blog site using session based auth **Golang, Redis, MongoDB, Gin**

## How to run

```bash
go run *.go
```

## Usage

```bash
# Redis is required to store the session of the user. The project uses the default redis port :6379

# Database used: Mongo Db Atlas
add the following to .env file 

MONGODB_DB_USERNAME = <>
MONGODB_PASSWORD = <>
MONGODB_DB_NAME = <>

# db structure: 
## DB name: as defined above in .env 
### Collections : articles, users 
```