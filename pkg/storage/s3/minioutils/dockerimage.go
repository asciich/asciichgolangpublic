package minioutils

// MINIO_DOCKER_IMAGE is the container image used to start a local MinIO
// server, for example in tests and examples.
//
// Upstream MinIO stopped publishing Community Edition images in October 2025.
// Since then, anonymous pulls of minioutils.MINIO_DOCKER_IMAGE and minio/minio fail
// with "unauthorized". The upstream repository github.com/minio/minio is
// archived, and its last release is RELEASE.2025-10-15T17-29-55Z.
//
// This image is a community build of that final upstream release. It won't
// get any more updates, including security fixes, so only use it for local
// development and testing, never in production.
//
// The tag is pinned on purpose. Don't replace it with "latest".
//
// Alternatives if this image becomes unavailable:
//   - pgsty/minio: an actively maintained community fork
//   - a self-built image from the upstream Dockerfile, pushed to a registry
//     under our control
const MINIO_DOCKER_IMAGE = "coollabsio/minio:RELEASE.2025-10-15T17-29-55Z"
