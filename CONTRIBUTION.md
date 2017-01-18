# Cameradar Contribution

This file will give you guidelines on how to contribute if you want to, and will list known contributors to this repo.

If you're not into software development or not into C++, you can still help. Updating the dictionaries for example, would be a really cool contribution! Just make sure the ids and routes you add are **default constructor credentials** and not custom credentials.

If you have other cool ideas, feel free to share them with me at [brendan.leglaunec@etixgroup.com](mailto:brendan.leglaunec@etixgroup.com) !

## Version 2.0.0

- *Cameradar* will become the name of the library.
- *Cameraccess* will be the name of the binary that uses Cameradar to _hack_ the cameras.

This quite big refactoring comes from the fact that most users who want to access cameras either want to launch it with the basic cache manager, mostly using the docker image already provided in this repository, or will not use it because it does not integrate into their software solution without sharing their database with Cameradar, which would cause issues with database migrations for example.

Transforming it into a library will allow developers to use it directly in their own code exactly as they want, allowing for a greater flexibility. The Cameraccess binary will then provide a simple use example as well as maintaining the current simple way of using Cameradar for non-developers.

This is quite a huge task compared to the tiny changes I usually do on Cameradar, so it might take a long time.

If you want to contribute, note that the develop will stay in 1.x until the 2.0.0 is released. A new development branch will be created especially for the 2.0 version, called `2.0.0` from which all work on the 2.0.0 version will be done until the 2.0.0 version is ready to replace the 1.x on the master and develop branches. The rest of the workflow is exactly the same as for the rest of the repository.

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

#### C++

+ All C++ code has to be formatted using `clang-format`
+ The namespaces should be respected and new files should implement the same namespace structure as the other files
+ Forward declarations should be used as much as possible
+ Use smart pointers instead of raw pointers as much as possible
+ Each constructor with only one parameter which is not a copy or a move constructor must be marked explicit
+ Use C++11 specifiers as much as possible *(override, noexcept)*
+ Variable and function names must always be in *snake_case*.

#### Golang

+ All Golang code has to be formated using `gofmt`
+ Make sure you follow the Golang [best practices](https://golang.org/doc/effective_go.html)

#### Shell scripting

+ Just make sure Codacy does not trigger warnings on your code. I probably suck more than you in shell anyway, who would I be to give you guidelines on it?

## Contributors

+ **Brendan Le Glaunec** - [@Ullaakut](https://github.com/Ullaakut) - brendan.leglaunec@etixgroup.com : *Original developer & Maintainer*
+ **Jeremy Letang** - [@jeremyletang](https://github.com/jeremyletang) - letang.jeremy@gmail.com : *Idea of the project & Mentorship*
