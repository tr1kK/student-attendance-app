package router

import (
	"net/http"
	"student-attendance-app/pkg/config"
	"student-attendance-app/pkg/handlers"
	"student-attendance-app/pkg/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Swagger
	_ "student-attendance-app/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const swaggerJSON = `{
    "swagger": "2.0",
    "info": {
        "description": "Это сервер для приложения по учету посещаемости студентов.",
        "title": "API Студенческого Портала",
        "contact": {
            "name": "Поддержка API",
            "url": "https://github.com/TR1K/student-attendance-app",
            "email": "trik.py@proton.me"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/",
    "paths": {
        "/api/admin/groups": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получает список всех групп.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Получить все группы (Админ)",
                "responses": {
                    "200": {
                        "description": "Список групп",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Group"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/admin/users": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получает список всех пользователей.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Получить всех пользователей (Админ)",
                "responses": {
                    "200": {
                        "description": "Список пользователей",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Создает нового пользователя с указанными данными.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Создать пользователя (Админ)",
                "parameters": [
                    {
                        "description": "Объект пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Созданный пользователь",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/admin/users/{id}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Обновляет данные существующего пользователя.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Обновить пользователя (Админ)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID Пользователя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Объект пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Обновленный пользователь",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Удаляет пользователя по его ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Удалить пользователя (Админ)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID Пользователя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь успешно удален",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/lessons": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает список занятий. Для студентов - занятия их группы. Для преподавателей/администраторов - все занятия.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "lessons"
                ],
                "summary": "Получить занятия",
                "responses": {
                    "200": {
                        "description": "Список занятий",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Lesson"
                            }
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/student/attendance": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получает все записи о посещаемости для залогиненного студента.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "student"
                ],
                "summary": "Получить записи о посещаемости студента",
                "responses": {
                    "200": {
                        "description": "Список записей о посещаемости",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Attendance"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Студент отправляет код посещаемости для определенного занятия.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "student"
                ],
                "summary": "Отметить посещаемость",
                "parameters": [
                    {
                        "description": "Данные для отметки посещаемости",
                        "name": "attendance",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.SubmitAttendanceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Посещаемость успешно отмечена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Неверный или просроченный код",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/teacher/attendance/{lessonId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получает все записи о посещаемости для определенного занятия.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "teacher"
                ],
                "summary": "Получить посещаемость занятия",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID Занятия",
                        "name": "lessonId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список записей о посещаемости",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Attendance"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/teacher/lessons/{lessonId}/code": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Генерирует новый 5-значный код для занятия, который истекает через 5 минут.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "teacher"
                ],
                "summary": "Сгенерировать код посещаемости",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID Занятия",
                        "name": "lessonId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Сгенерированный код",
                        "schema": {
                            "$ref": "#/definitions/models.GeneratedCode"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Деактивирует активный код для занятия.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "teacher"
                ],
                "summary": "Деактивировать код посещаемости",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID Занятия",
                        "name": "lessonId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Код успешно деактивирован",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Активный код не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Аутентифицирует пользователя и возвращает JWT токен.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Вход пользователя",
                "parameters": [
                    {
                        "description": "Учетные данные для входа",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный вход",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "Неверные учетные данные",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Создает нового пользователя-студента.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Регистрация нового студента",
                "parameters": [
                    {
                        "description": "Данные для регистрации пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь успешно зарегистрирован",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "409": {
                        "description": "Конфликт (идентификатор или email уже существует)",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/groups": {
            "get": {
                "description": "Возвращает список всех студенческих групп.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public"
                ],
                "summary": "Получить все группы",
                "responses": {
                    "200": {
                        "description": "Список групп",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Group"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.LoginRequest": {
            "type": "object",
            "required": [
                "identifier",
                "password"
            ],
            "properties": {
                "identifier": {
                    "type": "string",
                    "example": "student1"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                }
            }
        },
        "handlers.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "group_id",
                "identifier",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "new@example.com"
                },
                "group_id": {
                    "type": "integer",
                    "example": 1
                },
                "identifier": {
                    "type": "string",
                    "example": "newstudent"
                },
                "name": {
                    "type": "string",
                    "example": "New Student"
                },
                "password": {
                    "type": "string",
                    "example": "securepassword"
                }
            }
        },
        "handlers.SubmitAttendanceRequest": {
            "type": "object",
            "required": [
                "code",
                "lesson_id"
            ],
            "properties": {
                "code": {
                    "type": "string",
                    "example": "12345"
                },
                "lesson_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "models.Attendance": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "lesson": {
                    "$ref": "#/definitions/models.Lesson"
                },
                "lesson_id": {
                    "type": "integer"
                },
                "student": {
                    "$ref": "#/definitions/models.User"
                },
                "student_id": {
                    "type": "integer"
                },
                "submitted_at": {
                    "type": "string"
                }
            }
        },
        "models.GeneratedCode": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "expires_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_active": {
                    "type": "boolean"
                },
                "lesson": {
                    "$ref": "#/definitions/models.Lesson"
                },
                "lesson_id": {
                    "type": "integer"
                }
            }
        },
        "models.Group": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.Lesson": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "day": {
                    "type": "string"
                },
                "group": {
                    "$ref": "#/definitions/models.Group"
                },
                "group_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "room": {
                    "type": "string"
                },
                "teacher": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "group": {
                    "$ref": "#/definitions/models.Group"
                },
                "group_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "identifier": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Swagger endpoint
	r.GET("/swagger.json", func(c *gin.Context) {
		c.String(http.StatusOK, swaggerJSON)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))

	// Public routes
	r.GET("/groups", func(c *gin.Context) {
		handlers.GetGroups(c, db)
	})

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", func(c *gin.Context) {
			handlers.Login(c, db, cfg)
		})
		authRoutes.POST("/register", func(c *gin.Context) {
			handlers.Register(c, db, cfg)
		})
	}

	// Authenticated routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		// Lesson routes (accessible to all authenticated users)
		api.GET("/lessons", func(c *gin.Context) {
			handlers.GetLessons(c, db)
		})

		// Student routes
		studentRoutes := api.Group("/student")
		studentRoutes.Use(middleware.RoleMiddleware("student"))
		{
			studentRoutes.POST("/attendance", func(c *gin.Context) {
				handlers.SubmitAttendance(c, db)
			})
			studentRoutes.GET("/attendance", func(c *gin.Context) {
				handlers.GetStudentAttendance(c, db)
			})
		}

		// Teacher routes
		teacherRoutes := api.Group("/teacher")
		teacherRoutes.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherRoutes.POST("/lessons/:lessonId/code", func(c *gin.Context) {
				handlers.GenerateCode(c, db)
			})
			teacherRoutes.DELETE("/lessons/:lessonId/code", func(c *gin.Context) {
				handlers.DeactivateCode(c, db)
			})
			teacherRoutes.GET("/attendance/:lessonId", func(c *gin.Context) {
				handlers.GetLessonAttendance(c, db)
			})
		}

		// Admin routes
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.RoleMiddleware("admin"))
		{
			adminRoutes.GET("/users", func(c *gin.Context) { handlers.AdminGetUsers(c, db) })
			adminRoutes.POST("/users", func(c *gin.Context) { handlers.AdminCreateUser(c, db) })
			adminRoutes.PUT("/users/:id", func(c *gin.Context) { handlers.AdminUpdateUser(c, db) })
			adminRoutes.DELETE("/users/:id", func(c *gin.Context) { handlers.AdminDeleteUser(c, db) })
			adminRoutes.GET("/groups", func(c *gin.Context) { handlers.AdminGetGroups(c, db) })
		}
	}
}
