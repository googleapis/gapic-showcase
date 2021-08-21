# server/genrest

This directory contains mostly auto-generated files used to implement a REST
endpoint for Showcase services. The `*_custom.go` files contain manually written
REST-specific handlers, which are useful for helping generators test
REST-specific functionality (such as invalid JSON).
