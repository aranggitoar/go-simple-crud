package simple_crud_test

import (
  "errors"
  "testing"
)

var (
	DuplicateRow    = errors.New("Duplicate row.")
	RowNotExist     = errors.New("Row doesn't exist.")
	RowUpdateFailed = errors.New("Row update failed.")
	RowDeleteFailed = errors.New("Row deletion failed.")
	TableNotExist   = errors.New("Table doesn't exist, created just now.")
)

func TestNewDriver(t *testing.T) {
}

func TestInitDB(t *testing.T) {
}

func TestCreateRow(t *testing.T) {
  t.Parallel()

  type args struct {
    tn string
    fn string
    vs string
  }

  tests := []struct {
    name string
    args args
    want error
  }{
    {
      name: "similar row exists",
      args: args{
        tn: "languages",
        fn: "name",
        vs: "indonesian",
      },
      want: DuplicateRow,
    },
  }

  for _, test := range tests {
    test := test
    t.Run(test.name, func(t *testing.T) {
      if got := CreateRow(test.args.tn, test.args.fn, test.args.vs); got != test.want {
        t.Errorf("CreateRow() = %s, want %s", got, test.want)
      }
    })
  }
}

func TestReadAllRow(t *testing.T) {
}

func TestReadRow(t *testing.T) {
}

func TestUpdateRow(t *testing.T) {
}

func TestDeleteRow(t *testing.T) {
}


/*

Go Simple CRUD is a simple database CRUD operation API with dynamic row
scanning for Go.
Copyright (C) 2022  Aranggi J. Toar

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; only version 2 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

*/
