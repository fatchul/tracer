package attribute

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	DatastoreCassandra     DatabaseProduct = "Cassandra"
	DatastoreDerby         DatabaseProduct = "Derby"
	DatastoreElasticsearch DatabaseProduct = "Elasticsearch"
	DatastoreFirebird      DatabaseProduct = "Firebird"
	DatastoreIBMDB2        DatabaseProduct = "IBMDB2"
	DatastoreInformix      DatabaseProduct = "Informix"
	DatastoreMemcached     DatabaseProduct = "Memcached"
	DatastoreMongoDB       DatabaseProduct = "MongoDB"
	DatastoreMySQL         DatabaseProduct = "MySQL"
	DatastoreMSSQL         DatabaseProduct = "MSSQL"
	DatastoreNeptune       DatabaseProduct = "Neptune"
	DatastoreOracle        DatabaseProduct = "Oracle"
	DatastorePostgres      DatabaseProduct = "Postgres"
	DatastoreRedis         DatabaseProduct = "Redis"
	DatastoreSolr          DatabaseProduct = "Solr"
	DatastoreSQLite        DatabaseProduct = "SQLite"
	DatastoreCouchDB       DatabaseProduct = "CouchDB"
	DatastoreRiak          DatabaseProduct = "Riak"
	DatastoreVoltDB        DatabaseProduct = "VoltDB"
	DatastoreDynamoDB      DatabaseProduct = "DynamoDB"
)

type (
	DatabaseProduct string

	DatabaseAttr struct {
		Collection         string                 `json:"db.collection"`
		Operation          string                 `json:"db.operation"`
		Product            DatabaseProduct        `json:"db.product"`
		ParameterizedQuery string                 `json:"db.parameterized_query"`
		QueryParameters    map[string]interface{} `json:"db.query_parameters"`
	}
)

func Database(attr DatabaseAttr) []attribute.KeyValue {
	return parseAttrToKV(attr)
}
