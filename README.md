# S3Utils Package

A Go package that provides utility functions for AWS S3 operations with enhanced credential management.

## Features

- Flexible AWS credential management with priority handling
- File upload to S3 with automatic session management
- File existence checking in S3 buckets
- Unique filename generation for S3 uploads
- Support for both temporary session tokens and profile-based credentials

## Credential Priority

The package implements a credential priority system:

1. Environment Variables (Session Token) - Highest Priority
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`
   - `AWS_SESSION_TOKEN`
2. AWS Profile (e.g., default, custom profiles) - Fallback Option

## Usage

### Initialize AWS Session

```go
// Create a new AWS session with automatic credential management
sess, err := s3utils.NewAWSSession("us-west-2", "default")
if err != nil {
    log.Fatal(err)
}
```

### Upload File to S3

```go
err := s3utils.UploadToS3(
    "us-west-2",      // AWS Region
    "default",        // AWS Profile (used if session token not available)
    "myfile.txt",     // Local file path
    "my-bucket",      // S3 bucket name
    "folder/path"     // S3 folder path
)
if err != nil {
    log.Fatal(err)
}
```

### Check File Existence in S3

```go
exists, err := s3utils.CheckS3FileExists(sess, "my-bucket", "folder/file.txt")
if err != nil {
    log.Fatal(err)
}
if exists {
    fmt.Println("File exists in S3")
}
```

### Generate Unique Filename

```go
uniqueName, err := s3utils.GenerateUniqueFileName(
    sess,
    "my-bucket",
    "folder/path",
    "original.txt"
)
if err != nil {
    log.Fatal(err)
}
```

## Using with Session Token

To use temporary credentials with a session token:

1. Set the required environment variables:
```bash
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"
export AWS_SESSION_TOKEN="your_session_token"
```

2. The package will automatically detect and use these credentials.

## Using with AWS Profile

If session token credentials are not available, the package will automatically fall back to using the specified AWS profile:

1. Ensure your AWS credentials are properly configured in `~/.aws/credentials`:
```ini
[default]
aws_access_key_id = your_access_key
aws_secret_access_key = your_secret_key

[custom-profile]
aws_access_key_id = your_access_key
aws_secret_access_key = your_secret_key
```

2. Use the profile name in your code:
```go
sess, err := s3utils.NewAWSSession("us-west-2", "custom-profile")
```

## Error Handling

The package provides detailed error messages for common issues:
- Missing or invalid credentials
- S3 bucket access issues
- File operation errors
- Session creation failures

## Dependencies

- github.com/aws/aws-sdk-go

## Best Practices

1. Always use temporary credentials (session tokens) when possible
2. Keep AWS profiles as a fallback option
3. Handle errors appropriately in your application
4. Use appropriate S3 bucket permissions
5. Clean up temporary credentials after use

## License

This package is distributed under the MIT license.