// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2022-12-11 17:15:09.868894726 +0300 +03 m=+0.036989016

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/login": {
            "get": {
                "description": "Try to authorize in system. JWT-Token send with success",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Try to authorize in system",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Student Number",
                        "name": "stud-number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseWithJWTMessage"
                        }
                    },
                    "403": {
                        "description": "Пароль введен неверно!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/rooms": {
            "get": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "View full information about rooms in dormitory.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rooms"
                ],
                "summary": "Get all rooms in dormitory",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page param for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Size param for pagination",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/objects.RoomResponseDTO"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/rooms/{room-id}": {
            "get": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "View information about room in dormitory.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rooms"
                ],
                "summary": "Get room information",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Room id",
                        "name": "room-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/objects.RoomResponseDTO"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Комната не найдена!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/student-live-acts": {
            "post": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Settle/evic student in certain room.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "Settle/evic student in dormitory",
                "parameters": [
                    {
                        "description": "Параметры запроса. Если roomID == 0, то студент выселяется.",
                        "name": "requestParams",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.StudentLiveActsRequestMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Данные о студенте успешно обновлены!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметр должен быть числом!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Студент не найден\" | \"Комната не найдена",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Студент уже живёт в другой комнате!\" | \"Студент уже нигде не живёт!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/student-things-acts": {
            "post": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Give thing to student without changing its location.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "Action with things by students",
                "parameters": [
                    {
                        "description": "Параметры запроса. У поля status 2 значения: give(выдать) и return (забрать)",
                        "name": "TransferParams",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.StudentThingsActsRequestMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Вещь передана!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметр должен быть числом!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Студент не найден\" | \"Вещь не найдена!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Вещь уже у другого студента!\" | \"Вещь и так не у студента.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/students": {
            "get": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "View full information about students have lived in dormitory.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "Get all students in dormitory",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page param for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Size param for pagination",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/objects.StudentResponseDTO"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Add new student in user and student base.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "Add new student to base",
                "parameters": [
                    {
                        "description": "student base information",
                        "name": "user-params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AddNewStudentRequestMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Операция успешно проведена!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметр должен быть числом!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Пользователь с таким логином уже существует!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/students/{stud-number}": {
            "get": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "View full information about student.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "View information about student",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Student Number",
                        "name": "stud-number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Данные о студенте успешно обновлены!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметр должен быть числом!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Студент не найден",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Change in database group information about student.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "Change student group for student",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Student Number",
                        "name": "stud-number",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New student group",
                        "name": "new-group",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ChangeStudentGroupRequestMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Данные о студенте успешно обновлены!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметр должен быть числом!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Студент не найден",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/things": {
            "get": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Get full information about different things.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "things"
                ],
                "summary": "Get things in dormitory with params",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page param for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Size param for pagination",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Status defines mode: all things, free things or student things. Possible values: all, free, student",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Student number for searching in Student mode",
                        "name": "stud-number",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ThingFullInfo"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметры указаны неверно",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Студент не найден",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Add new thing with params in base.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "things"
                ],
                "summary": "Add new thing",
                "parameters": [
                    {
                        "description": "body for buy service",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AddNewThingRequestMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Операция успешно проведена!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметр должен быть числом!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Вещь с таким же уникальным номером уже есть в базе",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/things/{mark-number}": {
            "get": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Get full information about thing by mark-number.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "things"
                ],
                "summary": "Get thing info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Mark number for thing",
                        "name": "mark-number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ThingFullInfo"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметры указаны неверно",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Вещь не найдена",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "JWT-Token": []
                    }
                ],
                "description": "Transfer thing to another room.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "things"
                ],
                "summary": "Transfer thing to another room.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Thing mark number",
                        "name": "mark-number",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Dst room in which thing will be transferred.",
                        "name": "room-id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TransferThingRequestMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Вещь успешно перемещена!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Параметр не должен быть пустой\" | \"Параметр должен быть числом!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "403": {
                        "description": "У вас нет достаточно прав!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Вещь не найдена\" | \"Комната не найдена",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Вещь уже находится в этой комнате!",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Проблемы на стороне сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.ShortResponseMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AddNewStudentRequestMessage": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "studentNumber": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "models.AddNewThingRequestMessage": {
            "type": "object",
            "properties": {
                "markNumber": {
                    "type": "integer"
                },
                "thingType": {
                    "type": "string"
                }
            }
        },
        "models.ChangeStudentGroupRequestMessage": {
            "type": "object",
            "properties": {
                "newGroup": {
                    "type": "string"
                }
            }
        },
        "models.ResponseWithJWTMessage": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "models.ShortResponseMessage": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string"
                }
            }
        },
        "models.StudentLiveActsRequestMessage": {
            "type": "object",
            "properties": {
                "roomID": {
                    "type": "integer"
                },
                "student-number": {
                    "type": "string"
                }
            }
        },
        "models.StudentThingsActsRequestMessage": {
            "type": "object",
            "properties": {
                "mark-number": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "student-number": {
                    "type": "string"
                }
            }
        },
        "models.ThingFullInfo": {
            "type": "object",
            "properties": {
                "thing": {
                    "type": "object",
                    "$ref": "#/definitions/objects.Thing"
                }
            }
        },
        "models.TransferThingRequestMessage": {
            "type": "object",
            "properties": {
                "room-id": {
                    "type": "integer"
                }
            }
        },
        "objects.RoomResponseDTO": {
            "type": "object",
            "properties": {
                "room-id": {
                    "type": "integer"
                },
                "room-number": {
                    "type": "integer"
                },
                "room-type": {
                    "type": "string"
                }
            }
        },
        "objects.StudentResponseDTO": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "roomID": {
                    "type": "integer"
                },
                "studentGroup": {
                    "type": "string"
                },
                "studentNumber": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "objects.Thing": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "markNumber": {
                    "type": "integer"
                },
                "ownerID": {
                    "type": "integer"
                },
                "roomID": {
                    "type": "integer"
                },
                "thingType": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8082",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "BMSTU-WEB API",
	Description: "Server for bmstu dormitory app.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
