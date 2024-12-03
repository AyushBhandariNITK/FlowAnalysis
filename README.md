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


# **Flow Analysis Service with In-Memory Map Approach**

## **Overview**

This repository defines a high-performance backend service for flow analysis that manages data using an in-memory map mechanism. The service uses two maps (primary and secondary) to efficiently store, process, and count unique entries. Every minute, the primary and secondary maps are swapped, ensuring that the system remains performant while processing data concurrently.

### **Key Features**
- **Concurrent Data Management**: Utilizes a thread-safe map to store and manage data entries, preventing any service calls for data operations.
- **Primary and Secondary Map**: Uses two maps to store and handle data. The primary map is actively used for read and write operations, while the secondary map is cleared and prepared for the next cycle.
- **Periodic Map Swap**: Every 60 seconds, the primary and secondary maps are swapped. This ensures that the active map is always available for new entries, while the inactive map can be processed for unique entry counting.
- **Logging**: Detailed logs at various levels (Info, Warn, Error, and File) for easy debugging and monitoring.

--- 

## **Why In-Memory Approach is Used**

The in-memory approach is implemented to efficiently manage high request counts while minimizing data storage overhead. This approach leverages concurrent maps to avoid service calls and offers several benefits:

- **Efficient Memory Usage**: The service stores only a small amount of data for each entry, and with two maps (primary and secondary), memory usage is controlled. These maps are rotated every minute to prevent memory overload or spikes, and the inactive map is cleared for the next cycle. This ensures consistent memory usage over time.
  
- **High Performance**: In-memory storage allows for fast data access and processing, enabling the service to handle high volumes of requests with low latency. By keeping the data in memory, we avoid the bottlenecks associated with disk-based storage or external service calls.
  
- **Data Rotation**: The primary and secondary map system allows the active data to be processed while the inactive map is cleared and prepared for the next cycle. This ensures that only the most recent data is kept active, while older data is discarded after each map swap.

### **Disadvantage**
- **Data Loss on Restart**: As all data is stored in memory, a service restart will result in data loss. There is no persistent storage for the maps, meaning that the request count and any ongoing sessions will be reset upon a restart.

---

This explains both the advantages of using the in-memory approach for high performance and memory efficiency, as well as the trade-off of potential data loss in case of a service restart.