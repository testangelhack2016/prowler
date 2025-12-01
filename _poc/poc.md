# Proof of Concept: The 10x Prowler System

## 1. Objective

This document outlines the plan to build a Proof of Concept (PoC) for the next-generation Prowler "10x system." The goal is to validate the foundational architectural pillars outlined in the `_docs/technical/10x system.md` document, specifically focusing on the **Dynamic Risk Graph** and **AI-Powered Remediation**.

All source code, documentation, and infrastructure related to this PoC must be located within the `_poc/` directory to ensure clear separation from the existing Prowler codebase.

---

## 2. User Personas and Use Cases

- **Persona:** **Alex, the Security Engineer**
  - **Needs:** To quickly identify and prioritize the most critical risks in a large, complex cloud environment. Alex is overwhelmed by long lists of individual findings and needs to focus on exploitable attack paths.
  - **Use Case:** Alex runs the PoC against their staging AWS account. The system identifies a publicly accessible S3 bucket containing sensitive data and presents this as a critical attack path. The system then generates a Terraform snippet, which Alex can use to immediately remediate the issue.

---

## 3. Engineering Requirements & Phased Implementation

The PoC will be built in three phases, starting with a foundational PoC and progressing toward an enterprise-grade solution.

### **Phase 1: Foundation - The Dynamic Risk Graph**

This phase focuses on validating our ability to ingest cloud data and identify a critical toxic combination using a graph model.

> **For a detailed breakdown of the architecture, use cases, and technical specifications for this phase, please see the [Technical Requirements: Dynamic Risk Graph PoC](docs/dynamic_risk_graph_requirements.md) document.**

- **Tasks:**
  - **1.1: Project Scaffolding:** Create `_poc/ingestor`, `_poc/engine`, `_poc/docs` directories.
  - **1.2: Graph Schema Definition:** Define document/edge collections for `S3Bucket`, `IAMRole`, `HAS_TAG`, `IS_PUBLIC` in `_poc/docs/schema.md`.
  - **1.3: Graph Ingestion Service:** Develop a Go service in `_poc/ingestor` to scan AWS S3 buckets and populate a local ArangoDB instance.
  - **1.4: Attack Path Engine v1:** Create a service in `_poc/engine` that executes a hardcoded AQL query to find S3 buckets tagged as `sensitivity: high` that are also public.

- **Non-Functional Requirements:**
  - **Performance:** Query execution should complete in under 60 seconds for an environment with up to 1,000 buckets.
  - **Environment:** The entire system must be runnable via Docker Compose on a local machine.
  - **Security:** The PoC will use read-only AWS credentials.

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
  - **3.1: Role-Based Access Control (RBAC):**
    - **Requirement:** Implement a basic RBAC system. Define `Admin` and `Viewer` roles. Only Admins can trigger or approve remediations.
    - **Implementation:** Create an internal user management module and protect API endpoints based on user roles.

  - **3.2: Compliance Framework Integration:**
    - **Requirement:** Map attack path findings to specific compliance controls (e.g., ISO 27001, CIS).
    - **Implementation:** Enrich the graph with compliance data. Modify the Attack Path Engine to report which compliance controls are failing due to a detected attack path (e.g., `ISO-27001-A.5.31`).

  - **3.3: Audit Logging:**
    - **Requirement:** All actions, especially remediation suggestions and approvals, must be logged for security and compliance audits.
    - **Implementation:** Create a dedicated, append-only log stream for all critical events. Each log entry must include a timestamp, the responsible user, and the action performed.

  - **3.4: Scalability & Performance Benchmarking:**
    - **Requirement:** The system must handle an environment with over 1 million resources and 5 million relationships.
    - **Implementation:** Optimize AQL queries, benchmark database performance under load, and implement horizontal scaling for the ingestor and engine services in a Kubernetes environment.

- **Non-Functional Requirements:**
  - **Availability (SLO):** The core services will have an availability SLO of 99.9%.
  - **Data Freshness:** Graph data must not be more than 4 hours older than the state of the live cloud environment.
  - **Legal & Regulatory:** All data handling must comply with relevant regulations identified during the compliance mapping task.
