package sqliterepo

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	functionservice "github.com/nilspolek/AstralFS/function-service"
	"github.com/nilspolek/AstralFS/repo"
	"github.com/nilspolek/goLog"
)

type sqliteRepo struct {
	db *sql.DB
}

func New(path string) (repo.Repo, error) {
	db, err := sql.Open("sqlite3", path)
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS functions (
			id VARCHAR(16) PRIMARY KEY,
			image VARCHAR(255),
			route VARCHAR(255),
			port INTEGER
		);
		`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return &sqliteRepo{
		db: db,
	}, err
}

func (repo *sqliteRepo) InsertFunction(fn functionservice.Function) error {
	goLog.Info("REPO: Inset")
	insertSql := "INSERT INTO functions (image, port, route, id) VALUES (?,?,?,?)"
	_, err := repo.db.Exec(insertSql, fn.Image, fn.Port, fn.Route, fn.Id)
	return err
}

func (repo *sqliteRepo) DeleteFunction(id uuid.UUID) error {
	goLog.Info("REPO: Delete")
	deleteSql := "DELETE FROM functions WHERE id = ?"
	_, err := repo.db.Exec(deleteSql, id.String())
	return err
}

func (repo *sqliteRepo) DeleteAllFunctions() error {
	goLog.Info("REPO: DeleteAll")
	deleteSql := "DELETE FROM functions"
	_, err := repo.db.Exec(deleteSql)
	return err
}

func (repo *sqliteRepo) GetFunctions() ([]functionservice.Function, error) {
	goLog.Info("REPO: GetAll")
	var (
		err error
		out []functionservice.Function
	)
	rows, err := repo.db.Query("SELECT image, port, route, id FROM functions")
	for rows.Next() {
		var (
			image string
			port  int
			route string
			id    string
		)
		err := rows.Scan(&image, &port, &route, &id)
		if err != nil {
			return nil, err
		}
		out = append(out, functionservice.Function{
			Image: image,
			Port:  port,
			Route: route,
			Id:    uuid.MustParse(id),
		})
	}
	return out, err
}
