package mongodb

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Ankr-network/ankr-chain/historystorage/types"
)

var once sync.Once

type MongoDB struct {
	dbName   string
	session  *mgo.Session
}

func NewMongoDB(host string, name string) *MongoDB{
	var connS *mgo.Session
	var err error

	once.Do(func() {
		info := mgo.DialInfo{
			Addrs:     []string{host},
			Timeout:   time.Duration(5) * time.Second,
			PoolLimit: 4096,
		}

		if connS, err = mgo.DialWithInfo(&info); err != nil {
			return
		}
		connS.SetMode(mgo.Monotonic, true)
		connS.SetSafe(&mgo.Safe{})
	})

	if connS != nil {
		return &MongoDB{name, connS}
	}

	return nil
}

func (db *MongoDB) collection(session *mgo.Session, colName string) *mgo.Collection {
	col := session.DB(db.dbName).C(colName)

	if colName == "account" {
		index := mgo.Index{
			Key: []string{"address"},
			Unique: true,
			DropDups: false,
			Background: true,
			Sparse: true,
		}
		col.EnsureIndex(index)
	}

	if colName == "transaction" {
		col.EnsureIndexKey("txhash")
		col.EnsureIndexKey("txtype")

		index := mgo.Index{
			Key: []string{"txhash", "txtype"},
			Unique: true,
			DropDups: false,
			Background: true,
			Sparse: true,
		}
		col.EnsureIndex(index)
	}

	return col
}

func (db *MongoDB) AddSendTx(tx *types.TransactionSendTx) error {
	session := db.session.Copy()
	defer session.Close()
	err := db.collection(session, "transaction").Insert(bson.M{"txhash": tx.TxHash, "txtype": tx.TxType, "height": tx.Height, "index": tx.Index, "time": tx.Time, "fromaddress": tx.FromAddress, "toaddress": tx.ToAddress, "amount": tx.Amount})
	return err
}

func (db *MongoDB) AddMetering(tx *types.TransactionMetering) error {
	session := db.session.Copy()
	defer session.Close()
	err := db.collection(session, "transaction").Insert(bson.M{"txhash": tx.TxHash, "txtype": tx.TxType, "height": tx.Height, "index": tx.Index, "time": tx.Time, "dc": tx.DC, "ns": tx.NS, "value": tx.Value})
	return err
}

func (db *MongoDB) AddSetBalanceTx(tx *types.TransactionSetBalanceTx) error {
	session := db.session.Copy()
	defer session.Close()
	err := db.collection(session, "transaction").Insert(bson.M{"txhash": tx.TxHash, "txtype": tx.TxType, "height": tx.Height, "index": tx.Index, "time": tx.Time, "address": tx.Address, "Amount": tx.Amount})
	return err
}

func (db *MongoDB) AddSetStakeTx(tx *types.TransactionSetStakeTx) error {
	session := db.session.Copy()
	defer session.Close()
	err := db.collection(session, "transaction").Insert(bson.M{"txhash": tx.TxHash, "txtype": tx.TxType, "height": tx.Height, "index": tx.Index, "time": tx.Time, "Amount": tx.Amount})
	return err
}

func (db *MongoDB) AddSetValidatorTx(tx *types.TransactionSetValidatorTx) error {
	session := db.session.Copy()
	defer session.Close()
	err := db.collection(session, "transaction").Insert(bson.M{"txhash": tx.TxHash, "txtype": tx.TxType, "height": tx.Height, "index": tx.Index, "time": tx.Time, "ValidatorPubkey": tx.ValidatorPubkey, "Power": tx.Power})
	return err
}

func (db *MongoDB) AddAccount(account *types.Account) error {
	session := db.session.Copy()
	defer session.Close()
	err := db.collection(session, "account").Insert(bson.M{"address": account.Address, "balance": account.Balance})
	return err
}

func (db *MongoDB) UpdateAccount(Address string, amount string) error {
	session := db.session.Copy()
	defer session.Close()
	err := db.collection(session, "account").Update(bson.M{"address": Address}, bson.M{"$set": bson.M{"balance": amount}})
	return err
}

func (db *MongoDB) GetAccount(Address string) (*types.Account, error) {
	session := db.session.Copy()
	defer session.Close()
	var accS [] *types.Account
	db.collection(session, "account").Find(bson.M{"address" :Address}).Limit(1).All(&accS)
	if len(accS) == 0 {
		return nil, fmt.Errorf("there is no account record of address %s", Address)
	}
	return accS[0], nil
}




