# pc-store
CRUD server providing the computer store API

### Run server
+ Set server settings in [configs/main.yml](./configs/main.yml)
+ Define some environment variables for connecting to the database in `.env` file
```bash
PC_HOST=localhost
PC_PORT=5432
PC_USERNAME=postgres
PC_PASSWORD=password
PC_SSLMODE=disable
PC_DBNAME=pcstore

PC_SERVICE_HASHSALT=test
PC_SERVICE_TOKENSECRET=test

AMQP_URI=amqp://guest:guest@localhost:5672/
```

+ Run server
```
make run
```

### API
After launch, you can look at the `api` on: [localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

### Author
**Baran Pavlo [GitHub](https://github.com/samurai-of-honor)**
