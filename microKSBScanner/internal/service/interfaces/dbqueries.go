package interfaces

const (
	//	"id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	sqlCreateTableAssets = `CREATE TABLE IF NOT EXISTS "assets" (
"asset_id" TEXT PRIMARY KEY NOT NULL,
"time" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
"value" TEXT
);`

	sqlInsertAsset = `INSERT OR IGNORE INTO "assets" ("time","value","asset_id") VALUES(?,?,?);`

	sqlUpdateAsset = `UPDATE "assets" SET "time"=?, "value"=? WHERE "asset_id"=?;`

	sqlSelectAssets = `SELECT "time","asset_id","value" FROM "assets";` // WHERE "time" > $1;`

	// это не нужно - нужно при отдаче сообщений искать максимальный stamp
	// это нужно - при старте узнать что в базе и начать с максимального, а не запрашивать с нуля
	// sqlMaxStamp = `SELECT max("time") as maxtime from "assets"`

	// TODO: нужно придумать как удалять записи
)
