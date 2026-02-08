package repositories

import (
	"database/sql"
	"task-crud-kategori/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetSummary(startDate, endDate string) (*models.ReportSummary, error) {
	summary := &models.ReportSummary{}

	// filter tanggal (optional)
	dateFilter := ""
	args := []interface{}{}

	if startDate != "" && endDate != "" {
		dateFilter = "WHERE DATE(created_at) BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	} else {
		dateFilter = "WHERE DATE(created_at) = DATE('now')"
	}

	// total revenue & transaksi
	err := r.db.QueryRow(`
		SELECT 
			IFNULL(SUM(total_amount), 0),
			COUNT(*)
		FROM transactions
		`+dateFilter,
		args...,
	).Scan(&summary.TotalRevenue, &summary.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// produk terlaris
	err = r.db.QueryRow(`
		SELECT 
			p.name,
			SUM(td.quantity) AS total_qty
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		`+dateFilter+`
		GROUP BY td.product_id
		ORDER BY total_qty DESC
		LIMIT 1
	`, args...).Scan(
		&summary.ProdukTerlaris.Nama,
		&summary.ProdukTerlaris.QtyTerjual,
	)

	// kalau belum ada transaksi
	if err == sql.ErrNoRows {
		summary.ProdukTerlaris = models.BestProduct{}
		return summary, nil
	}

	if err != nil {
		return nil, err
	}

	return summary, nil
}
