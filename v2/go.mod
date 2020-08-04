module gitub.com/spark451/snowblower/v2

go 1.13

replace github.com/spark451/snowman/v2 => ../../snowman/v2

require (
	github.com/aws/aws-sdk-go v1.33.4
	github.com/duncan/base64x v0.0.0-20150429221403-a119b4bf1ecd
	github.com/google/uuid v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/remeh/sizedwaitgroup v1.0.0
	github.com/spark451/snowman/v2 v2.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.0.0
	github.com/ua-parser/uap-go v0.0.0-20200325213135-e1c09f13e2fe
	github.com/xeipuuv/gojsonschema v1.2.0
	go.mongodb.org/mongo-driver v1.3.5
)
