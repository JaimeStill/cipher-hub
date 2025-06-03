# Planning

This document will capture key retrospective concepts as well as decision adjustments that will help initiate a successful and scalable start to this project.

## Retrospective

Too much emphasis on testing and documentation straight away. Package completion should be the main emphasis, with validation occurring through direct execution as opposed to testing. Once the package is in a state of completion, then testing and documentation can be addressed.

The initial workflow was too context and token intensive. It additionally did not scope units of effort towards singularly-focused tasks that LLMs tend to excel at. Focus should be 15-20 minute iterative processes that have rapid feedback loops with minimal extraneous variables.

Adopt a scientific-driven development mindset. During package development, verify by running the actual code as it is intended to run so you can see and feel how it operates. This will help you understand what's important to test when the package is complete. Metric for testing should shift from 90% or higher to testing only critical infrastructure. Don't overdo it.

## Strategy

Project Development Phase: this is where you identify the ONE realistically scoped thing that you want to work on for an extended period of time. AI integration is crucial at this phase and should be conversational, helping to establish the parameters of the project, ensuring it stays appropriately scoped, and is a feasible and viable project in the context of where your skills are and what you want to accomplish.

* The output of this stage should be a roadmap that outlines the primary stages of effort that will be required to accomplish the project. These stages should contain: an overview of what the stage accomplishes, a list of key tasks required to execute the stage (these are not specific, just cursory level to aid in planning when reaching that stage), any key design considerations or pitfalls prevalent at this stage, and links to any helpful resources associated with the features of this stage.
  * The roadmap output from your Project Development Phase could benefit from explicit dependency mapping between stages. This helps identify which experiments can be parallelized versus which must be sequential.

Planning Phase: you have a roadmap and are at a specific task and stage within the project. This is where you work with AI to devise an experiment with an expected outcome intended to yield progress along your roadmap. If the experiment is geared towards a novel concept that you are unfamiliar with, it should be geared towards standalone execution so that the amount of variables affecting the experiment are minimized. If the experiment is directly integrated into or dependent on the work you've already established within the project, you should gear the experiment towards being executed directly within the project.

* The output of this stage is an experiment with a hypothesis aimed towards a desired outcome, and the set of steps you will take to validate the hypothesis. These steps are more descriptive what action needs to be taken in sequential order. Pseudo-code can be provided if it provides sufficient insight into the step, but otherwise this is just letting you know what needs to be done, not precisely how to do it (again, optimizing the impact on context and token limitations).
* A step should not have a cost greater than a certain context or token threshold. I haven't built enough context yet to describe what that is in terms of tokens, but my intuition from previous experimentation says that a step should not affect more than a single file or user action.
* The overall experiment should be small enough to reasonably be developed within the context window of a single chat session, resulting in the output artifact.

Execution Phase: you have the experiment written out, now it's time to execute it. During this phase, pair program with the AI agent (Claude Code for instance) within your development environment to execute the steps. All steps prior to this were executed in a chat context, this is the first instance where we're leaning on an AI agent to assist with development. Focus on describing exactly what the outcome of the step should be so that the agent can generate accurate results. The steps in the experiment are guidelines to help you achieve a desirable outcome. If you encounter any issues during the experiment, it is okay to make dynamic adjustments if you can conceive of changes that would affect the outcome. Be sure to capture these deviations and their resulting impacts. If an issue invalidates or otherwise makes the experiment unfeasible, it is a failed experiment. Be sure to capture all of the details of the failure.

* The output of this phase is either the successful execution of the experiment, or the details surrounding its failure.
* If the experiment was successful, you can move on to the next task. If conducted external to the project, perform whatever steps are needed to integrate the feature into the project. Also be sure to generate documentation and keep the isolated experiment and project retained in a standards repository for future reference. Be sure to adjust the experiment for any deviations encountered from the original experiment.
* If the experiment was unsuccessful, use the details of the failed experiment to return to the planning phase and devise another experiment
