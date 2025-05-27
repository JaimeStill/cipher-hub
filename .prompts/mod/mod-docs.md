# Moderate Docs

In the Go code base, packages contain a `doc.go` document that provides package documentation. Evaluate the contents of the documentation against the current state of the code. Generate any corrections and optimizations, then output the adjusted document into an artifact that can be used to adjust the package documentation.

If a package is missing a `doc.go` file, analyze the state of the package and generate the `doc.go` contents into the artifact. Artifact only. No preamble or summary.