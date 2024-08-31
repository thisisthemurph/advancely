package application

import "strconv"

type DatabaseConfig struct {
	Name          string
	Password      string
	URI           string
	AutoMigrateOn bool
}

type SupabaseConfig struct {
	URL               string
	PublicKey         string
	ServiceRoleSecret string
}

type ResendConfig struct {
	Key string
}

type AppConfig struct {
	Host          string
	ClientBaseURL string

	Database DatabaseConfig
	Supabase SupabaseConfig
	Resend   ResendConfig
}

func NewAppConfig(get func(string) string) AppConfig {
	var autoMigrateOn bool
	autoMigrateOn, _ = strconv.ParseBool(get("AUTO_MIGRATE_ON"))

	return AppConfig{
		Host:          get("LISTEN_ADDRESS"),
		ClientBaseURL: get("CLIENT_BASE_URL"),

		Database: DatabaseConfig{
			Name:          get("DATABASE_NAME"),
			Password:      get("DATABASE_PASSWORD"),
			URI:           get("DATABASE_URI"),
			AutoMigrateOn: autoMigrateOn,
		},
		Supabase: SupabaseConfig{
			URL:               get("SUPABASE_URL"),
			PublicKey:         get("SUPABASE_PUBLIC_KEY"),
			ServiceRoleSecret: get("SUPABASE_SERVICE_ROLE_SECRET"),
		},
		Resend: ResendConfig{
			Key: get("RESEND_KEY"),
		},
	}
}
