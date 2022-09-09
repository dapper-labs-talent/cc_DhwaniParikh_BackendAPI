package repositories

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"server/internal/entities"
)

type IAppUserRepository interface {
	Create(ctx context.Context, user *entities.AppUserEntity) (*entities.AppUserEntity, error)
	GetByEmail(ctx context.Context, email string) (*entities.AppUserEntity, error)
	GetById(ctx context.Context, uuid uuid.UUID) (*entities.AppUserEntity, error)
	GetAllUsers(ctx context.Context) ([]*entities.UserDetailsEntity, error)
	Update(ctx context.Context, user *entities.AppUserEntity) error
}

type AppUserRepo struct {
	db     *pg.DB
	logger *zerolog.Logger
}

func NewAppUserRepo(db *pg.DB, l *zerolog.Logger) *AppUserRepo {
	return &AppUserRepo{
		db:     db,
		logger: l,
	}
}

func (r *AppUserRepo) Create(ctx context.Context, user *entities.AppUserEntity) (*entities.AppUserEntity, error) {
	result, err := r.db.WithContext(ctx).Model(user).Returning("*").Insert()
	fmt.Println(result)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *AppUserRepo) GetById(ctx context.Context, uuid uuid.UUID) (*entities.AppUserEntity, error) {
	user := new(entities.AppUserEntity)
	err := r.db.WithContext(ctx).Model(user).Where("t.id = ?", uuid).Select()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AppUserRepo) GetByEmail(ctx context.Context, email string) (*entities.AppUserEntity, error) {
	// TODO: Do we need to do anything with the returned result?
	// TODO: Don't return everything, just return the id
	user := new(entities.AppUserEntity)
	err := r.db.WithContext(ctx).Model(user).Where("email = ?", email).Select()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AppUserRepo) GetAllUsers(ctx context.Context) ([]*entities.UserDetailsEntity, error) {
	var users []*entities.UserDetailsEntity
	err := r.db.WithContext(ctx).Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *AppUserRepo) Update(ctx context.Context, user *entities.AppUserEntity) error {
	// TODO: Do we need to do anything with the returned result?
	_, err := r.db.WithContext(ctx).Model(user).Returning("*").Where("id = ?", user.ID).Update()
	if err != nil {
		return err
	}
	return nil
}
