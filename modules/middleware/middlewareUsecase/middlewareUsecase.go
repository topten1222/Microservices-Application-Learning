package middlewareUsecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/topten1222/hello_sekai/config"
	middlewarerepository "github.com/topten1222/hello_sekai/modules/middleware/middlewareRepository"
	"github.com/topten1222/hello_sekai/pkg/jwtauth"
	"github.com/topten1222/hello_sekai/pkg/rbac"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(echo.Context, *config.Config, string) (echo.Context, error)
		RbacAuthorization(echo.Context, *config.Config, []int) (echo.Context, error)
		PlayerIdParamValidation(echo.Context) (echo.Context, error)
	}
	middlewareUsecase struct {
		middlewarerepository middlewarerepository.MiddlewareRepositoryHandlerService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewarerepository.MiddlewareRepositoryHandlerService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewarerepository: middlewareRepository}
}

func (u *middlewareUsecase) JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error) {
	context := c.Request().Context()
	claims, err := jwtauth.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return nil, err
	}
	if err := u.middlewarerepository.AccessTokenSearch(context, cfg.Grpc.AuthUrl, accessToken); err != nil {
		return nil, err
	}
	c.Set("player_id", claims.Id)
	c.Set("role_code", claims.RoleCode)
	return c, nil
}

func (u *middlewareUsecase) RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error) {
	ctx := c.Request().Context()
	playerRoleCode := c.Get("role_code").(int)
	fmt.Println("playerRoleCode::: ", playerRoleCode)
	roleCount, err := u.middlewarerepository.RolesCount(ctx, cfg.Grpc.AuthUrl)
	if err != nil {
		return nil, err
	}
	fmt.Println("RoleCount:: ", roleCount)

	playerRoleBinary := rbac.IntToBinary(playerRoleCode, int(roleCount))
	fmt.Println("RoleBinary:: ", playerRoleBinary)
	fmt.Println("Expected:: ", expected)

	for i := 0; i < int(roleCount); i++ {
		if playerRoleBinary[i]&expected[i] == 1 {
			return c, nil
		}
	}
	return nil, errors.New("Error: Permission denind")
}

func (u *middlewareUsecase) PlayerIdParamValidation(c echo.Context) (echo.Context, error) {

	playerIdReq := c.Param("player_id")
	fmt.Println("PlayerIDReq:::: ", playerIdReq)
	playerIdToken := c.Get("player_id").(string)

	fmt.Println("playerIdToken:::: ", playerIdToken)

	if playerIdToken == "" {
		log.Println("Error: PlayerId is empty in token")

		return nil, errors.New("Error: PlayerId is empty in token")
	}
	if playerIdReq != playerIdToken {
		log.Printf("Error: PlayerId is not match. playerIdReq: %s, playerIdToken: %s", playerIdReq, playerIdToken)

		return nil, errors.New("Error: PlayerId is not match")
	}
	return c, nil
}
