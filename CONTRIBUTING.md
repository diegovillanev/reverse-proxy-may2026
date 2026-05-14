# Contributing to this project

We would love for you to contribute to this project and help make it even better than it is today! As a contributor, here are the guidelines we would like you to follow:

1. [Code of Conduct](#coc)
2. [Question or Problem?](#question)
3. [Issues and Bugs](#issue)
4. [Feature Requests](#feature)
5. [Submission Guidelines](#submit)
6. [Coding Rules](#rules)
7. [Commit Message Guidelines](#commit)


## 1. Code of Conduct

This project is governed by a Code of Conduct aimed at promoting a respectful environment for all participants. We encourage constructive and professional collaboration regardless of contributors’ level of experience or technical knowledge. We strongly recommend reading the [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.


## 2. Questions and Issues

If you encounter difficulties or have questions about the project’s operation, we suggest first reviewing the existing documentation and the Frequently Asked Questions (FAQ). If you don’t find a solution there, open an issue in the repository describing your inquiry in detail. Include relevant information such as your working environment, software version, and any specific error messages.

Providing a detailed report will help maintainers and other contributors address your inquiry more efficiently, reducing possible delays and improving the quality of responses.


## 3. Bug Reports

- Fill out the provided issue template, making sure to include clear and precise details. This helps reproduce the problem and identify its root cause.
- Include a thorough description of the error along with the steps required to reproduce it. This allows maintainers to quickly locate the cause of the issue.
- If possible, provide screenshots, error logs, and code samples related to the reported problem. These additional details are crucial for efficiently diagnosing the issue and proposing effective solutions.


## 4. Feature Proposals

- Explain how the new feature will contribute to the project’s progress. Describe its relevance and expected benefits. Try to include a brief analysis of the potential impact on end users.
- Provide a detailed outline of the proposed implementation, including examples, dependencies, and any possible impact on system performance. Explain how this feature complements existing ones.
- If applicable, compare your proposal with existing alternatives to highlight its advantages. This helps reviewers better understand its value and contribution to the project.


## 5. Contribution Guidelines

1. Create a new descriptive branch using a command such as `git checkout -b my-new-feature`. Make sure the branch name is clear and representative of the change.
2. Make the necessary changes and write clear, concise commit messages following the predefined standards.
3. Submit a well-structured and documented Pull Request (PR). Include a clear description of the change, its purpose, and any other relevant information for reviewers. Be collaborative when addressing any feedback you receive.


## 6. Coding Standards

- Follow the project’s defined style conventions to ensure consistency across the codebase. A uniform style improves readability and lowers entry barriers for new contributors.
- Write clear, efficient, and properly documented code. The goal is for other developers to quickly understand your contributions, even if they are not familiar with every detail of the project.
- Include tests that validate the changes you’ve made. Tests help ensure the system continues to function correctly after your modifications are implemented.
- Review your code to identify and fix errors before submitting the PR. Use static analysis tools or linters to ensure code quality.


## <a name="commit"></a> 7. Commit Message Guidelines

### @dievilz' custom Git Commit Message Guidelines

This is intended to extend/work with the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.<br>
You can download @smallhadroncollider's [cmt](https://github.com/smallhadroncollider/cmt) Haskell tool, @emac's [gitx](https://github.com/emac/gitx) shell script and/or my tool: [gocmt](https://github.com/dievilz/gocmt), to enforce your team to use these guidelines.<br>
Forked from various articles, see the [links](#commit-links).

### Structure and Composition of the Commit Message

A commit message should be structured in three sections separated by a blank line. This structure helps distinguish between the general purpose of the change, technical details, and additional considerations.

1. **Header**
2. **Body** (optional)
3. **Footer** (optional).

```
[<ticket>] <type>(<scope>): <subject> #<meta>           F
-------------- blank line -------------------             O
Body                                                        R
-------------- blank line ----------------------------        M
Footer                                                          A
------------ last line (EOF) -------------------------            T
Header limit: ------------- 50-72 chars ---->|
Body & Footer limit: ------------- 72-100 chars ----->|
```

### Header

The header is a concise line that summarizes the change made. It should be clear, short, and directly communicate the purpose of the commit. Ideally, it summarizes the essence of the change in clear and simple terms.

It consists of 5 parts, which will be explained below:
`[<ticket>] <type>(<scope>): <subject> #<meta>`

#### Header's Type (mandatory)

The **type** of a commit message should be a single word or abbreviation drawn from an ontology, according to the nature of the project.
This document specifies a programming ontology, with the following elements:

* **`mod`:**          The everyday changes to the code that are not fixes of any kind, style changes, refactoring or anything else. This should be the most common type showed in the commit graph. If you don't like it, `edit` may also be used, but use only one of both consistently.
* **`blob`:**         A newly created file. Better when just creating an empty file rather than an almost finished file. May include _#wip_ metatag. May be used when changing file permissions.
* **`tree`:**         Like above, but for directories. May include _#wip_ metatag. May be used when changing dir permissions.
* **`fix`:**          A hot/bug fix.
* **`style`:**        Changes to the text **_Styling/Formatting_** - that do not affect the meaning of the code (white-spaces, comments, indentation, missing punctuation, typos, etc).
* **`refactor`:**     Changes to the code **_Composition/Organization_** - improving readability/performance and does not fix bugs or add features (changes to code lines/blocks, functions, structuring of the file, etc).
* **`feat`:**         A new feature (a whole unit of functionalities). Use this only when merging a feature branch or a whole set of files that only have feature changes. May be used for a squash commit.
* **`docs`:**         Documentation-only changes.
* **`test`:**         Adding missing tests or correcting existing ones.
* **`build`:**        Changes to the build/compilation/packaging process or auxiliary tools such as documentation generation or external dependencies.
* **`devops/ci-cd`:** Changes in the **_Continuous Integration/Delivery_** setup, files, scripts, etc. `devops` may be used for a broader spectrum of this type of changes or may be used instead of `ci-cd` consistently.
* **`notice`:**       Changes to announce/warn anything related to: files, code blocks, etc.
* **`chore`:**        For any other repetitive and periodic tasks (like cleanups of deprecated bits, or bumping versions of things). If it's something a bot could have done instead of the devs, it's likely a chore. May be used instead as meta hashtag for repetitive refactors, i.e. always adding aliases to an .aliases file.

* **`revert`:**       If the commit reverts a previous commit, it should begin with _"revert:"_, followed by the subject of the reverted commit. In the body it should say: "This reverts commit `<hash>`.", where the hash is the SHA of the commit being reverted and explain the reason(s), and footer say: "Reverts `<hash>`". For example:

    ```
    revert: include more details in command-line help text:
    ---------- blank line -----------
    This reverts commit 5b233b5a because of...
    --------- blank line -----------
    Reverts 5b233b5a
    ```

Header line may be prefixed for continuous integration purposes.

> For example, [Jira](https://bigbrassband.com/git-for-jira.html)
> requires ticket in the beggining of commit message:<br><br>
> `[LHJ-16] fix(compile): add unit tests for Windows`
> <br>
<br>

This is a custom ontology. For more information on a more commonly used ontology, see the full specification of Conventional Commits (https://www.conventionalcommits.org/en/v1.0.0/).

#### Header's Scope (optional)

Usually it is convenient to mention exactly which part of the code base changed.
The **scope** token is responsible for providing that information. While the granularity of the scope can vary, it is important for it to be a part of the "common language" spoken in the project.
Please notice that in some cases the scope is naturally too broad, and therefore not worthy to mention. `<TYPE>` and `<SCOPE>` may be mutually exclusive.

```
feat(auth): introduce sign-in via GitHub
```

#### Header's Subject (mandatory)

The **subject** token should contain a succinct description of the change(s).

* Soft limit: **50** chars. Hard limit: **72** chars.
* Use the infinitive tense to mainly describe the behavior of the program after the commit, i.e. _“change”_. Avoid describe your _coding behavior_.
* May be prefixed for CI/CD purposes.
* Do not capitalize the subject line. (Non-standard)
* Do not end the subject line with a period.

```
refactor: move folder structure to `src` directory layout
```

#### Header's Meta (optional)

The end of subject-line may contain **hashtags** to facilitate changelog generation and bissecting:

* **`#wip`**: The feature being implemented is not complete yet. Should not be included in changelogs (just the last commit for a feature goes to the changelog).

* **`#irrelevant`**: The commit does not add useful information. Used when fixing typos, etc. Should not be included in changelogs.

```
blob: add TODO markdown file #wip #irrelevant
```

### Body (optional)

Includes motivation for the change and contrasts with previous behavior in order to illustrate the impact of the change.

* Soft limit: **72** chars. Hard limit: **100** chars.
* Use infinitive, present tense: “change”, not “changed” nor “changes”
* Use the body to explain _What_ and _Why_, not _How_.
* Keep it Simple, Future-proof, Junior-dev friendly.
* Markup syntax as Markdown can be applied here.
  * For simple headers, type a space (So git couldn't parse it as comment), then use `H5` or `H6`: ` ##### <header>`
* See [1](#commit-link-1) and [2](#commit-link-2) for more info.

Optional directives to write the body.

* Why you made the change in the first place:
  * The way things worked before the change (and what was wrong with that), the way they work now, and why  you decided to solve it the way you did.<br>
    See [3](#commit-link-3).
* You may explain the same changes in 4 different perspectives (optional):
  * From the user’s perspective: A description of how a user would see incorrect program behavior, steps to reproduce a bug, user goals that the commit is addressing, what they can see, who is affected.
  * From a manager’s perspective: Design choices, your creativity, why you made the changes.
  * From the code’s perspective: A line-by-line, function-by-function, or file-by-file summary.
  * From git’s perspective: Any related commits in this or another repository, especially if you are reverting earlier changes; related GitHub issues.<br>
    See [4](#commit-link-4).

```
feat($browser): add onUrlChange event (popstate/hashchange/polling)

##### New $browser event:
 - forward popstate event if available
 - forward hashchange event if popstate not available
 - do polling when neither popstate nor hashchange available

Breaks $browser.onHashChange, which was removed (use onUrlChange instead)
```

### Footer (optional)

* All breaking changes or deprecations have to be mentioned in footer with the description of the change, justification and migration notes.

  ```
  Breaks $browser.onHashChange, which was removed (use onUrlChange instead)
  ```

* Referencing issues: closed bugs should be listed on a separate line in the footer prefixed with "Closes" keyword.

  ```
  Closes #123
  Closes #123, #245, #992
  ```

### Extras

#### Generating `CHANGELOG.md`

Changelogs may contain three sections: **new features**, **bug fixes**, **breaking changes**.
This list could be generated by script when doing a release, along with links to related commits. Of course you can edit this change log before actual release, but it could generate the skeleton.

* List of all subjects (first lines in commit message) since last release:

  ```
  git log <last tag> HEAD --pretty=format:%s
  ```

* New features in this release:

  ```
  git log <last release> HEAD --grep feat
  ```

### <a name="commit-links"></a> Links

These guidelines are based directly from @abravalheri's [gist](https://gist.github.com/abravalheri/34aeb7b18d61392251a2), which in turn is an extended version from @stephenparish's [gist](https://gist.github.com/stephenparish/9941e89d80e2bc58a153).<br>
_(Links [1](#commit-link-1) and [2](#commit-link-2) are referenced in both gists and I borrowed some things from [3](#commit-link-3) and [4](#commit-link-4))_

1. @abizern's (365git) article: <a name="commit-link-1"></a>[http://365git.tumblr.com/post/3308646748/writing-git-commit-messages](http://365git.tumblr.com/post/3308646748/writing-git-commit-messages)
2. @tpope's article: <a name="commit-link-2"></a>[http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
3. @chris.beams' article: <a name="commit-link-3"></a>[https://chris.beams.io/posts/git-commit/#why-not-how](https://7chris.beams.io/posts/git-commit/#why-not-how)
4. @joshuatauberer's article: <a name="commit-link-4"></a>[https://medium.com/@joshuatauberer/write-joyous-git-commit-messages-2f98891114c4](https://medium.com/@joshuatauberer/write-joyous-git-commit-messages-2f98891114c4).

@abravalheri's and @stephenparish's guidelines are both based on AngularJS project's [Commit Guidelines](https://github.com/angular/angular.js/blob/master/CONTRIBUTING.md).<br>

---

We deeply appreciate your willingness to improve this **project**. This project is a collaborative effort, and every contribution, no matter how small, brings us closer to a more robust, efficient, and valuable product. Together, we can take this project to new heights, benefiting a wider community.