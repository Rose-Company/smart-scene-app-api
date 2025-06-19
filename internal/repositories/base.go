package repositories

import (
	"context"
	"log"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model interface {
}

type Clause func(tx *gorm.DB)

type BaseRepository[M Model] interface {
	List(ctx context.Context, params models.QueryParams, clauses ...Clause) ([]*M, error)
	GetByID(ctx context.Context, id interface{}) (*M, error)
	Count(ctx context.Context, params models.QueryParams, clauses ...Clause) (int64, error)
	Create(ctx context.Context, o *M) (*M, error)
	Update(ctx context.Context, id interface{}, o *M, clauses ...Clause) (*M, error)
	UpdateColumns(ctx context.Context, id interface{}, columns map[string]interface{}, clauses ...Clause) (*M, error)
	GetByIDSelected(ctx context.Context, id interface{}, fields []string) (data *M, err error)
	GetIDsByConditions(ctx context.Context, clauses ...Clause) ([]int, error)
	GetDetailByConditions(ctx context.Context, clauses ...Clause) (*M, error)
	Delete(ctx context.Context, clauses ...Clause) error
	CreatesMultiple(ctx context.Context, o []*M) error
	UpdatesByConditions(ctx context.Context, o *M, clauses ...Clause) error
	CountWithGroup(ctx context.Context, params models.QueryParams, groupBy string, clauses ...Clause) (map[string]int64, error)
	UpdatesColumnsByConditions(ctx context.Context, columns map[string]interface{}, clauses ...Clause) error
	ExecRaw(ctx context.Context, raw string) error
}

type baseRepository[M Model] struct {
	model *M
	db    *gorm.DB
}

func NewBaseRepository[M Model](db *gorm.DB) BaseRepository[M] {
	return &baseRepository[M]{
		model: new(M),
		db:    db,
	}
}

func (b *baseRepository[M]) List(ctx context.Context, params models.QueryParams, clauses ...Clause) ([]*M, error) {
	var oList []*M

	// Get table name from model
	var tx *gorm.DB
	if tabler, ok := interface{}(b.model).(interface{ TableName() string }); ok {
		tableName := tabler.TableName()
		log.Printf("DEBUG List: Using table name: %s", tableName)
		tx = b.db.Table(tableName).Offset(params.Offset)
	} else {
		log.Printf("DEBUG List: TableName() method NOT found, using Model approach")
		tx = b.db.Model(b.model).Offset(params.Offset)
	}

	if params.Limit > 0 {
		tx = tx.Limit(params.Limit)
	}
	if len(params.QuerySort.Parse()) > 0 {
		tx = tx.Order(params.QuerySort.Parse())
	}

	if params.Selected != nil {
		tx.Select(params.Selected)
	}
	if params.Preload != nil {
		for _, p := range params.Preload {
			common.ApplyPreload(tx, p)
		}
	}
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Find(&oList).Error
	if err != nil {
		log.Printf("DEBUG List: GORM error: %v", err)
		return nil, err
	}

	log.Printf("DEBUG List: Success - found %d records", len(oList))
	return oList, nil
}

func (b *baseRepository[M]) Count(ctx context.Context, params models.QueryParams, clauses ...Clause) (int64, error) {
	var count int64

	// Get table name from model
	var tableName string
	if tabler, ok := interface{}(b.model).(interface{ TableName() string }); ok {
		tableName = tabler.TableName()
		log.Printf("DEBUG Count: Using table name: %s", tableName)
	} else {
		log.Printf("DEBUG Count: TableName() method NOT found, using Model approach")
		tx := b.db.Model(b.model)
		for _, f := range clauses {
			f(tx)
		}
		err := tx.Count(&count).Error
		if err != nil {
			log.Printf("DEBUG Count: GORM error: %v", err)
			return 0, err
		}
		return count, nil
	}

	// Use Table() instead of Model() to avoid generic type issues
	tx := b.db.Table(tableName)
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Count(&count).Error
	if err != nil {
		log.Printf("DEBUG Count: GORM error with Table(): %v", err)
		return 0, err
	}

	log.Printf("DEBUG Count: Success with Table() - count: %d", count)
	return count, nil
}

func (b *baseRepository[M]) GetByID(ctx context.Context, id interface{}) (*M, error) {
	var o *M
	err := b.db.Model(b.model).First(&o, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (b *baseRepository[M]) Create(ctx context.Context, o *M) (*M, error) {
	err := b.db.Model(b.model).Create(o).Error
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (b *baseRepository[M]) Update(ctx context.Context, id interface{}, o *M, clauses ...Clause) (*M, error) {
	updatedObj := new(M)
	tx := b.db.Model(updatedObj).Clauses(clause.Returning{})
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Where("id = ?", id).Updates(o).Error
	if err != nil {
		return nil, err
	}
	return updatedObj, nil
}

func (b *baseRepository[M]) UpdateColumns(ctx context.Context, id interface{}, columns map[string]interface{}, clauses ...Clause) (*M, error) {
	updatedObj := new(M)
	tx := b.db.Model(updatedObj).Clauses(clause.Returning{})
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Where("id = ?", id).Updates(columns).Error
	if err != nil {
		return nil, err
	}

	return updatedObj, nil
}

func (b *baseRepository[M]) GetByIDSelected(ctx context.Context, id interface{}, fields []string) (data *M, err error) {
	tb := b.db.Model(b.model)
	tb.Select(fields)
	err = tb.First(&data, "id = ? ", id).Error
	return
}

func (b *baseRepository[M]) GetIDsByConditions(ctx context.Context, clauses ...Clause) ([]int, error) {
	var ids []int
	tx := b.db.Model(b.model)
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Select("id").Find(&ids).Error
	return ids, err
}

func (b *baseRepository[M]) GetDetailByConditions(ctx context.Context, clauses ...Clause) (*M, error) {
	var o *M
	tx := b.db.Model(b.model)
	for _, f := range clauses {
		f(tx)
	}
	err := tx.First(&o).Error
	return o, err
}

func (b *baseRepository[M]) Delete(ctx context.Context, clauses ...Clause) error {
	var o *M
	tx := b.db.Model(b.model)
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Delete(&o).Error
	return err
}

func (b *baseRepository[M]) CreatesMultiple(ctx context.Context, o []*M) error {
	return b.db.Model(b.model).Create(o).Error
}

func (b *baseRepository[M]) UpdatesByConditions(ctx context.Context, o *M, clauses ...Clause) error {
	updatedObj := new(M)
	tx := b.db.Model(updatedObj).Clauses(clause.Returning{})
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Updates(o).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *baseRepository[M]) CountWithGroup(ctx context.Context, params models.QueryParams, groupBy string, clauses ...Clause) (map[string]int64, error) {
	var results []struct {
		GroupField string `gorm:"column:group_field"`
		Count      int64  `gorm:"column:count"`
	}

	tx := b.db.Model(b.model)
	for _, f := range clauses {
		f(tx)
	}

	tx = tx.Select(groupBy + " as group_field, COUNT(*) as count").Group(groupBy)

	err := tx.Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.GroupField] = result.Count
	}

	return counts, nil
}

func (b *baseRepository[M]) UpdatesColumnsByConditions(ctx context.Context, columns map[string]interface{}, clauses ...Clause) error {
	updatedObj := new(M)
	tx := b.db.Model(updatedObj).Clauses(clause.Returning{})
	for _, f := range clauses {
		f(tx)
	}
	err := tx.Updates(columns).Error
	return err
}

func (b *baseRepository[M]) ExecRaw(ctx context.Context, raw string) error {
	return b.db.Exec(raw).Error
}
