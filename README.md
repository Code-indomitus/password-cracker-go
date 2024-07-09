# password-cracker-go

This project is an implementation of a password cracker for the MD5 hash algorithm. The md5 hash algorithm has also been implemented from scratch. You can brute force a four-character password as well as use a word list to crack an MD5 hashed password using this password cracker.


## Usage

### Command-Line Arguments

To use the password cracker, run the program with the following command-line arguments:

#### Brute force passwords of length four
```bash
go run passcracker.go <hash>

# Example: 
go run .\passcracker.go 2bdb742fc3d075ec6b73ea414f27819a
```

#### Use word list to 
```bash
go run passcracker.go -l <hash> <wordlist_path>

# Example:
go run .\passcracker.go -l 2bdb742fc3d075ec6b73ea414f27819a .\word-lists\realhuman_phill.txt
```