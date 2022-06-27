CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sshcrack
rm -rf sshpwd_crack
mkdir sshpwd_crack
cp sshcrack reScanStrs.conf sshpwd_crack/
touch sshpwd_crack/ip_list.txt
touch sshpwd_crack/user.txt
touch sshpwd_crack/pass.txt