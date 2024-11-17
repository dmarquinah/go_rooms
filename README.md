# Go Rooms

Go Rooms is a backend application written in Golang. It aims to create WebSocket rooms for users to communicate in real-time.

## Requirements

- Golang 1.22.2 or later

## Getting Started

1. **Clone the Repository**  
   ```bash
   git clone https://github.com/your-username/go-rooms.git
   cd go-rooms
   ```
2. **Run the Application**  
   `go run main.go`

## Features

- Create and Login Users.
- Authentication and Authorization.
- Real-time Server communication with Front-End.

## To-Do

### Core Functionality

- [ ] Set up basic WebSocket server
- [ ] Create and manage WebSocket rooms
- [ ] Allow users to join and leave rooms
- [ ] Broadcast messages to all users within a room
- [ ] Implement unique room identifiers

### User Management

- [ ] Track user connections and disconnections
- [ ] Handle user reconnections gracefully
- [ ] Allow users to send private messages within a room
- [ ] Implement username or ID-based user identification

### Security

- [ ] Validate incoming messages and data
- [ ] Prevent unauthorized access to rooms
- [ ] Add rate-limiting for message frequency
- [ ] Implement basic input sanitization

### Error Handling

- [ ] Handle unexpected WebSocket closures
- [ ] Return meaningful error messages to users
- [ ] Log server-side errors for debugging

### Optimization

- [ ] Optimize WebSocket performance for multiple rooms
- [ ] Test server with high loads and concurrent connections
- [ ] Minimize memory usage for each room and user connection

### Testing

- [ ] Create unit tests for room creation logic
- [ ] Test user connection and disconnection handling
- [ ] Test broadcasting and private messaging functionality
- [ ] Implement integration tests for the entire server flow

## Database Modeling
![model](img/db_model_0_1.png)

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.