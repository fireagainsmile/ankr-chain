package mysql

import (
	"encoding/json"
	"errors"

	"github.com/Ankr-network/ankr-chain/historystorage/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	TableSQLSendTx = "CREATE TABLE IF NOT EXISTS `sendtx`" +
		"(`txhash` varchar(255) NOT NULL," +
		"`txtype` varchar(255) default NULL," +
		"`height` bigint(20) default 0," +
		"`index` int(11) default 0," +
		"`time` varchar(255) default NULL," +
		"`fromaddress` varchar(255) default NULL," +
		"`toaddress` varchar(255) default NULL," +
		"`amount` varchar(255) default NULL," +
		"PRIMARY KEY (`txhash`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"

	TableSQLMetering = "CREATE TABLE IF NOT EXISTS `metering`" +
		"(`txhash` varchar(255) NOT NULL," +
		"`txtype` varchar(255) default NULL," +
		"`height` bigint(20) default 0," +
		"`index` int(11) default 0," +
		"`time` varchar(255) default NULL," +
		"`dc` varchar(255) default NULL," +
		"`ns` varchar(255) default NULL," +
		"`value` varchar(255) default NULL," +
		"PRIMARY KEY (`txhash`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"

	TableSQLSetBalanceTx = "CREATE TABLE IF NOT EXISTS `setbalancetx`" +
		"(`txhash` varchar(255) NOT NULL," +
		"`txtype` varchar(255) default NULL," +
		"`height` bigint(20) default 0," +
		"`index` int(11) default 0," +
		"`time` varchar(255) default NULL," +
		"`address` varchar(255) default NULL," +
		"`amount` varchar(255) default NULL," +
		"PRIMARY KEY (`txhash`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"

	TableSQLSetStakex = "CREATE TABLE IF NOT EXISTS `setstaketx`" +
		"(`txhash` varchar(255) NOT NULL," +
		"`txtype` varchar(255) default NULL," +
		"`height` bigint(20) default 0," +
		"`index` int(11) default 0," +
		"`time` varchar(255) default NULL," +
		"`amount` varchar(255) default NULL," +
		"PRIMARY KEY (`txhash`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"

	TableSQLSetValidatorTx = "CREATE TABLE IF NOT EXISTS `setvalidatortx`" +
		"(`txhash` varchar(255) NOT NULL," +
		"`txtype` varchar(255) default NULL," +
		"`height` bigint(20) default 0," +
		"`index` int(11) default 0," +
		"`time` varchar(255) default NULL," +
		"`validatorpubkey` varchar(255) default NULL," +
		"`power` varchar(255) default NULL," +
		"PRIMARY KEY (`txhash`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"

	TableSQLAccount = "CREATE TABLE IF NOT EXISTS `account`" +
		"(`address` varchar(255) NOT NULL," +
		"`balance` varchar(255) default NULL," +
		"PRIMARY KEY (`address`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
)

type  MySql struct {
	dbOper *DBOperation
	logHis log.Logger
}

func NewMySql(dbHost string, dbName string, logHis log.Logger) *MySql {
	dbo := NewDBOperation(dbHost, dbName, logHis)
	if dbo != nil {
		err := dbo.createTable(TableSQLSendTx)
		if err != nil {
			return nil
		}

		err = dbo.createTable(TableSQLMetering)
		if err != nil {
			return nil
		}

		err = dbo.createTable(TableSQLSetBalanceTx)
		if err != nil {
			return nil
		}

		err = dbo.createTable(TableSQLSetStakex)
		if err != nil {
			return nil
		}

		err = dbo.createTable(TableSQLSetValidatorTx)
		if err != nil {
			return nil
		}

		err = dbo.createTable(TableSQLAccount)
		if err != nil {
			return nil
		}

		return &MySql{dbo, logHis}
	}

	return nil
}

func (db *MySql) AddSendTx(tx *types.TransactionSendTx) error {
	sql := "INSERT INTO sendtx(txhash,txtype,height,sendtx.index,time,fromaddress,toaddress,amount) values(?,?,?,?,?,?,?,?)"
	if db.dbOper != nil{
		return db.dbOper.set(sql, tx.TxHash, tx.TxType, tx.Height, tx.Index, tx.Time, tx.FromAddress, tx.ToAddress, tx.Amount)
	}

	db.logHis.Error("db.dbOper nil, can't AddSendTx")

	return errors.New("db.dbOper nil")
}

