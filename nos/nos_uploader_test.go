package nos

import (
	"fmt"
	"github.com/NetEase-Object-Storage/nos-golang-sdk/config"
	"github.com/NetEase-Object-Storage/nos-golang-sdk/nosclient"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	. "github.com/wq1019/go-image_uploader"
	"log"
	"os"
	"testing"
	"time"
)

var (
	Endpoint   = "nos-eastchina1.126.net"
	AccessKey  = "5fdde629014441f080b10d3f0299f85e"
	SecretKey  = "db67cb47b58d4bcb832684a17c97f29a"
	BucketName = "zm-dev"
	uploader   Uploader
)

func TestMain(m *testing.M) {
	nosClient, err := nosclient.New(&config.Config{
		Endpoint:  Endpoint,
		AccessKey: AccessKey,
		SecretKey: SecretKey,
	})
	if err != nil {
		log.Fatalf("nos client 创建失败! error: %+v", err)
	}
	store := NewDBStore(setupGorm())

	uploader = NewNosUploader(
		HashFunc(MD5HashFunc),
		store,
		nosClient,
		BucketName,
		Hash2StorageNameFunc(TwoCharsPrefixHash2StorageNameFunc),
	)
	m.Run()
}

func TestNosUploader_Upload(t *testing.T) {
	filename := "./../testdata/Go-Logo_Black.jpg"
	fi, err := os.Stat(filename)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	image, err := uploader.Upload(FileHeader{Filename: file.Name(), Size: fi.Size(), File: file})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(image.Hash)
}

func setupGorm() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 10; i++ {
		// db, err = gorm.Open("sqlite3", "file::memory:?cache=shared")
		db, err = gorm.Open("sqlite3", "cloud-images.db")
		if err == nil {
			autoMigrate(db)
			return db
		}
		log.Println(err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("数据库链接失败！ error: %+v", err)
	return nil
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Image{},
	).Error
	if err != nil {
		log.Fatalf("AutoMigrate 失败！ error: %+v", err)
	}
}
