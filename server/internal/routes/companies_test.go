package routes_test

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"testing"

	"advancely/internal/routes"
	"advancely/internal/store"
	"advancely/internal/tests"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func newTestCompaniesHandler(db *sqlx.DB) routes.CompaniesHandler {
	return routes.CompaniesHandler{
		CompanySettingsStore: store.NewPostgresCompanySettingsStore(db),
		Logger:               tests.NewDefaultLogger(),
		EnsurePermission:     routes.EnsurePermissionsFnFactory(store.NewPostgresPermissionsStore(db)),
	}
}

func TestAdminUserCanCreateNewAllowedDomain(t *testing.T) {
	testCases := []struct {
		name               string
		domain             string
		expectedStatusCode int
	}{
		{
			name:               "Google",
			domain:             "google.com",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Invalid domain",
			domain:             "someemail@domain.com",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Unknown domain",
			domain:             fmt.Sprintf("%s.com", uuid.NewString()),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	db, user, companyId := setUpTestAdminUserAndCompany(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := routes.AddAllowedDomainRequest{
				Domain:              tc.domain,
				AllowUnknownDomains: false,
			}

			c, rec := tests.NewRequestRecorder(t, http.MethodPost, "/admin/companies/add-domain", model)
			tests.SaveSessionInContext(c, user.ID, companyId)

			// Act
			handler := newTestCompaniesHandler(db)
			err := handler.HandleAddAllowedDomain()(c)
			if err != nil {
				c.Error(err)
			}

			// Assert
			require.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedStatusCode == http.StatusCreated {
				var domains []string
				_ = db.Select(&domains,
					"SELECT domain FROM allowed_email_domains where company_id = $1;", companyId)
				require.Len(t, domains, 1)
				require.Equal(t, tc.domain, domains[0])
			}
		})
	}
}
