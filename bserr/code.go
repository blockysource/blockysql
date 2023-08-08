// Copyright 2023 The Blocky Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bserr

// Code is a sql driver independent error code, used to determine the type of the error.
type Code uint32

const (
	// OK is returned when there is no error.
	OK Code = 0

	// NotFound is returned when a row is not found.
	NotFound Code = 1

	// UniqueViolation is returned when a unique constraint is violated.
	UniqueViolation Code = 2

	// TableNotFound is returned when a table is not found.
	TableNotFound Code = 3

	// DataException is returned when a data exception is detected.
	DataException Code = 4

	// ConcurrentUpdate is returned when a concurrent update is detected.
	ConcurrentUpdate Code = 5

	// AuthenticationFailed is returned when authentication failed.
	AuthenticationFailed Code = 6

	// InternalError is returned when an internal error occurred.
	InternalError Code = 7

	// ForeignKeyViolation is returned when a foreign key constraint is violated.
	ForeignKeyViolation Code = 8

	// ConstraintViolation is returned when a not null, check or exclusion constraint is violated.
	// It is not used for unique or foreign key constraints.
	ConstraintViolation Code = 9

	// InvalidInputSyntax is returned when an invalid input syntax is detected.
	InvalidInputSyntax Code = 10

	// PermissionDenied is returned when authorization failed.
	PermissionDenied Code = 11

	// OutOfDisk is returned when the disk is full.
	OutOfDisk Code = 12

	// OutOfMemory is returned when the memory is full.
	OutOfMemory Code = 13

	// TooManyConnections is returned when the maximum number of connections is reached.
	TooManyConnections Code = 14

	// TxDone is returned when the transaction is already closed.
	TxDone Code = 15

	// Timeout is returned when the query timed out.
	Timeout Code = 16

	// Unknown is returned when the error is unknown.
	Unknown Code = 17
)
