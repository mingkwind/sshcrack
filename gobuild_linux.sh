CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sshcrack
rm -rf sshcrack_linux
mkdir sshcrack_linux
mv sshcrack sshcrack_linux/
cp reScanStrs.conf sshcrack_linux/
touch sshcrack_linux/ip_list.txt
touch sshcrack_linux/user.txt
touch sshcrack_linux/pass.txt