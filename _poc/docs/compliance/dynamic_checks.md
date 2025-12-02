
# Technical Specification: Dynamic and Parameterized Checks

This document outlines the technical specification for evolving Prowler's compliance checks from static, single-purpose scripts to dynamic, reusable, and parameterized components.

## 1. Executive Summary

The current compliance framework relies on a large number of checks where the logic is hardcoded. This leads to significant code duplication, particularly for checks that perform similar actions but target different resources or parameters (e.g., checking for open ports). A new, parameterized approach will treat checks as generic functions that accept arguments.

This proposal introduces a system where a single, generic check (e.g., `securitygroup_allow_ingress_to_tcp_port`) can be called with different parameters (e.g., `port: 22` or `port: 3389`), dramatically reducing code duplication, simplifying maintenance, and increasing the framework's flexibility.

## 2. Current State Analysis

- **High Redundancy:** Numerous checks, such as those for open TCP ports, share 95% of their code, with only the port number being the variable.
- **Maintenance Complexity:** A logic update requires changes in every duplicated check, a process that is inefficient and prone to error.
- **Limited Flexibility:** Creating a new check for a minor variation (e.g., a new port) requires writing a completely new check file, increasing the codebase's size and complexity.

## 3. Proposed Architecture

### 3.1. Parameterized Check Definitions

Canonical check definitions will be extended to support an `input_parameters` block, defining the arguments a check can accept, including their type and a default value.

**Example: `prowler/compliance/checks/aws/ec2/securitygroup_allow_ingress_to_tcp_port.yaml`**
```yaml
Id: securitygroup_allow_ingress_to_tcp_port
Provider: AWS
Service: EC2
Description: "Checks if a specific TCP port is open to the internet in any security group."
InputParameters:
  - name: port
    type: int
    description: "The TCP port to check."
  - name: allowed_cidrs
    type: list[string]
    description: "A list of CIDR blocks that are allowed. Defaults to any."
    default: ["0.0.0.0/0"]
Severity: high
Remediation:
  Description: "Restrict ingress to the specified TCP port to only trusted IP addresses."
```

### 3.2. Frameworks with Parameterized Invocations

Compliance frameworks will now be able to call these generic checks and pass parameters to them. This makes the framework definitions more expressive and concise.

**Example: `prowler/compliance/frameworks/aws/cis_v1.4_aws.yaml`**
```yaml
Framework: CIS-AWS-v1.4
...
Requirements:
  - Id: "4.1"
    Name: "Ensure no security groups allow ingress from 0.0.0.0/0 to port 22"
    Checks:
      - check: securitygroup_allow_ingress_to_tcp_port
        parameters:
          port: 22

  - Id: "4.2"
    Name: "Ensure no security groups allow ingress from 0.0.0.0/0 to port 3389"
    Checks:
      - check: securitygroup_allow_ingress_to_tcp_port
        parameters:
          port: 3389
```

### 3.3. Refactoring the Check Execution Engine

The Prowler engine will be updated to:
1.  **Parse Parameters:** Identify when a check in a framework includes a `parameters` block.
2.  **Dynamic Invocation:** Instantiate and execute the corresponding check, passing the specified parameters as arguments to its underlying Python function.
3.  **Default Values:** If a parameter is not provided in the framework, the engine will use the `default` value from the check's canonical definition.

## 4. Benefits

- **Massive Code Reduction:** A single check can now replace dozens of nearly identical files.
- **Simplified Maintenance:** Bug fixes or logic improvements are made once in the generic check and are instantly propagated to all requirements that use it.
- **Increased Agility:** Adding a new requirement for a different port becomes a simple, one-line change in a framework's YAML file, with no new Python code needed.
- **Improved Consistency:** All checks for a given type of vulnerability will now use the exact same logic, ensuring consistent and reliable findings.

## 5. Migration Strategy

1.  **Identify Generic Check Candidates:** Analyze the existing codebase to identify groups of checks that can be consolidated into a single, parameterized check.
2.  **Develop Generic Implementations:** Create the new parameterized check files and refactor the underlying Python code to accept arguments.
3.  **Update Frameworks:** Modify the framework YAML files to replace the old, specific check names with the new parameterized invocations.
4.  **Extend the Engine:** Upgrade the Prowler engine to handle the new `parameters` syntax.
5.  **Deprecate and Remove:** Once the new system is validated, the redundant check files can be safely deleted.
