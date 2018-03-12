# Cameradar Contribution

This file will give you guidelines on how to contribute if you want to, and will list known contributors to this repo.

If you're not into software development or not into Golang, you can still help. Updating the dictionaries for example, would be a really cool contribution! Just make sure the credentials and routes you add are **default constructor credentials** and not custom credentials.

If you have other cool ideas, feel free to share them with me at [brendan.leglaunec@etixgroup.com](mailto:brendan.leglaunec@etixgroup.com) or to directly [create an issue](https://github.com/Ullaakut/cameradar/issues)!

## Version 2.0.0

*Cameradar* is the name of the Golang library and the binary that serves as an example of its use, as well as the docker image that runs the binary.

The 2.0.0 version was a complete refactorring of the Cameradar C++ tool, which came from the fact that most users who want to access cameras either wanted to launch it with the basic cache manager, mostly using the docker image already provided in this repository, or did not use it because it did not integrate into their software solution easily.

Transforming it into a library allowed developers to use it directly in their own code exactly as they want, allowing for a greater flexibility. The Cameradar binary also provides a simple use example as well as maintains the old simple way of using Cameradar for non-developers.

## Workflow

### Branches & issues

If you want to work on an issue, make sure you create a specific branch for this issue using the format `issue_number-solution_explanation`. Examples are:

If issue `#64` is `Improve network scan performance`, the branch to fix it should be something like: `64-improve-network-scan-performance`. Note that it should always start with a verb conjugated in the infinitive form, and describe what the commits's effects will be on the codebase. One branch should only be for one change. If your branch fixes multiple things, you're doing it wrong.

Always make sure you're not working on the same issue as someone else, by asking on the issue thread to be assigned to it.

### Commit names

The name of the commits should always be #[issue number] [effect of the issue] (ex: `#343 Improve test coverage`).

When working on your local branch, you can do as many commits as you want, obviously. The most important is that you squash your commits before creating your pull request, or at least before it is merged.

In case you're not familiar with squashing, here is a simple way to do it :

- `git fetch origin` will make sure that you have a local version of the origin repository that is up to date (will not overwrite anything on your branch, no worries)
- `git rebase -i origin/master` will start the process of rebasing your branch
- This will open a file letting you decide what to do with the commits. You want to keep the first `pick` and write `s` or `squash` instead of `pick` for all other commits below.
- If there are conflicts, you will fix them step by step by following what git tells you, it's pretty straight-forward.
- If there are no conflicts or if they are resolved, git will let you edit the commit names. Don't forget to comment the commit names of the commits you squashed if they are not relevant by adding a # character in front of the commit message, and make sure that the commit message you left follows the aforementioned guidelines.
- Now run `git log`, you should see only one commit by the name you chose during the rebase.
- You can now `git push -f` if you already pused your branch on origin or simply push without the `-f` if it's your first push on origin. The reason for the `-f` is that when you squash your commits, you create a new one that will conflict with the state of your branch on origin. If you pull, it will overwrite your local state, so don't do that except if you messed up your rebase.

### Pull Requests

When your pull request is created, GitHub will first check for conflicts, Codacy will check the shell and C++ code's quality and then Travis CI will try to build and launch functional tests of your versions of Cameradar.

If GitHub reports conflicts with the develop branch, you should resolve them by yourself using your git command-line interface. The easiest and cleanest way is to use `git rebase -i origin/develop` and follow git's instructions.
If Codacy reports new issues, they will be added in the comments of the PR to let you know what you should fix.
If Travis CI reports errors, you should be able to view the logs [by clicking here](https://travis-ci.org/Ullaakut/cameradar/builds) and you should fix it. No PR will be merged before all tests are passing correctly.

When creating your pull request, our hooks will make sure that your code:

- Builds
- Has 100% passing unit tests
- Can actually access a camera using a functional test
- Still has equivalent or higher test coverage (using coveralls)

Make sure to write in the PR description what issue it fixes. GitHub will intepret it and automatically close the issue once your pull request is closed. Just write Fixes #IssueNumber in the description.

When your pull request is created, GitHub will first check for conflicts and then your code will be reviewed by the maintainers of this repository.

If GitHub reports conflicts with the `master` branch, you should resolve them by yourself using your git command-line interface. The easiest and cleanest way is to use `git rebase -i origin/master` and follow git's instructions. If we report issues with your code, you should resolve them and then ping the person that reported them to notify them that you did the requested changes.

Once everything is in order, we will merge your pull request.

### Coding guidelines

Your code should just

- Not decrease the results of Cameradar on https://goreportcard.com/report/github.com/Ullaakut/cameradar
- Pass the code review

#### Golang

- All Golang code has to be formated using `gofmt` or `goreturns`.
- Make sure you follow the Golang [best practices](https://golang.org/doc/effective_go.html)

## Contributors

- **Brendan Le Glaunec** - [@Ullaakut](https://github.com/Ullaakut) - brendan.leglaunec@etixgroup.com : *Original developer & Maintainer*
- **Jeremy Letang** - [@jeremyletang](https://github.com/jeremyletang) - letang.jeremy@gmail.com : *Idea of the project & Mentorship*
- **ishanjain28** - [@ishanjain28](https://github.com/ishanjain28) - ishanjain28@gmail.com : *Implemented the environment variables support*