# SFTP

A low &amp; high level golang SFTP (version 3) client library.

## Low level

Sending packets over SFTP and receiving a response is the lowest level of abstraction sftp libraries can provide. Existing 
solutions do not provide low level access to the SFTP protocol's raw packets, which makes them difficult to use in low level 
FUSE programs. This SFTP library provides these mechanisims since it is being designed to plug into an existing FUSE project.

## High level

Eventually, this library will provide a high level SFTP implementation that mimics golang's "os" package.
