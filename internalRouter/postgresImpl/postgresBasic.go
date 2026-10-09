package postgresimpl

import(
	"log/slog"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	//rtr "sshGateway/addressResolve"

)

type envOwnership struct{
	userId string
	envName string
	targetAddr string
}

type PostgresImpl struct{
	meta envOwnership
	db *sql.DB
}

//one-off function to start postgres
func (p *PostgresImpl) DbOpen(connStr string) {

	var err error
	p.db, err = sql.Open("pgx", connStr)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		panic(err)
	}
}

//verifies presence of entry in the DB
func (p *PostgresImpl) isPresent(key string, envName string) bool {
	row := p.db.QueryRow("SELECT 1 FROM envOwnership w WHERE (w.userId = $1 AND w.envName = $2)", key, envName)
	var exists int
	
	err := row.Scan(&exists)
	
	if err != nil {
		slog.Error("DB query bad, error ", err)
		if err == sql.ErrNoRows {
			slog.Debug("No entry found in database", "userId", key, "envName", envName)
			return false
		}
		return false
	}

	slog.Debug("Entry exists in database", "userId", key, "envName", envName)
	return true
}


func (p *PostgresImpl) getRef(key string, envName string) string {
	return ""
}