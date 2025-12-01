# Prowler Configuration: A Technical Deep-Dive

**Version:** 1.15.0

This document provides a comprehensive technical overview of the configuration system in Prowler. It is intended for advanced users and developers who need to understand how to customize Prowler's behavior and integrate it into their own systems.

---

## 1. Introduction to Prowler's Configuration

Prowler's configuration is managed through a set of YAML files. This design choice allows for a high degree of flexibility and enables users to define complex configurations in a human-readable format. The core of the configuration system lies in the `config.yaml` file, which acts as the main entry point for all settings.

However, the system is not monolithic. It is designed to be modular, allowing users to split their configurations across multiple files and directories. This is particularly useful for large-scale deployments where managing a single, massive configuration file would be impractical.

---

## 2. The Core `config.yaml` File

The `config.yaml` file is the heart of Prowler's configuration. It defines the providers, accounts, and other global settings that Prowler will use during its execution. Here is a breakdown of the key sections within this file.

### Example `config.yaml`

```yaml
providers:
  aws:
    accounts:
      - id: "123456789012"
        role: "arn:aws:iam::123456789012:role/prowler"
        regions:
          - "us-east-1"
          - "eu-west-1"

checks:
  - s3_bucket_public_access
  - ec2_instance_no_public_ip

output:
  format: "json"
  path: "/prowler/output"
```

### Key-Value Explanations

| Key                              | Type     | Description                                                                                        |
| :------------------------------- | :------- | :------------------------------------------------------------------------------------------------- |
| `providers`                      | `map`    | Defines the cloud providers and accounts to be scanned.                                            |
| `providers.aws`                  | `map`    | Specifies the configuration for the AWS provider.                                                  |
| `providers.aws.accounts`         | `list`   | A list of AWS accounts to be scanned. Each account is defined by its ID and the role to assume.    |
| `providers.aws.accounts.id`      | `string` | The 12-digit AWS account ID.                                                                       |
| `providers.aws.accounts.role`    | `string` | The ARN of the IAM role that Prowler should assume to perform the security assessment.             |
| `providers.aws.accounts.regions` | `list`   | A list of regions to be scanned within the account. If not specified, all regions will be scanned. |
| `checks`                         | `list`   | A list of specific checks to be executed. If not specified, all checks will be run.                |
| `output`                         | `map`    | Defines the output format and location for the scan results.                                       |
| `output.format`                  | `string` | The format of the output file. Supported formats include `json`, `csv`, `html`, and `junit`.       |
| `output.path`                    | `string` | The directory where the output files will be saved.                                                |

---

## 3. Design Analysis: The Good, The Bad, and The Ugly

This section provides a critical analysis of Prowler's configuration system, highlighting its strengths and weaknesses.

### The Good (Pros)

*   **Flexibility and Expressiveness:** The YAML-based configuration allows for a high degree of flexibility. Users can define complex provider and check configurations with ease.
*   **Human-Readable:** YAML is a human-readable format, which makes it easy for users to understand and edit the configuration files.
*   **Modularity:** The ability to split configurations across multiple files is a major advantage for large-scale deployments. It allows for better organization and management of the configuration.
*   **Extensibility:** The design makes it relatively easy to add new providers and checks without having to modify the core Prowler code.

### The Bad (Cons)

*   **Complexity for Beginners:** The sheer number of configuration options can be overwhelming for new users. The learning curve can be steep, especially for those who are not familiar with YAML or cloud security concepts.
*   **Lack of Validation:** The current system lacks robust validation for the configuration files. This can lead to runtime errors if the configuration is not well-formed.
*   **Potential for Misconfiguration:** The flexibility of the system also means that there is a higher potential for misconfiguration. A small typo in the configuration file can lead to unexpected behavior.

### The Ugly (The Inelegant Parts)

*   **Role and Credential Management:** While the configuration allows for specifying roles, the management of credentials is still a bit clunky. The system could benefit from a more streamlined way of handling authentication, perhaps by integrating with a secret management system like HashiCorp Vault.
*   **Verbosity:** The YAML configuration can become verbose, especially for large deployments with many accounts and regions. This can make it difficult to get a high-level overview of the configuration.

---

## 4. Conclusion

Prowler's configuration system is a powerful and flexible tool that allows for a high degree of customization. However, it is not without its flaws. The complexity of the system can be a barrier for new users, and the lack of validation can lead to runtime errors.

Despite these weaknesses, the modular and extensible design of the configuration system is a major strength. It allows Prowler to be adapted to a wide range of use cases and environments. With some improvements in the areas of validation and credential management, Prowler's configuration system could be even more robust and user-friendly.
