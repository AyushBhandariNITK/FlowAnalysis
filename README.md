Here's a README for your project that includes the details about the Docker Compose setup:

---

# **Flow Analysis Service with Docker Compose**

## **Overview**

This repository defines a multi-container Docker Compose setup to run a PostgreSQL database, Kafka message broker, and a Flow Analysis backend service. The service processes data entries and switches between in-memory or persistent data storage based on configuration. This setup is designed to provide a scalable and efficient environment for your backend services.

## **Docker Compose Setup**

The Docker Compose configuration contains three services:

1. **`db`**: A PostgreSQL database container to store and manage application data.
2. **`flowanalysis`**: A backend service that processes flow analysis, interacting with the PostgreSQL database.
3. **`kafka`**: A Kafka container for message brokering between services.

### **Key Features**
- **PostgreSQL** for data storage.
- **Kafka** for message brokering.
- **Flow Analysis** service that can switch between in-memory storage or persistent storage based on the `INMEMORY` configuration.

## **Getting Started**

### **Requirements**
- **Docker**: Make sure Docker is installed on your system.
- **Docker Compose**: Ensure you have Docker Compose installed.

### **Usage**

To start the application using Docker Compose, run the following command:

```bash
docker-compose up -d
```

This command will:
- Pull the necessary images.
- Create and start the containers in detached mode.

To stop the services and remove the containers:

```bash
docker-compose down
```

### **Environment Configuration**
The services are configured via environment variables defined in the Docker Compose file:

- **`INMEMORY`**: Set to `true` for in-memory data management (default). Set it to `false` for persistent storage. This variable determines whether data is stored in memory or on the filesystem.
  
```yaml
INMEMORY: true  # Set to true for in-memory mode, false for persistent
```

### **Service Details**

- **PostgreSQL Service (`db`)**:
  - The PostgreSQL database service is configured with the following settings:
    - User: `postgres`
    - Password: `password`
    - Database: `mydb`
    - Port: Exposed on `5432`.
  
  Volumes are mounted to ensure data persistence:
  - `./volumes/postgress/postgress_data:/var/lib/postgresql/data`
  - `./volumes/postgress/init.sql:/docker-entrypoint-initdb.d/init.sql`

- **Flow Analysis Service (`flowanalysis`)**:
  - Connects to the PostgreSQL service and performs data analysis.
  - The `INMEMORY` variable determines whether data is stored in memory or on disk.

- **Kafka Service (`kafka`)**:
  - Configures a Kafka instance with custom settings for KRaft mode and security protocols.
  - Exposes ports `9092` and `9094` for internal and external communication.

### **Volumes and Networking**
- Volumes are mounted using relative paths for persistent data storage:
  - PostgreSQL: `./volumes/postgress/postgress_data`
  - Kafka: `./volumes/kafka/kafka_data`

- Docker’s default network is used to ensure communication between the services.

---
