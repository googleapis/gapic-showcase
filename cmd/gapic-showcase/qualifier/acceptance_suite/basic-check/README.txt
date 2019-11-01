This " basic-check" scenario also showcases some features of the
"qualifer" acceptance harness. In particular, it shows how files and
directories can be included at various points in the sandbox hierarchy
by means of include files:

* basic-check/include.google includes the subdirectory of api-common
  protos needed by the echo service
  
* basic-check/include.example.just_a_file is just an example showing
  how to include just a file instead of a directory

* basic-chec/foo/bar/include.example.an_arbitrarily_nested_file is
  just an example showing how to include a file or directory at any
  point in a sandbox directory tree.
