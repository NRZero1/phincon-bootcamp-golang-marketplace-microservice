package handler

import "github.com/gin-gonic/gin"

type UserHandlerInterface interface {
	UserSave
	UserFindById
	UserGetAll
	UserLogin
	UserFindByUsername
	UserGetAllOrFindByUsername
}

type UserSave interface {
	Save(context *gin.Context)
}

type UserFindById interface {
	FindById(context *gin.Context)
}

type UserGetAll interface {
	GetAll(context *gin.Context)
}

type UserLogin interface {
	Login(context *gin.Context)
}

type UserFindByUsername interface {
	FindByUsername(context *gin.Context)
}

type UserGetAllOrFindByUsername interface {
	GetAllOrFindByName(context *gin.Context)
}
