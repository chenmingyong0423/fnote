package repository

import (
	"context"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/repository/dao"
	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeAssetDAO struct {
	assets []*dao.Asset
}

func (f *fakeAssetDAO) FindById(context.Context, bson.ObjectID) (*dao.Asset, error) {
	if len(f.assets) == 0 {
		return nil, nil
	}
	return f.assets[0], nil
}

func (f *fakeAssetDAO) FindByIds(context.Context, []bson.ObjectID) ([]*dao.Asset, error) {
	return f.assets, nil
}

func (f *fakeAssetDAO) FindByFileID(context.Context, string) (*dao.Asset, error) {
	if len(f.assets) == 0 {
		return nil, nil
	}
	return f.assets[0], nil
}

func (f *fakeAssetDAO) Add(context.Context, *dao.Asset) (bson.ObjectID, error) {
	return bson.NilObjectID, nil
}

func (f *fakeAssetDAO) DeleteById(context.Context, bson.ObjectID) (int64, error) {
	return 0, nil
}

func TestFindByIdsPreservesFolderOrder(t *testing.T) {
	firstID := bson.NewObjectID()
	secondID := bson.NewObjectID()
	repository := NewAssetRepository(&fakeAssetDAO{assets: []*dao.Asset{
		{Model: mongox.Model{ID: secondID}, Content: "second"},
		{Model: mongox.Model{ID: firstID}, Content: "first"},
	}})

	assets, err := repository.FindByIds(context.Background(), []string{firstID.Hex(), secondID.Hex()})
	if err != nil {
		t.Fatal(err)
	}
	if len(assets) != 2 || assets[0].Id != firstID.Hex() || assets[1].Id != secondID.Hex() {
		t.Fatalf("unexpected asset order: %#v", assets)
	}
}
