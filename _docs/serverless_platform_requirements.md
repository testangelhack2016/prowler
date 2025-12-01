# Technical Requirements: Prowler as a Serverless-Native Platform

**Author:** Angel

This document outlines the technical requirements for re-architecting Prowler into a fully serverless-native platform, as proposed in the "10x Improvements" document. This initiative will transform Prowler from a command-line tool into a highly scalable, resilient, and cost-efficient cloud security platform, paving the way for a viable SaaS offering.

---

## 1. Introduction

- **Current State:** Prowler is a powerful CLI tool. The associated Prowler App consists of a Django backend, a Next.js frontend, and a PostgreSQL database, typically containerized with Docker. This architecture requires users to manage and scale server infrastructure, which can be complex and costly.

- **Proposed Future State:** We will re-architect the Prowler backend into a fully serverless-native platform, primarily using event-driven compute services like AWS Lambda and workflow orchestration like AWS Step Functions. This will eliminate the need for server management, enable massive parallelization, and introduce significant operational efficiencies.

## 2. Return on Investment (ROI)

Adopting a serverless architecture provides a significant ROI through cost savings, performance gains, and new business opportunities.

*   **Drastic Reduction in Operational Cost:**
    *   **Pay-per-Use:** Eliminates costs associated with idle servers. Instead of paying for a 24/7 running API and database server, we only pay for the compute time consumed during a scan.
    *   **Zero Infrastructure Management:** DevOps and engineering time is freed from patching, scaling, and managing servers, allowing the team to focus on core product features and security check development.

*   **10x Performance and Scalability:**
    *   **Massive Parallelism:** A serverless architecture allows us to break a large scan into thousands of parallel micro-tasks (e.g., one Lambda per check per resource). This can reduce scan times for large enterprise environments from hours to minutes.
    *   **Infinite Scalability:** The platform will automatically scale to handle any number of concurrent scans without manual intervention or performance degradation, a key requirement for a multi-tenant SaaS application.

*   **Business Enablement:**
    *   **Foundation for SaaS:** This architecture is the essential technical foundation for a scalable, resilient, and profitable Prowler SaaS offering.
    *   **Enhanced User Experience:** Faster scan times and a more resilient platform directly lead to higher customer satisfaction.

## 3. Use Cases

1.  **On-Demand Enterprise Scans:** A security engineer from a large enterprise logs into the Prowler web app and triggers a full scan on their 50+ AWS accounts. The platform executes the scan across all accounts in parallel, delivering a comprehensive report in under 30 minutes.

2.  **Scheduled Compliance Audits:** A compliance officer configures a recurring weekly scan for PCI-DSS across all production environments. The results are automatically generated and archived for audit purposes without any manual intervention.

3.  **Event-Driven Security Response:** An AWS CloudTrail event shows that an S3 bucket's public access policy has been changed. An EventBridge rule triggers a targeted Prowler scan via a Lambda function on only that specific S3 bucket. A finding is generated and a notification is sent to the security team within seconds of the event.

## 4. System Diagram

This diagram illustrates the proposed serverless architecture using AWS services.

```mermaid
graph TD
    subgraph User Interaction
        A[Web UI / CLI] --> B{API Gateway}
    end

    subgraph Prowler Serverless Platform
        B --> C[Scan Orchestrator Lambda]
        C --> D{Step Functions Workflow}
        D --> E[Dispatcher Lambda]
        E --> F((Executor Lambdas))
        F --> G[Results Database (DynamoDB)]
    end

    subgraph Event-Driven Scans
        H[Cloud Provider Events (e.g., CloudTrail)] --> I{EventBridge}
        I --> C
    end

    F -- Check Results --> G
    B -- Query Results --> G

    style F fill:#f9f,stroke:#333,stroke-width:2px
```

**Workflow Explanation:**

1.  A scan is initiated either via the **API Gateway** (user-triggered) or **EventBridge** (event-driven).
2.  The **Scan Orchestrator Lambda** validates the request and starts the **Step Functions Workflow**.
3.  The **Dispatcher Lambda** fetches the list of resources to be scanned and breaks the job into thousands of small tasks (e.g., `check_s3_public_access` on `my-bucket`).
4.  Each task is executed by a separate, concurrent **Executor Lambda**. These are lightweight functions containing the logic for a single Prowler check.
5.  Results are written directly from the Executor Lambdas to the **Results Database** (e.g., DynamoDB) for high-throughput, scalable persistence.
6.  The API Gateway can query this database to serve results to the user.

## 5. Functional Requirements

*   The platform MUST support triggering scans via a secure RESTful API.
*   The platform MUST be able to orchestrate scans for AWS, GCP, and Azure.
*   The system MUST support on-demand, scheduled, and event-driven scan triggers.
*   Scan results MUST be stored in a persistent, queryable database.
*   The platform MUST provide a mechanism for securely storing and managing customer cloud credentials (e.g., using AWS Secrets Manager).
*   The system MUST be designed for multi-tenancy, with strict logical and/or physical data separation between tenants.

## 6. Non-Functional Requirements

*   **Scalability:** The platform must scale horizontally to support thousands of concurrent checks without a degradation in performance.
*   **Performance:** Full scans of large environments (>100 accounts, >10,000 resources) should complete in under 30 minutes.
*   **Cost-Effectiveness:** The architecture MUST leverage pay-per-use services to minimize cost during idle periods.
*   **Security:** All data in transit and at rest must be encrypted. The platform itself must adhere to security best practices to prevent unauthorized access.
*   **Maintainability:** The codebase must be modular and well-documented. Components should be decoupled to allow for independent updates.
*   **Usability:** The API must be well-documented (e.g., OpenAPI spec) and easy to use.

## 7. Out of Scope

*   This project focuses solely on the re-architecture of the backend scanning engine. It does not include the development of a new frontend UI, though it will provide the necessary APIs to support one.
*   Migration of existing data from the current Prowler App database is not in the initial scope.
