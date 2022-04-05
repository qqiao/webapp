// Copyright 2022 Qian Qiao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*

Package rememberme contains all functionalities regarding the storage and
validation of the rememberme tokens.

Basics

A rememberme token is a way for the server side application to "remember" the
user logged in.

A user could have an unlimited number of rememberme tokens associated to his
or her account, each with a different unique identifier.

The purpose of the unique identifier is to identify different tokens from
different client device and/or browser instances and to facilitate token
revocation.

If a rememberme token is already available to a client on a particular device
or browser, the client should re-use this token instead of attempting to
obtain a new one. Reusing valid avoids additional tokens being created
unnecessarily and hence the storage burden that is associated to it.

Expiration

A rememberme token itself does not automatically expire. However, it would be
good practice for applications to periodically clean up unused tokens as a good
security measure.

The Purge method is designed for this purpose. The frequency of calling this
function and the cut-off point from which older tokens are deleted are entirely
determined by the application.

Revocation

A token can also be revoked. Once done so, subsequent ValidateToken calls from
any client would fail.

Token revocation is usually done as a security measure, when a user no longer
have access to the browser instance and/or the device a particular token was
initially issued to.

This action is usually initiated by a user, but applications are strongly
recommended to provide the necessary UI to facilitate the token revocation.

*/
package rememberme
