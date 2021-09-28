package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

const (
	host = "localhost"
	port = 5432
	user = "anousone"
	password = "smpidpmr"
	dbname = "gorm_db"
)

type SqlLogger struct {
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error){
	sql, _ := fc()
	fmt.Printf("%v\n===================\n", sql)
}

var db *gorm.DB

func main(){

	DSN := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=Asia/Bangkok", host, port, user, password, dbname)
	gormDirector := postgres.Open(DSN)
	var err error
	db, err = gorm.Open(gormDirector, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false, // if true not perform in database but display sql statement
	})
	if err != nil {
		panic(err)
	}

	//? create table
	// err = db.Migrator().CreateTable(Customer{}) // this case need implement mannual
	// err = db.AutoMigrate(Gender{}, Test{}) // create table && auto migration(check table)
	// if err != nil {
	// 	fmt.Printf("error = %v", err)
	// }
	//? table gender
	//* insert into table
	// CreateGender("male")
	// CreateGender("female")
	//* get
	// GetGenders()
	// GetGender(3)
	// GetGenderByName("sone")
	//* update
	// UpdateGender(4, "daky") // sql 2 statement
	// UpdateGenderV2(4, "")  //! field not null or empty
	UpdateGenderV2(5, "males")
	//* delete
	// DeleteGender(1)
	// DeleteGender(3)

	//? table test
	//* create
	// CreateTest(2, "yourname")
	// CreateTest(0, "ans")
	// CreateTest(0, "phong")
	//* delete
	// DeleteTest(3)
	//* get
	// GetTests()

	//? table customer & relationsive
	//* create
	// CreateCustomer("anousone", 5) // genderId 5 is male, 6 is female
	// CreateCustomer("daky", 5)
	// CreateCustomer("girlname", 6)
	//* get
	GetCustomers()
}

func GetCustomers() {
	customers := []Customer{}
	// tx := db.Preload("Gender").Find(&customers) // manual relation to gender
	tx := db.Preload(clause.Associations).Find(&customers) // auto relation to gender
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	for _, v := range customers {
		fmt.Printf("%v|%v|%v\n",v.ID, v.Name, v.Gender.Name)	
	}
}

func CreateCustomer(name string, genderId uint) {
	customer := Customer{
		Name: name,
		GenderID: genderId,
	}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

func CreateTest(code int, name string) {
	test := Test{Code: code, Name: name}
	tx := db.Create(&test)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Created")
	fmt.Println(test)
}

func DeleteTest(id int) {
	tx := db.Delete(&Test{}, id)
	if tx.Error != nil {
		fmt.Println("tx.Error")
		return
	}
	fmt.Println("Deleted")
	GetTest(id)
}

func GetTests() {
	tests := []Test{}
	tx := db.Find(&tests)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	for _, v := range tests {
		fmt.Println(v)
	}
}

func GetTest(id int) {
	test := Test{}
	tx := db.First(&test, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(test)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println("tx.Error")
		return
	}
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id) // get
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name // update
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
}

func UpdateGenderV2(id uint, name string) {
	gender := Gender{Name: name}
	tx := db.Model(&Gender{}).Where("id=@myid", sql.Named("myid", id)).Updates(gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
}

func GetGenders(){
	gender := []Gender{}
	tx := db.Order("id").Find(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGender(id uint){
	gender := Gender{}
	tx := db.First(&gender, id) // default where by id
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenderByName(name string){
	gender := Gender{}
	tx := db.Where("name=?", name).First(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func CreateGender(name string){
	gender := Gender{Name: name}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

type Gender struct {
	ID uint
	Name string `gorm:"unique;size(10)"`
}

type Test struct {
	gorm.Model
	Code int
	Name string `gorm:"type:varchar(50);unique;default:Hello"`
}

type Customer struct {
	ID uint
	Name string
	Gender Gender
	GenderID uint
}