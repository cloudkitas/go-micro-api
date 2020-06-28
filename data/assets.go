package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type Asset struct {
	ID          int
	Name        string
	Description string
}

func (a *Asset) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

type assets []*Asset

func GetAssets() assets {
	return assetList
}

func (a *assets) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func AddAssets(a *Asset) {
	a.ID = getNextAssetID()
	assetList = append(assetList, a)
}

func getNextAssetID() int {
	la := assetList[len(assetList)-1]
	return la.ID + 1
}

func UpdateAsset(id int, a *Asset) error {
	_, pos, err := findAsset(id)
	if err != nil {
		return err
	}

	a.ID = id
	assetList[pos] = a
	return nil

}

var ErrAssetNotFound = fmt.Errorf("Asset Not Found")

func findAsset(id int) (*Asset, int, error) {
	for i, a := range assetList {
		if a.ID == id {
			return a, i, nil
		}
	}

	return nil, -1, ErrAssetNotFound
}

var assetList = []*Asset{
	&Asset{
		ID:          1,
		Name:        "FTSE100",
		Description: "UK 100 companies",
	},
	&Asset{
		ID:          2,
		Name:        "DOW Jones Index",
		Description: "US 30",
	},
}
