# Afro Vintage Backend 

This is the backend service for the Afro Vintage platform — built with Go (Gin), MongoDB, and Clean Architecture. It powers both the mobile (Flutter) and web (Next.js) apps.

## Tech Stack

- Go + Gin
- MongoDB
- Clean Architecture
- JWT Auth
- Stripe (simulated)

## Running with Docker

To build and run the application using Docker, follow these steps:

1. **Build the Docker Image**:
   ```bash
   docker build -t afro-vintage-backend .
   ```

2. **Run the Docker Container**:
   ```bash
   docker run -p 8080:8080 --env-file .env afro-vintage-backend
   ```

3. **Using Docker Compose**:
   If you have a `docker-compose.yml` file, you can use the following commands:
   - Build and start the services:
     ```bash
     docker-compose up --build
     ```
   - Stop the services:
     ```bash
     docker-compose down
     ```

4. **Access the Application**:
   - The application will be available at `http://localhost:8080`.
   - Ensure MongoDB is running on `localhost:27017` if not using Docker Compose.

## Running the Project Locally

To run this project on your machine, follow these steps:

### Prerequisites
1. Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/install/).
2. Clone the repository:
   ```bash
   git clone <repository-url>
   cd backend/afro-vintage-backend
   ```

### Steps to Run
1. Start the services using Docker Compose:
   ```bash
   docker-compose up
   ```
2. The application will be accessible at `http://localhost:8080`.

### Stopping the Services
To stop the services, press `Ctrl+C` in the terminal and run:
```bash
docker-compose down
```

### Notes
- The MongoDB database is initialized with the name `afro_vintage`.
- Ensure that port `8080` and `27017` are not in use by other applications.
