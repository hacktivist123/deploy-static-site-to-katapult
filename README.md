# Deploy a Static Site to Katapult Object Storage

This is a simple script written in Go that deploys a static site to [Katapult Object Storage](https://katapult.io/products/object-storage/) from it's build folder.

>[!NOTE]
> - Make sure to create a .env file to hold your `KATAPULT_ENDPOINT`, `KATAPULT_ACCESS_KEY_ID`, `KATAPULT_SECRET_KEY` & `BUCKET_NAME`.
> - Make sure to enable your bucket to serve static sites when creating a new bucket. 
