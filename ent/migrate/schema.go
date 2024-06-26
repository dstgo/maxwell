// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ContainersColumns holds the columns for the "containers" table.
	ContainersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "cid", Type: field.TypeString, Unique: true},
		{Name: "image", Type: field.TypeString},
		{Name: "cpu", Type: field.TypeInt64},
		{Name: "memory", Type: field.TypeInt64},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "updated_at", Type: field.TypeInt64},
	}
	// ContainersTable holds the schema information for the "containers" table.
	ContainersTable = &schema.Table{
		Name:       "containers",
		Comment:    "container info table",
		Columns:    ContainersColumns,
		PrimaryKey: []*schema.Column{ContainersColumns[0]},
	}
	// JobsColumns holds the columns for the "jobs" table.
	JobsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString},
		{Name: "job_type", Type: field.TypeInt},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "updated_at", Type: field.TypeInt64},
	}
	// JobsTable holds the schema information for the "jobs" table.
	JobsTable = &schema.Table{
		Name:       "jobs",
		Comment:    "container cron job tables",
		Columns:    JobsColumns,
		PrimaryKey: []*schema.Column{JobsColumns[0]},
	}
	// MountsColumns holds the columns for the "mounts" table.
	MountsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "host", Type: field.TypeString},
		{Name: "bind", Type: field.TypeString},
		{Name: "is_dir", Type: field.TypeBool},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "updated_at", Type: field.TypeInt64},
	}
	// MountsTable holds the schema information for the "mounts" table.
	MountsTable = &schema.Table{
		Name:       "mounts",
		Comment:    "container mount record table",
		Columns:    MountsColumns,
		PrimaryKey: []*schema.Column{MountsColumns[0]},
	}
	// NodesColumns holds the columns for the "nodes" table.
	NodesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "uid", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "address", Type: field.TypeString, Unique: true},
		{Name: "note", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "updated_at", Type: field.TypeInt64},
	}
	// NodesTable holds the schema information for the "nodes" table.
	NodesTable = &schema.Table{
		Name:       "nodes",
		Comment:    "remote node info table",
		Columns:    NodesColumns,
		PrimaryKey: []*schema.Column{NodesColumns[0]},
	}
	// PortsColumns holds the columns for the "ports" table.
	PortsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "host", Type: field.TypeInt, Unique: true},
		{Name: "bind", Type: field.TypeInt, Comment: "bind port in container"},
		{Name: "protocol", Type: field.TypeInt, Comment: "network protocol"},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "updated_at", Type: field.TypeInt64},
	}
	// PortsTable holds the schema information for the "ports" table.
	PortsTable = &schema.Table{
		Name:       "ports",
		Comment:    "container ports bind info table",
		Columns:    PortsColumns,
		PrimaryKey: []*schema.Column{PortsColumns[0]},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "uid", Type: field.TypeString, Unique: true},
		{Name: "username", Type: field.TypeString, Unique: true},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "password", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeInt64},
		{Name: "updated_at", Type: field.TypeInt64},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Comment:    "user info table",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// ContainerMountsColumns holds the columns for the "container_mounts" table.
	ContainerMountsColumns = []*schema.Column{
		{Name: "container_id", Type: field.TypeInt},
		{Name: "mount_id", Type: field.TypeInt},
	}
	// ContainerMountsTable holds the schema information for the "container_mounts" table.
	ContainerMountsTable = &schema.Table{
		Name:       "container_mounts",
		Columns:    ContainerMountsColumns,
		PrimaryKey: []*schema.Column{ContainerMountsColumns[0], ContainerMountsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "container_mounts_container_id",
				Columns:    []*schema.Column{ContainerMountsColumns[0]},
				RefColumns: []*schema.Column{ContainersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "container_mounts_mount_id",
				Columns:    []*schema.Column{ContainerMountsColumns[1]},
				RefColumns: []*schema.Column{MountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// ContainerPortsColumns holds the columns for the "container_ports" table.
	ContainerPortsColumns = []*schema.Column{
		{Name: "container_id", Type: field.TypeInt},
		{Name: "port_id", Type: field.TypeInt},
	}
	// ContainerPortsTable holds the schema information for the "container_ports" table.
	ContainerPortsTable = &schema.Table{
		Name:       "container_ports",
		Columns:    ContainerPortsColumns,
		PrimaryKey: []*schema.Column{ContainerPortsColumns[0], ContainerPortsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "container_ports_container_id",
				Columns:    []*schema.Column{ContainerPortsColumns[0]},
				RefColumns: []*schema.Column{ContainersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "container_ports_port_id",
				Columns:    []*schema.Column{ContainerPortsColumns[1]},
				RefColumns: []*schema.Column{PortsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// ContainerJobsColumns holds the columns for the "container_jobs" table.
	ContainerJobsColumns = []*schema.Column{
		{Name: "container_id", Type: field.TypeInt},
		{Name: "job_id", Type: field.TypeInt},
	}
	// ContainerJobsTable holds the schema information for the "container_jobs" table.
	ContainerJobsTable = &schema.Table{
		Name:       "container_jobs",
		Columns:    ContainerJobsColumns,
		PrimaryKey: []*schema.Column{ContainerJobsColumns[0], ContainerJobsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "container_jobs_container_id",
				Columns:    []*schema.Column{ContainerJobsColumns[0]},
				RefColumns: []*schema.Column{ContainersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "container_jobs_job_id",
				Columns:    []*schema.Column{ContainerJobsColumns[1]},
				RefColumns: []*schema.Column{JobsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// NodeContainersColumns holds the columns for the "node_containers" table.
	NodeContainersColumns = []*schema.Column{
		{Name: "node_id", Type: field.TypeInt},
		{Name: "container_id", Type: field.TypeInt},
	}
	// NodeContainersTable holds the schema information for the "node_containers" table.
	NodeContainersTable = &schema.Table{
		Name:       "node_containers",
		Columns:    NodeContainersColumns,
		PrimaryKey: []*schema.Column{NodeContainersColumns[0], NodeContainersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "node_containers_node_id",
				Columns:    []*schema.Column{NodeContainersColumns[0]},
				RefColumns: []*schema.Column{NodesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "node_containers_container_id",
				Columns:    []*schema.Column{NodeContainersColumns[1]},
				RefColumns: []*schema.Column{ContainersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UserContainersColumns holds the columns for the "user_containers" table.
	UserContainersColumns = []*schema.Column{
		{Name: "user_id", Type: field.TypeInt},
		{Name: "container_id", Type: field.TypeInt},
	}
	// UserContainersTable holds the schema information for the "user_containers" table.
	UserContainersTable = &schema.Table{
		Name:       "user_containers",
		Columns:    UserContainersColumns,
		PrimaryKey: []*schema.Column{UserContainersColumns[0], UserContainersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "user_containers_user_id",
				Columns:    []*schema.Column{UserContainersColumns[0]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "user_containers_container_id",
				Columns:    []*schema.Column{UserContainersColumns[1]},
				RefColumns: []*schema.Column{ContainersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ContainersTable,
		JobsTable,
		MountsTable,
		NodesTable,
		PortsTable,
		UsersTable,
		ContainerMountsTable,
		ContainerPortsTable,
		ContainerJobsTable,
		NodeContainersTable,
		UserContainersTable,
	}
)

func init() {
	ContainersTable.Annotation = &entsql.Annotation{}
	JobsTable.Annotation = &entsql.Annotation{}
	MountsTable.Annotation = &entsql.Annotation{}
	NodesTable.Annotation = &entsql.Annotation{}
	PortsTable.Annotation = &entsql.Annotation{}
	UsersTable.Annotation = &entsql.Annotation{}
	ContainerMountsTable.ForeignKeys[0].RefTable = ContainersTable
	ContainerMountsTable.ForeignKeys[1].RefTable = MountsTable
	ContainerPortsTable.ForeignKeys[0].RefTable = ContainersTable
	ContainerPortsTable.ForeignKeys[1].RefTable = PortsTable
	ContainerJobsTable.ForeignKeys[0].RefTable = ContainersTable
	ContainerJobsTable.ForeignKeys[1].RefTable = JobsTable
	NodeContainersTable.ForeignKeys[0].RefTable = NodesTable
	NodeContainersTable.ForeignKeys[1].RefTable = ContainersTable
	UserContainersTable.ForeignKeys[0].RefTable = UsersTable
	UserContainersTable.ForeignKeys[1].RefTable = ContainersTable
}
