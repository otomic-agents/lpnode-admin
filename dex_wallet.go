package adminapiservice

import (
	dexwallet "admin-panel/gen/dex_wallet"
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// dexWallet service example implementation.
// The example methods log the requests and return zero values.
type dexWalletsrvc struct {
	logger *log.Logger
}

// NewDexWallet returns the dexWallet service implementation.
func NewDexWallet(logger *log.Logger) dexwallet.Service {
	return &dexWalletsrvc{logger}
}

// ListDexWallet implements listDexWallet.
func (s *dexWalletsrvc) ListDexWallet(ctx context.Context) (res *dexwallet.ListDexWalletResult, err error) {
	res = &dexwallet.ListDexWalletResult{}
	dws := service.NewDexWalletLogicService()
	ret, findErr := dws.ListAll(bson.M{})

	if findErr != nil {
		err = findErr
		return
	}
	logger.System.Debug(ret)
	res.Result = make([]*dexwallet.WalletRow, 0)
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	for _, v := range ret {
		res.Result = append(res.Result, &dexwallet.WalletRow{
			ID:        ptr.String(v.ID.Hex()),
			ChainID:   v.ChainId,
			ChainType: v.ChainType,
			Address:   ptr.String(v.Address),
			AccountID: ptr.String(v.AccountId),
			// PrivateKey: v.PrivateKey,
			SignServiceEndpoint: ptr.String(v.SignServiceEndpoint),
			WalletType:      v.WalletType,
			WalletName:      v.WalletName,
			VaultHostType:   ptr.String(v.VaultHostType),
			VaultName:       ptr.String(v.VaultName),
			VaultSecertType: ptr.String(v.VaultSecertType),
		})
	}

	logger.System.Info(len(ret))
	s.logger.Print("dexWallet.listDexWallet")
	return
}

// CreateDexWallet implements createDexWallet.
func (s *dexWalletsrvc) CreateDexWallet(ctx context.Context, p *dexwallet.WalletRow) (res *dexwallet.CreateDexWalletResult, err error) {
	fmt.Println(p.WalletName, "_______", ptr.ToString(&p.WalletName))
	res = &dexwallet.CreateDexWalletResult{Result: &struct{ ID *string }{}}

	address := ""
	vaultHostType := ""
	vaultSecertType := ""
	vaultName:=""
	storeId := ptr.ToString(p.StoreID)
	dwls := service.NewDexWalletLogicService()
	if p.WalletType == "storeId" {
		if storeId == "" {
			err = errors.WithMessage(utils.GetNoEmptyError(err), "storeIdType is required")
			return
		}
		vault, getVaultErr := dwls.GetVault(storeId)
		if getVaultErr != nil {
			err = errors.WithMessage(getVaultErr, "get vault failed")
			return
		}
		if vault == nil {
			err = errors.WithMessage(utils.GetNoEmptyError(err), "vault not found")
			return
		}
		address = vault.Address
		vaultHostType = vault.HostType
		vaultName = vault.Name
		vaultSecertType = vault.SecertType
	}

	if p.WalletType == "privateKey" {
		if ptr.ToString(p.Address) == "" {
			err = errors.WithMessage(utils.GetNoEmptyError(err), "address cannot be empty")
			return
		}
		vaultName = ""
		address = ptr.ToString(p.Address)
	}
	if p.WalletType == "privateKey" {
		p.WalletType = "secretVault"
	}

	if address == "" {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "address cannot be empty")
		return
	}
	filter := bson.M{
		"chainId": p.ChainID,
		"$or": []bson.M{
			{"addressLower": strings.ToLower(address)},
			{"walletName": p.WalletName},
		},
	}
	baseDatasrvc := service.NewBaseDataLogicService()
	chainRow, err := baseDatasrvc.GetChainRowById(p.ChainID)
	if err != nil {
		return
	}
	if chainRow.ID.Hex() == types.MongoEmptyIdHex {
		err = errors.New("no chain")
		return
	}

	ret, err := dwls.FindOneByFilter(filter)
	if err != nil {
		return
	}
	fmt.Println("ret.ID.Hex()", ret.ID.Hex())
	if ret.ID.Hex() != types.MongoEmptyIdHex {
		err = errors.New("wallet is already exist")
		return
	}

	createData := &types.DBWalletRow{
		ID:                  primitive.NewObjectID(),
		WalletName:          p.WalletName,
		PrivateKey:          ptr.ToString(p.PrivateKey),
		SignServiceEndpoint: ptr.ToString(p.SignServiceEndpoint),
		Address:             address,
		ChainType:           p.ChainType,
		ChainId:             p.ChainID,
		AccountId:           ptr.ToString(p.AccountID),
		AddressLower:        strings.ToLower(address),
		StoreId:             ptr.ToString(p.StoreID),
		WalletType:          p.WalletType,
		VaultHostType:       vaultHostType,
		VaultName:           vaultName,
		VaultSecertType:     vaultSecertType,
	}
	err = dwls.CreateByBsonMap(createData)
	if err != nil {
		err = errors.WithMessage(err, "create wallet failed")
		return
	}
	res.Code = ptr.Int64(0)
	res.Result.ID = ptr.String(createData.ID.Hex())
	s.logger.Print("dexWallet.createDexWallet")
	return
}

