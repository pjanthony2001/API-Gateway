# API-Gateway

# 1 	Introduction

## 1.1 Project Overview

The purpose of this design document is to outline the comprehensive design of the API Gateway, implemented using CloudWeGo Projects. This document aims to provide a detailed understanding of the architecture, components, and implementation plan for the API Gateway project.

## 1.2 Objectives

The primary objectives of the API Gateway project include:

Accepting HTTP requests with JSON-encoded bodies: The API Gateway will be responsible for accepting incoming HTTP requests from clients. These requests will include JSON-encoded bodies containing the necessary data for processing.

Using the Generic-Call feature in Kitex: The API Gateway will leverage the Generic-Call feature provided by the Kitex framework. This feature allows translating JSON requests received by the API Gateway into the Thrift binary format. This translation enables efficient communication and interoperability with the backend RPC servers.

Integrating a load balancing mechanism: To ensure optimal resource utilization and high availability, a load balancing mechanism will be integrated into the API Gateway. This mechanism will distribute incoming requests among multiple backend RPC servers, effectively managing the workload and preventing any single server from becoming a bottleneck.

Integrating a service registry and discovery mechanism: The API Gateway and RPC servers will be integrated with a service registry and discovery mechanism. This integration will enable dynamic service registration and discovery, allowing the API Gateway to locate and communicate with the available RPC servers efficiently.

Develop backend RPC servers using Kitex: This will allow us to test API Gateway's ability to handle and process requests, validate the integration with the backend services, and assess the overall performance and responsiveness of the system. This testing phase will help identify any potential issues or bottlenecks early on and allow for necessary improvements and optimizations to be made.



# 2 System Architecture

## 2.1 High Level
The API Gateway accepts JSON-encoded HTTP requests, translates them into Thrift binary format using Kitex's Generic-Call feature, and forwards them to backend RPC servers, with the help  of the service registry that enables service discovery. The load balancing mechanism ensures an even distribution of requests.

## 2.2 Components

### Hertz Framework: 
Hertz is a Golang HTTP framework. It serves as the entry of the API Gateway, that accepts JSON-encoded HTTP requests.

### Kitex Framework:
Kitex is a high-performance, scalable RPC framework that supports Thrift-based communication. The API Gateway will leverage:

Kitex's Generic-Call feature for JSON to Thrift translation. 
Kitex’s Round Robin Load Balancing feature to distribute requests amongst backend servers. 
Kitex’s Service Registry feature to seamlessly integrate new servers into the network

## 2.3 Interaction Flow

The interaction flow within the system is as follows:

1. An HTTP request with a JSON-encoded body is received by the API Gateway.

2. The API Gateway uses the Generic-Call feature of Kitex to translate the JSON request into Thrift binary format.

3. The API Gateway consults the load balancer to determine the appropriate backend RPC server to forward the request.

4. The load balancer selects a backend server based on the configured load balancing strategy and forwards the request to that server.

5. The backend RPC server receives the request, processes it, and sends the response back to the API Gateway.

6. The API Gateway translates the response from Thrift binary format to JSON and sends it back as an HTTP response to the client.



# 3 How to run 

## 3.0 Prerequisites
1) Install Golang https://go.dev/doc/install
2) Install Hertz https://www.cloudwego.io/docs/hertz/getting-started/
3) Install Kitex https://www.cloudwego.io/docs/kitex/getting-started/

## 3.1 Set up the Hertz Server
1) Open a terminal window and navigate to `hertz-server` directory
2) `go run .` to start the server on `localhost:8080`

## 3.2 Set up the Kitex Server
1) Open a terminal window and navigate to the `kitex_server` directory
2) `go run .` to start the server on `localhost:8888`

## 3.3 Send a HTTP Request
1) Send a GET request with the following command `curl -X GET localhost:8080 -d <JSON HERE>`
2) JSON should have a component `Msg : <string>`, following the structure in the `example_service.thrift` idl

# 4 Design Document
The design document, which includes the projected timeline at the end, can be accessed here : 
https://docs.google.com/document/d/19cSJfAP8_TKRUjOMC_lAO5g-itLNmyKLcib8Qp2RyqQ/edit
