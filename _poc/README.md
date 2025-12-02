# Proof of Concept (PoC): Automated Security Finding Remediation

This Proof of Concept (PoC) demonstrates a simplified, automated workflow for detecting and remediating a security misconfiguration in a cloud environment.

## Overview

The PoC consists of three core services orchestrated with Docker Compose:

*   **Ingestor**: A Go-based service that simulates the discovery of an S3 bucket and its metadata, storing this information in an ArangoDB database.
*   **Engine**: Another Go-based service that periodically queries the database for security findings. In this PoC, it specifically looks for a publicly accessible S3 bucket with a `sensitivity` tag set to `high`.
*   **Remediation**: A Python Flask application that receives findings from the Engine. It simulates a call to a Large Language Model (LLM) to generate a Terraform HCL block to remediate the finding and saves it to a local file.

## How to Run

1.  **Navigate to the PoC Directory**:

    ```bash
    cd _poc
    ```

2.  **Start the Docker Compose Environment**:

    ```bash
    docker-compose up --build
    ```

This command will build the Docker images for each service and start the containers. The `--build` flag ensures that any changes to the source code are included.

## Workflow

1.  **Infrastructure Setup**: The `docker-compose.yml` file defines the services. It also includes a service that creates a public S3 bucket in a LocalStack container to simulate the insecure resource.
2.  **Ingestion**: The `ingestor` service discovers the S3 bucket and its tags, then writes this information to the ArangoDB database.
3.  **Detection**: The `engine` service queries the database, identifies the public S3 bucket as a finding, and sends the finding to the `remediation` service.
4.  **Remediation Generation**: The `remediation` service receives the finding, simulates an LLM call to generate a Terraform remediation, and saves the code to the `_poc/remediations` directory on your local machine.

## Expected Outcome

After running the PoC, you will find a new Terraform file named `prowler-poc-bucket.tf` inside the `_poc/remediations` directory. This file contains the HCL code to block public access to the S3 bucket.
