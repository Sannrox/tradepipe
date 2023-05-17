package scylla

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type Scylla struct {
	Session      *gocql.Session
	ContactPoint string
	Keyspace     string
}

type KeyspaceConnection interface {
	All(string) ([]map[string]interface{}, error)
	GetById(string, string) (map[string]interface{}, error)
	InsertOne(string, map[string]interface{}) error
	DeleteOne(string, string) error
}

func NewScyllaKeySpaceConnection(contactPoint string, keyspace string) (*Scylla, error) {
	if err := CreateKeyspace(contactPoint, keyspace); err != nil {
		return nil, err
	}
	return NewScyllaDbWithPool(contactPoint, keyspace, 10)
}

func TryToConnectWithRetry(contactPoint string, attempts int, timeout time.Duration) error {
	cluster := gocql.NewCluster(contactPoint)
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

func NewScyllaDbWithPool(contactPoint string, keyspace string, poolSize int) (*Scylla, error) {
	cluster := gocql.NewCluster(contactPoint)
	cluster.Keyspace = keyspace
	cluster.PoolConfig.HostSelectionPolicy = gocql.RoundRobinHostPolicy()
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.NumConns = poolSize

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	logrus.Infof("established connection to %s keyspace", keyspace)

	return &Scylla{
		Session:  session,
		Keyspace: keyspace,
	}, nil
}

func CreateKeyspace(contactPoint string, keyspace string) error {
	clusterConfig := gocql.NewCluster(contactPoint)
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

func (s *Scylla) CreateTable(tablename, schema string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tablename, schema)
	logrus.Debug(query)
	return s.Session.Query(query).Exec()
}

func (s *Scylla) CreateTablePath(tableName string, prefix string) string {

	return s.Keyspace + "." + strings.TrimSpace(prefix) + strings.ReplaceAll(tableName, "-", "_")
}

func (s *Scylla) TableExists(keyspace, table string) bool {
	var tableName string
	err := s.Session.Query("SELECT table_name FROM system_schema.tables WHERE keyspace_name = ? AND table_name = ?", keyspace, table).Scan(&tableName)
	return err == nil
}

func (s *Scylla) Close() {
	s.Session.Close()
}

func (s *Scylla) Insert(tableName string, model interface{}) error {
	fields := getFields(model)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(fields, ", "), getValuesPlaceholder(model))
	logrus.Debug("QUERY:", query, " Len of Placeholders:", len(fields), " VALUES: ", getValues(model), "Len of Values", len(getValues(model)))
	return s.Session.Query(query, getValues(model)...).Exec()
}

func (s *Scylla) Query(query string, model ...interface{}) error {
	logrus.Info(s.Session.Query(query, model...).String())
	return s.Session.Query(query, model...).Exec()
}
func getFields(model interface{}) []string {
	t := reflect.TypeOf(model)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "ID" {
			fields = append(fields, "id")
		} else {
			fields = append(fields, strings.ToLower(field.Name))
		}
	}

	return fields
}

// func getFields(model interface{}) []string {
// 	t := reflect.TypeOf(model)

// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}

// 	var fields []string
// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		if field.Name == "ID" {
// 			fields = append(fields, "id")
// 		} else if field.Type.Kind() == reflect.Struct {
// 			// recursively call getFields on inner struct type
// 			innerFields := getFields(reflect.New(field.Type).Interface())
// 			fields = append(fields, innerFields...)
// 		} else {
// 			fields = append(fields, strings.ToLower(field.Name))
// 		}
// 	}

// 	return fields
// }

func getValuesPlaceholder(model interface{}) string {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var placeholders []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			structPlacesholders := getValuesPlaceholder(field.Addr().Interface())
			placeholders = append(placeholders, fmt.Sprintf("(%s)", structPlacesholders))
		} else {
			placeholders = append(placeholders, "?")
		}
	}
	return strings.Join(placeholders, ", ")
}

func getValues(model interface{}) []interface{} {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var values []interface{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			structValues := getValues(field.Addr().Interface())
			values = append(values, structValues...)
		} else {
			values = append(values, field.Interface())
		}
	}
	return values
}

func (s *Scylla) GetByID(tableName string, id gocql.UUID, model interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	logrus.Debug(query)
	return s.Session.Query(query, id).Scan(model)
}

func (s *Scylla) GetByKey(tableName string, key string, value string, model interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? ", tableName, key)
	logrus.Debug(query)
	return s.Session.Query(query, value).Scan(model)
}

func (s *Scylla) Update(tableName string, id gocql.UUID, model interface{}) error {
	fields := getFields(model)

	setValues := make([]string, len(fields)-1)
	v := reflect.ValueOf(model).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Name != "ID" {
			setValues[i-1] = fmt.Sprintf("%s = ?", strings.ToLower(field.Name))
		}
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, strings.Join(setValues, ", "))
	logrus.Debug(query)
	logrus.Debug(getValues(model)[1:])
	queryObj := s.Session.Query(query, getValues(model)[1:], id)
	logrus.Debug(queryObj)
	return queryObj.Exec()
}

func (s *Scylla) Delete(tableName string, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	logrus.Debug(query)
	return s.Session.Query(query, id).Exec()
}

func (s *Scylla) GetAllValues(tableName string) *gocql.Iter {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	logrus.Debug(query)

	return s.Session.Query(query).Iter()

}
