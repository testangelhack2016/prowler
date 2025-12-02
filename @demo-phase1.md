# Phase 2 Demo: Dynamic Check Definitions

This document provides the script and verification steps to demonstrate the successful completion of Phase 2 of the Dynamic Risk Graph PoC.

**Objective:** To demonstrate a more flexible and scalable system where "toxic combination" checks are defined as data within the graph database, rather than being hardcoded in the `engine` service.

---

## 1. Setup and Configuration

The entire PoC is containerized and managed by Docker Compose, requiring minimal setup.

**Prerequisites:**
- Docker and Docker Compose are installed.
- The project codebase is checked out.

**Execution:**

To start the system, run the following command from the root of the project directory:

```bash
docker-compose -f _poc/docker-compose.yml up --build -d
```
*(Note: The `-d` flag runs the containers in the background. You can omit it to see all logs in your terminal.)*

This single command will:
1.  **Build** the container images for the `ingestor` and `engine` services.
2.  **Start** all the necessary services in the correct order: `arangodb`, `localstack`, `create-bucket`, `ingestor`, and `engine`.
3.  **Create Bucket:** The `create-bucket` service will automatically create a test S3 bucket named `prowler-poc-bucket` in the LocalStack container, make it public, and tag it as `sensitivity:high`.
4.  **Ingest Data:** The `ingestor` service will connect to LocalStack, find the newly created bucket, and populate the ArangoDB database with its metadata, tags, and relationships.
5.  **Ingest Check Definition:** The `ingestor` will also create and insert a "Dynamic Check Definition" into the `Checks` collection.
6.  **Detect Risk:** The `engine` service will query the `Checks` collection, and then execute the AQL queries from each check to find any toxic combinations.

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

The output should show that `arangodb`, `localstack`, and `engine` are in the `running` state. The `create-bucket` and `ingestor` services will have the status `exited (0)` as they are designed to run once and exit.

### Step 2: Confirm Bucket and Check Ingestion

Next, inspect the logs to prove that the setup and data ingestion steps ran correctly.

**Action:**

```bash
# Check the logs of the ingestor service
docker-compose -f _poc/docker-compose.yml logs ingestor
```

**Expected Outcome:**
- The `ingestor` logs should show the message: `Successfully ingested default check definition.`
- It should also show: `Found toxic combination: Public and sensitive bucket 'prowler-poc-bucket'. Ingesting into ArangoDB.`

### Step 3: Show the Finding

Finally, check the logs of the `engine` service. This is where the success of the dynamic check is demonstrated.

**Action:**

```bash
docker-compose -f _poc/docker-compose.yml logs engine
```

**Expected Outcome:**

The `engine` logs will show it loading and executing the dynamic check. You will see output similar to this, confirming that the `engine` found the risk by using the check definition from the database:

```
Executing check: Public Sensitive S3 Bucket
2025/12/02 17:21:25 [FINDING] Check 'Public Sensitive S3 Bucket': map[message:S3 bucket 'prowler-poc-bucket' is public and tagged as sensitive. resource_id:prowler-poc-bucket]
```
*(Note: You might see an initial, harmless error message `collection or view not found: Checks` if the engine starts before the ingestor has created the collection. The engine will retry and succeed on the next cycle.)*

---

## 3. Advanced Verification: Inspecting the Graph Manually

For a deeper dive, you can directly access the ArangoDB database to see the raw graph data, including the dynamic check itself.

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

Now that you are connected, you can run AQL (ArangoDB Query Language) queries to inspect the ingested data. Because `arangosh` is a JavaScript shell, you must wrap the AQL query in `db._query()`.

**Action 1: View the Ingested Bucket**

```javascript
db._query('FOR bucket IN S3Bucket RETURN bucket');
```

**Expected Outcome 1:**

The query will return a JSON object for the ingested bucket.

```json
[
  {
    "_key": "prowler-poc-bucket",
    "_id": "S3Bucket/prowler-poc-bucket",
    "_rev": "...",
    "is_public": true,
    "name": "prowler-poc-bucket"
  }
]
```

**Action 2: View the Dynamic Check Definition**

```javascript
db._query('FOR check IN Checks RETURN check');
```

**Expected Outcome 2:**

This query returns the check document that the `engine` uses to find the toxic combination. Notice the `aql` field, which contains the exact query the engine runs.

```json
[
  {
    "_key": "public-sensitive-s3-bucket",
    "_id": "Checks/public-sensitive-s3-bucket",
    "_rev": "...",
    "name": "Public Sensitive S3 Bucket",
    "description": "Finds S3 buckets that are public and have a 'sensitivity:high' tag.",
    "aql": "
FOR bucket IN S3Bucket
    FILTER bucket.is_public == true
    FOR tag, edge IN 1..1 OUTBOUND bucket has_tag
        FILTER tag.name == "sensitivity:high"
        RETURN {
            "resource_id": bucket.name,
            "message": CONCAT("S3 bucket '", bucket.name, "' is public and tagged as sensitive.")
        }
"
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

## 5. Phase 2 Summary and Next Steps

Phase 2 successfully addressed the major concern from Phase 1 by removing the hardcoded logic from the `engine`.

- **Completed: Dynamic Check Definition**
  - The logic to find the toxic combination is no longer hardcoded in the `engine`. It is now stored as a document in the `Checks` collection in ArangoDB.
  - The `engine` now dynamically loads and executes these checks, making the system extensible without requiring code changes for new checks.

The following concerns from Phase 1 remain and will be the focus of the next phase:

- **Concern: Manual Verification**
  - The results of the scan are still only visible by manually inspecting the container logs.
  - **Next Step (Phase 3):** We will integrate the engine's output with Prowler's standard reporting formats and eventually a dedicated UI.

- **Concern: Limited Scope**
  - The current PoC only covers a single, simple use case.
  - **Next Step (Phase 3):** We will expand the graph schema and our check definitions to model more complex, multi-hop attack paths involving multiple services and resources.
