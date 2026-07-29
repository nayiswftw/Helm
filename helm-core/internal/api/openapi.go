//go:build linux

package api

import (
	"net/http"
)

const openAPISpec = `{
  "openapi": "3.1.0",
  "info": {
    "title": "Helm Control API",
    "description": "Zero-dependency Linux server management and container control API.",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "/api/v1",
      "description": "Version 1 API base path"
    }
  ],
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer"
      },
      "apiKeyHeader": {
        "type": "apiKey",
        "in": "header",
        "name": "X-API-Key"
      }
    }
  },
  "security": [
    { "bearerAuth": [] },
    { "apiKeyHeader": [] }
  ],
  "paths": {
    "/dashboard": {
      "get": {
        "summary": "Get system metrics snapshot",
        "responses": {
          "200": { "description": "System metrics snapshot including CPU, Memory, Disk, Uptime, Load Average, and Temperature" }
        }
      }
    },
    "/devices": {
      "get": {
        "summary": "List managed devices",
        "responses": {
          "200": { "description": "Array of managed devices" }
        }
      }
    },
    "/devices/{id}": {
      "get": {
        "summary": "Get device by ID",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": {
          "200": { "description": "Device details" },
          "404": { "description": "Device not found" }
        }
      }
    },
    "/actions": {
      "get": {
        "summary": "List administrative actions",
        "responses": {
          "200": { "description": "Array of administrative actions" }
        }
      }
    },
    "/actions/{id}/execute": {
      "post": {
        "summary": "Execute administrative action",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": {
          "202": { "description": "Action accepted for execution" },
          "404": { "description": "Action not found" }
        }
      }
    },
    "/containers": {
      "get": {
        "summary": "List Docker containers",
        "responses": {
          "200": { "description": "Array of Docker containers" }
        }
      }
    },
    "/containers/{id}/stats": {
      "get": {
        "summary": "Get live CPU and Memory usage for a Docker container",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": {
          "200": { "description": "Container resource statistics" },
          "404": { "description": "Container not found" }
        }
      }
    },
    "/containers/{id}/logs": {
      "get": {
        "summary": "Get recent stdout/stderr log lines for a Docker container",
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } },
          { "name": "tail", "in": "query", "required": false, "schema": { "type": "integer", "default": 100 } }
        ],
        "responses": {
          "200": { "description": "Container log lines" },
          "404": { "description": "Container not found" }
        }
      }
    },
    "/containers/{id}/start": {
      "post": {
        "summary": "Start a Docker container",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Container started" } }
      }
    },
    "/containers/{id}/stop": {
      "post": {
        "summary": "Stop a Docker container",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Container stopped" } }
      }
    },
    "/containers/{id}/restart": {
      "post": {
        "summary": "Restart a Docker container",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Container restarted" } }
      }
    },
    "/dokploy/projects": {
      "get": {
        "summary": "List Dokploy projects",
        "responses": { "200": { "description": "Array of Dokploy projects" }, "503": { "description": "Dokploy not configured" } }
      }
    },
    "/dokploy/applications/{id}": {
      "get": {
        "summary": "Get Dokploy application details",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Application details" }, "503": { "description": "Dokploy not configured" } }
      }
    },
    "/dokploy/applications/{id}/deploy": {
      "post": {
        "summary": "Trigger deployment of a Dokploy application",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "202": { "description": "Deployment triggered" } }
      }
    },
    "/dokploy/applications/{id}/redeploy": {
      "post": {
        "summary": "Trigger redeployment of a Dokploy application",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "202": { "description": "Redeployment triggered" } }
      }
    },
    "/dokploy/applications/{id}/deployments": {
      "get": {
        "summary": "List deployment history for a Dokploy application",
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Deployment history" } }
      }
    },
    "/notifications/test": {
      "post": {
        "summary": "Send test webhook notification",
        "responses": { "200": { "description": "Test notification dispatched" } }
      }
    },
    "/openapi.json": {
      "get": {
        "summary": "Get OpenAPI 3.1 JSON specification",
        "responses": { "200": { "description": "OpenAPI specification document" } }
      }
    }
  }
}`

// handleOpenAPI returns the OpenAPI 3.1 JSON specification.
func handleOpenAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(openAPISpec))
	}
}
