//+build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"
	"github.com/Ullaakut/disgo"
	"github.com/Ullaakut/disgo/style"
)

var supportedPlatforms = map[string]string{
	"linux/amd64": "ullaakut/cameradar:amd64",
	"linux/386": "ullaakut/cameradar:386",
	"linux/arm64": "ullaakut/cameradar:arm64",
	//"linux/riscv64": "ullaakut/cameradar:riscv64", // UNSUPPORTED.
	//"linux/ppc64le": "ullaakut/cameradar:ppc64le", // UNSUPPORTED.
	//"linux/s390x": "ullaakut/cameradar:s390x", // UNSUPPORTED.
	"linux/arm/v7": "ullaakut/cameradar:armv7",
	//"linux/arm/v6": "ullaakut/cameradar:armv6", // UNSUPPORTED.
}

var Default = Build

// Follows https://www.docker.com/blog/multi-platform-docker-builds/.
func Build() error {
	term := disgo.NewTerminal(disgo.WithColors(true))

	term.StartStep("Building images for all platforms")
	term.Infof("Builds planned for %v\n", supportedPlatforms)
	for platform, name := range supportedPlatforms {
		term.Infoln("Building image for", platform, "at", name)

		// docker buildx build --platform linux/arm/v7 -t ullaakut/cameradar:armv7 .
		if err := sh.Run("docker", "buildx", "build", "--platform", platform, "-t", name, "../../"); err != nil {
			return term.FailStepf("unable to build image: %v", err)
		}
	}

	term.Infoln(style.Success("Cross-platform docker build successful."))

	return nil
}

func Publish() error {
	term := disgo.NewTerminal(disgo.WithColors(true))

	term.StartStep("Pushing images to DockerHub")
	term.Infoln("Pushing ullaakut/cameradar:latest")
	if err := sh.Run("docker", "push", "ullaakut/cameradar:latest"); err != nil {
		return term.FailStepf("unable to push latest docker images to docker hub: %v", err)
	}

	if version, exists := os.LookupEnv("CAMERADAR_VERSION"); exists {
		term.Infoln("Pushing ullaakut/cameradar:"+version)
		if err := sh.Run("docker", "push", "ullaakut/cameradar:"+version); err != nil {
			return term.FailStepf("unable to push versionned docker images to docker hub: %v", err)
		}
	}

	term.StartStep("Pushing images to GitHub Packages")
	term.Infoln("Pushing docker.pkg.github.com/ullaakut/cameradar/cameradar:latest")
	if err := sh.Run("docker", "tag", "ullaakut/cameradar:latest", "docker.pkg.github.com/ullaakut/cameradar/cameradar:latest"); err != nil {
		return term.FailStepf("unable to push latest docker images to docker hub: %v", err)
	}
	if err := sh.Run("docker", "push", "docker.pkg.github.com/ullaakut/cameradar/cameradar:latest"); err != nil {
		return term.FailStepf("unable to push latest docker images to docker hub: %v", err)
	}

	if version, exists := os.LookupEnv("CAMERADAR_VERSION"); exists {
		term.Infoln("Pushing docker.pkg.github.com/ullaakut/cameradar/cameradar:"+version)
		if err := sh.Run("docker", "tag", "ullaakut/cameradar:"+version, "docker.pkg.github.com/ullaakut/cameradar/cameradar:"+version); err != nil {
			return term.FailStepf("unable to push latest docker images to docker hub: %v", err)
		}
		if err := sh.Run("docker", "push", "ullaakut/cameradar:"+version); err != nil {
			return term.FailStepf("unable to push versionned docker images to docker hub: %v", err)
		}
	}

	term.StartStep("Creating manifest(s) for cross platform builds")

	var manifestImages []string
	for _, image := range supportedPlatforms {
		manifestImages = append(manifestImages, image)
	}

	args := []string{"manifest", "create", "--amend", "ullaakut/cameradar:latest"}
	args = append(args, manifestImages...)

	// docker manifest create ullaakut/cameradar:latest ullaakut/cameradar:amd64 ullaakut/cameradar:armv7 [...]
	if err := sh.Run("docker", args...); err != nil {
		return term.FailStepf("unable to create manifest: %v", err)
	}

	if version, exists := os.LookupEnv("CAMERADAR_VERSION"); exists {
		args = []string{"manifest", "create", "--amend", "ullaakut/cameradar:"+version}
		args = append(args, manifestImages...)

		if err := sh.Run("docker", args...); err != nil {
			return term.FailStepf("unable to create manifest: %v", err)
		}
	}
	term.EndStep()

	term.Infoln(style.Success("Images published successfully."))

	return nil
}