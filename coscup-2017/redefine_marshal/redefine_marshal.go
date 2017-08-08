package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

type UnixTime struct {
	//anonymous field
	time.Time `bson:",inline" json:",inline"`
}

func (unixTime UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(unixTime.Time.Unix())
}

func (unixTime *UnixTime) UnmarshalJSON(data []byte) error {
	var unixTimeInt int64

	err := json.Unmarshal(data, &unixTimeInt)
	if err != nil {
		return err
	}

	unixTime.Time = time.Unix(unixTimeInt, 0)
	return nil
}

// AccountCoreData

type AccountCoreData struct {
	Name      string
	UUID      string
	CreatedAt UnixTime
	UpdatedAt UnixTime
	//CreatedAt time.Time
	//UpdatedAt time.Time
}

/*
// avoid inheritance!!!
type _AccountCoreDataAlias AccountCoreData

// AccountCoreData reflection redefinition
type _AuxAccountCoreData struct {
	//*AccountCoreData       `bson:",inline" json:",inline"`
	*_AccountCoreDataAlias `bson:",inline" json:",inline"`
	CreatedAt              int64
	UpdatedAt              int64
}

func (accountCoreData AccountCoreData) MarshalJSON() ([]byte, error) {
	return json.Marshal(_AuxAccountCoreData{
		//AccountCoreData:       &accountCoreData,
		_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountCoreData),
		CreatedAt:             accountCoreData.CreatedAt.Unix(),
		UpdatedAt:             accountCoreData.UpdatedAt.Unix(),
	})
}

func (accountCoreData *AccountCoreData) UnmarshalJSON(data []byte) error {
	aux := &_AuxAccountCoreData{
		//AccountCoreData: accountCoreData,
		_AccountCoreDataAlias: (*_AccountCoreDataAlias)(accountCoreData),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	accountCoreData.CreatedAt = time.Unix(aux.CreatedAt, 0)
	accountCoreData.UpdatedAt = time.Unix(aux.UpdatedAt, 0)
	return nil
}
*/

// AccountInfo

type AccountInfo struct {
	AccountCoreData `bson:",inline" json:",inline"`
	BirthYear       int64
	Description     string
}

/*
// avoid inheritance!!!
//type _AccountInfoAlias AccountInfo

// AccountInfo reflection redefinition
type _AuxAccountInfo struct {
	*_AuxAccountCoreData `bson:",inline" json:",inline"`
	BirthYear            *int64
	Description          *string
}

func (accountInfo AccountInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(_AuxAccountInfo{
		_AuxAccountCoreData: &_AuxAccountCoreData{
			//AccountCoreData:       &accountInfo.AccountCoreData,
			_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountInfo.AccountCoreData),
			CreatedAt:             accountInfo.AccountCoreData.CreatedAt.Unix(),
			UpdatedAt:             accountInfo.AccountCoreData.UpdatedAt.Unix(),
		},
		BirthYear:   &accountInfo.BirthYear,
		Description: &accountInfo.Description,
	})
}

func (accountInfo *AccountInfo) UnmarshalJSON(data []byte) error {
	aux := &_AuxAccountInfo{
		_AuxAccountCoreData: &_AuxAccountCoreData{
			//AccountCoreData: &accountInfo.AccountCoreData,
			_AccountCoreDataAlias: (*_AccountCoreDataAlias)(&accountInfo.AccountCoreData),
		},
		BirthYear:   &accountInfo.BirthYear,
		Description: &accountInfo.Description,
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	accountInfo.AccountCoreData.CreatedAt = time.Unix(aux._AuxAccountCoreData.CreatedAt, 0)
	accountInfo.AccountCoreData.UpdatedAt = time.Unix(aux._AuxAccountCoreData.UpdatedAt, 0)
	return nil
}
*/

func main() {
	accountCoreData := AccountCoreData{
		Name:      "Rasmus Faber",
		UUID:      "123e4567-e89b-12d3-a456-426655440000",
		CreatedAt: UnixTime{Time: time.Now()},
		UpdatedAt: UnixTime{Time: time.Unix(1501998000, 0)},
		//CreatedAt: time.Now(),
		//UpdatedAt: time.Unix(1501998000, 0),
	}

	accountInfo := AccountInfo{
		AccountCoreData: accountCoreData,
		BirthYear:       1979,
		Description:     "Swedish pianist, DJ, remixer, composer, record producer, sound engineer, and founder of the record label Farplane Records.",
	}

	fmt.Println("account core data:")
	fmt.Println("original:")
	fmt.Println(accountCoreData)
	accountCoreDataJson, err := json.Marshal(accountCoreData)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON:")
	fmt.Println(string(accountCoreDataJson))
	var accountCoreDataParsed AccountCoreData
	json.Unmarshal(accountCoreDataJson, &accountCoreDataParsed)
	fmt.Println("parsed:")
	fmt.Println(accountCoreDataParsed)

	accountCoreDataXML, err := xml.Marshal(accountCoreData)
	if err != nil {
		panic(err)
	}
	fmt.Println("XML:")
	fmt.Println(string(accountCoreDataXML))

	fmt.Println("account info:")
	fmt.Println("original:")
	fmt.Println(accountInfo)
	accountInfoJson, err := json.Marshal(accountInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON:")
	fmt.Println(string(accountInfoJson))
	var accountInfoParsed AccountInfo
	json.Unmarshal(accountInfoJson, &accountInfoParsed)
	fmt.Println("parsed:")
	fmt.Println(accountInfoParsed)

	accountInfoXML, err := xml.Marshal(accountInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("XML:")
	fmt.Println(string(accountInfoXML))
}
