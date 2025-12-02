
# Technical Specification: Structured Remediation Metadata

This document provides the technical specification for integrating structured, actionable remediation metadata into Prowler's compliance checks.

## 1. Executive Summary

Current remediation guidance in Prowler is largely descriptive text. While helpful, it is not structured in a way that allows for automation or easy integration with other tools. This proposal outlines a new model where remediation steps are defined as structured data, including code snippets for various IaC tools and CLIs.

By embedding actionable, machine-readable remediation metadata directly into the check definitions, we can lay the foundation for a new class of features, including the automated generation of remediation scripts and direct integration with remediation-as-code platforms.

## 2. Current State Analysis

- **Descriptive Text:** Remediation guidance is provided as free-form text, which is not machine-readable.
- **No Automation:** The lack of structured data makes it impossible to automate remediation actions.
- **Manual Effort:** Users must manually interpret the remediation text and translate it into a series of actions.

## 3. Proposed Architecture

The `Remediation` block within the canonical check definition will be expanded to include structured, code-based guidance.

### 3.1. Enhanced Remediation Block

The new `Remediation` block will have the following structure:

```yaml
Remediation:
  Description: "A human-readable description of the remediation steps."
  Recommendation:
    Text: "Title of the official documentation."
    Url: "URL to the official documentation."
  Code:
    IaC:
      - terraform: "<Terraform code snippet>"
      - cloudformation: "<CloudFormation code snippet>"
      - ansible: "<Ansible code snippet>"
    CLI:
      - aws-cli: "<AWS CLI command>"
      - azure-cli: "<Azure CLI command>"
      - gcloud-cli: "<gcloud CLI command>"
    Scripts:
      - python: "<Python script>"
      - bash: "<Bash script>"
```

### 3.2. Example

**`prowler/compliance/checks/aws/s3/s3_bucket_object_versioning.yaml`**

```yaml
Id: s3_bucket_object_versioning
...
Remediation:
  Description: "Enable versioning on the S3 bucket to protect against accidental deletion."
  Recommendation:
    Text: "Using versioning in S3 buckets"
    Url: "https://docs.aws.amazon.com/AmazonS3/latest/userguide/Versioning.html"
  Code:
    IaC:
      - terraform: |
          resource "aws_s3_bucket_versioning" "example" {
            bucket = aws_s3_bucket.example.id
            versioning_configuration {
              status = "Enabled"
            }
          }
    CLI:
      - aws-cli: "aws s3api put-bucket-versioning --bucket <bucket-name> --versioning-configuration Status=Enabled"
```

## 4. Benefits

- **Actionable Guidance:** Remediation steps are no longer just descriptive; they are actionable code snippets that can be directly used by developers and SREs.
- **Enables Automation:** The structured format allows for the automated generation of remediation scripts, dramatically reducing the time to fix findings.
- **Improved Accuracy:** By providing the exact code needed, we reduce the risk of human error during the remediation process.
- **Integration Hub:** Prowler becomes a central hub for security and compliance, providing not just the "what" (the finding) but also the "how" (the remediation).

## 5. Implementation

1.  **Update Schema:** Update the schema for the canonical check definition to include the new `Remediation` block.
2.  **Populate Data:** Begin populating the new `Remediation` block for all existing checks. This can be a community-driven effort.
3.  **Enhance Reporting:** The reporting engine will be updated to include these structured remediation steps in the output, allowing users to easily copy and paste the code they need.
4.  **Future Development (Remediation Engine):** This new structure paves the way for a future `prowler remediate` command that could automatically apply these fixes.
