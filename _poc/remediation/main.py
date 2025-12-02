
from flask import Flask, request, jsonify
import os

app = Flask(__name__)

REMEDIATIONS_DIR = "/app/remediations"

def call_llm(prompt):
    """
    Simulates a call to a Large Language Model.
    In a real-world scenario, this would make an API request to an LLM service.
    """
    print(f"Simulating LLM call with prompt:\n{prompt}")

    # Realistic placeholder for AI-generated HCL code
    # In a real implementation, this would be the dynamic response from the LLM
    resource_id = "prowler-poc-bucket" # Extract from prompt or pass as argument
    hcl_code = f'''
resource "aws_s3_bucket_public_access_block" "{resource_id}_access_block" {{
  bucket = "{resource_id}"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls    = true
  restrict_public_buckets = true
}}
'''
    return hcl_code

def generate_remediation(finding):
    """
    Constructs a prompt, calls the simulated LLM, and saves the remediation.
    """
    resource_id = finding.get("resource_id", "unknown_resource")
    message = finding.get("message", "No message provided.")

    prompt = f"""
    A security finding has been detected:
    - Resource ID: {resource_id}
    - Description: {message}

    Please generate the Terraform HCL code to remediate this finding.
    The remediation should enforce that the S3 bucket is not public.
    """

    # Call the simulated LLM to get the remediation code
    remediation_code = call_llm(prompt)

    # Save the remediation to a file
    if not os.path.exists(REMEDIATIONS_DIR):
        os.makedirs(REMEDIATIONS_DIR)
    
    file_path = os.path.join(REMEDIATIONS_DIR, f"{resource_id}.tf")
    with open(file_path, "w") as f:
        f.write(remediation_code)
    
    print(f"Saved remediation to: {file_path}")

    return remediation_code

@app.route("/", methods=['POST'])
def remediate():
    """
    Receives a finding from the engine, generates a remediation, 
    and returns it.
    """
    finding_data = request.json
    print(f"Received finding: {finding_data}")
    
    remediation_code = generate_remediation(finding_data)
    
    print(f"Generated remediation:\n{remediation_code}")
    return remediation_code, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)

