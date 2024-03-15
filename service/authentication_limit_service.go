package service

import (
	database "admin-panel/mongo_database"
	types "admin-panel/types"

	"go.mongodb.org/mongo-driver/bson"
)

type AuthenticationLimiterService struct{}

func (als *AuthenticationLimiterService) Set(value string) {
	valueSet := bson.M{"data": value}
	database.FindOneAndUpdate("main", "authenticationLimiters", bson.M{}, bson.M{"$set": valueSet})
}
func (als *AuthenticationLimiterService) Del() (count int64, err error) {
	count, err = database.DeleteOne("main", "authenticationLimiters", bson.M{})
	return
}
func (als *AuthenticationLimiterService) Get() (row types.AuthenticationLimiterRow,err error){
	row = types.AuthenticationLimiterRow{}
	err = database.FindOne("main", "authenticationLimiters", bson.M{}, &row)
	if err!=nil{
		return
	}
	return
}
func NewAuthenticationLimiterService() *AuthenticationLimiterService {
	return &AuthenticationLimiterService{}
}
