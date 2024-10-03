# Clean Architecture CLI

Clean Architecture CLI is a command-line tool for quickly generating Go projects with a Clean Architecture structure, using Fiber as the web framework and MongoDB Atlas as the database.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Dependencies](#dependencies)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Features

- Generates a complete project structure following Clean Architecture principles
- Sets up a basic user registration and login system
- Integrates Fiber web framework
- Configures MongoDB Atlas as the database
- Implements JWT-based authentication
- Provides a solid foundation for building scalable Go applications

## Installation

To install the Clean Architecture CLI, follow these steps:

1. Clone the repository:
```
git clone https://github.com/ankit-apk/cleanarch-cli.git
```


2. Change to the project directory:
```
cd cleanarch-cli
```


3. Build the CLI tool:
```
go build -o cleanarch-cli cmd/cleanarch-cli.go
```


4. (Optional) Move the binary to a directory in your PATH for global access:
```
sudo mv cleanarch-cli /usr/local/bin/
```


## Usage

To generate a new project, use the following command:

```
cleanarch-cli -name -module
```


Example:

```
cleanarch-cli -name myproject -module github.com/username/myproject
```


This will create a new directory `myproject` with the Clean Architecture structure and all necessary boilerplate code.

## Project Structure

The generated project will have the following structure:

```
myproject/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── user.go
│   ├── usecase/
│   │   └── user_usecase.go
│   ├── repository/
│   │   └── user_repository.go
│   └── handler/
│       └── user_handler.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── mongodb.go
│   └── auth/
│       └── jwt.go
├── go.mod
└── .env

```


## Dependencies

The generated project uses the following main dependencies:

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - MongoDB driver for Go
- [jwt-go](https://github.com/dgrijalva/jwt-go) - JWT implementation for Go
- [godotenv](https://github.com/joho/godotenv) - For loading environment variables
- [bcrypt](https://golang.org/x/crypto/bcrypt) - For password hashing

## Configuration

The generated project uses a `.env` file for configuration. Make sure to update the following variables in the `.env` file:

```
MONGO_URI=your_mongodb_atlas_connection_string JWT_SECRET=your_jwt_secret
```


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
