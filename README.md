# Greeter CLI

A simple command-line application written in Go that prompts the user for their name and greets them multiple times based on the provided input.

This project was built as part of learning how to design, implement, and test CLI applications in Go.

---

## 🚀 Features

- Parse command-line arguments
- Validate user input
- Interactive input using stdin
- Configurable output repetition
- Clean separation of concerns
- Fully unit tested (table-driven tests)

---

## 🧠 How It Works

1. User provides a number as an argument
2. CLI prompts for user's name
3. Prints greeting message `n` times

---

## 📦 Usage

`go run main.go <number>`

### Example

`go run main.go 2`

Output:

    Your Name please? Press the Enter key when done.
    Ankur
    Nice to meet you Ankur
    Nice to meet you Ankur

---

## ❗ Help Flag

`go run main.go -h`

Displays usage instructions.

---

## ⚠️ Error Handling

The application handles:

- Missing arguments
- Invalid number format
- Negative or zero values
- Empty name input

---

## 🏗️ Project Structure

    .
    ├── main.go        # Entry point
    ├── main_test.go   # Unit tests

---

## 🧪 Testing

Run tests using:

`go test ./...`

### Test Coverage Includes

- Argument parsing
- Input validation
- CLI execution flow
- User input handling

---

## 🛠️ Design Decisions

- Used `io.Reader` and `io.Writer` for better testability
- Followed table-driven testing (standard Go practice)
- Separated parsing, validation, and execution logic

---

## 📚 What I Learned

- Building CLI applications in Go
- Writing clean and testable code
- Table-driven testing patterns
- Handling user input via stdin
- Structuring small projects properly

---

## 🔮 Future Improvements

- Add support for flags using `flag` package
- Improve help/usage formatting
- Add colored output
- Package as installable CLI tool

---

## 👨‍💻 Author

Ankur Dwivedi