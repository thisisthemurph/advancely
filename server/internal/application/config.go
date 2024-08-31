package application

type DatabaseConfig struct {
	Password string
	URI      string
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
	Host     string
	Database DatabaseConfig
	Supabase SupabaseConfig
	Resend   ResendConfig
}

func NewAppConfig(get func(string) string) AppConfig {
	return AppConfig{
		Host: get("LISTEN_ADDRESS"),
		Database: DatabaseConfig{
			Password: get("DATABASE_PASSWORD"),
			URI:      get("DATABASE_URI"),
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
