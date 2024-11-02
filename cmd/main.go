package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id    int    `gorm:"primaryKey;column:id" json:"id"`
	Name  string `gorm:"column:name" json:"name"`
	Email string `gorm:"column:email" json:"email"`
}

var db *gorm.DB

func initDatabase() {
	var err error
	dsn := "root:toor@tcp(127.0.0.1:3306)/practica_7?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error al conectar con la base de datos")
	}
	// Migrar la estructura User a la base de datos
	db.AutoMigrate(&User{})
}

func main() {
	// Inicializar la conexión a la base de datos
	initDatabase()

	router := gin.Default()

	// Cargar los archivos de la carpeta template
	router.LoadHTMLGlob("templates/*")

	// Ping para verificar que el servidor esté en funcionamiento
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Página principal
	router.GET("/", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.HTML(200, "index.html", gin.H{
			"title":       "Main website",
			"total_users": len(users),
			"users":       users,
		})
	})

	// Obtener usuarios
	router.GET("/api/users", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.JSON(200, users)
	})

	// Crear usuario
	router.POST("/api/users", func(c *gin.Context) {
		var user User
		if c.BindJSON(&user) == nil {
			db.Create(&user)
			c.JSON(200, user)
		} else {
			c.JSON(400, gin.H{
				"error": "Invalid payload",
			})
		}
	})

	// Eliminar usuario
	router.DELETE("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return
		}

		db.Delete(&User{}, idParsed)
		c.JSON(200, gin.H{
			"message": "User Deleted",
		})
	})

	// Actualizar usuario
	router.PUT("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return
		}
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid payload",
			})
			return
		}

		db.Model(&User{}).Where("id = ?", idParsed).Updates(User{Name: user.Name, Email: user.Email})
		c.JSON(200, user)
	})

	// Ejecutar el servidor en el puerto 8001
	router.Run(":8001")
}
