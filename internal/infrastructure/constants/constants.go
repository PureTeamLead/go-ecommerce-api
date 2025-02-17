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
	MigrationsPath = "postgres-migs"
	HashCost       = 12
	ExpTime        = time.Minute
)
