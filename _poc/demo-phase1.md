# Phase 1 Demo: Dynamic Risk Graph

This document provides the script and verification steps to demonstrate the successful completion of Phase 1 of the Dynamic Risk Graph PoC.

**Objective:** To show a working, end-to-end data pipeline that can ingest cloud asset data, store it in a graph database, and run a query to identify a pre-defined "toxic combination" (a sensitive, publicly accessible S3 bucket).

---

## 1. Setup and Configuration

The entire PoC is containerized and managed by Docker Compose, requiring minimal setup.

**Prerequisites:**
- Docker and Docker Compose are installed.
- The project codebase is checked out.

**Execution:**

To start the system, run the following command from the root of the project directory:

```bash
docker-compose -f _poc/docker-compose.yml up --build
```

This single command will:
1.  **Build** the container images for the `ingestor` and `engine` services.
2.  **Start** all the necessary services in the correct order: `arangodb`, `localstack`, `create-bucket`, `ingestor`, and `engine`.
3.  **Create Bucket:** The `create-bucket` service will automatically create a test S3 bucket named `prowler-poc-bucket` in the LocalStack container.
4.  **Ingest Data:** The `ingestor` service will connect to LocalStack, find the newly created bucket, and populate the ArangoDB database with its metadata.
5.  **Detect Risk:** The `engine` service will repeatedly query the database to find the toxic combination we defined.

---

## 2. Demo Script

Follow these steps to run the demo and verify the outcome.

### Step 1: Verify Container Status

First, confirm that all the necessary services are up and running as expected.

**Action:**

```bash
docker-compose -f _poc/docker-compose.yml ps
```

**Expected Outcome:**

The output should show that `arangodb`, `localstack`, and `engine` are running. The `create-bucket` and `ingestor` services will have already run to completion and will not be in the running list.

### Step 2: Confirm Bucket Creation and Ingestion

Next, inspect the logs to prove that the setup and data ingestion steps ran correctly.

**Action:**

```bash
# Check the logs of the bucket creation service
docker-compose -f _poc/docker-compose.yml logs create-bucket

# Check the logs of the ingestor service
docker-compose -f _poc/docker-compose.yml logs ingestor
```

**Expected Outcome:**
- The `create-bucket` logs should show the message: `make_bucket: prowler-poc-bucket`.
- The `ingestor` logs should show a clean run without any connection errors.

### Step 3: Show the Finding

Finally, check the logs of the `engine` service. This is where the PoC's success is demonstrated.

**Action:**

```bash
docker-compose -f _poc/docker-compose.yml logs engine
```

**Expected Outcome:**

The `engine` logs will repeatedly print the results of its query. You will see output showing that it has found the `prowler-poc-bucket`, confirming that the entire pipeline is working.

---

## 3. Advanced Verification: Inspecting the Graph Manually

For a deeper dive, you can directly access the ArangoDB database to see the raw graph data.

### Step 1: Access the ArangoDB Container

First, open a shell session inside the running `arangodb` container.

**Action:**
```bash
docker-compose -f _poc/docker-compose.yml exec arangodb /bin/sh
```

### Step 2: Connect to the Database

Once inside the container, use the ArangoDB shell (`arangosh`) to connect to the database. The password is `prowler`.

**Action:**
```bash
arangosh --server.password prowler
```

### Step 3: Query the Data

Now that you are connected, you can run an AQL (ArangoDB Query Language) query to inspect the ingested data. Because `arangosh` is a JavaScript shell, you must wrap the AQL query in `db._query()`.

**Action:**

```javascript
db._query('FOR bucket IN S3Bucket RETURN bucket');
```

**Expected Outcome:**

The query will return a JSON object representing the ingested bucket, confirming that the `ingestor` successfully populated the database.

```json
[
  {
    "_key": "prowler-poc-bucket",
    "_id": "S3Bucket/prowler-poc-bucket",
    "is_public": true,
    "name": "prowler-poc-bucket",
    "tags": []
  }
]
```

---

## 4. Manual Triggering (for Debugging)

While `docker-compose up` automates the entire process, you can manually run the `create-bucket` and `ingestor` steps using `docker-compose run`. This is useful for debugging or re-running the pipeline without restarting the database and LocalStack.

### Step 1: Start Core Services

If they are not already running, start the `arangodb` and `localstack` services in the background.

**Action:**
```bash
docker-compose -f _poc/docker-compose.yml up -d arangodb localstack
```

### Step 2: Manually Create the Bucket

Run the `create-bucket` service as a one-off task.

**Action:**
```bash
docker-compose -f _poc/docker-compose.yml run create-bucket
```

### Step 3: Manually Trigger the Ingestion Scan

Run the `ingestor` service as a one-off task.

**Action:**
```bash
docker-compose -f _poc/docker-compose.yml run ingestor
```

After these commands complete, you can check the `engine` logs or query ArangoDB to verify the results.

---

## 5. Follow-up Tasks and Concerns

This PoC successfully validates our core hypothesis, but it also highlights several areas for improvement that will be addressed in the next phases.

- **Concern: Hardcoded Logic**
  - The query to find the toxic combination is currently hardcoded directly into the `engine`'s source code. This is inflexible and not scalable.
  - **Next Step (Phase 3):** We will develop a **Dynamic Check Definition**, which will allow us to define new toxic combinations and attack paths as data, without requiring code changes.

- **Concern: Manual Verification**
  - The results of the scan are only visible by manually inspecting the container logs. This is not a user-friendly or practical way to consume findings.
  - **Next Step (Phase 3):** We will integrate the engine's output with Prowler's standard reporting formats and eventually a dedicated UI.

- **Concern: Limited Scope**
  - The current PoC only covers a single, simple use case (one S3 bucket). A real-world environment is far more complex.
  - **Next Step (Phase 3):** We will expand the graph schema and our check definitions to model more complex, multi-hop attack paths involving multiple services and resources.
