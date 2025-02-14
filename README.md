# gRPC Web Chat

## Overview
This is a web-based chat application that utilizes gRPC for communication between the backend services and an API gateway for handling client requests.

## Key Features
- User Authentication - Register, login, and manage user accounts securely.

- Real-time Messaging - Communicate with other users via private or group chats.

- Group Chat Management - Create, join, and manage chat groups with assigned roles.

- Contact Management - Add, remove, and organize contacts for easier communication.
## Demo

- https://youtu.be/IpN2K1RcIFg
## Technologies Used

- Frontend: HTML, CSS, JavaScript
- Backend: Go, Echo Framework, gRPC, WebSocket
- Database: My SQL
## Installation
Before setting up the project, ensure you have the following installed:

- **Go** (>=1.2) [Download here](https://go.dev/dl/)
- **MySQL**  [Download here](https://dev.mysql.com/downloads/)

📥 Clone the Repository

```sh
git clone https://github.com/quanbin27/gRPC-Web-Chat
cd yourproject
```
Open MySQL and create the database:
```sh
CREATE DATABASE IF NOT EXISTS yourdatabase;
```
Update the .env file with your database credentials:
```sh
DSN=root:12345678@tcp(127.0.0.1:3306)/yourdatabase?charset=utf8mb4&parseTime=True&loc=Local
```

▶️ Run the Project
```sh
go run cmd/grpc_server/main.go
go run cmd/api_gateway/main.go
```
