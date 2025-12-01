# Prowler AWS Integration Details

This document provides details on the AWS credentials, checks, APIs used, and the data returned by Prowler's AWS integration.

## AWS Credentials

Prowler supports several methods for authenticating with AWS. The credentials must belong to a user or IAM role with the necessary permissions to perform the security checks.

### 1. IAM Role

You can configure an IAM role for Prowler to assume. This is the recommended approach for secure and manageable access.

```yaml
- type: object
  title: AWS Assume Role
  properties:
    role_arn:
      type: string
      description: The Amazon Resource Name (ARN) of the role to assume.
    external_id:
      type: string
      description: An identifier to enhance security for role assumption.
```

### 2. IAM User with Access Keys

You can use an IAM user's access keys (`aws_access_key_id` and `aws_secret_access_key`).

```yaml
- type: object
  title: AWS Credentials
  properties:
    aws_access_key_id:
      type: string
    aws_secret_access_key:
      type: string
```

### 3. Temporary Credentials

Prowler also supports temporary credentials, which include a session token.

```yaml
- type: object
  title: AWS Temporary Credentials
  properties:
    aws_access_key_id:
      type: string
    aws_secret_access_key:
      type: string
    aws_session_token:
      type: string
```

For more detailed information, please refer to the [Authentication](/user-guide/providers/aws/authentication) page.

## Prowler AWS Checks

Prowler comes with a wide range of predefined checks for AWS, covering various services and security best practices. Each check is designed to assess a specific security configuration.

**Example Check:** `s3_bucket_public_access`

This check verifies that S3 buckets do not have public access.

## AWS APIs Used

Prowler interacts with various AWS APIs to perform its security checks. The specific APIs used depend on the checks being executed. Some of the key APIs include:

*   **IAM API:** For checks related to IAM policies, roles, and users.
*   **S3 API:** For checks related to S3 buckets, access controls, and encryption.
*   **EC2 API:** For checks related to EC2 instances, security groups, and VPCs.
*   **RDS API:** For checks related to RDS instances and their configurations.
*   **VPC API:** For checks related to VPCs, subnets, and network ACLs.

## Returned Data

Prowler's checks return findings that provide information about the security posture of your AWS environment. Each finding includes:

*   **Status:** `PASS`, `FAIL`, or `MANUAL`.
*   **Severity:** `Critical`, `High`, `Medium`, `Low`, or `Informational`.
*   **Service Name:** The AWS service the check belongs to.
*   **Check Title:** A descriptive title of the check.
*   **Resource ID:** The identifier of the resource that was checked.
*   **Region:** The AWS region where the resource is located.
*   **Description:** A detailed description of the finding.
*   **Remediation:** Steps to remediate the finding.

Prowler can generate reports in various formats, including:

*   JSON
*   CSV
*   HTML
