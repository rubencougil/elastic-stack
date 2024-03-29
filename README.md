# 📶 Elastic Stack Docker + Sample Go App

[![Actions Status](https://github.com/rubencougil/elastic-stack/workflows/Build/badge.svg)](https://github.com/rubencougil/elastic-stack/actions)

Elastic Stack which includes **Elasticsearch**, **Kibana**, **Filebeat** and **Metricbeat**. It comes with a very simple Go application that publishes a tiny HTTP API. **Go** app connects to a **MySQL** database. Also there's an **Nginx** that acts as the load balancer of the Go App to allow scale up/down the app containers.

![Clase 1 - Practica 3 (1)](https://user-images.githubusercontent.com/1073799/202543773-a0ba29fb-a558-417c-974f-60f5f1c662f6.jpg)

## How to run it?

For spinning up the stack:

`docker-compose up -d`

After all services are running, you can use following Go application endpoints to generate random data (:8080 is the exposed port of th Nginx connected to Go application instances):

- `curl http://localhost:8080/`
- `curl http://localhost:8080/create -d {}`

To scale up/down Go application:

`docker-compose up -d --scale app=3`

To change where to tog change `LOG_TO` env var in `.env`:

- `elasticsearch` to send logs directly to ES
- `stdout` to send logs to standard output

## ES Clustering

`docker-compose-es-cluster.yml` file will allow us to spin up an ES cluster con 3 nodes.

```
docker-compose -f docker-compose-es-cluster.yml up -d --remove-orphans
```

Once done, we can check the status of the cluster:

```
curl http://localhost:9200/_cluster/health?pretty
```

Now we create a new index with replication set to 2 and we add a new document:

```
curl -H "Content-Type: application/json" -XPUT http://localhost:9200/test -d '{"settings" : {"index" : {"number_of_shards" : 3, "number_of_replicas" : 2 }}}'
curl -H "Content-Type: application/json" -XPUT http://localhost:9200/test/docs/1 -d '{"name": "ruben"}'
```

For getting index and document distribution through the Cluster:

```
curl http://127.0.0.1:9200/_cat/indices?v
```

There's an extra service with an UI for admin the cluster:

`http://localhost:9100`

