# docker-swarm-service-status
[![](https://images.microbadger.com/badges/image/albertogviana/docker-swarm-service-status.svg)](https://microbadger.com/images/albertogviana/docker-swarm-service-status "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/albertogviana/docker-swarm-service-status.svg)](https://microbadger.com/images/albertogviana/docker-swarm-service-status "Get your own version badge on microbadger.com")
[![Go Report Card](https://goreportcard.com/badge/github.com/albertogviana/docker-swarm-service-status)](https://goreportcard.com/report/github.com/albertogviana/docker-swarm-service-status)

## Motivation

The idea behind of this project was to provide an easy way to get the service deployment information on Docker Swarm Cluster. This service would be used with Jenkins where in a pipeline and I would be able to call my service and guarateer that my service was properly deployed or if for some reason failed or was rollback.

## TODO
- [ ] Add Prometheus metrics
- [ ] Improve documentation  

## Endpoint

### Deployment Status (/v1/docker-swarm-service-status/deployment-status/{service}/{image})

The Deployment Status endpoint is available on `/v1/docker-swarm-deployment-status/{service}/{image}` and it requires the parameters:
- `service` is related to the service name on Docker
- `image` is the image deployed in the cluster.
    - The `image` parameter must be sent enconded with `base64`.

You can easily create a base64 using the command line, with the command below: 
```
echo -n "albertogviana/docker-routing-mesh:1.0.0" | base64
```

### Service Status (/v1/docker-swarm-service-status/service-status/{service})

The Deployment Status endpoint is available on `/v1/docker-swarm-service-status/{service}` and it requires the parameters:
- `service` is related to the service name on Docker
