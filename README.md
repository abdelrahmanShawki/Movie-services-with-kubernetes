# MovieApp

MovieApp is a microservices-based application that provides a robust platform for managing and browsing movie-related data. This project demonstrates how to design, build, and deploy a scalable application using kubernetes , 
in addition to that : i use protocol buffers for service communication 

## Table of Contents

- [To DO](#TO-DO)
- [Project Overview](#project-overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [How to Run the Services with Docker](#how-to-run-the-services-with-docker)
- [License](#license)


## TO DO

1 - replacing http with grpc communication model

2 - using mysql to store services data

3 - deployment with kubernetes

4 - Unit and integration testing



## Project Overview

MovieApp is designed to offer users a seamless experience when browsing movie details, ratings, and reviews. It leverages microservices architecture to ensure modularity, scalability, and maintainability.

## Features

- **Microservices Architecture:** Each service is containerized and runs independently.
- **Service Discovery:** Utilizes Consul for service discovery and configuration.
- **Dockerized Environment:** All services can be easily deployed using Docker.
- **RESTful APIs:** Provides robust API endpoints for various functionalities.
- **Scalability:** Designed to adapt.

## Architecture

The application is composed of several microservices including:
- **Movie Service:** Handles movie data and details through communication with rating and metadata services.
- **rating Service:** Processes movie ratings.
- **metadata Service:** aggregate movie metadata.
- **Service Discovery:** Uses Hashicorp Consul to register and manage services for local testing , kubernetes for production .

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/abdelrahmanShawki/Movie-services-with-kubernetes
   ```

## How to Run the Services with Docker

To facilitate service discovery and inter-service communication, the project uses **Hashicorp Consul**. Follow these steps to set up Consul using Docker:

1. **Pull the Consul Docker image:**

   ```bash
   docker pull hashicorp/consul
   ```

2. **Run the Consul container:**

   ```bash
   docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
   ```

3. **View work done:**

   Open [http://localhost:8500](http://localhost:8500) to check service discovery and status.


## License

This project is licensed under the [MIT License](LICENSE).









