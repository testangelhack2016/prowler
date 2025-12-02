# Technical Specification: A Unified and Decoupled Compliance Data Model

This document provides a detailed technical specification for transforming the Prowler compliance framework from its current monolithic structure to a unified, decoupled, and relational data model. This new model will serve as the foundation for future 10x improvements, including dynamic reporting, automated remediation, and graph-based analysis.

## 1. Executive Summary

The current compliance framework, while functional, suffers from significant data duplication and a tight coupling between framework definitions and check logic. Each compliance JSON file is a self-contained monolith, leading to inconsistencies and a high maintenance burden.

This proposal advocates for a "database-as-code" approach. We will decompose the monolithic files into their core entities—`Frameworks`, `Requirements`, and `Checks`—and manage them in a structured, relational manner using YAML files. `Checks` will have a single, canonical definition, and `Frameworks` will reference them by a unique ID, eliminating redundancy and creating a single source of truth.

## 2. Current State Analysis ("As-Is")

The existing model is characterized by:
- **Monolithic JSON Files:** Each framework (e.g., `cis_v1.4_aws.json`) contains full definitions for all its requirements and the checks they use.
- **Data Duplication:** Critical metadata, such as a check's description, service, and remediation steps, is copied into every framework that uses it.
- **High Maintenance Overhead:** Updating a single check's description requires finding and editing every framework JSON file that includes it, a process that is both tedious and error-prone.
- **Inconsistency:** Without a single source of truth, the same check can have different descriptions or attributes across different frameworks.

## 3. Proposed Architecture ("To-Be")

We will separate the data into a hierarchical and relational structure using YAML for improved readability and maintainability.

### 3.1. Directory Structure

The `compliance` directory will be restructured to treat definitions as code, separating canonical check definitions from the frameworks that consume them.

```
prowler/
└── compliance/
    ├── checks/
    │   ├── aws/
    │   │   ├── ec2/
    │   │   │   └── ec2_ebs_volume_encryption.yaml
    │   │   └── s3/
    │   │       └── s3_bucket_object_versioning.yaml
    │   └── azure/
    │       └── ...
    └── frameworks/
        ├── aws/
        │   ├── aws_audit_manager_control_tower_guardrails_aws.yaml
        │   └── cis_v1.4_aws.yaml
        └── azure/
            └── ...
```

### 3.2. Canonical Check Definition

Each check will be defined **once** in its own YAML file within the `compliance/checks/` directory. This file becomes the single source of truth for all metadata related to that check.

**Example: `prowler/compliance/checks/aws/ec2/ec2_ebs_volume_encryption.yaml`**
```yaml
Id: ec2_ebs_volume_encryption
Provider: AWS
Service: EC2
Description: "Checks whether EBS volumes that are in an attached state are encrypted."
Severity: high
Remediation:
  Description: "To enable encryption for an EBS volume, you must create a new, encrypted volume and migrate the data. You cannot encrypt an existing unencrypted volume directly."
  Recommendation:
    Text: "Encrypting EBS Volumes"
    Url: "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html"
  Code:
    IaC:
      - terraform: "resource '''aws_ebs_volume''' '''example''' { encrypted = true }"
      - cloudformation: "EBSVolume: { Type: '''AWS::EC2::Volume''', Properties: { Encrypted: True } }"
    CLI: "aws ec2 create-snapshot --volume-id <volume-id> && aws ec2 create-volume --snapshot-id <snapshot-id> --encrypted"
```

### 3.3. Decoupled Framework Definition

Framework files will be significantly slimmed down. They will define the framework's metadata and its structure of requirements, but instead of containing full check details, they will simply **reference check IDs**.

**Example: `prowler/compliance/frameworks/aws/aws_audit_manager_control_tower_guardrails_aws.yaml`**
```yaml
Framework: AWS-Audit-Manager-Control-Tower-Guardrails
Name: AWS Audit Manager Control Tower Guardrails
Version: "1.0"
Provider: AWS
Description: "AWS Control Tower is a management and governance service that you can use to navigate through the setup process and governance requirements that are involved in creating a multi-account AWS environment."
Requirements:
  - Id: "1.0.3"
    Name: "Enable encryption for EBS volumes attached to EC2 instances"
    Checks:
      - ec2_ebs_default_encryption
      - ec2_ebs_volume_encryption
  - Id: "5.1.1"
    Name: "Disallow S3 buckets that are not versioning enabled"
    Checks:
      - s3_bucket_object_versioning
```
Notice this file is now much smaller, more readable, and contains no redundant information about the checks themselves.

## 4. Benefits of the New Model

This decoupled model represents a 10x improvement over the current state:

- **Single Source of Truth:** All check metadata is defined in one place, ensuring consistency and accuracy across all compliance frameworks.
- **Dramatically Reduced Maintenance:** To update a check's description or add remediation steps, you only need to edit one file.
- **Enhanced Readability:** YAML is inherently more human-friendly than JSON. The new, smaller framework files are significantly easier to read, review, and manage.
- **Scalability:** Adding new frameworks or checks is simplified. There is no longer a need to copy and paste large blocks of JSON.
- **Foundation for Automation:** With structured, canonical data, we can build powerful new features, including:
    - A generic reporting engine that works with any framework.
    - Automated generation of remediation steps.
    - CI/CD validation to ensure the integrity of compliance definitions.

## 5. Migration Strategy

1.  **Develop Tooling:** Create a script to parse all existing framework JSON files and generate the new directory structure and canonical `check` YAML files. The script will identify all unique checks and extract their metadata.
2.  **Convert Frameworks:** The same script will then generate the new, lean `framework` YAML files that reference the check IDs.
3.  **Refactor Prowler's Engine:** Update the Prowler scanning engine to read and assemble this new data model at runtime. The engine will first load all checks from the `compliance/checks` directory and then use the `framework` file to map requirements to those checks.
4.  **Testing & Validation:** Rigorously test the new system to ensure that compliance reports are generated correctly and match the output of the old system.
5.  **Deprecate and Remove:** Once the new model is validated, the old JSON files can be safely removed from the codebase.
