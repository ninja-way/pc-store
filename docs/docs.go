// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/computer": {
            "post": {
                "description": "Add new pc from request body to database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "computer"
                ],
                "summary": "Add Computer",
                "parameters": [
                    {
                        "description": "computer and its accessories",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PC"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/transport.IDResponse"
                        }
                    },
                    "400": {
                        "description": "bad request body"
                    },
                    "500": {
                        "description": "add pc to database error"
                    }
                }
            }
        },
        "/computer/{id}": {
            "get": {
                "description": "Get pc from database by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "computer"
                ],
                "summary": "Get Computer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Computer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PC"
                        }
                    },
                    "400": {
                        "description": "bad id passed"
                    },
                    "404": {
                        "description": "pc with passed id not found"
                    }
                }
            },
            "put": {
                "description": "Update existing pc in database by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "computer"
                ],
                "summary": "Update Computer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Computer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "new computer or some new accessories",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PC"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "pc updated"
                    },
                    "400": {
                        "description": "pc with passed id not found"
                    }
                }
            },
            "delete": {
                "description": "Delete pc from database by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "computer"
                ],
                "summary": "Delete Computer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Computer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "pc deleted"
                    },
                    "400": {
                        "description": "pc with passed id not found"
                    }
                }
            }
        },
        "/computers": {
            "get": {
                "description": "Get all pc from database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "computers"
                ],
                "summary": "Get Computers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.PC"
                            }
                        }
                    },
                    "500": {
                        "description": "get pcs from database error"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.PC": {
            "type": "object",
            "properties": {
                "added_at": {
                    "type": "string",
                    "example": "2023-01-01T00:00:00.00000Z"
                },
                "cpu": {
                    "type": "string",
                    "example": "i9"
                },
                "data_storage": {
                    "type": "string",
                    "example": "ssd 1tb"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Super PC"
                },
                "price": {
                    "type": "integer",
                    "example": 79999
                },
                "ram": {
                    "type": "integer",
                    "example": 32
                },
                "videocard": {
                    "type": "string",
                    "example": "RTX"
                }
            }
        },
        "transport.IDResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v1.6.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Computer store API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
