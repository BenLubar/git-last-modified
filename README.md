# git-last-modified

Sets files in a Git repository to their last-modified date.

    usage: git last-modified [<options>] [[--] <path>...]
      -commit-date
            Use the commit date for the last commit this file was involved in instead of the author date.
      -n    Dry run. Implies -v. Don't modify any file modification times.
      -q    Quiet. Don't warn about files specified on the command line that are not in Git.
      -v    Verbose. Print each filename and modification time as they are processed.

[![Build Status](https://travis-ci.org/BenLubar/git-last-modified.svg?branch=master)](https://travis-ci.org/BenLubar/git-last-modified)
[![CodeFactor](https://www.codefactor.io/repository/github/benlubar/git-last-modified/badge)](https://www.codefactor.io/repository/github/benlubar/git-last-modified)
[![Maintainability](https://api.codeclimate.com/v1/badges/139bdab6b8b2bb5ffd17/maintainability)](https://codeclimate.com/github/BenLubar/git-last-modified/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/139bdab6b8b2bb5ffd17/test_coverage)](https://codeclimate.com/github/BenLubar/git-last-modified/test_coverage)
