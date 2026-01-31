# Go Repo Template

This repository serves as a template for creating new Go projects. It includes a helper script to bootstrap a new repository from this template.

## Getting Started

To create a new repository based on this template, you will need:

1. A local copy of this `go-repo-template` repository.
2. A newly created, empty repository on your git host (e.g., GitHub, GitLab).

### Initialization Script

The `new_repo_init.sh` script automates the process of copying the template, initializing the git remote, and pushing the initial commit.

#### Usage

Run the script from the root of this template directory, providing the URL of your new repository:

```bash
./new_repo_init.sh <new_repo_url>
```

**Example (SSH):**

```bash
./new_repo_init.sh git@github.com:username/my-new-service.git
```

**Example (HTTPS):**

```bash
./new_repo_init.sh https://github.com/username/my-new-service.git
```

#### What the script does:

1. Creates a new directory for your repo as a sibling to this template directory.
2. Copies all template files to the new directory.
3. Removes the `new_repo_init.sh` script from the new repository (cleanup).
4. Initializes the git remote to the provided URL.
5. Creates and checks out a new branch named `repo_init`.
6. Commits and pushes the initial state to the remote repository.

Once the script completes, your new repository will be ready in `../<repo-name>` and pushed to the remote.
