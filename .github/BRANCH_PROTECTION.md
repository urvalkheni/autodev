# 🛡️ AutoDev Branch Protection Setup Guide

To protect the `main` branch from direct pushes, prevent open-source contributors from merging their own pull requests, and ensure that all CI pipelines are run and approved by you, follow this step-by-step setup in the GitHub web interface:

## 1. Enable Branch Protection for `main`
1. Go to your repository on GitHub: **https://github.com/HEETMEHTA18/autodev**
2. Click on the **Settings** tab in the top navigation bar.
3. In the left sidebar, under **Code and automation**, click on **Branches**.
4. Under **Branch protection rules**, click **Add branch protection rule** (or edit the rule for `main` if it already exists).
5. Set the **Branch name pattern** to `main`.

## 2. Block Direct Merges and Require PR Approval
Configure the following settings in the protection rule:
* **Check: Require a pull request before merging**
  * **Check: Require approvals**
  * Set **Required number of approvals before merging** to `1` (or more).
  * This prevents any open-source contributor from merging their own PR. It will require an explicit approval review from you (`HEETMEHTA18`) before the "Merge" button is unlocked.
* **Check: Dismiss stale pull request approvals when new commits are pushed**
  * (Recommended) This ensures that if a contributor pushes new changes after your approval, the approval is reset and you must review it again.

## 3. Require CI Pipeline Status Checks to Pass
To ensure no broken code is merged into `main`:
* **Check: Require status checks to pass before merging**
  * Check **Require branches to be up to date before merging** (this ensures the PR is tested with the latest main code).
  * In the search box under **Status checks that must pass**, search for and select your CI jobs:
    * `Go Lint` (from `.github/workflows/ci.yml`)
    * `Go Tests` (matrix modules)
    * `Go Build`
    * `Website Lint & Type Check` (from `.github/workflows/ci.yml`)
    * `Cross Compile`
    * `Analyze Code Quality / Analyze (go)` (from CodeQL)
    * `Analyze Code Quality / Analyze (javascript-typescript)` (from CodeQL)

## 4. Save the Rule
1. Scroll down to the bottom of the page.
2. Click **Create** (or **Save changes**).

---

Once configured:
1. Open-source contributors cannot push directly to `main` or merge their own PRs.
2. You will be notified to review every PR.
3. The merge button will remain locked until your approval is submitted and all automated checks (Linter, Tests, Build, CodeQL) pass cleanly.
