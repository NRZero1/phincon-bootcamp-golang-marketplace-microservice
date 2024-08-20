package handler

import "github.com/gin-gonic/gin"

type BalanceHandlerInterface interface {
	BalanceFindByID
	BalanceDeduct
	BalanceAddBalance
	BalanceGetAll
}

type BalanceFindByID interface {
	FindByID(context *gin.Context)
}

type BalanceGetAll interface {
	GetAll(context *gin.Context)
}

type BalanceDeduct interface {
	Deduct(context *gin.Context)
}

type BalanceAddBalance interface {
	AddBalance(context *gin.Context)
}

type BalanceGetAllOrFindByName interface {
	GetAllOrFindByName(context *gin.Context)
}
