package helper

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func SetupS3Uploader(file *multipart.FileHeader) (string, string, error) {
	openFile, err := file.Open()
	if err != nil {
		return "", "", err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", "", err
	}

	// Buat koneksi ke AWS S3
	client := s3.NewFromConfig(cfg)

	// membuat objek uploader yang akan digunakan untuk mengunggah file ke bucket S3
	uploader := manager.NewUploader(client)

	// Upload file ke bucket S3
	uploadOutput, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("manajemen-tugas"),
		Key:    aws.String(file.Filename),
		Body:   openFile,
		ACL:    "public-read",
	})

	if err != nil {
		return "", "", err
	}

	return uploadOutput.Location, file.Filename, nil
}

func SetupS3Delete(filename string) error {
	if filename == "" {
		return errors.New("Parameter key diperlukan") // kembalikan error yang bukan bertipe error
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("Error loading AWS config: %w", err) // kembalikan error yang bertipe error
	}

	// Buat koneksi ke AWS S3
	svc := s3.NewFromConfig(cfg)

	// Hapus file dari bucket S3
	_, err = svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String("manajemen-tugas"),
		Key:    aws.String(filename),
	})

	if err != nil {
		return fmt.Errorf("Error deleting object: %w", err) // kembalikan error yang bertipe error
	}

	return nil
}

func SetupS3DeleteAll() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("Error loading AWS config: %w", err)
	}

	// Buat koneksi ke AWS S3
	svc := s3.NewFromConfig(cfg)

	// Panggil ListObjectsV2 API untuk mendapatkan daftar objek dalam bucket
	resp, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("manajemen-tugas"),
	})
	if err != nil {
		return fmt.Errorf("Error listing objects: %w", err)
	}

	// Iterasi setiap objek dan hapus
	for _, obj := range resp.Contents {
		_, err := svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String("manajemen-tugas"),
			Key:    obj.Key,
		})
		if err != nil {
			return fmt.Errorf("Error deleting object: %w", err)
		}
	}

	return nil
}

//func SetupS3GetAllFiles() (map[string]string, map[string]string, error) {
//	cfg, err := config.LoadDefaultConfig(context.TODO())
//	if err != nil {
//		return nil, nil, fmt.Errorf("Error loading AWS config: %w", err)
//	}
//
//	// Buat koneksi ke AWS S3
//	svc := s3.NewFromConfig(cfg)
//
//	// Panggil metode ListObjectsV2 untuk mendapatkan daftar objek di bucket
//	response, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
//		Bucket: aws.String("manajemen-tugas"),
//	})
//
//	if err != nil {
//		return nil, nil, fmt.Errorf("Error listing objects: %w", err)
//	}
//
//	urls := make(map[string]string)
//	fileNames := make(map[string]string)
//
//	// Loop melalui daftar objek dan tambahkan nama dan URL file ke dalam map files
//	for _, item := range response.Contents {
//		fileName := *item.Key
//		// Ubah spasi dalam nama file menjadi "%20" dalam URL
//		encodedFileName := url.PathEscape(fileName)
//		url := fmt.Sprintf("https://manajemen-tugas.s3.ap-southeast-3.amazonaws.com/%s", encodedFileName)
//		urls[fileName] = url
//		fileNames[url] = fileName
//	}
//
//	return urls, fileNames, nil
//}
