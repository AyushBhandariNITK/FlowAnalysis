### **To run the project, follow these steps**  

### **Prerequisites**  
Ensure the following are installed on your system:  
1. **Docker**  
2. **Docker Compose**  

---

### **Steps to Run the Project**  

1. **Clone the repository**  
   Clone the project from the GitHub link:  
   ```bash
   git clone https://github.com/AyushBhandariNITK/FlowAnalysis
   cd FlowAnalysis
   ```
   **GitHub Repository**: [GitHub Link](https://github.com/AyushBhandariNITK/FlowAnalysis)  

2. **Prepare the images**  
   The project consists of three components:  
   - **Kafka**  
   - **Postgres**  
   - **FlowAnalysis (Main Application)**  

   You can prepare the images in the following ways:  

   - **For Kafka and Postgres**:  
     Option 1: Download from the provided Drive link (recommended for large images).  
       **Google Drive Link**: [Download Images](https://drive.google.com/drive/folders/1nFYBQRvLIpa1sbXEMVABziwy8V2h2Gf5)  
     Option 2: Pull the images directly:  
       ```bash
       docker pull postgres:latest
       docker pull bitnami/kafka:3.5.0
       ```  

   - **For FlowAnalysis**:  
     Option 1: Load the image from the provided tar file:  
       ```bash
       docker load -i flowanalysis_v1.0.tar
       ```  
     Option 2: Build the image from the `Dockerfile` (located in the root directory):  
       ```bash
       docker build -t analysis:v1.0 .
       ```  

3. **Start the project**  
   Once all images are ready on your system, start the project using Docker Compose:  
   ```bash
   docker-compose up -d
   ```  

   The `docker-compose.yml` file is located in the root directory.  

4. **Access the application**  
   The setup is now ready for testing. The application is exposed on port **5010**.  

---

### **Happy Testing!**  
If you encounter any issues, feel free to reach out for support.  
