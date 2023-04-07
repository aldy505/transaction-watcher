# Transaction Watcher

This is a repository to harden your SQL skill, learn Kafka (or Redpanda) and learn Docker.

## Your Tasks

This project contains a few tasks for you to work with.

### #1 The Watcher

Create a subdirectory called `watcher` that will do a SQL query periodically to the `transactions` table, 
and for every new row on that table, the `watcher` will produce a message to `transactions` topic on Redpanda.
You will also need to create the `transactions` topic yourself.

The schema of the message is defined on the `kafka-schemas` directory. Search for `transaction.json` file.

You must create a Dockerfile for your application. You can choose any language.

You don't need to create a HTTP API for this one. Just create a single function that will run when the 
program is executed. The application *must* work without any interference from you. It must work without
having anyone (including you) to trigger the run or consume function.

You can see that your `watcher` is emitting correct message to Redpanda through Redpanda Console that is
running on your local machine on port `8080`.

There will be environment variable available for you when you run it though Docker Compose:

```yaml
DATABASE_URL: "postgresql://watcher:password@postgres:5432/watcher?sslmode=disable"
KAFKA_ADDRESSES: "kafka:9092"
```

### #2 The Swimmer

Create a subdirectory called `swimmer` that will consume the `balance` topic from Redpanda, and do an event sourcing
of a customer's current balance. You must expose a HTTP API that handles a single endpoint of:

```http request
GET /current-balance?customer_id=123
Accept: application/json
```

With a response schema of:
```json
{
  "customer_id": 123,
  "current_balance": 123456
}
```

It's called swimmer because you will swim through the Redpanda `balance` topic records, and get an answer from that.
I will not give a clue on how you should consume your topic and produce a result based on the customer ID.

For the `balance` topic schema, you can look for `balance.json` on `kafka-schemas` directory.

You must create a Dockerfile for your application. You can choose any language.

Please expose the HTTP API at port 3000.

There will be environment variable available for you when you run it though Docker Compose:

```yaml
KAFKA_ADDRESSES: "kafka:9092"
```

### #3 The Frontend

This is an optional task.

Create a frontend (on `frontend` directory) that shows list of customer ID from `customer-list` service. 
Then, show the amount or balance that each customer have. You can create a full SPA page, or an SSR website.

You can hit these two endpoints:

```http request
GET http://customer-list:7201/customers
Accept: application/json

# Returns list of customer ID
[1,2,3,4,5]
```

```http request
GET http://swimmer:3000/current-balance?customer_id=123
Accept: application/json

# Returns:
{
  "customer_id": 123,
  "current_balance": 123456
}
```

Hey, the `swimmer` one is the one you made!

If you are developing it from your local machine, you can replace both `customer-list` and `swimmer` hostname
into `localhost`.

You can create a Dockerfile and add a new service schema on `docker-compose.yml` file to state your frontend container
configuration.

## How do I know if I did it?

Easy, run the Docker Compose. Build everything. See the frontend if it works. See the Redpanda Console if you
published and consume your topics and messages flawlessly.

```bash
docker compose up -d
```

To stop it, use

```bash
docker compose down
```

If you'd like to run just a few services for your local development purposes, for example only `postgres` and `kafka`
services, you can specify it on the `up` command, like so:

```bash
docker compose up -d postgres kafka
```

## Words of affirmation

Good luck, you can do it!