func (db *MySql) AddMetering(tx *types.TransactionMetering) error {
	sql := "INSERT INTO metering(txhash,txtype,height,metering.index,time,dc,ns,value) values(?,?,?,?,?,?,?,?)"
	if db.dbOper != nil{
		return db.dbOper.set(sql, tx.TxHash, tx.TxType, tx.Height, tx.Index, tx.Time, tx.DC, tx.NS, tx.Value)
	}

	db.logHis.Error("db.dbOper nil, can't AddMetering")

	return errors.New("db.dbOper nil")
}

func (db *MySql) AddSetBalanceTx(tx *types.TransactionSetBalanceTx) error {
	sql := "INSERT INTO setbalancetx(txhash,txtype,height,setbalancetx.index,time,address,amount) values(?,?,?,?,?,?,?)"
	if db.dbOper != nil{
		return db.dbOper.set(sql, tx.TxHash, tx.TxType, tx.Height, tx.Index, tx.Time, tx.Address, tx.Amount)
	}

	db.logHis.Error("db.dbOper nil, can't AddSetBalanceTx")

	return errors.New("db.dbOper nil")
}

func (db *MySql) AddSetStakeTx(tx *types.TransactionSetStakeTx) error {
	sql := "INSERT INTO setstaketx(txhash,txtype,height,setstaketx.index,time,amount) values(?,?,?,?,?,?)"
	if db.dbOper != nil{
		return db.dbOper.set(sql, tx.TxHash, tx.TxType, tx.Height, tx.Index, tx.Time, tx.Amount)
	}

	db.logHis.Error("db.dbOper nil, can't AddSetStakeTx")

	return errors.New("db.dbOper nil")
}

func (db *MySql) AddSetValidatorTx(tx *types.TransactionSetValidatorTx) error {
	sql := "INSERT INTO setvalidatortx(txhash,txtype,height,setvalidatortx.index,time,validatorpubkey, power) values(?,?,?,?,?,?,?)"
	if db.dbOper != nil{
		return db.dbOper.set(sql, tx.TxHash, tx.TxType, tx.Height, tx.Index, tx.Time, tx.ValidatorPubkey, tx.Power)
	}

	db.logHis.Error("db.dbOper nil, can't AddSetValidatorTx")

	return errors.New("db.dbOper nil")
}

func (db *MySql) AddAccount(account *types.Account) error {
	sql := "INSERT INTO account(address, balance) values(?,?)"
	if db.dbOper != nil{
		return db.dbOper.set(sql, account.Address, account.Balance)
	}

	db.logHis.Error("db.dbOper nil, can't AddAccount")

	return errors.New("db.dbOper nil")
}
func (db *MySql) UpdateAccount(Address string, amount string) error {
	sql := "UPDATE account set balance = ? where address = ?"
	if db.dbOper != nil{
		return db.dbOper.set(sql, amount, Address)
	}

	db.logHis.Error("db.dbOper nil, can't UpdateAccount")

	return errors.New("db.dbOper nil")
}

func (db *MySql) GetAccount(Address string) (*types.Account, error) {
	sql := "SELECT * from account where address = ?"
	if db.dbOper != nil{
		jsonObj, err := db.dbOper.query(sql, Address)
		if err == nil {
			mapAcc := make(map[string]interface{}, 1)
			err = json.Unmarshal([]byte(jsonObj), &mapAcc)
			if err == nil {
				addr, _    := mapAcc["address"].(string)
				balance, _ := mapAcc["balance"].(string)
				return &types.Account{addr, balance}, nil
			}else {
				db.logHis.Error("jsonObj Unmarshal failed", "err", err, "jsonObj", jsonObj)
			}
		}else {
			db.logHis.Error("db.dbOper.query failed", "err", err, "sql", sql, "address", Address)
		}
	}

	db.logHis.Error("db.dbOper nil, can't GetAccount")

	return nil, nil
}


