# Flight Itinerary Reconstruction API
A Go-based web API that reconstructs complete flight itineraries from a collection of flight ticket pairs. Given an array of source-destination airport code pairs, the API returns the complete travel sequence from the first departure to the final destination.

## Features

- **Single Endpoint**: POST `/api/v1/itinerary/reconstruct` accepts JSON payload with flight tickets
- **Itinerary Reconstruction**: Automatically determines the correct travel sequence
- **Error Handling**: Comprehensive validation and error reporting
- **Health Check**: GET `/health` endpoint for service monitoring
- **CORS Support**: Enabled for cross-origin requests during development
- **Comprehensive Testing**: Unit tests for core functionality and HTTP handlers

## API Specification

### Reconstruct Itinerary

**Endpoint**: `POST /api/v1/itinerary/reconstruct`

**Request Body**:
```json[
    ["LAX", "DXB"],
    ["JFK", "LAX"], 
    ["SFO", "SJC"],
    ["DXB", "SFO"]
  ]
```

**Response** (Success):
```json
   ["JFK", "LAX", "DXB", "SFO", "SJC"]
```

**Response** (Error):
```json
{
  "error": "Error description here"
}
```

### Health Check

**Endpoint**: `GET /api/v1//health/status`

**Response**:
```json
{
  "status": "healthy",
  "service": "flight-itinerary-go"
}
```

### Installation Steps

#### Option 1: Local Development

1. **Clone or create the project directory**:
   ```bash
   git clone https://github.com/code-elixrs/flight-itinerary-go.git
   mkdir flight-itinerary-go
   cd flight-itinerary-go
   ```

2. **Setup dependencies on local**:
   ```bash
   make deps
   ```
3. **Run the application on local**:
   ```bash
   make run
   ```
   

4. **Setup on docker**:
   ```bash
   make docker-build
   ```

5. **Run the application on docker**:
   ```bash
   make docker-run
   ```

   The server will start on port 8080. You should see output similar to:
   ```
   {"level":"info","timestamp":"2025-06-29T06:39:09.207+0400","caller":"cmd/main.go:40","msg":"Initializing..."} http server started on [::]:8080  
   ```

#### Option 2: Docker (Recommended)

1. **Clone or create the project directory**:
   ```bash
   mkdir flight-itinerary-go
   cd flight-itinerary-go
   ```

2. **Build and run with Docker Compose** (easiest):
   ```bash
   docker-compose up --build
   ```

3. **Or build and run with Docker directly**:
   ```bash
   # Build the image
   docker build -t flight-itinerary-go .
   
   # Run the container
   docker run -p 8080:8080 --name flight-api flight-itinerary-go
   ```

4. **Run in detached mode** (background):
   ```bash
   docker-compose up -d --build
   ```

6. **View logs**:
   ```bash
   docker-compose logs -f flight-api
   ```

7. **Stop the application**:
   ```bash
   docker-compose down
   ```

### Docker Management Commands

**View running containers**:
```bash
docker ps
```

**Stop and remove container**:
```bash
docker stop flight-api
docker rm flight-api
```

**Remove image**:
```bash
docker rmi flight-itinerary-go
```

**Check container health**:
```bash
docker inspect --format='{{.State.Health.Status}}' flight-api
```

## Testing

### Run Unit Tests

**Local Development**:
```bash
make test-unit # for unit test
make test-integration # for integration tests
make test # to test both
```


This will run comprehensive tests covering:
- Core itinerary reconstruction logic
- HTTP handler functionality  
- Error handling scenarios
- Edge cases (empty input, invalid formats, disconnected flights)

### Manual Testing with curl

**Test the main endpoint**:
```bash
curl -X POST http://localhost:8080/reconstruct-itinerary \
  -H "Content-Type: application/json" \
  -d '[
      ["LAX", "DXB"],
      ["JFK", "LAX"],
      ["SFO", "SJC"], 
      ["DXB", "SFO"]
    ]'
```

**Expected Response**:
```json
{
  ["JFK", "LAX", "DXB", "SFO", "SJC"]
}
```

**Test health check**:
```bash
curl http://localhost:8080/api/v1/health/status
```

**Test error handling**:
```bash
curl -X POST http://localhost:8080/api/v1itinerary/reconstruct \
  -H "Content-Type: application/json" \
  -d '[
      ["NYC", "LAX"],
      ["SFO", "DEN"]]'
```
## Project Structure

```
flight-itinerary-api/
├── cmd
  ├── main.go              # Main application code
├── internal
  ├── handler
    ├── itinerary_handler.go
    ├── itinerary_handler_test.go
  ├── logger
    ├── logger.go
  ├── middleware
    ├── logger.go
    ├── validator.go
  ├── model
    ├── itinerary.go
    ├── itinerary_test.go
  ├── service
    ├── itinerary_service.go
    ├── itinerary_service_test.go
├── pkg
  ├── errors
    ├── error.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── Makefile
├── README.md
├── integration_test.go
```

## Error Handling

The API handles various error conditions:

- **Invalid JSON**: Malformed request payload
- **Missing Fields**: Required `tickets` field not provided
- **Invalid Ticket Format**: Tickets without exactly 2 elements
- **Disconnected Flights**: Tickets that don't form a continuous path
- **Circular Routes**: Tickets that form cycles without clear starting point

All errors return appropriate HTTP status codes and descriptive error messages.
