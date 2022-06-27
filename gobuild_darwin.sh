go build -o sshcrack
rm -rf sshcrack_darwin
mkdir sshcrack_darwin
mv sshcrack sshcrack_darwin/
cp reScanStrs.conf sshcrack_darwin/
touch sshcrack_darwin/ip_list.txt
touch sshcrack_darwin/user.txt
touch sshcrack_darwin/pass.txt