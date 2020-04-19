- Post sample image:  curl -X POST -H "Content-Type: image/png" --data-binary "@maru.png" https://mcs95iikje.execute-api.ap-northeast-1.amazonaws.com/staging/PostImage 

- update lambda function: aws lambda update-function-code --function-name InsertMovie --zip-file fileb://./deployment.zip --region ap-northeast-1