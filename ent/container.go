// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/dstgo/maxwell/ent/container"
)

// container info table
type Container struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Cid holds the value of the "cid" field.
	Cid string `json:"cid,omitempty"`
	// Image holds the value of the "image" field.
	Image string `json:"image,omitempty"`
	// CPU holds the value of the "cpu" field.
	CPU int64 `json:"cpu,omitempty"`
	// Memory holds the value of the "memory" field.
	Memory int64 `json:"memory,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt int64 `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt int64 `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ContainerQuery when eager-loading is set.
	Edges        ContainerEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ContainerEdges holds the relations/edges for other nodes in the graph.
type ContainerEdges struct {
	// Owner holds the value of the owner edge.
	Owner []*User `json:"owner,omitempty"`
	// Node holds the value of the node edge.
	Node []*Node `json:"node,omitempty"`
	// Mounts holds the value of the mounts edge.
	Mounts []*Mount `json:"mounts,omitempty"`
	// Ports holds the value of the ports edge.
	Ports []*Port `json:"ports,omitempty"`
	// Jobs holds the value of the jobs edge.
	Jobs []*Job `json:"jobs,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading.
func (e ContainerEdges) OwnerOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Owner, nil
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// NodeOrErr returns the Node value or an error if the edge
// was not loaded in eager-loading.
func (e ContainerEdges) NodeOrErr() ([]*Node, error) {
	if e.loadedTypes[1] {
		return e.Node, nil
	}
	return nil, &NotLoadedError{edge: "node"}
}

// MountsOrErr returns the Mounts value or an error if the edge
// was not loaded in eager-loading.
func (e ContainerEdges) MountsOrErr() ([]*Mount, error) {
	if e.loadedTypes[2] {
		return e.Mounts, nil
	}
	return nil, &NotLoadedError{edge: "mounts"}
}

// PortsOrErr returns the Ports value or an error if the edge
// was not loaded in eager-loading.
func (e ContainerEdges) PortsOrErr() ([]*Port, error) {
	if e.loadedTypes[3] {
		return e.Ports, nil
	}
	return nil, &NotLoadedError{edge: "ports"}
}

// JobsOrErr returns the Jobs value or an error if the edge
// was not loaded in eager-loading.
func (e ContainerEdges) JobsOrErr() ([]*Job, error) {
	if e.loadedTypes[4] {
		return e.Jobs, nil
	}
	return nil, &NotLoadedError{edge: "jobs"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Container) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case container.FieldID, container.FieldCPU, container.FieldMemory, container.FieldCreatedAt, container.FieldUpdatedAt:
			values[i] = new(sql.NullInt64)
		case container.FieldName, container.FieldCid, container.FieldImage:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Container fields.
func (c *Container) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case container.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case container.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case container.FieldCid:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field cid", values[i])
			} else if value.Valid {
				c.Cid = value.String
			}
		case container.FieldImage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field image", values[i])
			} else if value.Valid {
				c.Image = value.String
			}
		case container.FieldCPU:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field cpu", values[i])
			} else if value.Valid {
				c.CPU = value.Int64
			}
		case container.FieldMemory:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field memory", values[i])
			} else if value.Valid {
				c.Memory = value.Int64
			}
		case container.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				c.CreatedAt = value.Int64
			}
		case container.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				c.UpdatedAt = value.Int64
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Container.
// This includes values selected through modifiers, order, etc.
func (c *Container) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryOwner queries the "owner" edge of the Container entity.
func (c *Container) QueryOwner() *UserQuery {
	return NewContainerClient(c.config).QueryOwner(c)
}

// QueryNode queries the "node" edge of the Container entity.
func (c *Container) QueryNode() *NodeQuery {
	return NewContainerClient(c.config).QueryNode(c)
}

// QueryMounts queries the "mounts" edge of the Container entity.
func (c *Container) QueryMounts() *MountQuery {
	return NewContainerClient(c.config).QueryMounts(c)
}

// QueryPorts queries the "ports" edge of the Container entity.
func (c *Container) QueryPorts() *PortQuery {
	return NewContainerClient(c.config).QueryPorts(c)
}

// QueryJobs queries the "jobs" edge of the Container entity.
func (c *Container) QueryJobs() *JobQuery {
	return NewContainerClient(c.config).QueryJobs(c)
}

// Update returns a builder for updating this Container.
// Note that you need to call Container.Unwrap() before calling this method if this Container
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Container) Update() *ContainerUpdateOne {
	return NewContainerClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Container entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Container) Unwrap() *Container {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Container is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Container) String() string {
	var builder strings.Builder
	builder.WriteString("Container(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("cid=")
	builder.WriteString(c.Cid)
	builder.WriteString(", ")
	builder.WriteString("image=")
	builder.WriteString(c.Image)
	builder.WriteString(", ")
	builder.WriteString("cpu=")
	builder.WriteString(fmt.Sprintf("%v", c.CPU))
	builder.WriteString(", ")
	builder.WriteString("memory=")
	builder.WriteString(fmt.Sprintf("%v", c.Memory))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(fmt.Sprintf("%v", c.CreatedAt))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(fmt.Sprintf("%v", c.UpdatedAt))
	builder.WriteByte(')')
	return builder.String()
}

// Containers is a parsable slice of Container.
type Containers []*Container
