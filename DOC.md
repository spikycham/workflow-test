# CI/CD Practice

## What I'm Going to Do?

I'm trying to make a CI/CD pipeline from scratch, I never did this before. By the way, CI stands for Continuous Integration, and CD stands for Continuous Deployment. Here I'm going to analyze what I need to do first.

Our goal of building a CI/CD pipeline is to make the project deployment easier. Like we just need to push our code to remote repository, the project is going to be tested and deployed automatically. And this is how it looks like overall: `Push -> Deploy`. Let's split this goal into steps. Right after we push the code to remote repository (GitHub), what should happen next.

GitHub Actions
1. install deps
2. test
3. build
4. docker build
5. docker push
6. send webhook

Server
1. receive webhook
2. docker pull
3. stop old container
4. run new container

The GitHub Actions steps are the CI pipeline, while the Server steps are the CD pipeline. And these are just my assumptions, and then I'm going to implement them from scratch.

## The CI Pipeline

Our current goal is to get the docker image on docker hub as long as we run the command `git push`. To achieve this, a GitHub action is required, we can set up a workflow by adding a YAML file to the `.github/workflow/` directory, GitHub will recognize the YAML file automatically. I created a `.github/workflow/build.yml` here. And my project contents are under the `web` directory.

```text
repository/
├── .github/workflow/
│   └── build.yml
└── web/
```

This is a minial template for a workflow YAML file. It's easy to know what each field does.

```yml
name: CI Practice

on:
  push:
    branches: ["main"]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: checkout
      uses: actions/checkout@v6

    - name: use pnpm
      uses: pnpm/action-setup@v6
      with:
        version: 12
    - name: install deps
      working-directory: ./web
      run: pnpm install
    - name: build project
      working-directory: ./web
      run: pnpm build

    - name: build image
      env:
        SHA: ${{ github.sha }}
      run: docker build ./web -t <username>/ci-test-web:${SHA:0:4}
    - name: push image
      env:
        SHA: ${{ github.sha }}
      run: |
        docker login -u <username> -p ${{ secrets.DOCKER_HUB_TOKEN }}
        docker push <username>/ci-test-web:${SHA:0:4}
```

The checkout should be the first step because the GitHub workflow server should have the repository itself. And after that, we start installing dependencies as we planed. Then build the project, build the docker image, push it to the docker hub. Everything looks fine and is just the same as we planed.

But here is a thing, we need a docker access token to let the GitHub workflow server access our own docker account for pushing the image. After generating the access token, if we don't have one, we should add it to the environment of our current repository through `Settings -> Secrets and variables -> Actions -> Repository secrets`.

After all these done, the pipeline will run well. You should see a new docker image in your docker hub for each git push.

## Moving on to CD
