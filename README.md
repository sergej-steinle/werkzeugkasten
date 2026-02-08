# Werkzeugkasten

A Go utility library providing common helper functions for web applications.

## Installation

```bash
go get github.com/sergej-steinle/werkzeugkasten
```

## Features

### Werkzeug

The `Werkzeug` struct provides file and string utilities:

- **RandomString(n int)** - Generates a random string of length `n`
- **UploadOneFile / UploadFiles** - Handles multipart file uploads with optional renaming and content-type validation
- **CreateDir(path string)** - Creates a directory (including parents) if it doesn't exist
- **DownloadStaticFile** - Serves a file as a downloadable attachment

```go
w := werkzeugkasten.Werkzeug{
    MaxFileSize:      10 * 1024 * 1024, // 10 MB
    AllowedFileTypes: []string{"image/png", "image/jpeg"},
}

// Generate a random string
s := w.RandomString(16)

// Upload a single file from an HTTP request
file, err := w.UploadOneFile(r, "./uploads")
```

### Validator

The `Validator` struct provides form/input validation:

- **Valid()** - Returns `true` if there are no field errors
- **AddFieldError(key, message)** - Adds a validation error for a field
- **CheckField(ok, key, message)** - Conditionally adds a field error

Standalone validation helpers:

- **NotBlank(value)** - Checks that a string is not empty/whitespace
- **MaxChars(value, n)** - Checks that a string does not exceed `n` characters
- **PermittedValue(value, ...permitted)** - Checks that a value is in an allowed set
- **Matches(value, regexp)** - Checks that a string matches a regex pattern
- **IsEmail(value)** - Validates an email address

```go
v := werkzeugkasten.Validator{}

v.CheckField(werkzeugkasten.NotBlank(name), "name", "Name is required")
v.CheckField(werkzeugkasten.MaxChars(name, 100), "name", "Name must be at most 100 characters")
v.CheckField(werkzeugkasten.IsEmail(email), "email", "Must be a valid email address")

if !v.Valid() {
    // handle v.FieldErrors
}
```

## License

This project is licensed under the [MIT License](LICENSE).
