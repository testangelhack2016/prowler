
# Technical Specification: Graph-Based Compliance Dependency Model

This document outlines the technical specification for re-architecting Prowler's compliance data model into a graph-based dependency model.

## 1. Executive Summary

The current Prowler compliance framework is hierarchical and lacks the ability to represent the complex, many-to-many relationships between resources, controls, and compliance requirements. A graph-based model, where these entities are treated as nodes and their relationships as edges, will provide a far more powerful and flexible foundation for analysis and visualization.

This proposal advocates for adopting a graph database model (either by using a dedicated graph database or by structuring our data in a graph-like manner) to map out the entire compliance ecosystem. This will enable advanced queries, impact analysis, and the visualization of compliance posture in a way that is simply not possible with the current model.

## 2. Current State Analysis

- **Hierarchical and Siloed:** Data is structured in a tree-like manner within each framework, with no easy way to see connections across frameworks or between resources and the various requirements they impact.
- **Limited Querying:** Answering complex questions like "What is the full compliance impact of this single S3 bucket?" or "Which requirements are affected by a failure in this specific check?" is extremely difficult and computationally expensive.
- **Static Views:** The current model only allows for static, pre-defined views of the data.

## 3. Proposed Architecture

We will model our compliance data as a labeled property graph, consisting of nodes and directed, labeled edges.

### 3.1. Node Types

- **`Resource`:** A specific cloud asset (e.g., an S3 bucket, an EC2 instance). Properties would include `id`, `provider`, `region`, etc.
- **`Check`:** A Prowler check. Properties would include `id`, `description`, `severity`, etc.
- **`Requirement`:** A specific compliance requirement (e.g., "CIS AWS v1.4 - 4.1").
- **`Framework`:** A compliance framework (e.g., "CIS AWS v1.4").

### 3.2. Edge Types

- **`HAS_CHECK`:** Connects a `Requirement` to a `Check`.
- **`HAS_REQUIREMENT`:** Connects a `Framework` to a `Requirement`.
- **`AFFECTS`:** Connects a `Check` to a `Resource` (with a `status` property of `PASS` or `FAIL`).

### 3.3. Example Graph

An S3 bucket (`my-bucket`) fails the `s3_bucket_object_versioning` check. This check is part of requirement `5.1.1` in the `AWS-Audit-Manager-Control-Tower-Guardrails` framework. The graph would look like this:

```
(Framework {name: "AWS-Audit-Manager-Control-Tower-Guardrails"})
  -[:HAS_REQUIREMENT]-> (Requirement {id: "5.1.1"})
                          -[:HAS_CHECK]-> (Check {id: "s3_bucket_object_versioning"})
                                           -[:AFFECTS {status: "FAIL"}]-> (Resource {id: "my-bucket"})
```

### 3.4. Technology Choice

- **Option A: Managed Graph Database:** Use a managed graph database like Amazon Neptune or Neo4j Aura. This provides a powerful query language (like openCypher) out of the box but introduces a new piece of infrastructure.
- **Option B: In-Memory Graph Library:** Use a library like `networkx` in Python to build and query the graph in memory during the Prowler execution. This is simpler to implement but may have performance limitations with very large datasets.

## 4. Benefits

- **Powerful, Flexible Queries:** A graph model allows for complex, multi-dimensional queries that are impossible with the current system. We can easily traverse the graph to answer questions like:
    - "Show me all failing resources for the CIS framework."
    - "What is the blast radius of this single misconfigured EC2 instance?"
    - "Which requirements are covered by the fewest checks?"
- **Impact and Root Cause Analysis:** The graph model makes it trivial to perform impact and root cause analysis.
- **Advanced Visualizations:** The graph can be used to generate rich, interactive visualizations of the compliance posture, providing a far more intuitive understanding of the data than a simple table.

## 5. Implementation

1.  **Select Technology:** Choose between a managed graph database or an in-memory library.
2.  **Build the Graph Loader:** Create a process that takes the standardized results JSON and populates the graph.
3.  **Develop a Query Service:** Create a service or library that exposes a simple API for querying the graph.
4.  **Integrate with Reporting:** Update the reporting engine to use the query service to fetch data for its reports.
