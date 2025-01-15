openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config csr.cnf
openssl x509 -req -in server.csr -signkey server.key -out server.crt -days 365 -extfile cert.cnf -extensions x509_ext
cat server.crt server.key > cert.pem
cp cert.pem ../cert.pem 
cp server.key ../key.pem
