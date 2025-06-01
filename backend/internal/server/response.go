package server

import "github.com/gin-gonic/gin"

func successResponse(data interface{}) gin.H {
	return gin.H{
		"success": true,
		"data":    data,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{
		"success": false,
		"error":   err.Error(),
	}
}
