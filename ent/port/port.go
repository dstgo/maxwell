// Code generated by ent, DO NOT EDIT.

package port

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the port type in the database.
	Label = "port"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldHost holds the string denoting the host field in the database.
	FieldHost = "host"
	// FieldBind holds the string denoting the bind field in the database.
	FieldBind = "bind"
	// FieldProtocol holds the string denoting the protocol field in the database.
	FieldProtocol = "protocol"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// Table holds the table name of the port in the database.
	Table = "ports"
	// OwnerTable is the table that holds the owner relation/edge. The primary key declared below.
	OwnerTable = "container_ports"
	// OwnerInverseTable is the table name for the Container entity.
	// It exists in this package in order to avoid circular dependency with the "container" package.
	OwnerInverseTable = "containers"
)

// Columns holds all SQL columns for port fields.
var Columns = []string{
	FieldID,
	FieldHost,
	FieldBind,
	FieldProtocol,
	FieldCreatedAt,
	FieldUpdatedAt,
}

var (
	// OwnerPrimaryKey and OwnerColumn2 are the table columns denoting the
	// primary key for the owner relation (M2M).
	OwnerPrimaryKey = []string{"container_id", "port_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() int64
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() int64
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() int64
)

// OrderOption defines the ordering options for the Port queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByHost orders the results by the host field.
func ByHost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHost, opts...).ToFunc()
}

// ByBind orders the results by the bind field.
func ByBind(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBind, opts...).ToFunc()
}

// ByProtocol orders the results by the protocol field.
func ByProtocol(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProtocol, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByOwnerCount orders the results by owner count.
func ByOwnerCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOwnerStep(), opts...)
	}
}

// ByOwner orders the results by owner terms.
func ByOwner(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOwnerStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newOwnerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OwnerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, OwnerTable, OwnerPrimaryKey...),
	)
}