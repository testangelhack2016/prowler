# Technical Requirements: Real-time IaC Scanning in the IDE

**Author:** Angel

This document provides the technical requirements for the "Real-time IaC Scanning in the IDE" feature, a core component of the proposed 10x improvements for Prowler. The goal is to shift security left by integrating Prowler's powerful scanning engine directly into the developer's primary workspace, the IDE.

---

## 1. Introduction

- **Current State:** Prowler scans deployed cloud resources or local Infrastructure as Code (IaC) files from the command line. This is effective but represents a late-stage feedback loop in the development lifecycle.

- **Proposed Future State:** We will create an IDE extension (e.g., for VS Code) that provides developers with real-time security feedback as they write IaC. This extension will analyze files (Terraform, CloudFormation, Kubernetes, etc.) on the fly, highlight misconfigurations, and provide actionable remediation advice directly in the editor, preventing security vulnerabilities before they are ever committed to version control.

## 2. Return on Investment (ROI)

Integrating security scanning into the IDE offers a compelling ROI by fundamentally changing when and how security issues are addressed.

*   **Exponentially Reduced Remediation Cost:** Fixing a security misconfiguration during development is orders of magnitude cheaper than fixing it in production. By catching issues at the source, we eliminate the downstream costs associated with CI/CD failures, security reviews, and emergency production patches.

*   **Increased Developer Velocity:** Developers receive instant, contextual feedback without leaving their editor. This eliminates the context-switching and delays associated with waiting for a pipeline to fail or a security team to review their code, allowing them to code faster and more securely.

*   **Scalable Security Governance:** This empowers every developer to be a first line of security defense. Security teams can scale their impact by embedding policy and best practices directly into the development workflow, rather than acting as a bottleneck. It also serves as a continuous learning tool, educating developers on security best practices.

## 3. Use Cases

1.  **Real-time Terraform Validation:** A DevOps engineer is writing a Terraform file to provision an AWS S3 bucket. As they save the file, the Prowler IDE extension immediately places a yellow squiggly line under the resource block. Hovering over it reveals a warning: "S3 bucket does not have public access blocks enabled." The extension offers a "quick fix" option that, when clicked, automatically inserts the correct `aws_s3_bucket_public_access_block` resource.

2.  **Kubernetes Manifest Hardening:** A developer is creating a Kubernetes `Deployment.yaml`. They configure a container to run with a `securityContext` that includes `privileged: true`. The extension immediately flags this with a critical severity error, explaining the security risks of privileged containers and suggesting alternative capabilities.

3.  **Pre-Commit Confidence:** Before committing code, a developer checks the IDE's "Problems" panel, which lists all Prowler findings for the files in their changeset. They quickly fix a few minor issues and commit their code with high confidence that it meets the organization's security standards.

## 4. System Diagram

This diagram illustrates the proposed architecture, which leverages the Language Server Protocol (LSP) for a clean separation of concerns between the IDE and the Prowler analysis engine.

```mermaid
graph TD
    subgraph IDE (e.g., VS Code)
        A[IaC File Editor (.tf, .yaml)] -- on change --> B(Prowler IDE Extension)
        B -- sends diagnostics --> A
    end

    subgraph Backend Process
        C(Language Server)
        D[Prowler Core Engine (as library)]
        E[Prowler Check Definitions]
    end

    B -- forwards file content --> C
    C -- receives diagnostics --> B
    C -- calls for analysis --> D
    D -- uses --> E
    D -- returns findings --> C

    style D fill:#f9f,stroke:#333,stroke-width:2px
```

**Workflow Explanation:**

1.  The developer opens and edits an IaC file in their IDE.
2.  The **Prowler IDE Extension** detects the file change and sends the full content of the document to the **Language Server**.
3.  The Language Server invokes the **Prowler Core Engine**, which has been packaged as a callable library, passing the IaC file content.
4.  The Prowler Engine runs all relevant **Check Definitions** against the in-memory file content.
5.  Findings (issue description, severity, line number, remediation advice) are returned to the Language Server.
6.  The Language Server translates these findings into LSP-standard diagnostic messages and sends them back to the IDE Extension.
7.  The IDE renders these diagnostics as inline squiggles, entries in the "Problems" panel, and hover/tooltip information.

## 5. Functional Requirements

| Requirement | Input | Output | Inferred/Assumed |
| :--- | :--- | :--- | :--- |
| **Real-time File Scanning** | IaC file content (on open/save/change) | A list of security findings with line numbers | Inferred |
| **Display Findings** | A list of findings from the engine | Visual indicators (squiggles, highlights) in the editor | Assumed |
| **Provide Remediation** | A single finding | A description of the issue, its severity, and a suggested code fix | Inferred |
| **Apply Quick Fix** | User clicks a "quick fix" action | The problematic code is automatically replaced with the suggested fix | Inferred |
| **Support IaC Languages** | Terraform, CloudFormation, Kubernetes YAML files | Correct parsing and analysis for each language | Assumed |

## 6. Non-Functional Requirements

*   **Performance:** The analysis of a single file must complete in under 500ms to provide a real-time, non-blocking user experience.
*   **Security:** All processing MUST happen locally on the user's machine. The IDE extension MUST NOT transmit file contents or credentials to any external service.
*   **Usability:** Findings must be clear, concise, and actionable. Quick fixes should be reliable and accurately reflect the suggested change.
*   **Maintainability:** The core analysis logic must be reused from the main Prowler codebase to avoid duplicating checks. The LSP architecture decouples the IDE-specific code from the analysis engine.
*   **Extensibility:** The system should allow for new checks to be added to the core Prowler library and be automatically picked up by the IDE extension without requiring an extension update.
