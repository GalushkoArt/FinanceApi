// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/api/v1/symbols": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": [
                            "client",
                            "admin"
                        ]
                    }
                ],
                "description": "Get all available latest symbols",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbols"
                ],
                "summary": "GetSymbols",
                "operationId": "get-symbols",
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Symbol"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "404": {
                        "description": "Data not found",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": [
                            "admin"
                        ]
                    }
                ],
                "description": "Update symbol data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbols"
                ],
                "summary": "UpdateSymbols",
                "operationId": "update-symbols",
                "parameters": [
                    {
                        "description": "Update symbol data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UpdateSymbol"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Add successfully",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "400": {
                        "description": "Client request errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": [
                            "client",
                            "admin"
                        ]
                    }
                ],
                "description": "Add new symbol data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbols"
                ],
                "summary": "AddSymbols",
                "operationId": "add-symbols",
                "parameters": [
                    {
                        "description": "New symbol data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Symbol"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Add successfully",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "400": {
                        "description": "Client request errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/symbols/{symbol}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": [
                            "client",
                            "admin"
                        ]
                    }
                ],
                "description": "Get latest data for particular symbol",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbols"
                ],
                "summary": "GetSymbol",
                "operationId": "get-symbol",
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Symbol"
                            }
                        }
                    },
                    "400": {
                        "description": "Client request error",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "404": {
                        "description": "Client request error",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": [
                            "admin"
                        ]
                    }
                ],
                "description": "Delete data for symbol",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbols"
                ],
                "summary": "DeleteSymbol",
                "operationId": "delete-symbol",
                "responses": {
                    "200": {
                        "description": "Deleted successfully",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "400": {
                        "description": "Client request errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "404": {
                        "description": "Client request errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "get": {
                "description": "Refresh auth token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Refresh",
                "operationId": "refresh-token",
                "responses": {
                    "200": {
                        "description": "Response with jwt token",
                        "schema": {
                            "$ref": "#/definitions/model.SuccessfulAuthentication"
                        }
                    },
                    "400": {
                        "description": "Wrong refresh token",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            }
        },
        "/auth/signin": {
            "put": {
                "description": "Authenticate user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignIn",
                "operationId": "sign-in",
                "parameters": [
                    {
                        "description": "Authentication user data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SignIn"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Response with jwt token",
                        "schema": {
                            "$ref": "#/definitions/model.SuccessfulAuthentication"
                        }
                    },
                    "400": {
                        "description": "Wrong user data",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "401": {
                        "description": "Wrong credentials",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "put": {
                "description": "Register new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignUp",
                "operationId": "sign-up",
                "parameters": [
                    {
                        "description": "New user data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SignUp"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "New user created successfully",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "400": {
                        "description": "Wrong user data",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server errors",
                        "schema": {
                            "$ref": "#/definitions/handler.CommonResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.CommonResponse": {
            "type": "object",
            "properties": {
                "authErrors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.AuthError"
                    }
                },
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "model.AuthError": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "rule": {
                    "type": "string"
                }
            }
        },
        "model.Exchange": {
            "type": "object",
            "properties": {
                "country": {
                    "type": "string"
                },
                "mic_code": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "timezone": {
                    "type": "string"
                }
            }
        },
        "model.Price": {
            "type": "object",
            "properties": {
                "close": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "high": {
                    "type": "string"
                },
                "low": {
                    "type": "string"
                },
                "open": {
                    "type": "string"
                },
                "volume": {
                    "type": "string"
                }
            }
        },
        "model.SignIn": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 3
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 6
                }
            }
        },
        "model.SignUp": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 6
                },
                "username": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 3
                }
            }
        },
        "model.SuccessfulAuthentication": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "model.Symbol": {
            "type": "object",
            "required": [
                "symbol"
            ],
            "properties": {
                "currency": {
                    "type": "string"
                },
                "currency_base": {
                    "type": "string"
                },
                "currency_quote": {
                    "type": "string"
                },
                "exchanges": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Exchange"
                    }
                },
                "name": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "values": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Price"
                    }
                }
            }
        },
        "model.UpdateSymbol": {
            "type": "object",
            "required": [
                "symbol"
            ],
            "properties": {
                "currency": {
                    "type": "string"
                },
                "currency_base": {
                    "type": "string"
                },
                "currency_quote": {
                    "type": "string"
                },
                "exchanges": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Exchange"
                    }
                },
                "name": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "values": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Price"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "scopes": {
                "admin": " Grants read and write access to resources",
                "client": " Grants read access to resources"
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
	Title:            "Finance API",
	Description:      "Finance REST API for equities, fx and crypto rates.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
