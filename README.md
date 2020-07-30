# GO FEEDER
Read document from Google Drive and sent to API server

### Prerequisites
- Docker

### Installing

Download repository
```
git clone https://github.com/amelendres/go-feeder
```

Set your env vars
```
cp .env.dist .env
```

Build container
```
make build
```

### HOW TO RUN 

**ENDPOINTS**

* Parse a Docx from Google Drive
1. Set your GOOGLE_API_KEY in your .env
2. Enable to share link in your drive
3. Copy the fileId in your payload as fileUrl field
4. Execute the curl request
```
curl --location --request POST 'http://localhost:8050/devotionals/parse' \
--header 'Content-Type: application/json' \
--data-raw '{
    "fileUrl": "1OA90lU_VuOStjvDKrjb2hJtFSKcLZCmq",
    "planId": "23a63256-f264-4d94-b7ed-8ce60f744ae3",
    "authorId": "9158becf-6f89-4366-9541-ae5b99689cc2",
    "publisherId": "2e62bcd1-b639-49fd-950b-9c2a937b07a5"
}'
```

* Import Devotionals from document
1. Set your DEVOM_API_URL
2. Execute the curl request updating your fileId in your payload as fileUrl field 
```
curl --location --request POST 'http://localhost:8050/devotionals/import' \
--header 'Content-Type: application/json' \
--data-raw '{
    "fileUrl": "1OA90lU_VuOStjvDKrjb2hJtFSKcLZCmq",
    "planId": "23a63256-f264-4d94-b7ed-8ce60f744ae3",
    "authorId": "9158becf-6f89-4366-9541-ae5b99689cc2",
    "publisherId": "2e62bcd1-b639-49fd-950b-9c2a937b07a5"
}'
```


## Authors

* **Alfredo Melendres** -  alfredo.melendres@gmail.com

[MIT license](LICENSE.md)
