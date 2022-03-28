/*

Package rememberme contains all functionalities regarding the storage and
validation of the rememberme tokens.

Basics

A user could have an unlimited number of rememberme tokens associated to his
account, each with a different unique identifier.

The purpose of the unique identifier is to identify different tokens from
different browser instances and to facilitate token revocation.

If a rememberme token is already available to a client on a particular device
or browser, the client should re-use it token instead of attempting to
obtaining a new one. Reusing avoids additional tokens being created
unnecessarily and hence the storage burden that is associated to it.

Revocation

A token can also be revoked. Once done so, subsequent ValidateToken calls from
any clients will fail, and clients should re-authenticate and obtain new
tokens.

Token revocation is usually done because a user no longer has access to the
browser instance and/or the device a particular token is issued to.

*/
package rememberme
