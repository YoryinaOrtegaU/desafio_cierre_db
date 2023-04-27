package invoices

import (
	"database/sql"

	"github.com/bootcamp-go/desafio-cierre-db.git/internal/domain"
)

type Repository interface {
	Create(invoices *domain.Invoices) (int64, error)
	ReadAll() (*[]domain.Invoices, error)
	CalculeTotal() (*[]domain.Invoices, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(invoices *domain.Invoices) (int64, error) {
	query := `INSERT INTO invoices (id, datetime, customer_id, total) VALUES (?, ?, ?, ?)`
	row, err := r.db.Exec(query, &invoices.Id, &invoices.Datetime, &invoices.CustomerId, &invoices.Total)
	if err != nil {
		return 0, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) ReadAll() (*[]domain.Invoices, error) {
	query := `SELECT id, customer_id, datetime, total FROM invoices`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	invoices := []domain.Invoices{}
	for rows.Next() {
		invoice := domain.Invoices{}
		err := rows.Scan(&invoice.Id, &invoice.CustomerId, &invoice.Datetime, &invoice.Total)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return &invoices, nil
}

func (r *repository) CalculeTotal() (*[]domain.Invoices, error) {
	statement, err := r.db.Prepare(`
	UPDATE invoices a 
	SET total=(SELECT sum(p.price * s.quantity) 
	from sales s 
	inner join products p on s.product_id = p.id  
	WHERE s.invoice_id = a.id)
	`)

	_, err = statement.Exec()
	if err != nil {
		var vacio []domain.Invoices
		return &vacio, err
	}

	return r.ReadAll()
}
