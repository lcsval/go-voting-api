package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type userHandler struct {
	db         *sqlx.DB
	repository UserRepository
}

func newUserHandler(db *sqlx.DB) *userHandler {
	return &userHandler{
		db:         db,
		repository: UserRepository{db: *db},
	}
}

func (u *userHandler) getUser(c *gin.Context) {
	id := c.Param("id")

	if user, err := u.repository.Get(id); err != nil {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func (u *userHandler) getAllUsers(c *gin.Context) {
	if users, err := u.repository.GetAll(); err != nil {
		c.JSON(http.StatusOK, users)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func (u *userHandler) createUser(c *gin.Context) {
	user := &User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newUser, err := u.repository.CreateUser(*user)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (u *userHandler) updateUser(c *gin.Context) {
	user := &User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")
	if user.ID != "" && user.ID != id {
		c.JSON(http.StatusBadRequest, errors.New("User id is not equal URL ID"))
		return
	}

	user.ID = id
	if err := u.repository.Update(*user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *userHandler) deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := u.repository.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, err)
	}
}

func RegisterRoutes(db *sqlx.DB, router *gin.Engine) {
	handler := newUserHandler(db)

	// TODO - fazer retornar todos os usuários
	router.GET("/users", func(c *gin.Context) {
		handler.getAllUsers(c)
	})

	router.GET("/user/:id", func(c *gin.Context) {
		handler.getUser(c)
	})

	router.POST("/user", func(c *gin.Context) {
		handler.createUser(c)
	})

	router.PUT("/user/:id", func(c *gin.Context) {
		handler.updateUser(c)
	})

	router.DELETE("/user/:id", func(c *gin.Context) {
		handler.deleteUser(c)
	})
}
