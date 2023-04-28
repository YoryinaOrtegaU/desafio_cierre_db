package customers

import (
	"database/sql"

	"github.com/bootcamp-go/desafio-cierre-db.git/internal/domain"
)

type Repository interface {
	Create(customers *domain.Customers) (int64, error)
	ReadAll() ([]*domain.Customers, error)
	TotalSalesByConditionCustomer() (map[string]float64, error)
	TopCustomers() (map[string]float64, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(customers *domain.Customers) (int64, error) {
	query := "INSERT INTO customers (id, first_name, last_name, `condition`) VALUES (?, ?, ?, ?)"
	row, err := r.db.Exec(query, &customers.Id, &customers.FirstName, &customers.LastName, &customers.Condition)
	if err != nil {
		return 0, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) ReadAll() ([]*domain.Customers, error) {
	query := `SELECT id, first_name, last_name, customers.condition FROM customers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	customers := make([]*domain.Customers, 0)
	for rows.Next() {
		customer := domain.Customers{}
		err := rows.Scan(&customer.Id, &customer.FirstName, &customer.LastName, &customer.Condition)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}
	return customers, nil
}

func (r *repository) TotalSalesByConditionCustomer() (map[string]float64, error) {
	query := "SELECT c.`condition`, ROUND(SUM(i.total), 2)  FROM customers c INNER JOIN invoices i ON c.id = i.customer_id GROUP BY c.`condition`"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	response := make(map[string]float64)
	for rows.Next() {
		var condition string
		var total float64
		err := rows.Scan(&condition, &total)
		if err != nil {
			return nil, err
		}
		response[condition] = total
	}
	return response, nil
}

func (r *repository) TopCustomers() (map[string]float64, error) {
	query := "SELECT c.first_name, c.last_name, ROUND(i.total, 2) as Amount FROM customers c INNER JOIN invoices i ON c.id = i.customer_id where c.`condition` = 1 ORDER by Amount DESC limit 5"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	response := make(map[string]float64)
	for rows.Next() {
		var first_name string
		var last_name string
		var amount float64
		err := rows.Scan(&first_name, &last_name, &amount)
		if err != nil {
			return nil, err
		}
		customer := first_name + last_name
		response[customer] = amount
	}
	return response, nil
}
