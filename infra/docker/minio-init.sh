#!/bin/sh

# Wait for MinIO to be ready
until mc alias set myminio http://minio:9000 minioadmin minioadmin123; do
  echo "Waiting for MinIO to be ready..."
  sleep 5
done

# Create bucket if it doesn't exist
mc mb myminio/gogym-photos 2>/dev/null || echo "Bucket already exists"

# Set bucket policy to public read
mc anonymous set download myminio/gogym-photos

echo "MinIO setup completed"