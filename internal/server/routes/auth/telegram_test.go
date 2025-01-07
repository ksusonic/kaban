package auth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ksusonic/kanban/internal/logger"
	"github.com/ksusonic/kanban/internal/models"
	"github.com/ksusonic/kanban/internal/server/routes/auth"
)

func TestController_TelegramCallback(t *testing.T) {
	botCfg := models.BotCfg{
		Name:  "test bot",
		Token: "super-secret-token",
	}

	type fields struct {
		MockuserRepo   func(*gomock.Controller) *auth.MockuserRepo
		MockauthModule func(*gomock.Controller) *auth.MockauthModule
	}
	type args struct {
		query  string
		method string
	}

	tests := []struct {
		name         string
		args         args
		fields       fields
		expectedCode int
		expectedBody string
	}{
		{
			name: "success_login",
			args: args{
				query: "/auth/tg-callback?" +
					"next=%2f&id=123&first_name=Daniil&username=ksusonic&photo_url=pic.jpg&auth_date=123" +
					"&hash=d7590bb4172bdb9029a2f08c976aeb86167037fcf788ba6c6a2aad849c8b3d1b",
				method: http.MethodGet,
			},
			fields: fields{
				MockuserRepo: func(ctrl *gomock.Controller) *auth.MockuserRepo {
					mock := auth.NewMockuserRepo(ctrl)
					mock.EXPECT().GetByTelegramID(gomock.Any(), int64(123)).
						Return(&models.User{ID: 777}, nil)

					return mock
				},
				MockauthModule: func(ctrl *gomock.Controller) *auth.MockauthModule {
					mock := auth.NewMockauthModule(ctrl)
					mock.EXPECT().GenerateJWTToken(777).
						Return(&models.JWTToken{
							Token:   "JWT-TEST-TOKEN",
							Expires: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
						}, nil)

					mock.EXPECT().ValidateTelegramCallbackData(url.Values{
						"id":         []string{"123"},
						"first_name": []string{"Daniil"},
						"username":   []string{"ksusonic"},
						"auth_date":  []string{"123"},
						"photo_url":  []string{"pic.jpg"},
						"hash":       []string{"d7590bb4172bdb9029a2f08c976aeb86167037fcf788ba6c6a2aad849c8b3d1b"},
						"next":       []string{"/"},
					}).Return(true)

					return mock
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"expires":"2009-11-10T23:00:00Z","token":"JWT-TEST-TOKEN"}`,
		},
		{
			name: "success_register",
			args: args{
				query: "/auth/tg-callback?" +
					"next=%2f&id=123&first_name=Daniil&username=ksusonic&photo_url=pic.jpg&auth_date=123" +
					"&hash=d7590bb4172bdb9029a2f08c976aeb86167037fcf788ba6c6a2aad849c8b3d1b",
				method: http.MethodGet,
			},
			fields: fields{
				MockuserRepo: func(ctrl *gomock.Controller) *auth.MockuserRepo {
					mock := auth.NewMockuserRepo(ctrl)
					mock.EXPECT().GetByTelegramID(gomock.Any(), int64(123)).
						Return(nil, models.ErrNotFound)

					pic := "pic.jpg"
					mock.EXPECT().AddTelegramUser(
						gomock.Any(),
						"ksusonic",
						int64(123),
						"Daniil",
						&pic,
					).Return(777, nil)

					return mock
				},
				MockauthModule: func(ctrl *gomock.Controller) *auth.MockauthModule {
					mock := auth.NewMockauthModule(ctrl)
					mock.EXPECT().GenerateJWTToken(777).
						Return(&models.JWTToken{
							Token:   "JWT-TEST-TOKEN",
							Expires: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
						}, nil)

					mock.EXPECT().ValidateTelegramCallbackData(url.Values{
						"id":         []string{"123"},
						"first_name": []string{"Daniil"},
						"username":   []string{"ksusonic"},
						"auth_date":  []string{"123"},
						"photo_url":  []string{"pic.jpg"},
						"hash":       []string{"d7590bb4172bdb9029a2f08c976aeb86167037fcf788ba6c6a2aad849c8b3d1b"},
						"next":       []string{"/"},
					}).Return(true)

					return mock
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"expires":"2009-11-10T23:00:00Z","token":"JWT-TEST-TOKEN"}`,
		},
		{
			name: "invalid_hash",
			args: args{
				query: "/auth/tg-callback?" +
					"next=%2f&id=123&first_name=Daniil&username=ksusonic&photo_url=pic.jpg&auth_date=123" +
					"&hash=keklol",
				method: http.MethodGet,
			},
			fields: fields{
				MockuserRepo: auth.NewMockuserRepo,
				MockauthModule: func(ctrl *gomock.Controller) *auth.MockauthModule {
					mock := auth.NewMockauthModule(ctrl)
					mock.EXPECT().ValidateTelegramCallbackData(url.Values{
						"id":         []string{"123"},
						"first_name": []string{"Daniil"},
						"username":   []string{"ksusonic"},
						"auth_date":  []string{"123"},
						"photo_url":  []string{"pic.jpg"},
						"hash":       []string{"keklol"},
						"next":       []string{"/"},
					}).Return(false)

					return mock
				},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"validation error: telegram query data invalid"}`,
		},
		{
			name: "no query",
			args: args{
				query:  "/auth/tg-callback",
				method: http.MethodGet,
			},
			fields: fields{
				MockuserRepo:   auth.NewMockuserRepo,
				MockauthModule: auth.NewMockauthModule,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"validation error: telegram query data invalid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			mockCtrl := gomock.NewController(t)

			gin.SetMode(gin.TestMode)

			engine := gin.Default()
			engine.GET("/auth/tg-callback", auth.NewController(
				tt.fields.MockuserRepo(mockCtrl),
				tt.fields.MockauthModule(mockCtrl),
				logger.NewDisabled(),
				botCfg,
			).TelegramCallback)

			req, err := http.NewRequest(tt.args.method, tt.args.query, nil)
			require.NoError(t, err)

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
