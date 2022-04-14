<div id="top"></div>

<div align="center">

  <h3 align="center">REST API and gRPC example</h3>

  The purpose of this repo is to make some boilerplate code for a REST API(using Echo framework) and gRPC.

  The Fibonacci sequence is used as sample data.
  The APIs take two integers(X,Y) as parameters and return a part of the Fibonacci sequence from X to Y indices.

  Redis is used to cache the Fibonacci sequence to not recalculate it for each request.
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Summary</summary>
  <ol>
    <li>
      <a href="#about">About</a>
      <ul>
        <li><a href="#dependencies">Dependencies</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting started</a>
      <ul>
        <li><a href="#build-and-start">Build and start</a></li>
        <li><a href="#testing">Testing</a></li>
      </ul>
    </li>
    <li><a href="#project-structure">Project structure</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About

The project contains a REST API and a gRPC servers.

In response to a request consisting of two uint64s(X,Y) the services send an array of numbers with indices from X to Y in the Fibonacci sequence.

To build and start the project Docker and Docker-Compose are needed.
If you build the project using Docker the unit tests are launched automatically.
The integration testing for both REST and gRPC services is launched by a separate Dockerfile and is also written in Golang.

<p align="right">(<a href="#top">back to top</a>)</p>



### Dependencies

All you need to run the program is:
* [Docker](https://www.docker.com/)
* [Docker compose](https://docs.docker.com/compose/install/)
Make sure to install the latest version.

Inside the Docker containers following software is used:
* [Golang](https://golang.org/)
* [Redis](https://redis.io/)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting started

The main program is launched by docker-compose, the integration testing can be started by a separate script in the "tests" folder once the main program is started.

### Build and start

Before you start you will need Docker и Docker-Compose installed.
The links to download all the tools required are available here: <a href="#dependencies">Dependencies</a>

After you've installed all the required software you can simply run the `docker-compose yml` script:
   ```
   docker-compose up
   ```


### Testing

All the unit tests are started automatically before building and starting the program in Dockerfile.

The integration testing can be found in the **tests** directory.
To build and run the integration testing Docker container you can simply run run-tests.sh.
You can also launch the integration test from the **tests** directory by executing following commands:

   ```
   docker build -t fibonacci -f Dockerfile ..

   docker run --network fibonacci fibonacci

   ```

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- STRUCTURE -->
## Project structure
The root directory of the project contains `docker-compose.yml` and `Dockerfile` scripts to build and start the program as well as the main.go file that reads the config file and starts the APIs.

* **app** contains the golang code, the program itself
    * **databases** contains code to interact with databases(only Redis in this case)
    * **gprc** contains the gRPC API code
    * **rest** contains the REST API code
    * **tools** contains the Fibonacci sequence generator

* **config** contains the config file
* **proto** contains the .proto file for gRPC
* **tests** contains the integration test and the Dockerfile script to build and lauch it

<p align="right">(<a href="#top">back to top</a>)</p>

