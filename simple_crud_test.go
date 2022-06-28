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
