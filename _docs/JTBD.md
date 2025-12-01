
# Jobs to be Done (JTBD) Documentation for Prowler

## 1. Core Job Definition

The primary job that customers "hire" Prowler to do is to gain a clear, comprehensive, and continuous understanding of their cloud security and compliance posture.

> **When** I am responsible for my organization's cloud infrastructure, **I want to** proactively identify, understand, and report on security risks and compliance gaps, **so I can** effectively mitigate threats, pass audits, and prevent security incidents.

This core job revolves around moving from a state of uncertainty and potential risk to a state of confidence and control over the cloud environment.

## 2. Job Map

This diagram illustrates the main job and the surrounding functional and emotional jobs that customers are trying to accomplish.

```mermaid
graph TD
    subgraph "Main Job To Be Done"
        A("Continuously Understand and Improve Cloud Security Posture")
    end

    subgraph "Functional Jobs"
        B("Assess Security Posture")
        C("Ensure Compliance")
        D("Monitor for Threats")
        E("Remediate Vulnerabilities")
    end

    subgraph "Emotional/Social Jobs"
        F("Reduce Security Anxiety")
        G("Gain Confidence in Security")
        H("Demonstrate Competence to Leadership & Auditors")
    end

    A --> B
    A --> C
    A --> D
    A --> E
    A --> F
    A --> G
    A --> H
```

## 3. Related Jobs

### Supporting Jobs
These are the smaller, tactical tasks that need to be done to achieve the core job. Prowler is hired to handle these tasks efficiently.

- **Run security assessments:** Execute scans across multiple cloud providers (AWS, Azure, GCP, etc.) and accounts to find misconfigurations.
- **Check for compliance:** Validate the cloud environment against hundreds of controls from frameworks like CIS, NIST, GDPR, and HIPAA.
- **Visualize findings:** See a high-level overview of the security posture through a dashboard or detailed reports.
- **Generate reports:** Create audit-ready compliance reports and evidence.
- **Integrate into pipelines:** Embed security checks into the CI/CD process to catch issues before they reach production ("shift-left").
- **Manage findings:** Triage, prioritize, and mute findings to focus on what matters most.

### Complementary Jobs
These are jobs that are often done in conjunction with using Prowler.

- **User and access management:** Onboard team members (Security, DevOps) and assign roles and permissions within the Prowler App.
- **Workflow integration:** Connect Prowler to other tools like Jira or Slack to create tickets and send notifications based on findings.

## 4. Customer Journey and Desired Outcomes

This diagram shows the typical steps a customer takes to get their job done with Prowler, along with the desired outcomes at each stage.

### Progress Diagram

```mermaid
graph TD
    Start("Discover Prowler") --> Step1("Install & Configure");
    Step1 --> Step2("Connect Cloud Accounts");
    Step2 --> Step3("Run First Security Scan");
    Step3 --> Step4("Analyze Findings in Dashboard/Report");
    Step4 --> Step5{Finding is critical?};
    Step5 -- Yes --> Step6a("Prioritize & Assign Remediation");
    Step5 -- No --> Step6b("Mute or Acknowledge Finding");
    Step6a --> Step7("Remediate in Cloud Environment");
    Step7 --> Step8("Re-run Scan to Verify Fix");
    Step8 --> Step4;
    Step6b --> Step4;
```

### Desired Outcomes
These are the measurable results customers expect from "hiring" Prowler.

| Desired Outcome | Metric / Indicator of Success |
| :--- | :--- |
| **Minimize time to value** | - Set up and run the first scan in **under 15 minutes**. |
| **Increase operational efficiency** | - Reduce the time spent manually checking for cloud misconfigurations by **>90%**. <br> - Automate the generation of compliance evidence, saving hours of manual work per audit. |
| **Improve security posture** | - Decrease the time to detect critical vulnerabilities from days/weeks to **minutes**. <br> - Increase the coverage of security checks across all cloud assets to **over 95%**. |
| **Enhance clarity and actionability** | - Reduce the time needed to understand the impact and remediation steps for a finding. <br> - Decrease the number of false positive findings that require investigation. |
| **Improve cross-team collaboration** | - Provide a single, shared view of security posture for Security, DevOps, and Compliance teams. |

## 5. Outcome Hierarchy

This diagram shows how the lower-level metrics and desired outcomes contribute to the ultimate high-level goals of the customer.

```mermaid
graph TD
    subgraph "Ultimate Goal"
        Z("A Secure and Compliant Cloud Environment")
    end
    subgraph "Key Customer Goals"
        Y1("Reduced Risk of a Security Breach")
        Y2("Effortless & Successful Audits")
        Y3("Increased Developer Velocity & Confidence")
    end
    subgraph "Necessary Outcomes (Metrics)"
        X1("Faster Time to Detection for Vulnerabilities")
        X2("Comprehensive Security Check Coverage")
        X3("Automated, Audit-Ready Reporting")
        X4("Fewer False Positives to Investigate")
        X5("Clear and Actionable Remediation Guidance")
    end
    X1 --> Y1
    X2 --> Y1
    X3 --> Y2
    X4 --> Y2
    X4 --> Y3
    X5 --> Y3

    Y1 --> Z
    Y2 --> Z
    Y3 --> Z
```
