# Cameradar Contribution

This file will give you guidelines on how to contribute if you want to, and will list known contributors to this repo.

If you're not into software development or not into Golang, you can still help. Updating the dictionaries for example, would be a really cool contribution! Just make sure the credentials and routes you add are **default constructor credentials** and not custom credentials.

If you have other cool ideas, feel free to share them with me at [brendan.leglaunec@etixgroup.com](mailto:brendan.leglaunec@etixgroup.com) !

## Version 2.0.0

- *Cameradar* is the name of the Golang library.
- *Cameraccess* is the name of the binary that uses Cameradar to discover and access the cameras.

This quite big refactoring comes from the fact that most users who want to access cameras either want to launch it with the basic cache manager, mostly using the docker image already provided in this repository, or will not use it because it does not integrate into their software solution without sharing their database with Cameradar, which would cause issues with database migrations for example.

Transforming it into a library allows developers to use it directly in their own code exactly as they want, allowing for a greater flexibility. The Cameraccess binary also provides a simple use example as well as maintains the old simple way of using Cameradar for non-developers.

## Workflow

### Branches & issues

When an issue is opened, a branch will be automatically created. If you want to work on this issue, this is the branch you **have** to work on and create your pull request from.

**Always make sure you're not working on the same issue as someone else, by asking on the issue to be assigned to it.**

### Commit names

The name of the commits should always be `v[next version] : [name of the fixed issue]` (ex: `v1.1.4 : Removed unnecessary null pointer checks`), and each PR should only contain one single commit.

When working on your local branch, you can do as many commits as you want, obviously. The most important is that you **squash** your commits before creating your pull request.

In case you're not familiar with squashing, here is a simple way to do it :

+ On your branch, when everything is clean and working, launch `git log` and count the number of commits your branch is ahead from compared to the `develop` branch.
+ Then launch `git rebase -i HEAD~X`, X being the number of commits you want to squash. For example if I had 12 commits on my branch, I will squash all of them by writing `git rebase -i HEAD~12`.
+ This will open a file letting you decide what to do with the commits. You want to keep the first `pick` and write `s` instead of the other ones, s meaning squash.
+ If there are conflicts, you will fix them step by step by following what git tells you, it's pretty straight-forward.
+ If there are no conflicts or if they are resolved, git will let you edit the commit names. Don't forget to comment the commit names of the commits you squashed by adding a `#` character in front of the commit message.
+ Now launch `git log`, you should see only one commit by the name you chose during the rebase.

### Pull Requests

When your pull request is created, GitHub will first check for conflicts, Codacy will check the shell and C++ code's quality and then Travis CI will try to build and launch functional tests of your versions of Cameradar.

If GitHub reports conflicts with the develop branch, you should resolve them by yourself using your git command-line interface. The easiest and cleanest way is to use `git rebase -i origin/develop` and follow git's instructions.
If Codacy reports new issues, they will be added in the comments of the PR to let you know what you should fix.
If Travis CI reports errors, you should be able to view the logs [by clicking here](https://travis-ci.org/EtixLabs/cameradar/builds) and you should fix it. No PR will be merged before all tests are passing correctly.

### Coding guidelines

This part will tell you about what are the general coding guidelines I want to keep on this project.

#### Golang

+ All Golang code has to be formated using `gofmt`
+ Make sure you follow the Golang [best practices](https://golang.org/doc/effective_go.html)

#### Shell scripting

+ Just make sure Codacy does not trigger warnings on your code.

## Contributors

+ **Brendan Le Glaunec** - [@Ullaakut](https://github.com/Ullaakut) - brendan.leglaunec@etixgroup.com : *Original developer & Maintainer*
+ **Jeremy Letang** - [@jeremyletang](https://github.com/jeremyletang) - letang.jeremy@gmail.com : *Idea of the project & Mentorship*