// DeleteDexWallet implements deleteDexWallet.
func (s *dexWalletsrvc) DeleteDexWallet(ctx context.Context, p *dexwallet.DeleteFilter) (res *dexwallet.DeleteDexWalletResult, err error) {
	res = &dexwallet.DeleteDexWalletResult{}
	dwls := service.NewDexWalletLogicService()
	objectId, idErr := primitive.ObjectIDFromHex(p.ID)
	if idErr != nil {
		err = errors.WithMessage(err, "id format incorrect, unable to convert to mongoid")
		return
	}
	v := struct {
		Id         primitive.ObjectID `bson:"_id"`
		BridgeName string             `bson:"bridgeName"`
		AmmName    string             `bson:"ammName"`
	}{}
	err = database.FindOne("main", "bridges", bson.M{
		"wallet_id": objectId,
	}, &v)
	if v.Id.Hex() != types.MongoEmptyIdHex {
		err = errors.WithMessage(utils.GetNoEmptyError(err), fmt.Sprintf("bridge already using this token, amm:%s, bridge:%s", v.AmmName, v.BridgeName))
		return
	}
	delCount, delErr := dwls.DeleteById(p.ID)
	if delErr != nil {
		err = delErr
		return
	}
	res.Result = ptr.Int64(delCount)
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	s.logger.Print("dexWallet.deleteDexWallet", delCount)
	return
}
func (s *dexWalletsrvc) VaultList(cxt context.Context) (res *dexwallet.VaultListResult, err error) {
	res = &dexwallet.VaultListResult{Result: make([]*dexwallet.VaultRow, 0)}
	logger.System.Debug("listVault")
	// request accessToken
	dwls := service.NewDexWalletLogicService()
	accessToken, err := dwls.GetVaultAccessToken()
	if err != nil {
		err = errors.WithMessage(err, "get accesstoken error")
		return
	}
	vaultList, err := dwls.GetVaultList(accessToken)
	if err != nil {
		err = errors.WithMessage(err, "get vaultlist occur error")
		return
	}

	for _, item := range vaultList {
		res.Result = append(res.Result, &dexwallet.VaultRow{
			Address:    ptr.String(item.Address),
			HostType:   ptr.String(item.HostType),
			ID:         ptr.String(item.Id),
			Name:       ptr.String(item.Name),
			SecertType: ptr.String(item.SecertType),
		})
	}

	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	return
}
func (s *dexWalletsrvc) UpdateLpWallet(cxt context.Context) (res *dexwallet.UpdateLpWalletResult, err error) {
	res = &dexwallet.UpdateLpWalletResult{}

	dwls := service.NewDexWalletLogicService()
	update, err := dwls.RefreshLpWallet()
	if err != nil {
		return
	}
	if !update {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "update error:")
		return
	}
	res.Code = ptr.Int64(0)
	res.Result = ""
	res.Message = ptr.String("")
	logger.System.Debug("updateLpWallet")
	return
}
