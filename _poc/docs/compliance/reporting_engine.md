
# Technical Specification: Generic and Extensible Reporting Engine

This document details the technical design for a generic, extensible reporting engine for Prowler. This engine will replace the current system of custom per-framework Python scripts with a unified and more powerful solution.

## 1. Executive Summary

The current approach to reporting in Prowler is fragmented and inefficient. Each compliance framework requires a dedicated Python script to parse its output and format it into a report. This creates a tight coupling between the data and its presentation, making the system difficult to maintain and extend.

This proposal outlines a new architecture where a single, generic reporting engine can process the results of any compliance scan and generate reports in multiple formats (HTML, PDF, JSON, CSV). This will be achieved by standardizing the data format of scan results and using a template-based approach for rendering the output.

## 2. Current State Analysis

- **Fragmented Logic:** Reporting logic is scattered across numerous Python scripts (e.g., `aws_audit_manager_control_tower_guardrails_aws.py`).
- **Boilerplate Code:** Each script contains a significant amount of boilerplate code for data processing and formatting.
- **Difficult to Extend:** Adding a new output format (e.g., XML) requires modifying every single reporting script.
- **Inconsistent Output:** Without a centralized engine, ensuring consistent styling and structure across different reports is a manual and challenging process.

## 3. Proposed Architecture

### 3.1. Standardized Results Format

The foundation of the new reporting engine is a standardized, intermediate JSON format for all scan results. After a scan is complete, the Prowler engine will generate a single JSON file that contains all the findings, regardless of the compliance framework used. This format will be rich and structured, containing all the information needed for any type of report.

**Example: `results.json`**
```json
{
  "scan_metadata": {
    "timestamp": "2023-10-27T10:00:00Z",
    "provider": "AWS",
    "account_id": "123456789012"
  },
  "framework": {
    "name": "CIS-AWS-v1.4",
    "version": "1.4"
  },
  "results": [
    {
      "requirement_id": "4.1",
      "requirement_name": "Ensure no security groups allow ingress from 0.0.0.0/0 to port 22",
      "check_id": "securitygroup_allow_ingress_to_tcp_port",
      "status": "FAIL",
      "region": "us-east-1",
      "resource_id": "sg-0123456789abcdef0",
      "remediation": {
        "description": "Restrict ingress to the specified TCP port...",
        "cli": "aws ec2 revoke-security-group-ingress ..."
      }
    }
  ]
}
```

### 3.2. Template-Based Rendering

The reporting engine will use a powerful templating library (such as Jinja2 for Python) to generate the final reports. For each supported output format, there will be a corresponding template.

- **`report.html.j2`:** A template for generating a rich, interactive HTML report.
- **`report.csv.j2`:** A template for a simple CSV export.
- **`report.json.j2`:** A template to output a clean, customer-facing JSON report.

This approach completely decouples the data from its presentation. To add a new output format, we simply need to create a new template file.

### 3.3. The Reporting Engine CLI

A new Prowler CLI command will be introduced to drive the reporting engine:

```bash
prowler report --input results.json --format html --output report.html
prowler report --input results.json --format csv --output report.csv
```

This command will:
1.  Load the standardized `results.json` file.
2.  Select the appropriate template based on the `--format` flag.
3.  Render the template with the results data.
4.  Save the output to the specified file.

## 4. Benefits

- **Centralized Logic:** All reporting logic is now consolidated into a single, maintainable engine.
- **Radical Extensibility:** Adding a new output format is as simple as creating a new template file. There is no need to touch any Python code.
- **Guaranteed Consistency:** All reports will share a consistent structure and style, improving the user experience.
- **Decoupling:** The compliance frameworks and the scanning engine have no knowledge of the reporting formats, leading to a much cleaner and more modular architecture.

## 5. Migration Strategy

1.  **Define the Standard Format:** Finalize the schema for the intermediate `results.json` format.
2.  **Update the Engine:** Modify the Prowler scanning engine to generate the standardized results file.
3.  **Build the Reporting Engine:** Create the new `prowler report` command and the initial set of templates (HTML, CSV, JSON).
4.  **Deprecate Old Scripts:** Remove the per-framework reporting scripts.
5.  **Documentation:** Update the Prowler documentation to reflect the new, simplified reporting workflow.
