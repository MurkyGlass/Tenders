package repositories

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrInvalidInput = errors.New("invalid input data")
	ErrConflict     = errors.New("data conflict")
	ErrNotAllowed   = errors.New("operation not allowed")
)

type Repository interface {
	Users() UserRepository
	Tenders() TenderRepository
	Company() CompanyRepository
	Offer() OfferRepository
	Doc() DocRepository
	Log() LogRepository
	RoleInCompany() RoleInCompanyRepository
	Category() CategoryRepository
	District() DistrictRepository
	Right() RightRepository
	Role() RoleRepository
	Status() StatusRepository
	Refresh() RefreshRepository
	Reset() ResetRepository
	CategoryLink() CategoryLinkRepository
	BeginTx(ctx context.Context) (Transaction, error)

	LinkerDoc(idDoc int) LinkerDoc
	LinkerLog(idLog int) LinkerLog
	LinkerTCategory(idTender int) LinkerCategory
	LinkerRoleRight(idRole int) LinkerRight
}

type Transaction interface {
	Commit() error
	Rollback() error

	Users() UserRepository
	Tenders() TenderRepository
	Company() CompanyRepository
	Offer() OfferRepository
	Doc() DocRepository
	Log() LogRepository
	RoleInCompany() RoleInCompanyRepository
	Category() CategoryRepository
	District() DistrictRepository
	Right() RightRepository
	Role() RoleRepository
	Status() StatusRepository
	Refresh() RefreshRepository
	Reset() ResetRepository
	CategoryLink() CategoryLinkRepository

	LinkerDoc(idDoc int) LinkerDoc
	LinkerLog(idLog int) LinkerLog
	LinkerTCategory(idTender int) LinkerCategory
	LinkerRoleRight(idRole int) LinkerRight
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Users() UserRepository {
	return NewUserRepository(r.db)
}
func (r *repository) Tenders() TenderRepository {
	return NewTenderRepository(r.db)
}
func (r *repository) Company() CompanyRepository {
	return NewCompanyRepository(r.db)
}
func (r *repository) Offer() OfferRepository {
	return NewOfferRepository(r.db)
}
func (r *repository) Doc() DocRepository {
	return NewDocRepository(r.db)
}
func (r *repository) Log() LogRepository {
	return NewLogRepository(r.db)
}
func (r *repository) RoleInCompany() RoleInCompanyRepository {
	return NewRoleInCompanyRepository(r.db)
}
func (r *repository) Category() CategoryRepository {
	return NewCategoryRepository(r.db)
}
func (r *repository) District() DistrictRepository {
	return NewDistrictRepository(r.db)
}
func (r *repository) Right() RightRepository {
	return NewRightRepository(r.db)
}
func (r *repository) Role() RoleRepository {
	return NewRoleRepository(r.db)
}
func (r *repository) Status() StatusRepository {
	return NewStatusRepository(r.db)
}
func (r *repository) Refresh() RefreshRepository {
	return NewRefreshRepository(r.db)
}
func (r *repository) Reset() ResetRepository {
	return NewResetRepository(r.db)
}
func (r *repository) CategoryLink() CategoryLinkRepository {
	return NewCategoryLinkRepository(r.db)
}
func (r *repository) LinkerDoc(idDoc int) LinkerDoc {
	return NewLinkerDoc(idDoc, r.db, nil)
}
func (r *repository) LinkerLog(idLog int) LinkerLog {
	return NewLinkerLog(idLog, r.db, nil)
}
func (r *repository) LinkerTCategory(idTender int) LinkerCategory {
	return NewLinkerCategory(idTender, r.db)
}
func (r *repository) LinkerRoleRight(idRole int) LinkerRight {
	return NewLinkerRight(idRole, r.db)
}

func (r *repository) BeginTx(ctx context.Context) (Transaction, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

type transaction struct {
	tx *sqlx.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}
func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}
func (t *transaction) Users() UserRepository {
	return NewUserRepository(t.tx)
}
func (t *transaction) Tenders() TenderRepository {
	return NewTenderRepository(t.tx)
}
func (t *transaction) Company() CompanyRepository {
	return NewCompanyRepository(t.tx)
}
func (t *transaction) Offer() OfferRepository {
	return NewOfferRepository(t.tx)
}
func (t *transaction) Doc() DocRepository {
	return NewDocRepository(t.tx)
}
func (t *transaction) Log() LogRepository {
	return NewLogRepository(t.tx)
}
func (t *transaction) RoleInCompany() RoleInCompanyRepository {
	return NewRoleInCompanyRepository(t.tx)
}
func (t *transaction) Category() CategoryRepository {
	return NewCategoryRepository(t.tx)
}
func (t *transaction) District() DistrictRepository {
	return NewDistrictRepository(t.tx)
}
func (t *transaction) Right() RightRepository {
	return NewRightRepository(t.tx)
}
func (t *transaction) Role() RoleRepository {
	return NewRoleRepository(t.tx)
}
func (t *transaction) Status() StatusRepository {
	return NewStatusRepository(t.tx)
}
func (t *transaction) LinkerDoc(idDoc int) LinkerDoc {
	return NewLinkerDoc(idDoc, t.tx, nil)
}
func (t *transaction) LinkerLog(idLog int) LinkerLog {
	return NewLinkerLog(idLog, t.tx, nil)
}
func (t *transaction) LinkerTCategory(idTender int) LinkerCategory {
	return NewLinkerCategory(idTender, t.tx)
}
func (t *transaction) LinkerRoleRight(idRole int) LinkerRight {
	return NewLinkerRight(idRole, t.tx)
}
func (t *transaction) Refresh() RefreshRepository {
	return NewRefreshRepository(t.tx)
}
func (t *transaction) Reset() ResetRepository {
	return NewResetRepository(t.tx)
}
func (t *transaction) CategoryLink() CategoryLinkRepository {
	return NewCategoryLinkRepository(t.tx)
}
