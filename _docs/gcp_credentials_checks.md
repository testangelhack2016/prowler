# Prowler GCP Integration Details

This document provides details on the GCP credentials, checks, APIs used, and the data returned by Prowler's GCP integration.

## GCP Credentials

Prowler supports several methods for authenticating with Google Cloud Platform. The credentials must belong to a user or service account with the necessary permissions to perform the security checks.

### 1. Application Default Credentials (ADC)

You can use the `gcloud` CLI to set up your Application Default Credentials:

```console
gcloud auth application-default login
```

### 2. Service Account JSON Key File

You can use a service account key file by setting the `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

```console
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials.json"
```

To create a service account and generate a key, you can follow the official Google Cloud documentation:
- [Create a service account](https://cloud.google.com/iam/docs/creating-managing-service-accounts)
- [Generate a key for a service account](https://cloud.google.com/iam/docs/creating-managing-service-account-keys)

### 3. Service Account Impersonation

Prowler can impersonate a service account using the `--impersonate-service-account` flag:

```console
prowler gcp --impersonate-service-account <service-account-email>
```

For more detailed information, please refer to the [Authentication](/user-guide/providers/gcp/authentication) page.

## Prowler GCP Checks

Prowler comes with a wide range of predefined checks for GCP, covering various services and security best practices. Each check is designed to assess a specific security configuration.

**Example Check:** `iam_sa_no_user_managed_keys`

This check verifies that service accounts are not using user-managed keys.

## Google Cloud APIs Used

Prowler interacts with various Google Cloud APIs to perform its security checks. The specific APIs used depend on the checks being executed. Some of the key APIs include:

*   **Identity and Access Management (IAM) API:** For checks related to IAM policies, roles, and service accounts.
*   **Compute Engine API:** For checks related to virtual machines, firewalls, and networks.
*   **Cloud Storage API:** For checks related to storage buckets, access controls, and encryption.
*   **BigQuery API:** For checks related to BigQuery datasets and tables.
*   **Cloud SQL API:** For checks related to Cloud SQL instances and their configurations.
*   **Google Kubernetes Engine (GKE) API:** For checks related to GKE clusters and their security settings.

## Returned Data

Prowler's checks return findings that provide information about the security posture of your GCP environment. Each finding includes:

*   **Status:** `PASS`, `FAIL`, or `MANUAL`.
*   **Severity:** `Critical`, `High`, `Medium`, `Low`, or `Informational`.
*   **Service Name:** The GCP service the check belongs to.
*   **Check Title:** A descriptive title of the check.
*   **Resource ID:** The identifier of the resource that was checked.
*   **Region:** The GCP region where the resource is located.
*   **Description:** A detailed description of the finding.
*   **Remediation:** Steps to remediate the finding.

Prowler can generate reports in various formats, including:

*   JSON
*   CSV
*   HTML
