package userstorage

import (
	"app-invite-service/module/user/usermodel"
	"context"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNewSQLStore(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	var tests []struct {
		name string
		args args
		want *sqlStore
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLStore(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sqlStore_CreateUser(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		in0  context.Context
		data *usermodel.UserCreate
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sqlStore{
				db: tt.fields.db,
			}
			if err := s.CreateUser(tt.args.in0, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sqlStore_FindUser(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		in0        context.Context
		conditions map[string]interface{}
		moreInfo   []string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *usermodel.User
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sqlStore{
				db: tt.fields.db,
			}
			got, err := s.FindUser(tt.args.in0, tt.args.conditions, tt.args.moreInfo...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
