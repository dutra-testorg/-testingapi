// Package api GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package api

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health": {
            "get": {
                "description": "Health-check returning a dummy HealthCheckResponse (config)",
                "produces": [
                    "application/json"
                ],
                "summary": "Provide health-check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "\$ref": "#/definitions/pkg_rest.HealthCheckResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "\$ref": "#/definitions/ghandler.HTTPError"
                        }
                    }
                }
            }
        },
        "/v1/demo/{uid}": {
            "get": {
                "description": "demo endpoint returning a Demo struct",
                "produces": [
                    "application/json"
                ],
                "summary": "Get demo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "uuidv4 (UUIDv4)",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "\$ref": "#/definitions/internal_micro.Demo"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "\$ref": "#/definitions/ghandler.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "\$ref": "#/definitions/ghandler.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ghandler.HTTPError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "trace_id": {
                    "type": "string"
                }
            }
        },
        "github.com_gympass_scheduler_internal_micro.Demo": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "github.com_gympass_scheduler_pkg_rest.HealthCheckResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "service": {
                    "type": "string"
                }
            }
        },
        "internal_micro.Demo": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "pkg_rest.HealthCheckResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "service": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Gympass Go Example API",
	Description:      "This is a sample Golang Gympass server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
