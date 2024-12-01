# **Flow Analysis Service**

## **Overview**
The Flow Analysis Service is a high-performance application for handling concurrent data operations, maintaining thread safety, and periodically managing data maps. It features an efficient mechanism for managing unique data entries with clear separation of active and inactive states.

---

## **Features**
- **Thread-Safe Data Management**: Uses atomic operations and synchronization primitives to ensure safe concurrent access.
- **Periodic Data Refresh**: Automatically switches and clears inactive data maps every 60 seconds.
- **Logging**: Provides detailed logs at various levels (Info, Warn, Error, and File).
- **Custom Error Handling**: Predefined error types for better API responses.
- **HTTP API**: A flexible interface for accepting and processing data entries.

---

## **Key Components**

### **FlowMap**
A thread-safe structure for storing and managing key-value pairs. It ensures operations are synchronized with the active state of the service.

### **Map Switching Mechanism**
- **Purpose**: Alternates between two data maps (`MapA` and `MapB`) every 60 seconds.
- **Workflow**:
  1. Temporarily halts map operations.
  2. Switches the active and inactive maps.
  3. Logs the count of unique entries from the inactive map.
  4. Clears the inactive map for future use.

### **Logging**
Provides real-time insights and writes detailed logs to `unique_requests.log`. Supports multiple log levels:
- **Info**: General information.
- **Warn**: Warnings.
- **Error**: Critical issues.
- **File**: Appends to the log file.

### **HTTP API**
- **Endpoint**: `GET /api/verve/accept`
  - **Query Parameters**:
    - `id`: A unique identifier to store.
    - `endpoint` (optional): A URL to send the count of unique entries.
- **Functionality**:
  - Stores the `id` in the active map.
  - Sends the unique entry count to the specified endpoint if provided.

---

## **How It Works**

### **Initialization**
- Two `FlowMap` instances (`MapA` and `MapB`) are initialized.
- A separate goroutine handles periodic map switching.
- The HTTP server starts and listens on port `5010`.

### **Map Lifecycle**
1. **Active Map**: Used for storing and processing current data.
2. **Inactive Map**: Cleared and prepared for the next active cycle.

---

## **Getting Started**

### **Requirements**
- Go 1.18+ 
- Dependencies:
  - [Echo](https://github.com/labstack/echo): Web framework for HTTP APIs.
  - [Concurrent Map](https://github.com/orcaman/concurrent-map): Thread-safe map implementation.
  - [Klog](https://github.com/kubernetes/klog): Logging library.

### **Run the Application**
1. Clone the repository.
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build and run the application:
   ```bash
   go run main.go
   ```
4. Access the API at `http://localhost:5010/api/verve/accept`.

---

## **Usage**

### **API Example**
```bash
curl "http://localhost:5010/api/verve/accept?id=123&endpoint=http://example.com"
```

- **Without `endpoint`**: Adds the `id` to the active map.
- **With `endpoint`**: Sends the unique count of entries to the specified endpoint.

---

## **Error Handling**

### **Custom Errors**
- **Connection Refused**: Indicates the server connection failed.
- **Invalid Endpoint**: Triggered for 404 responses.
- **Bad Request**: Triggered for malformed requests.

---

## **Logs**
- Real-time logs are available in the console and written to `unique_requests.log`.
- Example log entry:
  ```
  2024-11-30 12:00:00 Number of unique entries: 100
  ```

---

## **System Flow**

### **Data Flow Diagram**
```plaintext
+--------------------------+
| HTTP API Request         |
+-----------+--------------+
            |
            v
+-----------+--------------+
| Store ID in Active Map   |
| - Optional: Send Unique  |
|   Count to Endpoint      |
+-----------+--------------+
            |
            v
+-----------+--------------+
| Map Switch Mechanism     |
| - Log Unique Count       |
| - Clear Inactive Map     |
+--------------------------+
```

### **Periodic Map Management**
- Executes every 60 seconds.
- Alternates between `MapA` and `MapB`.
