package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUserTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "createusertable",
		Migrate: func(d *gorm.DB) error {
			return d.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					username VARCHAR(255) NOT NULL,
					password VARCHAR(255) NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					deleted_at TIMESTAMPTZ
				)
			`).Error
		},
		Rollback: func(d *gorm.DB) error {
			return d.Exec("DROP TABLE IF EXISTS users").Error
		},
	}
}

func AddLevelToUser() *gormigrate.Migration {
    return &gormigrate.Migration{
        ID: "20250827_add_level_to_user",
        Migrate: func(tx *gorm.DB) error {
            return tx.Exec("ALTER TABLE users ADD COLUMN level INT NOT NULL DEFAULT 1").Error
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Exec("ALTER TABLE users DROP COLUMN level").Error
        },
    }
}

func AddBalanceToUser() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250827_add_balance_to_user",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE users ADD COLUMN balance INT NOT NULL DEFAULT 0").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE users DROP COLUMN balance").Error
		},
	}
}
