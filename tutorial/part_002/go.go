/*********************************************\
  ===========================================

         A didactic Go playing program       

                      by

               Code Monkey King

  ===========================================
\*********************************************/

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
  group []int;
  liberties []int;
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
        case BLACK|MARKER: fmt.Print(" #");
        case WHITE: fmt.Print(" O");
        case WHITE|MARKER: fmt.Print(" #");
        case LIBERTY: fmt.Print(" +");
      }
    };fmt.Println();
  };fmt.Print("   ");
  fmt.Print(" A B C D E F G H J K L M N O P Q R S T"[:board.size*2-4]);
  fmt.Print("\n\n");
  fmt.Print("    Side: ");
  if board.side == BLACK {
    fmt.Print("BLACK\n");
  } else { fmt.Print("WHITE\n"); }
  fmt.Print("      Ko: ");
  if board.ko == EMPTY { fmt.Print("EMPTY");
  } else { fmt.Print(board.square(board.ko)); }
  fmt.Print("\n\n");
  fmt.Print("    group: ", board.group, "\n");
  fmt.Print("    liberties: ", board.liberties);
  fmt.Print("\n\n");
}

func (board *Board) count(sq, color int) {
  stone := board.position[sq];
  if stone == OFFBOARD { return; }
  if stone > 0 && (stone & color) > 0 && (stone & MARKER) == 0 {
    board.position[sq] |= MARKER;
    board.group = append(board.group, sq);
    board.count(sq+1, color);
    board.count(sq-1, color);
    board.count(sq+board.size, color);
    board.count(sq-board.size, color);
  } else if stone == EMPTY {
    board.position[sq] |= LIBERTY;
    board.liberties = append(board.liberties, sq);
  }
}

func (board *Board) restore() {
  board.group = []int{};
  board.liberties = []int{};
  for sq := 0; sq < board.size*board.size; sq++ {
    if board.position[sq] != OFFBOARD {
      board.position[sq] &= 3;
    }
  }
}

func (board *Board) square(sq int) string {
  row := sq / board.size-1;
  col := sq % board.size-1;
  coord := make([]byte, 4);
  if col > 8 { coord[0] = 'A' + byte(col) + 1;
  } else { coord[0] = 'A' + byte(col); }
  copy(coord[1:], strconv.Itoa(board.size-2-row));
  return string(coord);
}

func main() {
  board := new(Board);
  board.init(19);
  board.position[100] = WHITE;
  board.position[101] = WHITE;
  board.position[121] = WHITE;
  board.show();
  board.count(100, WHITE);
  board.show();
  board.restore();
  board.show();
}
