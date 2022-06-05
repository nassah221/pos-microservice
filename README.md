# pos-microservice

Point of Sales APIs implemented using microservice architecture

## Motivation

I'm trying to follow domain driven design style with this particular project.
The goal is to write a single microservice for each domain context (i.e. order, sales, products-inventory etc).

This project is mainly an exercise for structuring multiple microservices, experimenting with go-kit, inter-service communication and mongodb.

I'll try to keep a flat-directory structure. Each microservice will reside in it's own directory at the root.

## References

- My own [pos-crud](https://github.com/nassah221/pos-crud/blob/main/sales-api/api/cashier.go) project which is written in REST style

- The [goddd](https://github.com/marcusolsson/goddd) sample project

- Some of the example code from go-kit [examples](https://github.com/go-kit/examples)