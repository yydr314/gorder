package persistent

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL() *MySQL {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &MySQL{db: db}
}

type StockModel struct {
	ID        int64
	ProductID string
	Quantity  int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (StockModel) TableName() string {
	return "o_stock"
}

func (m MySQL) BatchGetStockByID(ctx context.Context, ids []string) ([]StockModel, error) {
	var result []StockModel
	tx := m.db.WithContext(ctx).Where("product_id IN ?", ids).Find(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return result, nil
}
