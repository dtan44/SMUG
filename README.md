<img src="" alt="SMUG"/>

# SMUG
### Simple Microservice Universal Gateway

SMUG is a simple microservice structure implemented in Golang. It was designed to be easy to use and configure by having few dependencies while covering general usage of a microservice gateway.

## Getting Started

To get started, make sure you have [Golang](https://golang.org/doc/install) installed and configured, then just `go get github.com/dtan44/SMUG`.

### SMUG Structure

SMUG is organized into two main services: **Service Registration** and **Service Discovery**. Service registration provides backend functionality to register and deregister services. Service Discovery allows clients (i.e. web browser, mobile phone, etc...) to find and use the available services.

## Technologies Used

This project is implemented in Golang.
[Viper](https://github.com/spf13/viper) was used for 12-Factor app configuration. [Glide](https://github.com/Masterminds/glide) was used for dependency management.
[API Blueprint](https://apiblueprint.org/) was used to create the documentation.

## Authors

* **Dylan Tan** - [dtan44](https://github.com/dtan44)

## License

This project is licensed under the MIT License
