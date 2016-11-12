package atmi
/* 
** Plain text IPC buffer support
**
** @file typed_string.go
** 
** -----------------------------------------------------------------------------
** Enduro/X Middleware Platform for Distributed Transaction Processing
** Copyright (C) 2015, ATR Baltic, SIA. All Rights Reserved.
** This software is released under one of the following licenses:
** GPL or ATR Baltic's license for commercial use.
** -----------------------------------------------------------------------------
** GPL license:
** 
** This program is free software; you can redistribute it and/or modify it under
** the terms of the GNU General Public License as published by the Free Software
** Foundation; either version 2 of the License, or (at your option) any later
** version.
**
** This program is distributed in the hope that it will be useful, but WITHOUT ANY
** WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
** PARTICULAR PURPOSE. See the GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License along with
** this program; if not, write to the Free Software Foundation, Inc., 59 Temple
** Place, Suite 330, Boston, MA 02111-1307 USA
**
** -----------------------------------------------------------------------------
** A commercial use license is available from ATR Baltic, SIA
** contact@atrbaltic.com
** -----------------------------------------------------------------------------
*/

/*
#cgo pkg-config: atmisrvinteg

#include <xatmi.h>
#include <string.h>
#include <stdlib.h>
#include <ubf.h>

*/
import "C"
import "unsafe"

//UBF Buffer
type TypedString struct {
	Buf *ATMIBuf
}

//Return The ATMI buffer to caller
func (u *TypedString) GetBuf() *ATMIBuf {
	return u.Buf
}

//Allocate new string buffer
//@param s - source string
func (ac *ATMICtx) NewString(gs string) (*TypedString, ATMIError) {
	var buf TypedString

	c_val := C.CString(gs)
	defer C.free(unsafe.Pointer(c_val))

	size := int64(C.strlen(c_val) + 1) /* 1 for EOS. */

	if ptr, err := ac.TpAlloc("STRING", "", size); nil != err {
		return nil, err
	} else {
		buf.Buf = ptr
		C.strcpy(buf.Buf.C_ptr, c_val)
		buf.Buf.TpSetCtxt(ac)

		return &buf, nil
	}
}

//Get the String Handler from ATMI Buffer
func (ac *ATMICtx) CastToString(abuf *ATMIBuf) (TypedString, ATMIError) {
	var buf TypedString

	buf.Buf = abuf

	return buf, nil
}

//Get the string value out from buffer
//@return String value
func (s *TypedString) GetString() string {
	return C.GoString(s.Buf.C_ptr)
}

//Set the string to the buffer
//@param str 	String value
func (s *TypedString) SetString(gs string) ATMIError {
	c_val := C.CString(gs)
	defer C.free(unsafe.Pointer(c_val))

	new_size := int64(C.strlen(c_val) + 1) /* 1 for EOS. */

	if cur_size, err := s.Buf.Ctx.TpTypes(s.Buf, nil, nil); nil != err {
		return err
	} else {
		if cur_size >= new_size {
			C.strcpy(s.Buf.C_ptr, c_val)
		} else if err := s.Buf.TpRealloc(new_size); nil != err {
			return err
		} else {
			C.strcpy(s.Buf.C_ptr, c_val)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// Wrappers for memory management
///////////////////////////////////////////////////////////////////////////////////

func (u *TypedString) TpRealloc(size int64) ATMIError {
	return u.Buf.TpRealloc(size)
}
