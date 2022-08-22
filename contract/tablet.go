package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// 유저, 빌리고있는 여부, 빌린 기록
type UserRating struct{
	User string `json:"user"`
	IsLend bool `json:"islend"`
	Dates []Date `json:"dates"`
}

type Date struct{
	Tablet string  `json:"tablet"`
	Record float64 `json:"record"`
}

// 유저 등록
func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, username string) error {

	var user = UserRating{User: username, IsLend: false}
	userAsBytes, _ := json.Marshal(user)	

	return ctx.GetStub().PutState(username, userAsBytes)
}

// 태블릿 대여 기록 추가
func (s *SmartContract) AddDate(ctx contractapi.TransactionContextInterface, username string, tablet string, record string) error {
	
	// getState User 
	userAsBytes, err := ctx.GetStub().GetState(username)	

	if err != nil{
		return err
	} else if userAsBytes == nil{ // no State! error
		return fmt.Errorf("\"Error\":\"User does not exist: "+ username+"\"")
	}
	// state ok
	user := UserRating{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return err
	}
	// create date structure
	newRecord, _ := strconv.ParseFloat(record,64) 
	var Date = Date{Tablet: tablet, Record: newRecord}

	user.Dates=append(user.Dates, Date)

	// 대여 여부 true
	//user.IsLend = true;
	user.IsLend = true


	// update to User World state
	userAsBytes, err = json.Marshal(user);
	if err != nil {
		return fmt.Errorf("failed to Marshaling: %v", err)
	}	

	err = ctx.GetStub().PutState(username, userAsBytes)
	if err != nil {
		return fmt.Errorf("failed to AddRating: %v", err)
	}	
	return nil
}

// 태블릿 반납
func (s *SmartContract) ReturnTablet(ctx contractapi.TransactionContextInterface, username string, tablet string, record string) error {
	
	// getState User 
	userAsBytes, err := ctx.GetStub().GetState(username)	

	if err != nil{
		return err
	} else if userAsBytes == nil{ // no State! error
		return fmt.Errorf("\"Error\":\"User does not exist: "+ username+"\"")
	}
	// state ok
	user := UserRating{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return err
	}
	// create date structure
	newRecord, _ := strconv.ParseFloat(record,64) 
	var Date = Date{Tablet: tablet, Record: newRecord}

	user.Dates=append(user.Dates, Date)

	// 반납
	//user.IsLend = true;
	user.IsLend = false


	// update to User World state
	userAsBytes, err = json.Marshal(user);
	if err != nil {
		return fmt.Errorf("failed to Marshaling: %v", err)
	}	

	err = ctx.GetStub().PutState(username, userAsBytes)
	if err != nil {
		return fmt.Errorf("failed to returning: %v", err)
	}	
	return nil
}


// 태블릿 대여 기록보기
func (s *SmartContract) ReadDate(ctx contractapi.TransactionContextInterface, username string) (string, error) {

	UserAsBytes, err := ctx.GetStub().GetState(username)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if UserAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", username)
	}

	// user := new(UserRating)
	// _ = json.Unmarshal(UserAsBytes, &user)
	
	return string(UserAsBytes[:]), nil	
}



func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create teamate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting teamate chaincode: %s", err.Error())
	}
}
