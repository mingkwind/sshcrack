CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sshcrack
rm -rf sshpwd_crack
mkdir sshpwd_crack
cp sshcrack sshpwd_crack
touch sshpwd_crack/ip_list.txt
touch sshpwd_crack/user.txt
touch ssshpwd_crack/pass.txt