

```golang

package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
    accessKey = "添加自己的accessKey"
    secretKey = "添加自己的secretKey"
    region    = "添加自己的region"
)

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

// S3CreateBucket ...
func S3CreateBucket(bucket string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region),
    })

    svc := s3.New(sess)
    _, err := svc.CreateBucket(&s3.CreateBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Unable to create bucket %q, %v", bucket, err)
    }

    err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Error occurred while waiting for bucket to be created, %v", err)
    }

    fmt.Printf("Bucket %q successfully created\n", bucket)
}

// S3ListBuckets ...
func S3ListBuckets() {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region),
    })

    svc := s3.New(sess)
    result, err := svc.ListBuckets(nil)
    if err != nil {
        exitErrorf("Unable to list buckets, %v", err)
    }

    for _, b := range result.Buckets {
        fmt.Printf("%s created on %s\n", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
    }
}

// S3DeleteBucket ...
func S3DeleteBucket(bucket string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    svc := s3.New(sess)
    _, err := svc.DeleteBucket(&s3.DeleteBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Unable to delete bucket %q, %v", bucket, err)
    }

    err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Error occurred while waiting for bucket to be deleted, %v", err)
    }

    fmt.Printf("Bucket %q successfully deleted\n", bucket)
}

// S3PutObject ....
func S3PutObject(bucket string, key string, value string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    uploader := s3manager.NewUploader(sess)
    _, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
        Body:   strings.NewReader(value),
    })
    if err != nil {
        exitErrorf("Unable to upload %q to %q, %v", key, bucket, err)
    }

    fmt.Printf("Successfully uploaded %q to %q\n", key, bucket)
}


// S3PutObjectFile ....
func S3PutObjectFile(bucket string, filename string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    file, err := os.Open(filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	
	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(filename),
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}
	
	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}



// S3ListObjects ...
func S3ListObjects(bucket string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    svc := s3.New(sess)
    resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
    if err != nil {
        exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
    }

    for _, item := range resp.Contents {
        fmt.Println("Name:         ", *item.Key)
        fmt.Println("Last modified:", *item.LastModified)
        fmt.Println("Size:         ", *item.Size)
        fmt.Println("Storage class:", *item.StorageClass)
        fmt.Println("")
    }

    fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
}


// S3GetObjectAsBytes ...
func S3GetObjectAsBytes(bucket string, key string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    downloader := s3manager.NewDownloader(sess)
    buffer := aws.NewWriteAtBuffer([]byte{})
    _, err := downloader.Download(buffer,
        &s3.GetObjectInput{
            Bucket: aws.String(bucket),
            Key:    aws.String(key),
        })
    if err != nil {
        exitErrorf("Unable to download key %q, %v", key, err)
    }

    fmt.Println("Downloaded", string(buffer.Bytes()))
}


// S3GetObjectAsFile ...
func S3GetObjectAsFile(bucket string, key string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    file, err := os.Create(key)
    if err != nil {
        exitErrorf("Unable to open file %q, %v", err)
    }
    defer file.Close()

    downloader := s3manager.NewDownloader(sess)

    numBytes, err := downloader.Download(file,
    &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
	if err != nil {
	    exitErrorf("Unable to download key %q, %v", key, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}


// S3DeleteObject ....
func S3DeleteObject(bucket string, key string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    svc := s3.New(sess)
    _, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
    if err != nil {
        exitErrorf("Unable to delete object %q from bucket %q, %v", key, bucket, err)
    }

    err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        exitErrorf("Error occurred while waiting for object %q to be deleted, %v", key, err)
    }

    fmt.Printf("Object %q successfully deleted\n", key)
}


// S3DeleteObjects ...
func S3DeleteObjects(bucket string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    svc := s3.New(sess)
    iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
        Bucket: aws.String(bucket),
    })
    if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
        exitErrorf("Unable to delete objects from bucket %q, %v", bucket, err)
    }

    fmt.Printf("Deleted object(s) from bucket: %q\n", bucket)
}


// S3CopyObjects ...
func S3CopyObjects(bucket string, other string, key string) {
    sess, _ := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region)},
    )

    svc := s3.New(sess)
    source := bucket + "/" + key
    _, err := svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(other), CopySource: aws.String(source), Key: aws.String(key)})
    if err != nil {
        exitErrorf("Unable to copy key from bucket %q to bucket %q, %v", bucket, other, err)
    }

    err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(other), Key: aws.String(key)})
    if err != nil {
        exitErrorf("Error occurred while waiting for key %q to be copied to bucket %q, %v", bucket, key, other, err)
    }

    fmt.Printf("Key %q successfully copied from bucket %q to bucket %q\n", key, bucket, other)
}



func main() {
    S3CreateBucket("owb-test-1")
    S3CreateBucket("owb-test-2")
    fmt.Println("-------------------------------------------------")

    S3ListBuckets()
    fmt.Println("-------------------------------------------------")

    S3PutObject("owb-test-1", "key1", "value1")
    S3PutObject("owb-test-1", "key2", "value2")
    S3PutObject("owb-test-1", "key3", "value3")
    fmt.Println("-------------------------------------------------")

    S3GetObject("owb-test-1", "key1")
    S3GetObject("owb-test-1", "key2")
    S3GetObject("owb-test-1", "key3")
    fmt.Println("-------------------------------------------------")

    S3CopyObjects("owb-test-1", "owb-test-2", "key1")
    S3CopyObjects("owb-test-1", "owb-test-2", "key2")
    S3CopyObjects("owb-test-1", "owb-test-2", "key3")
    fmt.Println("-------------------------------------------------")

    S3ListObjects("owb-test-1")
    S3ListObjects("owb-test-2")
    fmt.Println("-------------------------------------------------")

    S3DeleteObject("owb-test-1", "key1")
    S3DeleteObject("owb-test-1", "key2")
    S3DeleteObject("owb-test-1", "key3")
    fmt.Println("-------------------------------------------------")

    S3DeleteObjects("owb-test-2")
    fmt.Println("-------------------------------------------------")

    S3DeleteBucket("owb-test-1")
    S3DeleteBucket("owb-test-2")
}


```