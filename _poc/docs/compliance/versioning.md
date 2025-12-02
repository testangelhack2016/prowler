
# Technical Specification: Versioning and Lifecycle Management for Compliance Frameworks

This document outlines the technical specification for implementing a robust versioning and lifecycle management system for Prowler's compliance frameworks and checks.

## 1. Executive Summary

The current compliance framework lacks a formal system for versioning and managing the lifecycle of its components. This makes it difficult to track changes, manage updates, and ensure that users are running the correct version of a compliance framework for their audits.

This proposal introduces a comprehensive versioning and lifecycle management system based on Semantic Versioning (SemVer). We will implement this system for both individual checks and the compliance frameworks that consume them. Tooling will be developed to automate the versioning process, create changelogs, and manage the state (e.g., `draft`, `active`, `deprecated`) of each component.

## 2. Current State Analysis

- **No Formal Versioning:** The `Version` field is a simple string and is not consistently applied or enforced.
- **Difficult to Track Changes:** There is no easy way to see what has changed between two versions of a compliance framework, a critical requirement for many audit and governance scenarios.
- **Lack of Lifecycle Management:** There is no formal process for introducing, updating, or deprecating frameworks and checks.

## 3. Proposed Architecture

### 3.1. Semantic Versioning (SemVer)

We will adopt Semantic Versioning (v2.0.0) for all compliance components. The version number `MAJOR.MINOR.PATCH` will be incremented as follows:

- **`MAJOR`** version when you make incompatible changes (e.g., removing a requirement).
- **`MINOR`** version when you add functionality in a backward-compatible manner (e.g., adding a new requirement).
- **`PATCH`** version when you make backward-compatible bug fixes (e.g., correcting a typo in a description).

### 3.2. Versioning in Definitions

Both check and framework YAML files will include a `version` field.

**Check Example:**
```yaml
Id: ec2_ebs_volume_encryption
Version: 1.1.0
...
```

**Framework Example:**
```yaml
Framework: CIS-AWS-v1.4
Version: 2.0.1
...
Requirements:
  ...
```

Frameworks will also specify the version of the checks they use, allowing for precise dependency management:
```yaml
    Checks:
      - check: ec2_ebs_volume_encryption@1.1.0
```

### 3.3. Lifecycle Status

Each component will have a `status` field to indicate its current state in the lifecycle:

- **`draft`:** Under development and not ready for production use.
- **`active`:** Stable and recommended for use.
- **`deprecated`:** No longer recommended. A `deprecated_by` field will point to the new version.

### 3.4. Automated Tooling

A new CLI tool will be developed to manage the versioning and lifecycle process:

- **`prowler compliance bump <component> --level <major|minor|patch>`:** This command will automatically increment the version number, update the `version` field in the YAML file, and create a new git tag.
- **`prowler compliance changelog <component> --from <version1> --to <version2>`:** This will generate a human-readable changelog detailing the differences between two versions.
- **`prowler compliance deprecate <component> --new <new-component>`:** This will mark a component as `deprecated` and set the `deprecated_by` field.

## 4. Benefits

- **Clear and Consistent Versioning:** SemVer provides a clear and unambiguous way to communicate the nature of changes to compliance components.
- **Traceability and Auditability:** The ability to generate changelogs and diffs between versions is essential for audit and compliance scenarios.
- **Improved Stability:** The lifecycle management process ensures that users are always aware of the status of the components they are using.
- **Automation:** The tooling will automate many of the manual and error-prone tasks associated with managing compliance-as-code.

## 5. Migration Strategy

1.  **Initial Versioning:** Start by assigning an initial version (e.g., `1.0.0`) to all existing checks and frameworks.
2.  **Develop Tooling:** Build the `prowler compliance` CLI with the `bump`, `changelog`, and `deprecate` commands.
3.  **Integrate with CI/CD:** Integrate the new tooling into the CI/CD pipeline to automate the versioning process and enforce best practices.
4.  **Documentation:** Update the developer and user documentation to explain the new versioning and lifecycle management system.
