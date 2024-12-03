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

## **Second Approach: Persistent Storage with PVC**

In this approach, instead of using in-memory data storage, we store the data persistently using Persistent Volume Claims (PVC). This ensures durability across service restarts and better long-term data retention.

### **Storage Options**

We have considered three storage options:

1. **Redis with TTL**: Redis is an in-memory data store with time-to-live (TTL) support, making it ideal for short-lived data.
   - **Disadvantages**: While Redis offers fast read/write access, it is not ideal for counting scenarios that involve high throughput and heavy write operations. Its performance can degrade with large datasets or high write operations. Additionally, managing TTL for counts becomes complex, and its persistence layer is not as straightforward as with SQL databases.

2. **SQL (PostgreSQL)**: PostgreSQL is an ACID-compliant relational database that can handle heavy write operations and structured data storage.
   - **Advantages**: PostgreSQL is suitable for storing request counts as it ensures data integrity and can handle heavy writes efficiently. It is ideal for systems requiring frequent updates and high consistency. With optimized schema design, PostgreSQL offers scalability for counting and indexing operations.
   - **Disadvantages**: Under extremely high write loads, SQL databases can face performance bottlenecks if not properly optimized. However, since we are not using batch processing and request counts are relatively small, PostgreSQL is an excellent choice for persistence.

3. **NoSQL (MongoDB)**: NoSQL databases offer flexibility and scalability but are generally not ideal for counting use cases.
   - **Disadvantages**: NoSQL databases lack built-in support for atomic operations like counting and aggregation, making them less efficient for our scenario. SQL or Redis would be better suited for efficient and precise counting operations.

---

## **Decision to Proceed with SQL Approach**

After evaluating the storage options, we proceeded with the **SQL (PostgreSQL)** approach. PostgreSQL’s capability to handle heavy write operations and structured data storage makes it a suitable choice for managing request counts and providing durability across service restarts.

---

## **Extensions Implemented**

### **Extension 1: Changing GET Request to POST**

The original approach used a `GET` request to accept new data entries. For better handling of larger payloads and to adhere to RESTful principles, the request method was changed to **POST**. This modification enables more efficient data handling, especially with a large number of incoming requests.

---

### **Extension 3: Data Dump to Kafka**

Initially, data was dumped into a log file. However, to enhance scalability and integrate with a more event-driven architecture, we transitioned to **Kafka** for the data dump. Kafka offers a distributed, fault-tolerant, high-throughput messaging system that processes real-time data streams more efficiently than writing to a log file.

This change enables better integration with other services, improved system monitoring, and enhanced scalability for future needs.

---

This updated approach combines persistent storage with PostgreSQL and introduces Kafka for better event streaming, making the system more scalable and efficient for handling large volumes of data.


## **Extension 2: Horizontal Scaling with Load Balancer**

To scale the service horizontally and handle increased traffic, we can deploy multiple instances of the service and use a **load balancer** to distribute incoming requests efficiently.

### **Solution:**
1. **Load Balancer (Nginx)**: 
   - Nginx can be used to distribute incoming traffic evenly across two or more instances of the service.
   - This ensures that the service can handle more traffic by load balancing requests between multiple service instances.

2. **Data Consistency**: 
   - All service instances will be polling from the same PostgreSQL database, ensuring that the count values remain consistent across instances.

3. **Kafka Message Duplication**:
   - Since the service sends a message once every minute, even if multiple instances send the same message due to the scheduler, the consumer will only process one message for that time window.
   - By handling Kafka message duplication at the consumer side, we ensure that only one message is processed, regardless of how many instances send the same data.

This solution enables horizontal scaling of the service with a load balancer while maintaining consistent data and preventing Kafka message duplication.

---