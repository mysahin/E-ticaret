package Controllers

import (
	database "ETicaret/Database"
	"ETicaret/Models"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

// FileController handles file operations.
type FileController struct {
	uploader   *s3manager.Uploader
	downloader *s3.S3
	bucketName string
}

var FileId uint = 0

func NewFileController(uploader *s3manager.Uploader, downloader *s3.S3, bucketName string) *FileController {
	if uploader == nil || downloader == nil {
		panic("uploader and downloader cannot be nil")
	}
	return &FileController{
		uploader:   uploader,
		downloader: downloader,
		bucketName: bucketName,
	}
}

// UploadFile uploads files to S3.
func (fc *FileController) UploadFile(c *fiber.Ctx, newFileName string) (string, error) {
	db := database.DB.Db
	form, err := c.MultipartForm()
	if err != nil {
		return "", c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var uploadedFile Models.Files
	files := form.File["files"]
	var uploadedURLs []string
	for _, file := range files {
		fileHeader := file
		f, err := fileHeader.Open()
		if err != nil {
			return "", c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer f.Close()

		uploadedURL, err := fc.saveFile(f, fileHeader.Filename, newFileName)
		if err != nil {
			return "", c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		uploadedURLs = append(uploadedURLs, uploadedURL)
		if err := db.Create(&uploadedFile).Error; err != nil {
			return "", err
		}
		FileId = uploadedFile.ID
	}

	return uploadedURLs[0], c.Status(http.StatusOK).JSON(fiber.Map{"urls": uploadedURLs})
}

// fixFileName replaces special characters in filenames and ensures it's an image file.
func fixFileName(filename string, productId string) (string, error) {
	if len(filename) >= 4 {
		lastFour := filename[len(filename)-4:]
		if lastFour == ".jpg" {
			filename = productId + ".jpg"
		} else if lastFour == ".png" {
			filename = productId + ".png"
		} else if lastFour == ".gif" {
			filename = productId + ".gif"
		} else if lastFour == ".bmp" {
			filename = productId + ".bmp"
		} else {
			return "", errors.New("Invalid file extension")
		}

	} else {
		fmt.Println("String uzunluğu 4'ten küçük.")
	}
	return filename, nil
}

// saveFile uploads a file to S3 and returns the URL.
func (fc *FileController) saveFile(fileReader io.Reader, filename string, newFileName string) (string, error) {
	newName, errr := fixFileName(filename, newFileName)
	if errr != nil {
		return "", errr
	}
	_, err := fc.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(fc.bucketName),
		Key:    aws.String(newName),
		Body:   fileReader,
	})
	if err != nil {
		return "", err
	}

	// Get the URL of the uploaded file
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", fc.bucketName, filename)

	return url, nil
}

// ListFiles lists all files in the S3 bucket.
func (fc *FileController) ListFiles(c *fiber.Ctx) error {
	resp, err := fc.downloader.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(fc.bucketName),
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var filenames []string
	for _, item := range resp.Contents {
		filenames = append(filenames, *item.Key)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"filenames": filenames})
}

// ShowFile retrieves and sends a file from S3.
func (fc *FileController) ShowFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	obj, err := fc.downloader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(fc.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer obj.Body.Close()

	// Determine the content type based on the file extension
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		// If the content type is not recognized, default to octet-stream
		contentType = "application/octet-stream"
	}

	// Set the content type header
	c.Set("Content-Type", contentType)

	content, err := io.ReadAll(obj.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).Send(content)
}

// DeleteFile deletes a file from S3.
func (fc *FileController) DeleteFile(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Delete the file from S3
	_, err := fc.downloader.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fc.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("File '%s' successfully deleted.\n", filename)

	return c.SendStatus(http.StatusOK)
}

func GetFileId() uint {
	fileId := FileId
	return fileId
}
