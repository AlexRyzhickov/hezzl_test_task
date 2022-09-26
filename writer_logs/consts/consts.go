package consts

const (
	TableName = "logs"
)

const Ddl = `
		CREATE TABLE logs (
		    l JSON
		) ENGINE = Memory
		`

const Set_allow_experimental_object_type = "Set allow_experimental_object_type = 1"

const Drop_table = "DROP TABLE IF EXISTS logs"
