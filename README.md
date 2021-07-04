# Jackcompiler

Building a compiler for Jack - High-Level Programming, simple, Java-like, object-based programming
language.

## Lexical Analyser

To execute the tests about the lexical analyse, run:

```
go test jackcompile/lexical_analysis
```

A file will be created in the folder */resources/tests/results*, containing the lexical analysis of */resources/Square.jack* class, in XML format.

To execute the specific test that generates the file, run:

```
go test -run ^TestNewJackTokenizerMustBuildTheExpectedXMLTree$ jackcompile/lexical_analysis
```
