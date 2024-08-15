package handler

import "github.com/gin-gonic/gin"

type ChannelHandlerInterface interface {
	ChannelFindById
	ChannelGetAll
	ChannelFindByName
	ChannelGetAllOrFindByName
}

type ChannelFindById interface {
	FindById(context *gin.Context)
}

type ChannelGetAll interface {
	GetAll(context *gin.Context)
}

type ChannelFindByName interface {
	FindByName(context *gin.Context)
}

type ChannelGetAllOrFindByName interface {
	GetAllOrFindByName(context *gin.Context)
}
