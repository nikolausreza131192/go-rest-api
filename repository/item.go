package repository

import (
	"log"
	"sort"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/nikolausreza131192/pos/entity"
)

// ItemRepo define the interface for item repository
type ItemRepo interface {
	GetAll() []entity.Item
	GetByID(id int) entity.Item
}

type itemRepo struct {
	db      *sqlx.DB
	itemMap *sync.Map
}

// ItemRepoParam will be used as repository parameter
type ItemRepoParam struct {
	DB *sqlx.DB
}

// NewItem initialize item repository
func NewItem(param ItemRepoParam) ItemRepo {
	itemMap := &sync.Map{}
	r := &itemRepo{
		db:      param.DB,
		itemMap: itemMap,
	}
	r.buildCache()

	return r
}

// Store all items from DB to sync map
func (m *itemRepo) buildCache() {
	rows, err := m.db.Queryx(`
		SELECT id, kategori_id, merk_id, kode, nama, panjang, lebar, initial_stok, initial_m2, keterangan, created_by, updated_by, created_at, updated_at
		FROM m_barang
	`)

	if err != nil {
		log.Println("GetAll Error query", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		item := entity.Item{}
		err = rows.StructScan(&item)
		if err != nil {
			log.Println("GetAll Error struct scan", err)
			continue
		}

		m.itemMap.Store(item.ID, item)
	}
}

// GetAll is used to fetch all items from sync map
func (m *itemRepo) GetAll() []entity.Item {
	items := []entity.Item{}
	m.itemMap.Range(func(k, v interface{}) bool {
		item := v.(entity.Item)
		items = append(items, item)
		return true
	})

	// Sort by name
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items
}

// Get item by ID
func (m *itemRepo) GetByID(id int) entity.Item {
	if v, ok := m.itemMap.Load(id); ok {
		return v.(entity.Item)
	}

	// If not found, then query from DB
	item := entity.Item{}
	row := m.db.QueryRowx(`
		SELECT id, kategori_id, merk_id, kode, nama, panjang, lebar, initial_stok, initial_m2, keterangan, created_by, updated_by, created_at, updated_at
		FROM m_barang
		WHERE id = ?
	`, id)
	if err := row.StructScan(&item); err != nil {
		log.Println("GetByID Error struct scan", id, err)
		return entity.Item{}
	}

	return item
}
