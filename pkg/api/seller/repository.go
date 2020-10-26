package seller

import (
	"database/sql"
	"fmt"
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) FindByUUID(uuid string) (*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller WHERE uuid = ?", uuid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	seller := &Seller{}

	err = rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)

	if err != nil {
		return nil, err
	}

	return seller, nil
}

func (r *Repository) list() ([]*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*Seller

	for rows.Next() {
		seller := &Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func (r *Repository) ListTop(top int) ([]*Seller, error) {
	rows, err := r.db.Query(GetTopSellers(top))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*Seller

	for rows.Next() {
		seller := &Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func GetTopSellers(top int)string{
	return fmt.Sprintf( `SELECT 
				id_seller, name, email, phone, uuid
			FROM
				seller
					INNER JOIN
				(SELECT 
					fk_seller, COUNT(*)
				FROM
					product
				GROUP BY fk_seller
				ORDER BY COUNT(*) DESC
				LIMIT %d) AS products ON seller.id_seller = products.fk_seller;`,top)
}