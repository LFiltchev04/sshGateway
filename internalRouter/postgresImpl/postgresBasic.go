package postgresimpl

import(
	"log/slog"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	rtr "sshGateway/internalRouter/addressResolve"

)


type PostgresImpl struct{
	db *sql.DB
}

func (p *PostgresImpl) DbOpen(connStr string) {

	var err error
	p.db, err = sql.Open("pgx", connStr)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		panic(err)
	}

	
}

func (p *PostgresImpl) isPresent(key string) bool {
	sql.Open("pgx", )
	return false
}

func (p *PostgresImpl) getRef(key string) string {
	return ""
}