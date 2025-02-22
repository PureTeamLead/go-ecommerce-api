package constants

import (
	"github.com/google/uuid"
	"time"
)

var (
	EmptyID = uuid.Nil
)

const (
	PostgresDriver = "postgres"
	CookieJWT      = "JWT token"
	MigrationsPath = "postgres-migs"
	HashCost       = 12
	ExpTime        = time.Minute
)
