package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

var md5hash = "77f62e3524cd583d698d51fa24fdff4f"
var sha256hash = "95a5e1547df73abdd4781b6c9e55f3377c15d08884b11738c2727dbd887d4ced"

var bCryptHash = "$2a$10$Zs3ZwsjV/nF.KuvSUE.5WuwtDrK6UVXcBpQrH84V8q3Opg1yNdWLu"

func makeBcryptHash(password string) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("hash = %s\n", hash)

	err = bcrypt.CompareHashAndPassword([]byte(bCryptHash), []byte(password))
	if err != nil {
		log.Println("[!] Authentication Failed")
		return
	}

	log.Println("[+] Authentication Successful")
}

func main() {
	f, err := os.Open("D:\\OffsecTools\\SecLists\\Passwords\\Leaked-Databases\\rockyou.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		password := scanner.Text()
		hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		if hash == md5hash {
			fmt.Printf("[+] Password Successfully Cracked (MD5): %s\n", password)
		}
		hash = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		if hash == sha256hash {
			fmt.Printf("[+] Password Found (SHA256): %s\n", password)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	makeBcryptHash("someC0mpl3xP@ssw0rd")
}
