## Prerequisites

Before you start the program, you need to generate a private key and a public key. These keys are used for generate token.

Follow these steps to generate the keys:

1. Open a terminal.
2. Run the following commands to generate a 4096-bit RSA private key and a corresponding public key:

```bash
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:4096
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

3. After generating the keys, create a new directory named `certs` at the root level of your project.
4. Move the `private_key.pem` and `public_key.pem` files into the `certs` directory.

## Dockerizing the Application

To run the application in a Docker container, follow these steps:

1. Ensure Docker is installed on your machine. If not, you can download and install Docker from [here](https://docs.docker.com/get-docker/).

2. Open a terminal and navigate to the root directory of the project.

3. Run the following command to build the Docker image:

```bash
docker-compose build
```

4. After the build completes, start the application with the following command:

```bash
docker-compose up
```
