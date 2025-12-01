# Third-Party Dependencies

This document outlines the third-party modules (libraries) and services that Prowler relies on to function.

## Third-Party Services

Prowler integrates with various third-party services to perform security assessments and for integrations. These include:

- **Cloud Providers:**
  - **AWS (Amazon Web Services):** Core provider for security assessments.
  - **Microsoft Azure:** Core provider for security assessments.
  - **GCP (Google Cloud Platform):** Core provider for security assessments.
  - **Oracle Cloud Infrastructure (OCI):** For security assessments.

- **Containerization & Orchestration:**
  - **Kubernetes:** To perform security checks within Kubernetes clusters.
  - **Docker:** Used in development and for running Prowler in a containerized environment.

- **Version Control & Code Hosting:**
  - **GitHub:** For version control and for performing checks on GitHub organizations and repositories.

- **Communication & Notifications:**
  - **Slack:** For sending scan notifications and results.

- **Security Scanning Services:**
  - **Shodan:** To gather information about internet-facing assets.

- **Other Services:**
  - **Microsoft 365 / Microsoft Graph:** For assessments related to M365 environments.

## Third-Party Modules (Libraries)

Prowler is built using several open-source Python libraries. These are managed using `poetry` and are listed in the `pyproject.toml` file.

### Core Dependencies

These packages are required for the core functionality of Prowler:

- `awsipranges`
- `alive-progress`
- `azure-identity`
- `azure-keyvault-keys`
- `azure-mgmt-applicationinsights`
- `azure-mgmt-authorization`
- `azure-mgmt-compute`
- `azure-mgmt-containerregistry`
- `azure-mgmt-containerservice`
- `azure-mgmt-cosmosdb`
- `azure-mgmt-databricks`
- `azure-mgmt-keyvault`
- `azure-mgmt-monitor`
- `azure-mgmt-network`
- `azure-mgmt-rdbms`
- `azure-mgmt-postgresqlflexibleservers`
- `azure-mgmt-recoveryservices`
- `azure-mgmt-recoveryservicesbackup`
- `azure-mgmt-resource`
- `azure-mgmt-search`
- `azure-mgmt-security`
- `azure-mgmt-sql`
- `azure-mgmt-storage`
- `azure-mgmt-subscription`
- `azure-mgmt-web`
- `azure-mgmt-apimanagement`
- `azure-mgmt-loganalytics`
- `azure-monitor-query`
- `azure-storage-blob`
- `boto3`
- `botocore`
- `colorama`
- `cryptography`
- `dash`
- `dash-bootstrap-components`
- `detect-secrets`
- `dulwich`
- `google-api-python-client`
- `google-auth-httplib2`
- `jsonschema`
- `kubernetes`
- `markdown`
- `microsoft-kiota-abstractions`
- `msgraph-sdk`
- `numpy`
- `pandas`
- `py-ocsf-models`
- `pydantic`
- `pygithub`
- `python-dateutil`
- `pytz`
- `schema`
- `shodan`
- `slack-sdk`
- `tabulate`
- `tzlocal`
- `py-iam-expand`
- `h2`
- `oci`

### Development Dependencies

These packages are used for development, testing, and code quality:

- `bandit`
- `black`
- `coverage`
- `docker`
- `flake8`
- `freezegun`
- `marshmallow`
- `mock`
- `moto`
- `openapi-schema-validator`
- `openapi-spec-validator`
- `pre-commit`
- `pylint`
- `pytest`
- `pytest-cov`
- `pytest-env`
- `pytest-randomly`
- `pytest-xdist`
- `safety`
- `vulture`
