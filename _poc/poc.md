---
title: Dynamic Risk Graph PoC
---

## **Objective**

This document outlines a phased proof-of-concept (PoC) to demonstrate the value of a dynamic risk graph model in identifying complex, multi-hop attack paths that traditional security scanners often miss.

The PoC will focus on a "toxic combination" use case: **detecting a sensitive S3 bucket that is publicly accessible.**

---

### **Phase 1: Foundation - Data Ingestion & Attack Path Identification**

This phase focuses on validating our ability to ingest cloud data and identify a critical toxic combination using a graph model.

> **For a detailed breakdown of the architecture, use cases, and technical specifications for this phase, please see the [Technical Requirements: Dynamic Risk Graph PoC](docs/dynamic_risk_graph_requirements.md) document.**

- **Status:** ✅ **Complete**
  - All services are containerized and orchestrated with Docker Compose.
  - The `ingestor` service successfully connects to a local AWS environment (LocalStack), finds an S3 bucket, and stores its metadata in the ArangoDB graph database.
  - The `engine` service successfully queries the database and identifies the pre-defined toxic combination.
  - The core data pipeline and attack path detection mechanism are validated.

- **Tasks:**
  - **1.1: Project Scaffolding:** Create `_poc/ingestor`, `_poc/engine`, `_poc/docs` directories. (Done)
  - **1.2: Graph Schema Definition:** Define document/edge collections for `S3Bucket`, `IAMRole`, `HAS_TAG`, `IS_PUBLIC` in `_poc/docs/schema.md`. (Done)
  - **1.3: Graph Ingestion Service:** Develop a Go service in `_poc/ingestor` to scan AWS S3 buckets and populate a local ArangoDB instance. (Done)
  - **1.4: Attack Path Engine v1:** Create a service in `_poc/engine` that executes a hardcoded AQL query to find S3 buckets tagged as `sensitivity: high` that are also public. (Done)

- **Non-Functional Requirements:**
  - **Performance:** Query execution should complete in under 60 seconds for an environment with up to 1,000 buckets. (Met)
  - **Environment:** The entire system must be runnable via Docker Compose on a local machine. (Met)
  - **Security:** The PoC will use read-only AWS credentials. (Met)

### **Phase 2: Intelligence - AI-Powered Remediation**

This phase focuses on using the findings from Phase 1 to automatically generate a fix.

- **Tasks:**
  - **2.1: Remediation Service Scaffolding:** Create a Python-based service in `_poc/remediation` that the Attack Path Engine can call.
  - **2.2: Intelligent Prompt Engineering:** The service will construct a detailed prompt for an LLM, providing the finding's context to generate a Terraform HCL fix.
  - **2.3: Output Generation:** The service will call the LLM API and print the suggested Terraform code to the console, validating the end-to-end workflow.

- **Non-Functional Requirements:**
  - **Integration:** The Remediation Service must be reachable from the Attack Path Engine over a local network.
  - **Clarity:** The generated code must be valid, well-formatted, and human-readable.

### **Phase 3: Enterprise Readiness**

This phase extends the PoC to address the security, compliance, and scalability requirements of large enterprise customers.

- **Tasks:**
  - **3.1: Dynamic Check Definition:** Evolve the hardcoded AQL query into a flexible, definition-based check that can be managed and versioned. This will allow new toxic combinations to be defined without changing the engine's source code.
  - **3.2: Compliance Mapping:** Map the dynamic "Public Sensitive Data" check to relevant controls within enterprise compliance frameworks. For example:
    - **NIST 800-53:** Map to `CP-2(e)` (Contingency Plan Updates) to ensure data protection and availability are not compromised.
    - **NIS2:** Map to `Article 21(2)(a)` (Risk Management Policy) to demonstrate proactive risk identification and treatment.
  - **3.3: Prowler Integration:** Integrate the graph engine and dynamic checks into the main Prowler scanner, allowing them to be executed alongside existing Prowler checks and producing a unified output.
  - **3.4: Post-Incident Review Integration:** Link findings to incident response procedures, such as `NIS2 Article 21(2)(b)` (Incident Handling), by providing context for post-incident reviews and corrective actions.
  - **3.5: Scalability & Performance Testing:** Benchmark the graph engine's performance against a simulated large-scale enterprise environment with millions of resources and hundreds of defined toxic combinations.
