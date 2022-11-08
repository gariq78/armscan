package interfaces

import (
	"encoding/json"
	"fmt"
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

type DataRepository struct {
	storage DbHandler
	logg    Logger
}

func NewDataRepository(dbhandler DbHandler, logg Logger) (*DataRepository, error) {
	rv := DataRepository{
		storage: dbhandler,
		logg:    logg,
	}

	err := rv.initDb()
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	return &rv, nil
}

func (dr *DataRepository) initDb() error {
	_, err := dr.storage.Execute(sqlCreateTableAssets)
	return err
}

// const sqltimelayout = `2006-01-02 15:04:05Z07:00`

// func (dr *DataRepository) GetLastTimeStamp() (time.Time, error) {
// 	var maxTime time.Time
// 	r, err := dr.storage.Query(sqlMaxStamp)
// 	if err != nil {
// 		return maxTime, fmt.Errorf("max stamp: %w", err)
// 	}
// 	defer r.Close()

// 	if r.Next() {
// 		var maxtimestr sql.NullString
// 		err = r.Scan(&maxtimestr) // 2020-10-20 15:16:51+00:00
// 		if err != nil {
// 			return maxTime, fmt.Errorf("row scan: %w", err)
// 		}

// 		if maxtimestr.Valid {
// 			maxTime, err = time.Parse(sqltimelayout, maxtimestr.String)
// 			if err != nil {
// 				return maxTime, fmt.Errorf("time.Parse: %w", err)
// 			}
// 		}
// 	} else {
// 		return maxTime, fmt.Errorf("max stamp return 0 rows")
// 	}

// 	return maxTime, nil
// }

func (dr *DataRepository) AddData(m structs.AddDataPacket) error {
	// id, err := getIDfromAsset(m.ДанныеПоАктиву)
	// if err != nil {
	// 	return fmt.Errorf("getIDfromAsset(%v) err: %w", m, err)
	// }

	jsonBlob, err := json.Marshal(m.Asset)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	dr.logg.Inf("DataRepository.AddData json asset: %s", jsonBlob)

	res, err := dr.storage.Execute(sqlUpdateAsset, m.ScanTime, jsonBlob, m.Asset.HostID)
	if err != nil {
		return fmt.Errorf("update asset: %w", err)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		dr.logg.Err(err, "DataRepository.RowsAffected sqlUpdateAsset")
	}
	dr.logg.Inf("DataRepository.sqlUpdateAsset affected %d", aff)

	res, err = dr.storage.Execute(sqlInsertAsset, m.ScanTime, jsonBlob, m.Asset.HostID)
	if err != nil {
		return fmt.Errorf("insert asset: %w", err)
	}

	aff, err = res.RowsAffected()
	if err != nil {
		dr.logg.Err(err, "DataRepository.RowsAffected sqlInsertAsset")
	}
	dr.logg.Inf("DataRepository.sqlInsertAsset affected %d", aff)

	return nil
}

func (dr *DataRepository) Datas() ([]asset.Asset, error) {
	r, err := dr.storage.Query(sqlSelectAssets)
	if err != nil {
		return nil, fmt.Errorf("select assets: %w", err)
	}
	defer r.Close()

	res := make([]asset.Asset, 0)

	for r.Next() {
		var value string
		var assetId string
		var rowstamp time.Time

		err := r.Scan(&rowstamp, &assetId, &value)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		var m asset.Asset
		err = json.Unmarshal([]byte(value), &m)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal message: %w", err)
		}

		// if rowstamp.After(maxTime) {
		// 	maxTime = rowstamp
		// }

		res = append(res, m)
	}

	return res, nil
}

// func getIDfromAsset(d domain.ДанныеПоАктиву) (string, error) {
// 	if d.Hardware.MAC == "" {
// 		if d.Software.HostName == "" {
// 			return "", fmt.Errorf("Cannot calculate asset ID")
// 		} else {
// 			return d.Software.HostName, nil
// 		}
// 	} else {
// 		return d.Hardware.MAC, nil
// 	}
// }
