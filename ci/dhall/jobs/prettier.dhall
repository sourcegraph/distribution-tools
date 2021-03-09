let GitHubActions = (../imports.dhall).GitHubActions

let Setup = ../setup.dhall

in  Setup.MakeJob
      Setup.JobArgs::{
      , name = "prettier-format"
      , additionalSteps =
        [ GitHubActions.Step::{
          , name = Some "Prettier formatting"
          , run = Some "ci/check-prettier.sh"
          }
        ]
      }
