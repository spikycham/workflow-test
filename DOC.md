# CI/CD Practice

## What I'm Going to Do?

I'm trying to make a CI/CD pipeline from scratch, I never did this before. Here I'm going to analyze what I need to do first.

Our goal of building a CI/CD pipeline is to make the project deployment easier. Like we just need to push our code to remote repository, the project is going to be tested and deployed automatically. And this is how it looks like overall: `Push -> Deploy`.

Let's split this goal into steps. So right after we push the code to remote repository (Github).

Github Actions
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

## The Core Things Required


