package command_db

import (
	"net"
	"net/http"
	"net/rpc"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CommandNewsModel struct {
	ID          int64     `gorm:"type:smallint ; primaryKey"`
	Version     int64     `gorm:"type:smallint ; primaryKey"`
	Title       string    `gorm:"type:varchar(200)"`
	Content     string    `gorm:"type:varchar(50)"`
	Author      string    `gorm:"type:varchar(50)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	PublishedAt time.Time `gorm:"autoUpdateTime"`
	NewsType    string    `gorm:"type:varchar(50) ; default:famous"`
	Tags        []string  `gorm:"type:text[]"`
}

func InitializeCommandDB(uri string) (*gorm.DB, error) {
	connection, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (commandNewsModel *CommandNewsModel) RpcCommandConnection() error {
	err := rpc.Register(commandNewsModel)
	if err != nil {
		return err
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		return err
	}
	http.Serve(listener, nil)
	return nil
}

func (commandNewsModel *CommandNewsModel) Add_news(connection gorm.DB) error {
	err := connection.Create(&commandNewsModel).Error
	if err != nil {
		return err
	}
	return nil
}
