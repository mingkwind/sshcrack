CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sshcrack
rm -rf sshpwd_crack
mkdir sshpwd_crack
cp sshcrack sshcrack/
touch sshcrack/ip_list.txt
touch sshcrack/user.txt
touch sshcrack/pass.txt