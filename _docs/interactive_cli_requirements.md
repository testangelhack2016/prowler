# Technical Requirements: Interactive CLI with Conversational AI

**Author:** Angel

This document outlines the technical requirements for creating an "Interactive CLI with Conversational AI," a key 10x improvement initiative for Prowler. This feature will transform the Prowler command-line interface from a flag-based tool into an intelligent, interactive shell that understands natural language.

---

## 1. Introduction

- **Current State:** The Prowler CLI is a powerful and scriptable tool, but it requires users to learn and remember a wide range of commands, providers, checks, and output formatting flags. This can present a steep learning curve for new users and can be cumbersome for even experienced users during exploratory analysis.

- **Proposed Future State:** We will develop a "Prowler Shell," an interactive mode for the CLI that is powered by a conversational AI. This will allow users to make requests in plain English. The AI will translate these natural language queries into the appropriate Prowler commands, execute them, and present the results in a user-friendly format, effectively creating a natural language interface for cloud security analysis.

## 2. Return on Investment (ROI)

*   **Democratization of Security:** Makes the full power of Prowler accessible to a much wider audience, including junior engineers, support staff, and managers who are not CLI power-users. This embeds security expertise directly into the tool.

*   **Drastic Increase in Efficiency:** Allows for rapid, iterative, and exploratory security analysis. Instead of manually constructing complex commands, users can simply ask questions, significantly speeding up investigations and daily tasks.

*   **Reduced Cognitive Load:** Users no longer need to memorize dozens of flags and command structures. They can focus on the security questions they want to answer, not on the syntax required to ask them.

## 3. Use Cases

1.  **Targeted Resource Query:** A junior cloud engineer needs to verify the security of a specific S3 bucket. They open the Prowler Shell and type: `"check the bucket named 'my-company-assets-123'"`. The AI translates this, runs `prowler aws --services s3 --resource-id my-company-assets-123`, and displays a summary of the findings.

2.  **Risk-Based Exploration:** A security manager wants a high-level overview of their environment's posture. They ask: `"what are the most critical risks in my 'prod' AWS account right now?"` The AI runs a scan, filters for `Critical` severity findings, and presents a summarized list, allowing the manager to quickly grasp the most pressing issues.

3.  **Compliance-Focused Investigation:** An auditor is assessing GDPR compliance. They query: `"show me all non-compliant GDPR controls related to data storage in Azure"`. The AI constructs the appropriate command, e.g., `prowler azure --compliance gdpr_europe --services storage blob cosmosdb sqldatabase`, executes it, and presents the failing controls.

## 4. System Diagram

This diagram illustrates the architecture for processing natural language commands. It involves a local CLI that communicates with a backend LLM for understanding and translation.

```mermaid
graph TD
    subgraph User's Terminal
        A[User types: "show me public s3 buckets"] --> B(Prowler Interactive Shell)
    end

    subgraph AI Processing Layer
        C{LLM API Endpoint (e.g., OpenAI, Vertex AI)}
        D[Prowler Command Library & Context]
    end

    subgraph Local Execution
        E[Prowler Core Engine]
        F[Command Execution Module]
    end

    B -- Sends secure prompt --> C
    C -- Ingests context --> D
    C -- Returns translated command --> B
    B -- Executes command via --> F
    F -- Invokes --> E
    E -- Returns results --> B
    B -- Displays formatted output --> A

    style C fill:#f9f,stroke:#333,stroke-width:2px
```

**Workflow Explanation:**

1.  The user enters a natural language query into the **Prowler Interactive Shell**.
2.  The shell sanitizes the input and embeds it into a carefully crafted prompt. This prompt includes context about Prowler's command structure (providers, services, flags) from the **Prowler Command Library & Context**.
3.  The prompt is sent to a secure **LLM API Endpoint**.
4.  The LLM processes the prompt and returns a structured JSON object containing the translated Prowler CLI command (e.g., `{"command": "prowler", "args": ["aws", "--services", "s3"], "filters": ["Status==FAIL", "CheckID==s3_bucket_public_access"]}`).
5.  The Prowler Shell validates the received command for safety and syntax.
6.  The validated command is passed to the local **Command Execution Module**, which runs the **Prowler Core Engine**.
7.  The results are returned to the shell, which may format them in a human-readable summary or table for display to the user.

## 5. Functional Requirements

| Requirement | Input | Output | Inferred/Assumed |
| :--- | :--- | :--- | :--- |
| **Natural Language Understanding** | A plain English sentence from the user. | A structured Prowler command. | Assumed |
| **Command Translation** | A user query about resources, checks, or compliance. | The correct combination of provider, services, filters, and flags. | Inferred |
| **Contextual Awareness** | A follow-up question like "now show me only the critical ones". | The previous command context is used to refine the new command. | Inferred |
| **Command Execution** | A valid Prowler command generated by the AI. | The standard output of the Prowler scan. | Assumed |
| **User-Friendly Display** | Raw Prowler output (e.g., JSON). | A summarized, formatted, and easy-to-read display of the results. | Inferred |
| **Safety Mechanism** | A command translated by the AI. | A confirmation prompt for any commands that are destructive or alter state (if ever supported). | Assumed |

## 6. Non-Functional Requirements

*   **Performance:** The round-trip time for AI translation (query -> translated command) should be under 3 seconds to feel interactive.
*   **Accuracy:** The AI's command translation must have a high degree of accuracy. The system should provide a confidence score and allow the user to see and approve the generated command before execution.
*   **Security:** No sensitive information (e.g., cloud credentials, scan results) should be sent to the external LLM. Only the sanitized natural language query and non-sensitive Prowler command context are sent.
*   **Cost:** The cost of LLM API calls must be managed. The system should implement caching for identical queries.
*   **Graceful Degradation:** If the LLM API is unavailable or fails, the interactive shell should gracefully fall back to the standard CLI mode, informing the user of the issue.
