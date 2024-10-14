# Data Analytics with Go, ClickHouse, Cube.js, and Terraform

This project demonstrates how to retrieve and visualize data analytics using Golang, ClickHouse as a database, Cube.js for data aggregation, and Terraform for infrastructure management.

## Project Setup

1. **Install Dependencies**  
   Make sure you have Docker, Terraform, and Cube.js installed on your system.

2. **Terraform Setup**  
   Run the following commands to set up ClickHouse using Docker via Terraform:

```bash
terraform init
```

### Step 2: Apply Terraform Configuration

This will create a custom Docker network and start Kafka and Zookeeper containers.

```bash
terraform apply
```

3. **Run Data Migration**  
Use the following command to generate fake data (100 invoices) in ClickHouse:

```bash
go run main.go --migration

```

This will insert 100 fake invoices into the `invoices` table with different amounts and timestamps.

4. **Fetch Analytics Data**  
To fetch analytics data such as total revenue and most profitable months, run:

```bash
go run main.go

```

5. **Cube.js Setup**  
To set up Cube.js for data aggregation and dashboard visualization, navigate to the Cube.js directory and start the development environment:
```bash
npm run dev

```

6. **Custom Port for Cube.js**  
To run Cube.js on a custom port, use:

CUBEJS_PORT=<PORT_NUMBER> npm run dev


7. **Access Cube.js Dashboard**  
Open the Cube.js Playground to access the dashboards and visualization:
http://localhost:4000


## Project Structure

- `main.go` – Contains the logic to either migrate data or fetch analytics from ClickHouse.
- `schema/` – Cube.js schema definitions for data models (invoices, etc.).
- `terraform/` – Terraform configuration files for setting up ClickHouse.

## Commands

- **Terraform**: 
- `terraform init` – Initializes Terraform in the project.
- `terraform apply` – Applies the Terraform configuration to set up the environment.

- **Golang**:
- `go run main.go --migration` – Runs the migration to insert fake invoice data into ClickHouse.
- `go run main.go` – Fetches and displays analytics data.

- **Cube.js**:
- `npm run dev` – Starts the Cube.js development environment.
- `CUBEJS_PORT=<PORT_NUMBER> npm run dev` – Starts Cube.js on a custom port.

## Authors
Amine Ameur



