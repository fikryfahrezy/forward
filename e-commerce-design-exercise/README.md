# E-Commerce Design Excercise

## Disclaimer

The author doesn't have any experience that relate to E-Commerce domain nor designing system architecture before.

## High-level architecture diagram

![Architecture Diagram](./architecture-diagram.png)

### Registration Flow

- User Register through `Auth Service`.
- `Auth Service` call `User Service` to create profile.
- `Auth Service` publish event to `Kafka`.
- `Notification Service` consume the user registartion event to create welcome email.
- `Analytic Service`consume the user registartion event for log track user registration metric.

### Login Flow

- `Auth Service` validate the request to `Auth DB`.
- `Auth Service` create session or auth token.

### Browse & Search Flow

- Fetch product list from `Product Service`.
- Handle search from `Search Service`.
- Show recommendation producs from `Recommendation Service`.

### Add to Cart Flow

- `Order Service` validates product to `Product Service`.
- `Order Service` store the cart on `Redis`.

### Checkout Flow

- `Order Service` validates cart item.
- `Order Service` call `Product Service` to reserve product.
- `Order Service` call `Transport Service` to calculate shipping cost.
- `Order Service` call `Payment Service` to process the payment.
- `Payment Service` call webhook to `Order Service`
- `Order Service` publish event to `Kafka`.

### Order Processing Flow

- `Product Service` consume order event from `Kafka` to deduct stock.
- `Notification Service` consume order event from `Kafka` to send confimation notification.
- `Transport Service` consume order event from `Kafka` to create shipping.
- `Analytic Service` consume order event from `Kafka` for processing.

### Shipping Flow

- `Transport Service` retreive websockt from shipping provider.
- `Transport Service` publish event to `Kafka`.
- `Order Service` consume event from `Kafka` to update the order status.

## Tech stack

Here are the tech stacks of choices, the reason why choose it, and the altenative(s).

### API Gateway

[Traefik](https://github.com/traefik/traefik)

- The most popular API gateway (https://ossinsight.io/collections/api-gateway/trends)
- Open source
- All-in-one for API Gateway and proxy
- Built-in auto-discovery of services with Kubernetes
- Observability follows OpenTelemetry semantic conventions (portable, no vendor lock-in)

Alternative(s): [Kong](https://github.com/Kong/kong), [Envoy](https://github.com/envoyproxy/envoy)

### Infastructure

[Kubernetes](https://github.com/kubernetes/kubernetes)

- Well known Container Orchestration
- Open source
- Vast ecosystem and third-party extensions

Alternative(s): [Nomad](https://github.com/hashicorp/nomad), [Docker Swarm](https://docs.docker.com/engine/swarm/)

### Database

[Postgres](https://github.com/postgres/postgres)

Alternative(s): [MariaDB](https://github.com/MariaDB/server)

[ClickHouse](https://github.com/ClickHouse/ClickHouse)

Alternative(s): [TimescaleDB](https://github.com/timescale/timescaledb)

[Redis](https://github.com/redis/redis)

Alternative(s): [Valkey](https://github.com/valkey-io/valkey), [Dragonfly](https://github.com/dragonflydb/dragonfly), [Memcached](https://github.com/memcached/memcached)

### Event Streaming

[Kafka](https://github.com/apache/kafka)

Alternative(s): [RabbitMQ](https://github.com/rabbitmq/rabbitmq-server), [Redpanda](https://github.com/redpanda-data/redpanda)

### Search Engine

[Elasticsearch](https://github.com/elastic/elasticsearch)

Alternative(s): [Meilisearch](https://github.com/meilisearch/meilisearch), [Typesense](https://github.com/typesense/typesense)

### Logging & Monitoring

[Loki](https://github.com/grafana/loki), [Grafana](https://github.com/grafana/grafana), [Tempo](https://github.com/grafana/tempo), [Mimir](https://github.com/grafana/mimir).

Alternative(s): [Elasticsearch](https://github.com/elastic/elasticsearch), [Logstash](https://github.com/elastic/logstash), [Kibana](https://github.com/elastic/kibana)

## Database Schema

## (Some) Potential Related API

## References

- https://medusajs.com/blog/ecommerce-architecture
- https://ngrok.com/blog/reverse-proxy-vs-api-gateway
- https://microservices.io/patterns
- https://konghq.com/blog/enterprise/why-kong-is-the-best-api-gateway
- https://www.qovery.com/blog/9-key-reasons-to-use-or-not-kubernetes-for-your-dev-environments