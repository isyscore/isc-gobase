package test

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	orm2 "github.com/isyscore/isc-gobase/extend/orm"
	"testing"
)

func TestGorm1(t *testing.T) {
	config.LoadYamlFile("./application-test1.yaml")
	orm2.AddGormHook(&GobaseOrmHookDemo{})
	db, _ := orm2.NewGormDb()

	// 删除表
	db.Exec("drop table isc_demo.gobase_demo")

	//新增表
	db.Exec("CREATE TABLE gobase_demo(\n" +
		"  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n" +
		"  `name` char(20) NOT NULL COMMENT '名字',\n" +
		"  `age` INT NOT NULL COMMENT '年龄',\n" +
		"  `address` char(20) NOT NULL COMMENT '名字',\n" +
		"  \n" +
		"  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n" +
		"  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n" +
		"\n" +
		"  PRIMARY KEY (`id`)\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")

	// 新增
	db.Create(&GobaseDemo{Name: "zhou", Age: 18, Address: "杭州"})
	db.Create(&GobaseDemo{Name: "zhou", Age: 11, Address: "杭州2"})

	// 查询：一行
	var demo GobaseDemo
	db.First(&demo).Where("name=?", "zhou")

	dd, _ := db.DB()
	dd.Query("select * from gobase_demo")

	// 查询：多行
	fmt.Println(demo)
}

type GobaseDemo struct {
	Id         uint64
	Name       string
	Age        int
	Address    string
}

func (GobaseDemo) TableName() string {
	return "gobase_demo"
}

type GobaseOrmHookDemo struct {
}

func (*GobaseOrmHookDemo) Before(ctx context.Context, parameters map[string]any) (context.Context, error){
	fmt.Println("before")
	fmt.Println(parameters)
	return ctx, nil
}

func (*GobaseOrmHookDemo) After(ctx context.Context, parameters map[string]any) (context.Context, error){
	fmt.Println("after")
	fmt.Println(parameters)
	return ctx, nil
}

func (*GobaseOrmHookDemo) Err(ctx context.Context, err error, parameters map[string]any) error {
	fmt.Println(parameters)
	return nil
}


