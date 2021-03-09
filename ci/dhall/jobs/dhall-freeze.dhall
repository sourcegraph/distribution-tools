let GitHubActions = (../imports.dhall).GitHubActions

let Setup = ../setup.dhall

in  Setup.MakeJob
      Setup.JobArgs::{
      , name = "dhall-freeze"
      , additionalSteps =
        [ GitHubActions.Step::{
          , name = Some "Check that dhall files are linted properly"
          , run = Some "just freeze-dhall"
          , env = Some (toMap { CHECK = "true" })
          }
        ]
      }
