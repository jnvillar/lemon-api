package userhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userconverter "lemonapp/domain/user/converter"
	usermodel "lemonapp/domain/user/model"
	usermock "lemonapp/domain/user/service/mocks"
	uservalidator "lemonapp/domain/user/validator"
	walletmodel "lemonapp/domain/wallet/model"
	walletmock "lemonapp/domain/wallet/service/mocks"
	customerrors "lemonapp/errors"
	"lemonapp/mocks"

	"github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	user := usermodel.NewUser("juan", "noli", "juaninv", "juan@gmail.com")
	userJSON, _ := json.Marshal(user)
	cases := []struct {
		responseMock      *mocks.ResponseWriter
		name              string
		userServiceMock   *usermock.Service
		walletServiceMock *walletmock.Service
		expectedResponse  interface{}
		request           *http.Request
	}{
		{

			name:             "it fails as expected when users service fails",
			expectedResponse: &customerrors.Error{Error: "internal server error"},
			request:          httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJSON))),
			userServiceMock: func() *usermock.Service {
				u := &usermock.Service{}
				u.
					On("Create", mock.Anything, mock.AnythingOfType("*usermodel.User")).
					Return(nil, fmt.Errorf("error creating user"))
				return u
			}(),
			responseMock: func() *mocks.ResponseWriter {
				r := &mocks.ResponseWriter{}
				r.
					On("WriteHeader", mock.Anything)
				r.
					On("Header").
					Return(http.Header{})
				r.
					On("Write", mock.Anything).
					Return(0, nil)
				return r
			}(),
		},
		{

			name:             "it fails as expected when wallet service fails",
			expectedResponse: &customerrors.Error{Error: "internal server error"},
			request:          httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJSON))),
			userServiceMock: func() *usermock.Service {
				u := &usermock.Service{}
				u.
					On("Create", mock.Anything, mock.AnythingOfType("*usermodel.User")).
					Return(user, nil)
				return u
			}(),
			walletServiceMock: func() *walletmock.Service {
				u := &walletmock.Service{}
				u.
					On("CreateWallets", mock.Anything, user.ID).
					Return(nil, fmt.Errorf("error creating wallets"))
				return u
			}(),
			responseMock: func() *mocks.ResponseWriter {
				r := &mocks.ResponseWriter{}
				r.
					On("Header").
					Return(http.Header{})
				r.
					On("WriteHeader", mock.Anything)
				r.
					On("Write", mock.Anything).
					Return(0, nil)
				return r
			}(),
		},
		{

			name:    "ok request",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJSON))),
			userServiceMock: func() *usermock.Service {
				u := &usermock.Service{}
				u.
					On("Create", mock.Anything, mock.AnythingOfType("*usermodel.User")).
					Return(user, nil)
				return u
			}(),
			walletServiceMock: func() *walletmock.Service {
				u := &walletmock.Service{}
				u.
					On("CreateWallets", mock.Anything, user.ID).
					Return([]walletmodel.Wallet{}, nil)
				return u
			}(),
			responseMock: func() *mocks.ResponseWriter {
				r := &mocks.ResponseWriter{}
				r.
					On("Header").
					Return(http.Header{})
				r.
					On("WriteHeader", mock.Anything)
				r.
					On("Write", mock.Anything).
					Return(0, nil)
				return r
			}(),
			expectedResponse: user,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			convey.Convey("Given a handler", t, func() {
				handler := NewUsersHandler(c.userServiceMock, userconverter.NewUsersConverter(), uservalidator.NewUserValidator(), c.walletServiceMock)
				convey.Convey("When we create an user", func() {
					handler.createUser(&gin.Context{
						Writer:  c.responseMock,
						Request: c.request,
					})
					if c.expectedResponse != nil {
						response, _ := json.Marshal(c.expectedResponse)
						r := string(response) + "\n"
						c.responseMock.AssertCalled(t, "Write", []byte(r))
					}
				})
			})
		})
	}
}
