package main

import (
	"fmt"
	"net/http"
	"protocol"
	"strings"
	"utils"

	"github.com/gin-gonic/gin"
)

//StartServer 启动服务器
func StartServer(port string) {
	r := gin.Default()
	api := r.Group("/api", globalHandler)
	api.GET("/items.json", projectBaobeiListHandler)
	api.GET("/favoies.json", getUserFavoriesHandler)
	api.POST("/favory.json", addUserFavoriesHandler)
	api.DELETE("/favory.json", rmUserFavoriesHandler)
	api.GET("/user/authenticate", userAuthenticateHandler)
	api.POST("/user/password.json", setUserPasswordHandler)
	r.Run(fmt.Sprintf(":%s", port))
}

var skipAuthApi = map[string]string{
	"/api/items.json":        http.MethodGet,
	"/api/user/authenticate": http.MethodGet,
}

func isSkipAuth(queryPath, method string) bool {
	if m, ok := skipAuthApi[queryPath]; ok && m == method {
		return true
	}
	return false
}
func globalHandler(context *gin.Context) {
	queryPath := context.Request.URL.Path
	method := context.Request.Method
	// gLogger.ILog("handler by globalHandler %s,[%s] %s", context.ClientIP(), method, queryPath)
	if !isSkipAuth(queryPath, method) {
		header := context.GetHeader("Authorization")
		tokenString := strings.TrimSpace(strings.Replace(header, "Bearer ", "", 1))
		token, err := utils.ParseJwt(tokenString, jwtSignKey)
		if err != nil {
			// gLogger.ELog("parse jwt error %s", err.Error())
			// context.AbortWithStatusJSON(http.StatusForbidden, map[string]string{"error_msg": "invalid token"})
			renderWithStatus(nil, protocol.NewError(ErrCodeInvalidToken, "invalid token"), context, http.StatusForbidden)
			return
		}
		if !token.Valid {
			// gLogger.ELog("parse jwt error %s", err.Error())
			renderWithStatus(nil, protocol.NewError(ErrCodeInvalidToken, "invalid token"), context, http.StatusForbidden)
			return
		}

		if claims, ok := token.Claims.(*UserClaims); !ok {
			renderWithStatus(nil, protocol.NewError(ErrCodeInvalidToken, "invalid token"), context, http.StatusForbidden)
			return
		} else {
			if user, err := getUserByUID(claims.UID); err != nil {
				renderJSON(nil, err, nil)
				// context.Abort()
			} else {
				context.Set("user_context", user)
			}
		}

		// context.Set()
	}
	// context.q
	// context.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"hahah": "oooo"})
}

type codeVersion struct {
	GitBranch string
	GitCommit string
}
type commonResponse struct {
	Info      interface{} `json:omitempty`
	Success   bool        `json:omitempty`
	ErrorCode int         `json:omitempty`
	ErrorMsg  string      `json:omitempty`
	Version   codeVersion
}

// func (resp *commonResponse) String() string {
// 	b, _ := json.Marshal(resp)
// 	return string(b)
// }
func renderWithStatus(obj interface{}, err error, c *gin.Context, status int) {
	resp := &commonResponse{
		Version: codeVersion{GitBranch, GitCommit},
	}
	if err != nil {
		var cper *protocol.Err
		var ok bool
		if cper, ok = err.(*protocol.Err); !ok {
			cper = protocol.NewUnkownErr(err)
		}
		resp.ErrorCode = cper.Code
		resp.ErrorMsg = cper.Msg
	} else {
		if obj == nil {
			obj = struct {
			}{}
		}
		resp.Success = true
		resp.Info = obj

	}
	c.AbortWithStatusJSON(status, resp)
}

func renderJSON(obj interface{}, err error, c *gin.Context) {
	if err != nil {
		renderWithStatus(nil, err, c, http.StatusBadRequest)
	} else {
		renderWithStatus(obj, nil, c, http.StatusOK)
	}
}
