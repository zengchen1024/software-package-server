// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/v1/cla": {
            "get": {
                "description": "verify cla",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "CLA"
                ],
                "summary": "verify cla",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.claSingedResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/sig": {
            "get": {
                "description": "list sigs",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Sig"
                ],
                "summary": "list sigs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/sigvalidator.Sig"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg": {
            "get": {
                "description": "list software packages",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "list software packages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "importer of the softwarePkg",
                        "name": "importer",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "phase of the softwarePkg",
                        "name": "phase",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "count per page",
                        "name": "count_per_page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "page num which starts from 1",
                        "name": "page_num",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.SoftwarePkgsDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            },
            "post": {
                "description": "apply a new software package",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "apply a new software package",
                "parameters": [
                    {
                        "description": "body of applying a new software package",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.softwarePkgRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/app.NewSoftwarePkgDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/:id": {
            "put": {
                "description": "update application of software package",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "update application of software package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body of updating software package application",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.softwarePkgRequest"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/{id}": {
            "get": {
                "description": "get software package",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "get software package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.SoftwarePkgReviewDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/{id}/review/abandon": {
            "put": {
                "description": "abandon software package",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "abandon software package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/{id}/review/approve": {
            "put": {
                "description": "approve software package",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "approve software package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/{id}/review/comment": {
            "post": {
                "description": "create a new software package review comment",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "create a new software package review comment",
                "parameters": [
                    {
                        "description": "body of creating a new software package review comment",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reviewCommentRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/{id}/review/comment/{cid}/translate": {
            "post": {
                "description": "translate review comment",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "translate review comment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "cid of review comment",
                        "name": "cid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body of translate review comment",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.translationCommentRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/app.TranslatedReveiwCommentDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/{id}/review/reject": {
            "put": {
                "description": "reject software package",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "reject software package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/softwarepkg/{id}/review/rerunci": {
            "put": {
                "description": "rerun ci of software package",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SoftwarePkg"
                ],
                "summary": "rerun ci of software package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of software package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.NewSoftwarePkgDTO": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "app.SoftwarePkgApplicationDTO": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                },
                "sig": {
                    "type": "string"
                },
                "spec_url": {
                    "type": "string"
                },
                "src_rpm_url": {
                    "type": "string"
                }
            }
        },
        "app.SoftwarePkgApproverDTO": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "is_tc": {
                    "type": "boolean"
                }
            }
        },
        "app.SoftwarePkgBasicInfoDTO": {
            "type": "object",
            "properties": {
                "applied_at": {
                    "type": "string"
                },
                "ci_status": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "importer": {
                    "type": "string"
                },
                "phase": {
                    "type": "string"
                },
                "pkg_name": {
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                },
                "repo_link": {
                    "type": "string"
                },
                "sig": {
                    "type": "string"
                }
            }
        },
        "app.SoftwarePkgOperationLogDTO": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "app.SoftwarePkgReviewCommentDTO": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "since_creation": {
                    "type": "integer"
                }
            }
        },
        "app.SoftwarePkgReviewDTO": {
            "type": "object",
            "properties": {
                "application": {
                    "$ref": "#/definitions/app.SoftwarePkgApplicationDTO"
                },
                "applied_at": {
                    "type": "string"
                },
                "approved_by": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/app.SoftwarePkgApproverDTO"
                    }
                },
                "ci_status": {
                    "type": "string"
                },
                "comments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/app.SoftwarePkgReviewCommentDTO"
                    }
                },
                "desc": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "importer": {
                    "type": "string"
                },
                "logs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/app.SoftwarePkgOperationLogDTO"
                    }
                },
                "phase": {
                    "type": "string"
                },
                "pkg_name": {
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                },
                "rejected_by": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/app.SoftwarePkgApproverDTO"
                    }
                },
                "repo_link": {
                    "type": "string"
                },
                "sig": {
                    "type": "string"
                }
            }
        },
        "app.SoftwarePkgsDTO": {
            "type": "object",
            "properties": {
                "pkgs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/app.SoftwarePkgBasicInfoDTO"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "app.TranslatedReveiwCommentDTO": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                }
            }
        },
        "controller.ResponseData": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        },
        "controller.claSingedResp": {
            "type": "object",
            "properties": {
                "signed": {
                    "type": "boolean"
                }
            }
        },
        "controller.reviewCommentRequest": {
            "type": "object",
            "required": [
                "comment"
            ],
            "properties": {
                "comment": {
                    "type": "string"
                }
            }
        },
        "controller.softwarePkgRequest": {
            "type": "object",
            "required": [
                "desc",
                "pkg_name",
                "platform",
                "reason",
                "sig",
                "spec_url",
                "src_rpm_url"
            ],
            "properties": {
                "desc": {
                    "type": "string"
                },
                "pkg_name": {
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                },
                "sig": {
                    "type": "string"
                },
                "spec_url": {
                    "type": "string"
                },
                "src_rpm_url": {
                    "type": "string"
                }
            }
        },
        "controller.translationCommentRequest": {
            "type": "object",
            "properties": {
                "language": {
                    "type": "string"
                }
            }
        },
        "sigvalidator.Sig": {
            "type": "object",
            "properties": {
                "en_feature": {
                    "type": "string"
                },
                "en_group": {
                    "type": "string"
                },
                "feature": {
                    "type": "string"
                },
                "group": {
                    "type": "string"
                },
                "sig_names": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
