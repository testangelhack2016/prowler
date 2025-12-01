# Technical Requirements: Graph-Based Attack Path Analysis

**Author:** Angel

This document provides the technical requirements for the "Graph-Based Attack Path Analysis" initiative. This feature will evolve Prowler from a tool that identifies individual misconfigurations into an intelligent platform that understands and visualizes how these issues can be chained together to create real-world attack paths.

---

## Executive Summary

**Purpose:** The goal is to create a system that models a cloud environment as a graph, connecting resources, identities, permissions, and network paths. This will enable Prowler to move beyond static, single-point-in-time checks and perform sophisticated analysis to uncover "toxic combinations"—chains of seemingly low-risk misconfigurations that create critical exploit paths.

**Strategic Value:** This initiative provides a true risk-based view of security posture. Instead of presenting a flat list of findings, Prowler will be able to prioritize the 2% of issues that pose 98% of the risk. This allows security teams to focus their limited resources on the vulnerabilities that matter most, dramatically improving the effectiveness of their security program.

**Audience:** This document is intended for project stakeholders, product managers, data scientists, and the engineering teams who will design and build this graph analysis platform.

## Target Architecture & System Diagram Description

The architecture is built on a high-performance, event-driven stack using Go for data collection, Kubernetes for orchestration, NATS for event streaming, and a dedicated Graph Database for analysis. This ensures the system can process vast amounts of cloud resource data in near real-time.

```mermaid
graph TD
    subgraph Data Collection (Prowler)
        A[Prowler Scanners] -- Resource & Finding Data --> B((Event Bus / NATS))
    end

    subgraph Kubernetes Cluster: Graph Platform
        B --> C[Graph Ingestion Service (Go)]
        C -- Transforms & Writes --> D[(Graph Database: Neo4j/Neptune)]
        E[Attack Path Engine (Go)] -- Queries --> D
        F[API Gateway] -- User Queries --> E
    end

    subgraph User Interface
        G[Prowler Web UI] -- Calls --> F
    end

    style D fill:#f9f,stroke:#333,stroke-width:2px
```

### System Diagram Flow Description

1.  **Data Ingestion:** The process begins with the **Prowler Scanners**, which collect detailed configuration data, relationships, and security findings from the target cloud environments. This data is published as structured events to the **Event Bus (NATS)**.

2.  **Transformation and Loading:** A **Graph Ingestion Service**, a scalable microservice written in Go, consumes the events from NATS. Its responsibility is to translate the raw resource data into a graph model of nodes (e.g., EC2 Instance, IAM User, S3 Bucket) and edges (e.g., `HAS_PERMISSION_TO`, `NETWORK_ACCESS_TO`, `ASSUMES_ROLE`). It then writes this graph structure into the **Graph Database** (e.g., Neo4j or Amazon Neptune).

3.  **Analysis and Querying:** Users interact with the system via the **Prowler Web UI** or a dedicated API. Their requests are sent to the **API Gateway**, which routes them to the **Attack Path Engine**. This Go-based service contains the core logic for translating user questions into complex graph queries (e.g., using Cypher for Neo4j).

4.  **Query Execution:** The Attack Path Engine executes these queries against the Graph Database. Example queries include shortest path algorithms to find exploit chains or pattern matching to identify known vulnerable configurations.

5.  **Visualization:** The results, representing one or more attack paths, are returned to the UI, which visualizes the chain of resources and permissions, highlighting the path an attacker could take.

## Key Use Cases

1.  **Discovering Public Exposure to Sensitive Data:** A security analyst needs to answer the question: "Show me all paths from a public-facing resource to any database tagged with `sensitivity: high`." The Attack Path Engine traverses the graph to find any chain of network access, permissions, and roles that could lead from the internet to a production database, even if it takes multiple hops.

2.  **Blast Radius Analysis for a Compromised Host:** The SOC team receives an alert that a specific EC2 instance may be compromised. They use the Attack Path Engine to ask, "What is the blast radius of this instance?" The system instantly queries the graph to show all other resources (e.g., S3 buckets, code repositories, other VMs) that the instance's IAM role has permissions to access, providing an immediate scope for the incident response.

3.  **Proactive "Choke Point" Remediation:** The system automatically runs a scheduled query to find the single permission or resource that, if remediated, would sever the most critical attack paths. This finding is then integrated with the "AI-Powered Remediation" system, which suggests a targeted fix (e.g., removing a single IAM policy) that provides the biggest security ROI.

## Return on Investment (ROI) Justification

*   **Cost Optimization:**
    *   **Reduced Analyst Fatigue:** Drastically reduces the manual effort and time required for security analysts to trace complex relationships in a cloud environment. What takes a human hours of manual investigation can be done by the graph engine in seconds.
    *   **Focused Remediation:** By prioritizing findings that are part of a valid attack path, engineering teams waste less time fixing low-risk, standalone issues. This optimizes expensive developer time.

*   **Business Value/Risk Improvement:**
    *   **True Risk Prioritization:** This is the core value. It moves the conversation from "We have 500 critical findings" to "We have 15 critical attack paths that expose our customer data." This allows the business to address actual risk, not theoretical severity.
    *   **Preventing Complex Breaches:** The majority of major cloud breaches are not the result of a single failure, but a chain of exploits. This system is specifically designed to find and help prevent exactly those scenarios, significantly reducing the likelihood of a major incident.

## Non-Functional Technical Requirements

1.  **Scalability:** The graph database must support storing and querying a graph with over 100 million nodes (resources/identities) and 500 million edges (relationships) to accommodate enterprise-scale cloud environments.

2.  **Query Latency:** The p95 latency for common attack path queries (e.g., shortest path between two nodes) must be under 15 seconds to enable an interactive user experience.

3.  **Data Freshness:** The data in the graph must not be more than 4 hours older than the state of the live cloud environment. The Graph Ingestion Service must be able to process a full environment refresh within this timeframe.

4.  **Security:** All data, both at rest in the graph database and in transit between services, must be encrypted. Access to the graph database and the Attack Path Engine API must be strictly controlled via IAM and network policies.

5.  **Operational Resilience (SLO):** The Graph Ingestion Service and the Attack Path Engine will have an availability SLO of 99.9%. Data completeness (the percentage of scanned resources successfully ingested into the graph) will have an SLI of >99.5%.

6.  **Extensibility:** The graph model (schema) must be designed to be easily extensible. Adding a new resource type (e.g., a new AWS service) or a new relationship type should not require a full database migration.

7.  **Completeness:** The system must be able to ingest and model all critical resource types relevant to security, including compute instances, identities (users/roles), permissions policies, network configurations (security groups, VPCs), and data stores.
