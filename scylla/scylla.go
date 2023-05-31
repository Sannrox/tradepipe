package scylla

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/sirupsen/logrus"
)

type Scylla struct {
	Session      *gocqlx.Session
	ContactPoint string
	Keyspace     string
}

func NewScyllaKeySpaceConnection(contactPoint string, port int, keyspace string) (*Scylla, error) {
	if err := CreateKeyspace(contactPoint, port, keyspace); err != nil {
		return nil, err
	}
	return NewScyllaDbWithPool(contactPoint, port, keyspace, 10)
}

func TryToConnectWithRetry(contactPoint string, port int, attempts int, timeout time.Duration) error {
	cluster := gocql.NewCluster(contactPoint)
	cluster.Port = port
	for i := 0; i < attempts; i++ {
		session, err := cluster.CreateSession()
		if err != nil {
			logrus.Warnf("could not connect to scylla, retrying in %s", timeout)
			time.Sleep(timeout)
			continue
		}

		defer session.Close()
		logrus.Infof("established connection to scylla")
		return nil
	}

	return fmt.Errorf("could not connect to scylla after %d attempts", attempts)
}

func NewScyllaDbWithPool(contactPoint string, port int, keyspace string, poolSize int) (*Scylla, error) {
	cluster := gocql.NewCluster(contactPoint)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.PoolConfig.HostSelectionPolicy = gocql.RoundRobinHostPolicy()
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.NumConns = poolSize

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, err
	}

	logrus.Infof("established connection to %s keyspace", keyspace)

	return &Scylla{
		Session:  &session,
		Keyspace: keyspace,
	}, nil
}

func CreateKeyspace(contactPoint string, port int, keyspace string) error {
	clusterConfig := gocql.NewCluster(contactPoint)
	clusterConfig.Port = port
	clusterConfig.Keyspace = "system"
	clusterConfig.ProtoVersion = 4
	session, err := clusterConfig.CreateSession()
	if err != nil {
		return err
	}

	defer session.Close()
	err = session.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}", keyspace)).Exec()
	if err != nil {
		logrus.Fatal(err)
	}
	return nil

}

func (s *Scylla) NewTable(schema table.Metadata) *table.Table {
	return table.New(schema)
}

func (s *Scylla) CreateTablePath(tableName string, prefix string) string {
	return s.Keyspace + "." + strings.TrimSpace(prefix) + strings.ReplaceAll(tableName, "-", "_")
}

func (s *Scylla) Close() {
	s.Session.Close()
}

func (s *Scylla) Insert(metaTable *table.Table, model interface{}) error {
	q := s.Session.Query(metaTable.Insert()).BindStruct(model)
	return q.ExecRelease()
}

func (s *Scylla) GetByKeys(metaTable *table.Table, model interface{}) *gocqlx.Queryx {
	return s.Session.Query(metaTable.Get()).BindStruct(model)
}

func (s *Scylla) GetAll(metaTable *table.Table) *gocqlx.Queryx {
	return s.Session.Query(metaTable.SelectAll())
}

func (s *Scylla) Delete(metaTable *table.Table, model interface{}) error {
	q := s.Session.Query(metaTable.Delete()).BindStruct(model)
	logrus.Info(q.Query)
	return q.ExecRelease()
}

func (s *Scylla) Update(metaTable *table.Table, model interface{}, columns ...string) error {
	q := s.Session.Query(metaTable.Update(columns...)).BindStruct(model)
	return q.ExecRelease()
}

func (s *Scylla) CheckIfTableExits(tableName string) (bool, error) {

	query, names := qb.Select("system_schema.tables").Columns("table_name").Where(qb.EqLit("keyspace_name", s.Keyspace)).Where(qb.EqLit("table_name", tableName)).Limit(1).ToCql()
	iter := s.Session.Query(query, names).Iter()

	exists := iter.Scan(nil)
	if iter.Close() != nil {
		return false, iter.Close()
	}

	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
func (s *Scylla) CreateTable(tablename, schema string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tablename, schema)
	logrus.Debug(query)
	return s.Session.Session.Query(query).Exec()
}
