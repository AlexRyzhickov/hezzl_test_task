package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/nats-io/nats.go"
	"hezzl_test_task/writer_logs/consts"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func createLogsTable(conn driver.Conn) error {
	if err := conn.Exec(context.Background(), consts.Drop_table); err != nil {
		return err
	}
	if err := conn.Exec(context.Background(), consts.Set_allow_experimental_object_type); err != nil {
		return err
	}
	if err := conn.Exec(context.Background(), consts.Ddl); err != nil {
		return err
	}
	return nil
}

func AsyncInsert(conn driver.Conn, log string) error {
	ctx := clickhouse.Context(context.Background(), clickhouse.WithStdAsync(false))
	return conn.Exec(ctx, fmt.Sprintf("INSERT INTO %s VALUES ('%s')", consts.TableName, log))
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	if err != nil {
		log.Fatal(err)
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = createLogsTable(conn)
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe("foo", func(m *nats.Msg) {
		err := AsyncInsert(conn, string(m.Data))
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start writing logs in ClickHouse")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	<-shutdown
}
