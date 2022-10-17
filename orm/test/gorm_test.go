package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/orm"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	config.LoadYamlFile("./application-test1.yaml")
	db, _ := orm.GetGormDb()
	db.Exec("drop table gobase_demo")
	//db.Exec("CREATE TABLE gobase_demo if not exist (\n  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n  `name` char(20) NOT NULL COMMENT '名字',\n  `age` INT NOT NULL COMMENT '年龄',\n  `address` char(20) NOT NULL COMMENT '名字',\n  \n  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")

	// 新增
	//db.Create(&GobaseDemo{Name: "zhou", Age: 18, Address: "杭州"})


	// 查询：一行
	var demo GobaseDemo
	db.Select("Name", "Age", "CreatedAt").Create(&demo)

	// 查询：多行
	fmt.Println(demo)

	// 查询：分页

	// 查询：一个值

	// 查询：个数
}

type GobaseDemo struct {
	Id         uint64
	Name       string
	Age        int
	Address    string
	CreateTime time.Time
	UpdateTime time.Time
}

func Test2(t *testing.T) {

}
