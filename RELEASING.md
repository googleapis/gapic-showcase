# Releasing
To make a new release of gapic-showcase the following steps should be performed.

1. Create a new branch of the version
```sh
git checkout master
git pull origin master
git checkout -b v{VERSION}
```

2. The semantic version of this package as well as the API is used across multiple files in gapic-showcase. This needs the be updated. Use the utility script found in util/cmd/bump_version to bump the version across all of the files. Pass the `--patch`, `--minor`, `--major`, `--api` flags to this script depending which version needs to be updated.
```sh
go run ./util/cmd/bump_version --{TYPE}
```

3. Update the CHANGELOG.md file with the changes made since the last release including the version in the heading for the changes. Please note that the version must start with the character `v` in order to match the git tag that will be pushed in a later step. You can get the list of PRs from your local git commit history on the `master` branch, or from GitHub by using the the query `is:pr merged:>=YYYY-MM-DD` (where the date refers to the _last_ release date).

4. Create a pull request for this release merging the version branch into master.

5. Change branches to master and pull.
```sh
git checkout master
git pull origin master
```

6. Create a tag for the version and push. The automated release will take over from here. Please note that version tags should start with the character `v`.
```sh
git tag v{VERSION}
git push origin v{VERSION}
```

7. Build a Docker image.
```sh
docker build -t gcr.io/gapic-images/gapic-showcase:latest .
```

8. Ensure your Docker daemon is authenticated with your GCP credentials.
```sh
gcloud auth login
gcloud auth configure-docker
```

9. Tag with the release version and push both the `latest` and the versioned images.
```sh
docker tag gcr.io/gapic-images/gapic-showcase:latest gcr.io/gapic-images/gapic-showcase:{VERSION}
docker push gcr.io/gapic-images/gapic-showcase:lastet
docker push gcr.io/gapic-images/gapic-showcase:{VERSION}
```
