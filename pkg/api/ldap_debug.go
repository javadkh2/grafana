package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/ldap"
	"github.com/grafana/grafana/pkg/services/multildap"
)

var (
	getLDAPConfig = multildap.GetConfig
	newLDAP       = multildap.New
)

// LDAPAttribute is a serializer for user attributes mapped from LDAP. Is meant to display both the serialized value and the LDAP key we received it from.
type LDAPAttribute struct {
	ConfigAttributeValue string `json:"cfgAttrValue"`
	LDAPAttributeValue   string `json:"ldapValue"`
}

// RoleDTO is a serializer for mapped roles from LDAP
type RoleDTO struct {
	OrgId   int64           `json:"orgId"`
	OrgRole models.RoleType `json:"orgRole"`
	GroupDN string          `json:"groupDN"`
}

// TeamDTO is a serializer for mapped Teams from LDAP
type TeamDTO struct {
	GroupDN string `json:"groupDN"`
	TeamId  int64  `json:"teamId"`
}

// LDAPUserDTO is a serializer for users mapped from LDAP
type LDAPUserDTO struct {
	Name           *LDAPAttribute `json:"name"`
	Surname        *LDAPAttribute `json:"surname"`
	Email          *LDAPAttribute `json:"email"`
	Username       *LDAPAttribute `json:"login"`
	IsGrafanaAdmin *bool          `json:"isGrafanaAdmin"`
	IsDisabled     bool           `json:"isDisabled"`
	OrgRoles       *[]RoleDTO     `json:"roles"`
	Teams          *[]TeamDTO     `json:"teams"`
}

// ReloadLDAPCfg reloads the LDAP configuration
func (server *HTTPServer) ReloadLDAPCfg() Response {
	if !ldap.IsEnabled() {
		return Error(400, "LDAP is not enabled", nil)
	}

	err := ldap.ReloadConfig()
	if err != nil {
		return Error(500, "Failed to reload ldap config.", err)
	}
	return Success("LDAP config reloaded")
}

// GetUserFromLDAP finds an user based on a username in LDAP. This helps illustrate how would the particular user be mapped in Grafana when synced.
func (server *HTTPServer) GetUserFromLDAP(c *models.ReqContext) Response {
	ldapConfig, err := getLDAPConfig()

	if err != nil {
		return Error(400, "Failed to obtain the LDAP configuration. Please ", err)
	}

	ldap := newLDAP(ldapConfig.Servers)

	username := c.Params(":username")

	if len(username) == 0 {
		return Error(http.StatusBadRequest, "Validation error. You must specify an username", nil)
	}

	user, serverConfig, err := ldap.User(username)

	if user == nil {
		return Error(
			404, "No user was found on the LDAP server(s)", err)
	}

	u := &LDAPUserDTO{
		Name:           &LDAPAttribute{serverConfig.Attr.Name, user.Name},
		Surname:        &LDAPAttribute{serverConfig.Attr.Surname, user.Name},
		Email:          &LDAPAttribute{serverConfig.Attr.Email, user.Email},
		Username:       &LDAPAttribute{serverConfig.Attr.Username, user.Login},
		IsGrafanaAdmin: user.IsGrafanaAdmin,
		IsDisabled:     user.IsDisabled,
	}

	orgRoles := []RoleDTO{}

	for _, g := range serverConfig.Groups {
		role := &RoleDTO{}

		if user.OrgRoles[g.OrgID] == g.OrgRole {
			role.OrgId = g.OrgID
			role.OrgRole = user.OrgRoles[g.OrgID]
			role.GroupDN = g.GroupDN

			orgRoles = append(orgRoles, *role)
		} else {
			role.GroupDN = g.GroupDN

			orgRoles = append(orgRoles, *role)
		}
	}

	u.OrgRoles = &orgRoles

	u.Teams = &[]TeamDTO{}

	return JSON(200, u)
}
