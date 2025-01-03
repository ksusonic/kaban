package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ksusonic/kanban/internal/controller/auth"
	"github.com/ksusonic/kanban/internal/logger"
	"github.com/ksusonic/kanban/internal/models"
)

func TestController_TelegramCallback(t *testing.T) {
	const (
		botName  = "Kanban"
		botToken = "super-secret-token"
	)

	botCfg := models.BotCfg{
		Name:  botName,
		Token: botToken,
	}

	type fields struct {
		mockUserRepo func(*gomock.Controller) *auth.MockuserRepo
	}
	type args struct {
		query string
	}

	tests := []struct {
		name         string
		args         args
		fields       fields
		expectedCode int
	}{
		{
			name: "success",
			args: args{
				query: "/auth/tg-callback?" +
					"next=%2f&id=123&first_name=Daniil&username=ksusonic&photo_url=pic.jpg&auth_date=123" +
					"&hash=d7590bb4172bdb9029a2f08c976aeb86167037fcf788ba6c6a2aad849c8b3d1b",
			},
			fields: fields{
				mockUserRepo: auth.NewMockuserRepo,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "no query",
			args: args{
				query: "/auth/tg-callback",
			},
			fields: fields{
				mockUserRepo: auth.NewMockuserRepo,
			},
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			mockCtrl := gomock.NewController(t)

			engine := gin.Default()
			engine.GET("/auth/tg-callback", auth.NewController(
				tt.fields.mockUserRepo(mockCtrl),
				botCfg,
				logger.NewDisabled(),
			).TelegramCallback)

			req, _ := http.NewRequest(http.MethodGet, tt.args.query, nil)
			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
