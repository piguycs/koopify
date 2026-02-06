# Koopify

A sassy little web app to transfer funds from your wallet to mine and transfer snake oil from my non-existant inventory to yours!

## Local Deployment

Docker compose is used for local deployments.

```
$ docker-compose --version
Docker Compose version 5.0.2
```

Links to locally hosted services:
- koopify-web: [https://localhost:8000/](https://localhost:8000/)
- koopify-api: [https://localhost:8080/](https://localhost:8000/)
- koopify-db: TODO

TLS is used, as the same container images are used for the real deployment on [https://koopify.piguy.nl](https://koopify.piguy.nl). The docker-compose manifest spins up an alpine container to generate locally signed certificates. The official deployment uses letsencrypt.

### License

All the code in this repository is licensed under Apache-2.0 license.
