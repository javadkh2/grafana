package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/ldap"
	"github.com/grafana/grafana/pkg/services/multildap"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type LDAPMock struct {
	mock.Mock
	Results []*models.ExternalUserInfo
}

type LDAPConfigMock struct {
	mock.Mock
}

var userSearchResult *models.ExternalUserInfo
var userSearchConfig *ldap.ServerConfig

func (c *LDAPConfigMock) Servers() []*ldap.ServerConfig {
	return nil
}

func (m *LDAPMock) Login(query *models.LoginUserQuery) (*models.ExternalUserInfo, error) {
	return &models.ExternalUserInfo{}, nil
}

func (m *LDAPMock) Users(logins []string) ([]*models.ExternalUserInfo, error) {
	s := []*models.ExternalUserInfo{}
	return s, nil
}

func (m *LDAPMock) User(login string) (*models.ExternalUserInfo, *ldap.ServerConfig, error) {
	return userSearchResult, userSearchConfig, nil
}

func getUserFromLDAPContext(t *testing.T, requestURL string) *scenarioContext {
	t.Helper()

	sc := setupScenarioContext(requestURL)

	hs := &HTTPServer{Cfg: setting.NewCfg()}

	sc.defaultHandler = Wrap(func(c *models.ReqContext) Response {
		sc.context = c
		return hs.GetUserFromLDAP(c)
	})

	sc.m.Get("/api/admin/ldap/:username", sc.defaultHandler)

	sc.resp = httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	sc.req = req
	sc.exec()

	return sc
}

func TestGetUserFromLDAPApiEndpoint_UserNotFound(t *testing.T) {
	getLDAPConfig = func() (*ldap.Config, error) {
		return &ldap.Config{}, nil
	}

	newLDAP = func(_ []*ldap.ServerConfig) multildap.IMultiLDAP {
		return new(LDAPMock)
	}

	userSearchResult = nil

	sc := getUserFromLDAPContext(t, "/api/admin/ldap/user-that-does-not-exist")

	require.Equal(t, sc.resp.Code, http.StatusNotFound)
	responseString, err := getBody(sc.resp)

	assert.Nil(t, err)
	assert.Equal(t, "{\"message\":\"No user was found on the LDAP server(s)\"}", responseString)
}

func TestGetUserFromLDAPApiEndpoint(t *testing.T) {
	isAdmin := true
	userSearchResult = &models.ExternalUserInfo{
		Name:           "John Doe",
		Email:          "john.doe@example.com",
		Login:          "johndoe",
		OrgRoles:       map[int64]models.RoleType{1: models.ROLE_ADMIN},
		IsGrafanaAdmin: &isAdmin,
	}

	userSearchConfig = &ldap.ServerConfig{
		Attr: ldap.AttributeMap{
			Name:     "ldap-name",
			Surname:  "ldap-surname",
			Email:    "ldap-email",
			Username: "ldap-username",
		},
		Groups: []*ldap.GroupToOrgRole{
			{
				GroupDN: "cn=admins,ou=groups,dc=grafana,dc=org",
				OrgID:   1,
				OrgRole: models.ROLE_ADMIN,
			},
		},
	}

	getLDAPConfig = func() (*ldap.Config, error) {
		return &ldap.Config{}, nil
	}

	newLDAP = func(_ []*ldap.ServerConfig) multildap.IMultiLDAP {
		return new(LDAPMock)
	}

	sc := getUserFromLDAPContext(t, "/api/admin/ldap/johndoe")

	require.Equal(t, sc.resp.Code, http.StatusOK)

	jsonResponse, err := getJSONbody(sc.resp)
	assert.Nil(t, err)

	expected := `
		{
		  "name": {
				"cfgAttrValue": "ldap-name", "ldapValue": "John Doe"
			},
			"surname": {
				"cfgAttrValue": "ldap-surname", "ldapValue": "John Doe"
			},
			"email": {
				"cfgAttrValue": "ldap-email", "ldapValue": "john.doe@example.com"
			},
			"login": {
				"cfgAttrValue": "ldap-username", "ldapValue": "johndoe"
			},
			"isGrafanaAdmin": true,
			"isDisabled": false,
			"roles": [
				{ "orgId": 1, "orgRole": "Admin", "groupDN": "cn=admins,ou=groups,dc=grafana,dc=org" }
			],
			"teams": []
		}
	`
	var expectedJSON interface{}
	_ = json.Unmarshal([]byte(expected), &expectedJSON)

	assert.Equal(t, jsonResponse, expectedJSON)
}
