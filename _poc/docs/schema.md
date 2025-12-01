# Graph Schema Definition

This document defines the graph database schema for the Dynamic Risk Graph PoC.

## Document Collections

- **S3Bucket:** Represents an S3 bucket in AWS.
  - **Example:** `{ "_key": "my-corp-bucket", "name": "my-corp-bucket", "is_public": true }`
- **Tag:** Represents a resource tag.
  - **Example:** `{ "_key": "sensitivity_high", "key": "sensitivity", "value": "high" }`

## Edge Collections

- **has_tag:** Connects an `S3Bucket` to a `Tag`.
  - **Example:** `_from: S3Bucket/my-corp-bucket`, `_to: Tag/sensitivity_high`
