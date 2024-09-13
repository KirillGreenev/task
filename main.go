package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type City struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	State string `db:"state"`
}

type SQLiteCityRepository struct {
	db *sqlx.DB
}

func NewSQLiteCityRepository(db *sqlx.DB) *SQLiteCityRepository {
	return &SQLiteCityRepository{db: db}
}

func initDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "./test.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS cities (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(30) NOT NULL,
        state VARCHAR(30) NOT NULL
    );`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return db, nil
}

func (repo *SQLiteCityRepository) Create(ctx context.Context, city City) (int, error) {
	result, err := repo.db.ExecContext(ctx, "INSERT INTO cities (name, state) VALUES (?, ?)",
		city.Name, city.State)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (repo *SQLiteCityRepository) Delete(ctx context.Context, id int) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM cities WHERE id = ?", id)
	return err
}

func (repo *SQLiteCityRepository) Update(ctx context.Context, city City) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE cities SET name = ?, state = ? WHERE id = ?",
		city.Name, city.State, city.ID)
	return err
}

func (repo *SQLiteCityRepository) List(ctx context.Context) ([]City, error) {
	var cities []City
	err := repo.db.SelectContext(ctx, &cities, "SELECT * FROM cities")
	return cities, err
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := NewSQLiteCityRepository(db)

	ctx := context.Background()

	city := City{Name: "Москва", State: "РФ"}
	id, err := repo.Create(ctx, city)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Добавлен город с id: %d\n", id)

	cities, err := repo.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Список городов:", cities)

	city.ID = id
	city.Name = "Санкт-Петербург"
	err = repo.Update(ctx, city)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Обновлён город:", city)

	err = repo.Delete(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Удалён город с ID:", id)
}

/* Написать репозиторий к базе данных test на mysql, используя библиотеку sqlx.CREATE TABLE cities (
id INTEGER NOT NULL PRIMARY KEY,
name VARCHAR(30) NOT NULL,
state VARCHAR(30) NOT NULL
);
Реализовать следующие методы:
Create
Delete
Update
List
*/
