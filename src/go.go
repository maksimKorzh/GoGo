/**********************************************\
  ============================================

      GoGo - A didactic Go playing program

                        by

                 Code Monkey King

  ============================================
\**********************************************/
package main

import (
  "fmt"
  "strconv"
)

const (
  EMPTY = 0;
  BLACK = 1;
  WHITE = 2;
  MARKER = 4;
  OFFBOARD = 7;
  LIBERTY = 8;
)

type Board struct {
  size int;
  position []int;
  side int;
  ko int;
}

func (board *Board) init(size int) {
  board.size = size+2;
  board.side = BLACK;
  board.ko = EMPTY;
  board.position = make([]int, board.size*board.size);
  for row := 0; row < board.size; row++ {
    for col := 0; col < board.size; col++ {
      sq := row * board.size + col;
      if row == 0 || row == board.size-1 ||
         col == 0 || col == board.size-1 {
        board.position[sq] = OFFBOARD;
      } else { board.position[sq] = EMPTY; }
    }
  }
}

func (board *Board) show() {
  for row := 0; row < board.size; row++ {
    for col := 0; col < board.size; col++ {
      if row > 0 && row < board.size-1 && col == 0 {
        fmt.Print(" ");
        if (board.size-1-row) < 10 { fmt.Print(" "); }
        fmt.Print(board.size-1-row);
      };sq := row * board.size + col;
      switch board.position[sq] {
        case EMPTY: fmt.Print(" .");
        case BLACK: fmt.Print(" X");
        case WHITE: fmt.Print(" O");
      }
    };fmt.Print("\n");
  };fmt.Print("   ");
  fmt.Print(" A B C D E F G H J K L M N O P Q R S T"[:board.size*2-4]);
  fmt.Print("\n\n");
  fmt.Print("    Side: ");
  if board.side == BLACK {
    fmt.Print("BLACK\n");
  } else { fmt.Print("WHITE"); }
  fmt.Print("      Ko: ", board.getSquare(board.ko), "\n\n");
}

func (board *Board) getSquare(sq int) string {
  row := sq/board.size-1;
  col := sq%board.size-1;
  c := make([]byte, 4);
	if col >= 8 { c[0] = 'A' + byte(col) + 1;
	} else { c[0] = 'A' + byte(col); }
	copy(c[1:], strconv.Itoa(board.size-2-row));
	return string(c);
}

func main() {
  board := new(Board);
  board.init(19);
  board.ko = 100;
  board.show();
}
