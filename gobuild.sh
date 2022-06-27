CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sshcrack
rm -rf sshpwd_crack_linux
mkdir sshpwd_crack_linux
cp sshcrack reScanStrs.conf sshpwd_crack_linux/
touch sshpwd_crack_linux/ip_list.txt
touch sshpwd_crack_linux/user.txt
touch sshpwd_crack_linux/pass.txt