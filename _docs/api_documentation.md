# Prowler API Documentation

### API Overview

| Detail | Value |
| :--- | :--- |
| **API Name** | Prowler API |
| **Version** | 1.15.0 |
| **Base URL** | `https://<your-prowler-api-domain>` |
| **Authentication** | Bearer Token (JWT) or API Key in the `Authorization` header. <br> e.g., `Authorization: Bearer <JWT_TOKEN>` or `Authorization: Api-Key <API_KEY>`|

----------

## 1. User Management

### 1.1. Retrieve User Details

**Endpoint:** `GET /api/v1/users/{id}`

| Detail | Description |
| :--- | :--- |
| **Purpose** | Retrieves the profile details for a specific user using their unique ID. |
| **Method** | `GET` |
| **URL** | `/api/v1/users/{id}` |

#### Request Parameters

| Parameter Name | Location | Data Type | Required | Description/Use Case |
| :--- | :--- | :--- | :--- | :--- |
| `id` | Path | string (uuid) | Yes | The **unique identifier** of the user to retrieve. |
| `fields[users]` | Query | array[string] | No | Comma-separated list of fields to include (e.g., `name,email`). If omitted, all standard fields are returned. |
| `include` | Query | array[string] | No | Include related resources in the response (e.g., `roles`, `memberships`). |

#### Success Response

**Status Code:** `200 OK`

##### Response Body Schema

| Parameter Name | Data Type | Description/Use Case |
| :--- | :--- | :--- |
| `id` | string (uuid) | The unique user identifier. |
| `type` | string | The resource type, which is `"users"`. |
| `attributes.name` | string | The user's full name. |
| `attributes.email`| string | The primary contact email address. |
| `attributes.company_name` | string | The name of the user's company. |
| `attributes.date_joined` | date-time | Timestamp of user creation. |
| `relationships` | object | Contains links to related resources like `roles` and `memberships`. |

##### Example (JSON)
```json
{
  "data": {
    "type": "users",
    "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
    "attributes": {
      "name": "Jane Doe",
      "email": "jane.doe@example.com",
      "company_name": "Prowler",
      "date_joined": "2023-11-28T10:00:00Z"
    },
    "relationships": {
        "roles": {
            "data": [
                {"type": "roles", "id": "r1b2c3d4-e5f6-7890-1234-567890abcdef"}
            ]
        }
    }
  }
}
```

#### Error Responses

| Status Code | Description | Error Response Example |
| :--- | :--- | :--- |
| **401 Unauthorized** | Token is missing or invalid. | `{"errors": [{"detail": "Authentication credentials were not provided."}]}` |
| **404 Not Found** | The specified `id` does not exist. | `{"errors": [{"detail": "Not found."}]}` |

----------

### 1.2. Create New User

**Endpoint:** `POST /api/v1/users`

| Detail | Description |
| :--- | :--- |
| **Purpose** | Creates a new user profile in the system. |
| **Method** | `POST` |
| **URL** | `/api/v1/users` |

#### Request Body Schema
The request body must be a **JSON** object following the JSON:API specification.

| Parameter Name | Data Type | Required | Description/Use Case |
| :--- | :--- | :--- | :--- |
| `data.type` | string | Yes | Must be `"users"`. |
| `data.attributes.name`| string | Yes | The full name of the new user. |
| `data.attributes.email`| string | Yes | The user's email address (must be unique). |
| `data.attributes.password` | string | Yes | The initial password for the user. |
| `data.attributes.company_name` | string | No | The name of the user's company. |

##### Example (JSON)
```json
{
  "data": {
    "type": "users",
    "attributes": {
      "name": "John Smith",
      "email": "john.smith@example.com",
      "password": "SecurePassword123",
      "company_name": "Prowler Corp"
    }
  }
}
```

#### Success Response
**Status Code:** `201 Created`

The response returns the newly created user object, including its assigned **ID**.

##### Response Body Schema
(Same as success schema for `GET /api/v1/users/{id}`)

#### Error Responses

| Status Code | Description | Error Response Example |
| :--- | :--- | :--- |
| **400 Bad Request** | Missing a required field (e.g., `email` is null) or invalid data type. | `{"errors": [{"detail": "email: This field may not be blank."}]}` |
| **409 Conflict** | The provided `email` already exists in the system. | `{"errors": [{"detail": "A user with that email already exists."}]}` |

----------

## 2. Compliance

### 2.1. Get Compliance Requirement Attributes

**Endpoint:** `GET /api/v1/compliance-overviews/attributes`

| Detail | Description |
| :--- | :--- |
| **Purpose** | Retrieve detailed attribute information for all requirements in a specific compliance framework. |
| **Method**| `GET` |
| **URL** | `/api/v1/compliance-overviews/attributes` |

#### Request Parameters

| Parameter Name | Location | Data Type | Required | Description/Use Case |
| :--- | :--- | :--- | :--- | :--- |
| `filter[compliance_id]`| Query | string | Yes | The ID of the compliance framework to get attributes for (e.g., `cis_v1.2_aws`). |
| `fields[compliance-requirements-attributes]` | Query | array[string] | No | Specify which fields to include in the response. |

#### Success Response
**Status Code:** `200 OK`

Returns a list of requirement attributes for the specified compliance framework.

##### Response Body Schema

| Parameter Name | Data Type | Description/Use Case |
| :--- | :--- | :--- |
| `id` | string | The unique identifier for the requirement attribute. |
| `type` | string | The resource type, `"compliance-requirements-attributes"`. |
| `attributes.compliance_name` | string | The human-readable name of the compliance framework. |
| `attributes.framework`| string | The machine-readable name of the framework (e.g., "cis"). |
| `attributes.version`| string | The version of the compliance framework. |
| `attributes.description` | string | A detailed description of the compliance requirement. |
| `attributes.attributes`| object | A JSON object containing the check IDs associated with the requirement. |

##### Example (JSON)
```json
{
  "data": [
    {
      "type": "compliance-requirements-attributes",
      "id": "cis_v1.2_aws_1.1",
      "attributes": {
        "id": "1.1",
        "compliance_name": "CIS AWS Foundations Benchmark v1.2.0",
        "framework_description": "CIS AWS Foundations Benchmark v1.2.0",
        "name": "Avoid the use of the 'root' account",
        "framework": "cis",
        "version": "1.2.0",
        "description": "The 'root' account is the most privileged user in an AWS account...",
        "attributes": {
          "checks": [
            "iam_avoid_root_usage"
          ]
        }
      }
    }
  ]
}
```

#### Error Responses

| Status Code | Description | Error Response Example |
| :--- | :--- | :--- |
| **400 Bad Request** | Required `compliance_id` filter is missing or invalid. | `{"errors": [{"detail": "Filter 'compliance_id' is required."}]}` |
| **401 Unauthorized** | Token is missing or invalid. | `{"errors": [{"detail": "Authentication credentials were not provided."}]}` |
