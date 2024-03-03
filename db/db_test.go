package db

import (
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type User struct {
	gorm.Model
	Profiles []Profile `gorm:"many2many:user_profiles;foreignKey:Refer;joinForeignKey:UserReferID;References:UserRefer;joinReferences:ProfileRefer"`
	Refer    uint      `gorm:"index:,unique"`
}

type Profile struct {
	gorm.Model
	Name      string
	UserRefer uint `gorm:"index:,unique"`
}

type Person struct {
	NodeID    int
	Name      string
	Addresses []Address `gorm:"many2many:person_addresses;"`
}

type Address struct {
	NodeID uint
	Name   string
}

type PersonAddress struct {
	PersonNodeID  int `gorm:"primaryKey;foreignKey:PersonNodeID;joinForeignKey:PersonNodeID"`
	AddressNodeID int `gorm:"primaryKey;foreignKey:AddressNodeID;joinForeignKey:AddressNodeID"`
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func TestCreateRepo(t *testing.T) {
	var err error
	DB, err = gorm.Open(mysql.Open("root:114514@tcp(localhost:3306)/openalysis?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(
		&model.Contributor{},
		&model.Group{},
		&model.Issue{},
		&model.Organization{},
		&model.PullRequest{},
		&model.Repository{},
		&model.GroupsOrganizations{},
		&model.GroupsRepositories{},
		&model.ContributorsIssues{},
		&model.ContributorsPullRequests{},
	)

	//err = DB.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{})
	//DB.AutoMigrate(&Person{}, &Address{})

	//DB.Create(&Profile{Name: "lorain", UserRefer: 13})
	//DB.Create(&Profile{Name: "lorain", UserRefer: 13})

	// 修改 Person 的 Addresses 字段的连接表为 PersonAddress
	// PersonAddress 必须定义好所需的外键，否则会报错
	if err != nil {
		fmt.Println(err)
	}

	//DB.Create(&Student{Name: "jack", Teachers: []Teacher{{Name: "mark"}}})
	//rows, err := CreateRepository(&model.Repository{
	//	Owner:            "cloudwego",
	//	Name:             "kitex",
	//	IssueCount:       500,
	//	PullRequestCount: 400,
	//	StarCount:        300,
	//	ForkCount:        200,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(rows)
}
