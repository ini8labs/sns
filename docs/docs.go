// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": [http,https],
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
        "/login/otp": {
            "post": {
                "tags": ["OTP APIs"],
                "parameters": [
                    {
                        "description": "enter a valid Phone Number",
                        "name": "PhoneNumber",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "phone_number":{
                                    "type": "string",
                                    "example": "+916379430684"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Status ok",
                        "schema": {
                            "type": "string",
                            "example": "Status ok"
                        }
                    }
                }
            }
        },
        "/login/verify": {
            "post": {
                "tags": ["OTP APIs"],
                "parameters": [
                    {
                        "description": "enter a valid OTP",
                        "name": "otp",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "otp":{
                                    "type": "string",
                                    "example": "123465"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Status ok",
                        "schema": {
                            "type": "string",
                            "example": "Status ok"
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "2.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "Customer-manager APIs",
	Description:      "This is Lottery SMS Notification Service API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
