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

// Code generated by "enumer -type=Code -text -output=code_string.go code.go"; DO NOT EDIT.

package bserr

import (
	"fmt"
	"strings"
)

const _CodeName = "OKNotFoundUniqueViolationTableNotFoundDataExceptionConcurrentUpdateAuthenticationFailedInternalErrorForeignKeyViolationConstraintViolationInvalidInputSyntaxPermissionDeniedOutOfDiskOutOfMemoryTooManyConnectionsTxDoneTimeoutUnknown"

var _CodeIndex = [...]uint8{0, 2, 10, 25, 38, 51, 67, 87, 100, 119, 138, 156, 172, 181, 192, 210, 216, 223, 230}

const _CodeLowerName = "oknotfounduniqueviolationtablenotfounddataexceptionconcurrentupdateauthenticationfailedinternalerrorforeignkeyviolationconstraintviolationinvalidinputsyntaxpermissiondeniedoutofdiskoutofmemorytoomanyconnectionstxdonetimeoutunknown"

func (i Code) String() string {
	if i >= Code(len(_CodeIndex)-1) {
		return fmt.Sprintf("Code(%d)", i)
	}
	return _CodeName[_CodeIndex[i]:_CodeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _CodeNoOp() {
	var x [1]struct{}
	_ = x[OK-(0)]
	_ = x[NotFound-(1)]
	_ = x[UniqueViolation-(2)]
	_ = x[TableNotFound-(3)]
	_ = x[DataException-(4)]
	_ = x[ConcurrentUpdate-(5)]
	_ = x[AuthenticationFailed-(6)]
	_ = x[InternalError-(7)]
	_ = x[ForeignKeyViolation-(8)]
	_ = x[ConstraintViolation-(9)]
	_ = x[InvalidInputSyntax-(10)]
	_ = x[PermissionDenied-(11)]
	_ = x[OutOfDisk-(12)]
	_ = x[OutOfMemory-(13)]
	_ = x[TooManyConnections-(14)]
	_ = x[TxDone-(15)]
	_ = x[Timeout-(16)]
	_ = x[Unknown-(17)]
}

var _CodeValues = []Code{OK, NotFound, UniqueViolation, TableNotFound, DataException, ConcurrentUpdate, AuthenticationFailed, InternalError, ForeignKeyViolation, ConstraintViolation, InvalidInputSyntax, PermissionDenied, OutOfDisk, OutOfMemory, TooManyConnections, TxDone, Timeout, Unknown}

var _CodeNameToValueMap = map[string]Code{
	_CodeName[0:2]:          OK,
	_CodeLowerName[0:2]:     OK,
	_CodeName[2:10]:         NotFound,
	_CodeLowerName[2:10]:    NotFound,
	_CodeName[10:25]:        UniqueViolation,
	_CodeLowerName[10:25]:   UniqueViolation,
	_CodeName[25:38]:        TableNotFound,
	_CodeLowerName[25:38]:   TableNotFound,
	_CodeName[38:51]:        DataException,
	_CodeLowerName[38:51]:   DataException,
	_CodeName[51:67]:        ConcurrentUpdate,
	_CodeLowerName[51:67]:   ConcurrentUpdate,
	_CodeName[67:87]:        AuthenticationFailed,
	_CodeLowerName[67:87]:   AuthenticationFailed,
	_CodeName[87:100]:       InternalError,
	_CodeLowerName[87:100]:  InternalError,
	_CodeName[100:119]:      ForeignKeyViolation,
	_CodeLowerName[100:119]: ForeignKeyViolation,
	_CodeName[119:138]:      ConstraintViolation,
	_CodeLowerName[119:138]: ConstraintViolation,
	_CodeName[138:156]:      InvalidInputSyntax,
	_CodeLowerName[138:156]: InvalidInputSyntax,
	_CodeName[156:172]:      PermissionDenied,
	_CodeLowerName[156:172]: PermissionDenied,
	_CodeName[172:181]:      OutOfDisk,
	_CodeLowerName[172:181]: OutOfDisk,
	_CodeName[181:192]:      OutOfMemory,
	_CodeLowerName[181:192]: OutOfMemory,
	_CodeName[192:210]:      TooManyConnections,
	_CodeLowerName[192:210]: TooManyConnections,
	_CodeName[210:216]:      TxDone,
	_CodeLowerName[210:216]: TxDone,
	_CodeName[216:223]:      Timeout,
	_CodeLowerName[216:223]: Timeout,
	_CodeName[223:230]:      Unknown,
	_CodeLowerName[223:230]: Unknown,
}

var _CodeNames = []string{
	_CodeName[0:2],
	_CodeName[2:10],
	_CodeName[10:25],
	_CodeName[25:38],
	_CodeName[38:51],
	_CodeName[51:67],
	_CodeName[67:87],
	_CodeName[87:100],
	_CodeName[100:119],
	_CodeName[119:138],
	_CodeName[138:156],
	_CodeName[156:172],
	_CodeName[172:181],
	_CodeName[181:192],
	_CodeName[192:210],
	_CodeName[210:216],
	_CodeName[216:223],
	_CodeName[223:230],
}

// CodeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func CodeString(s string) (Code, error) {
	if val, ok := _CodeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _CodeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Code values", s)
}

// CodeValues returns all values of the enum
func CodeValues() []Code {
	return _CodeValues
}

// CodeStrings returns a slice of all String values of the enum
func CodeStrings() []string {
	strs := make([]string, len(_CodeNames))
	copy(strs, _CodeNames)
	return strs
}

// IsACode returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Code) IsACode() bool {
	for _, v := range _CodeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalText implements the encoding.TextMarshaler interface for Code
func (i Code) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Code
func (i *Code) UnmarshalText(text []byte) error {
	var err error
	*i, err = CodeString(string(text))
	return err
}
