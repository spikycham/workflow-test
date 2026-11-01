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

This is a minial template for a workflow YAML file. It's easy to know what each field does. The `github.sha` is actually the new commit hash value.

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

If we want to know when the build is going to be over, we need a service to wait for the GitHub workflow's notification. And that's what GitHub webhooks does. Create a webhook in settings to listen to the workflow run event, and use a service to wait for the post request sent by the webhook. And the event is not that specified, like a workflow run event contains the `requested/in_progress/completed` stages.

We could check the action and the conclusion of a webhook request in our endpoint, pull and run the docker image when the action is `completed` and the conclusion is `success`. And remember to stop the original docker container. If you need the ability of rollback, you should keep the original docker container. Otherwise, you can just delete it. Here I'm not going to keep it.

```go
func runDocker(containerName, imageName string) {
	ctx, cacenl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cacenl()

	stop := exec.CommandContext(ctx, "docker", "stop", containerName)
	if err := stop.Run(); err != nil {
		log.Printf("[ERROR] failed to stop: %v", err)
	}
	rm := exec.CommandContext(ctx, "docker", "rm", containerName)
	if err := rm.Run(); err != nil {
		log.Printf("[ERROR] failed to remove: %v", err)
	}

    // I think there should be an approach to clean up the old images as well. I may try to design a method to handle this in the future.
	pull := exec.CommandContext(ctx, "docker", "pull", imageName)
	if err := pull.Run(); err != nil {
		log.Printf("[ERROR] failed to pull: %v", err)
		return
	}

	cmd := exec.CommandContext(ctx, "docker", "run", "--name", containerName, "-p", "20004:80", "-d", imageName)
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to run docker: %v", err)
		return
	}

	log.Printf("[INFO] success to run docker")
}
```

## Make the CD Service in A Container Affect the Host

Now we can actually run the new image as long as we push the code to remote repository. But here is still an important problem, we should actually run this service in a container, and the commands we write here only affects the container itself. So we need to explore how to connect the container commands to the host environment.

It's easy to know that we must have a docker CLI installed inside of the docker container for running the docker commands, otherwise the program may throw errors. So for a go service, we know that we can compile the program, and copy the artifact to another container with alpine image to run it. And the key is to install the docker inside of the alpine container, with a command in the dockerfile `RUN apk add --no-cache docker-cli`.

And after that, we just need to mount the host's docker socket by `-v /var/run/docker.sock:/var/run/docker.sock \`. This gives your service access to the host's docker daemon, and you don't need to change your service code for even one line.

> NOTE: Access to /var/run/docker.sock effectively grants root-level control over the host.

## Result

This is wonderful. After all these done, you can just code your new feature, and push to your remote repository (or create a pull request and merge it to main branch). Then the website is going to be updated automatically!

