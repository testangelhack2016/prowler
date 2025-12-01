# First Principles in Prowler

This document outlines the fundamental principles that form the foundation of Prowler's architecture and design. These principles ensure the solution is robust, scalable, maintainable, and extensible.

---

## 1. Modularity and Abstraction

**Principle:** Breaking down a complex system into smaller, independent, and interchangeable components (modules). Abstraction is used to hide the complexity of the implementation behind a simplified and consistent interface.

**Application in Prowler:**

- **Layered Architecture:** Prowler is built on a layered architecture that separates the core engine, the cloud provider interactions, and the individual security checks.
- **Provider Pattern:** The `Provider` (e.g., `AWSProvider`, `GCPProvider`, `AzureProvider`) is a key abstraction. It encapsulates all the specific logic required to interact with a particular cloud's SDK (Boto3 for AWS, Azure SDK for Python, etc.).
- **Benefit:** This design allows the `Prowler Core` and the `Checks` to remain cloud-agnostic. A check requests data (e.g., "list all storage buckets") from the provider, and the provider handles the specific API calls. This makes it possible to add new cloud providers without altering the core logic or the existing checks.

---

## 2. Separation of Concerns

**Principle:** Each component in the system should have a single, well-defined responsibility. This prevents components from becoming tightly coupled and overly complex.

**Application in Prowler:**

- **Prowler Core:** Its sole concern is orchestrating the scan. It manages configuration, loads the appropriate provider and checks, and handles reporting. It does not contain any specific cloud API logic.
- **Providers:** Each provider is only responsible for authenticating and communicating with its designated cloud provider's API.
- **Checks:** Each check file contains the logic for a single, atomic security control (e.g., `s3_bucket_public_access`). It is not concerned with authentication or how the data is fetched, only with analyzing the data provided to it.

---

## 3. Statelessness and Parallelism

**Principle:** Components should not retain information (state) from previous interactions. This allows them to be executed independently and concurrently, which is crucial for performance and scalability.

**Application in Prowler:**

- **Stateless Checks:** Prowler's checks are designed to be stateless. Each check is a self-contained unit that receives all the necessary information as input, performs its analysis, and returns a finding. It does not depend on the results or execution of other checks.
- **Parallel Execution:** Because checks are stateless and independent, Prowler can execute them in parallel. This dramatically reduces the total scan time, especially in large cloud environments with thousands of resources.

---

## 4. Configuration as Code

**Principle:** Managing and provisioning infrastructure and system configurations through machine-readable definition files, rather than manual configurations.

**Application in Prowler:**

- **YAML Configuration:** Prowler uses YAML files for managing its configuration, including:
    - **Checks to execute:** Users can define which checks or compliance frameworks to run.
    - **Mute-listing (Allowlisting):** Findings can be suppressed for specific resources via configuration files, allowing for auditable and version-controlled exceptions.
- **Benefit:** This approach allows for version control of Prowler's configuration, repeatable scans, and easier integration into CI/CD pipelines for continuous security monitoring.

---

## 5. Extensibility

**Principle:** The system should be designed to accommodate new functionality with minimal changes to the existing codebase.

**Application in Prowler:**

- **Custom Checks:** The modular design makes it straightforward for users to write their own custom checks. A new check can be added by creating a new Python file that adheres to the defined check structure, without needing to modify Prowler's core.
- **New Providers:** The provider abstraction layer makes it feasible to add support for entirely new cloud providers or services by creating a new provider class that implements the standard interface for interacting with the Prowler engine.

---

## 6. Automation of Security Auditing

**Principle:** Automating repetitive and manual tasks to improve efficiency, reduce human error, and ensure consistency.

**Application in Prowler:**

- **Codified Best Practices:** Prowler's primary function is the automation of security audits. It codifies well-established security benchmarks (like CIS, GDPR, ISO-27001, etc.) into automated, executable checks.
- **Continuous Monitoring:** By running Prowler in an automated fashion (e.g., via cron jobs or CI/CD pipelines), organizations can achieve continuous monitoring of their cloud security posture, moving from periodic manual audits to a state of constant readiness.
