# WeCredit API

WeCredit is a Go-based API for managing user accounts and authentication. It provides endpoints for user registration, login, and other related functionalities.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Database Migrations](#database-migrations)
- [Endpoints](#endpoints)
- [Twilio Configuration](#twilio-configuration)
- [License](#license)

## Features
- User registration and login
- JWT-based authentication
- OTP verification using Twilio
- Swagger API documentation
- Database migrations

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/weCredit.git
   cd weCredit
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create environment files:
   ```bash
   make setup
   ```

4. Run database migrations:
   ```bash
   make migrate
   ```

## Configuration

Before running the application, configure the environment variables in the `.env` file. You can use the `sample.env` as a reference.

### Example `.env` file
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5433
DB_USERNAME=db_username
DB_PASSWORD=your_password
DB_DATABASE_NAME=your_database_name

# Application Configuration
APP_PORT=8080
AUTH_SECRET=secret
AUTH_EXPIRY_PERIOD=3600

# Swagger Configuration
SWAGGER_HOST_SCHEME=http
SWAGGER_USERNAME=your_username
SWAGGER_PASSWORD=your_password

# Twilio Configuration
ACCOUNT_SID=your_twilio_account_sid
AUTH_TOKEN=your_twilio_auth_token
TWILIO_NUMBER=your_twilio_phone_number
```

## Usage

To start the application, run:
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`.

## API Documentation

API documentation is generated using Swagger. You can access it at `http://localhost:8080/swagger/index.html` after starting the server.

## Endpoints

### User Registration
- **POST** `/users`
  - **Description**: Create a new user with the provided details.
  - **Request Body**: 
    ```json
    {
      "full_name": "John Doe",
      "user_name": "+919876543210",
      "role": "USER"
    }
    ```
  - **Responses**:
    - `201 Created`: User successfully registered.
    - `400 Bad Request`: Validation errors.

### Initialize Login
- **POST** `/users/init/login`
  - **Description**: Initiate the login process by sending an OTP to the user's registered phone number.
  - **Request Body**: 
    ```json
    {
      "user_name": "+919876543210"
    }
    ```
  - **Responses**:
    - `200 OK`: OTP sent successfully.
    - `404 Not Found`: User not found.

### User Login
- **POST** `/users/login`
  - **Description**: Authenticate a user using provided credentials and OTP.
  - **Request Body**: 
    ```json
    {
      "user_name": "+919876543210",
      "otp": "123456"
    }
    ```
  - **Responses**:
    - `200 OK`: Successful login with JWT token.
    - `401 Unauthorized`: Invalid credentials or OTP.

### Get User by ID
- **GET** `/users/:id`
  - **Description**: Retrieve user details by ID.
  - **Responses**:
    - `200 OK`: User details.
    - `404 Not Found`: User not found.

## Twilio Configuration

To send OTPs using Twilio, you need to set up your Twilio account and obtain the following credentials:

1. **Account SID**: Your Twilio account identifier.
2. **Auth Token**: Your Twilio account authentication token.
3. **Twilio Phone Number**: The phone number you will use to send OTPs.

### Sending OTPs

The application uses Twilio's API to send OTPs to users during the login process. Ensure that you have configured the Twilio settings in your `.env` file as shown above.

### Example of Sending OTP
The OTP sending functionality is integrated into the user login process. When a user attempts to log in, an OTP will be generated and sent to their registered phone number using Twilio.

## Database Migrations

Database migrations are managed using Goose. To run migrations, use the following command:
```bash
sh script/migration.sh
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
