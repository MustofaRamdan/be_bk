package config

import "fmt"

// AddIndexes creates database indexes on frequently queried columns
// to speed up WHERE, ORDER BY, and COUNT operations.
func AddIndexes() {
	indexes := []struct {
		table string
		name  string
		cols  string
	}{
		{"alumni", "idx_alumni_status_pengajuan_status", "status_pengajuan, status"},
		{"alumni", "idx_alumni_created_at", "created_at"},
		{"karyas", "idx_karyas_status", "status"},
		{"karyas", "idx_karyas_createdAt", "createdAt"},
		{"posts", "idx_posts_published", "published"},
		{"posts", "idx_posts_createdAt", "createdAt"},
		{"guru", "idx_guru_createdAt", "createdAt"},
	}

	for _, idx := range indexes {
		sql := fmt.Sprintf(
			"CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
			idx.name, idx.table, idx.cols,
		)
		if err := DB.Exec(sql).Error; err != nil {
			// MySQL < 8.0 doesn't support IF NOT EXISTS for indexes,
			// so we silently ignore duplicate index errors
			fmt.Printf("Index %s: %v (may already exist)\n", idx.name, err)
		}
	}

	fmt.Println("Database indexes checked")
}
