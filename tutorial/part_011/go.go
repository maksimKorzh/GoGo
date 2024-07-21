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
  "bufio"
  "os"
  "strings"
  "math"
  "math/rand"
  "time"
)

/*********************************************\
  ===========================================
                     BOARD
  ===========================================
\*********************************************/

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
  lastMove int;
}

func (board *Board) init(size int) {
  board.size = size+2;
  board.side = BLACK;
  board.ko = EMPTY;
  board.lastMove = EMPTY;
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
    board.count(sq+board.size, color);
    board.count(sq-1, color);
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

func (board *Board) play(sq, color int) bool {
  if board.position[sq] != EMPTY { return false;
  } else if sq == board.ko { return false; }
  oldKo := board.ko;
  board.ko = EMPTY;
  board.position[sq] = color;
  for s := 0; s < board.size*board.size; s++ {
    stone := board.position[s];
    if stone == OFFBOARD { continue; }
    if stone & (3-board.side) > 0 {
      board.count(s, 3-color);
      if len(board.liberties) == 0 {
        if len(board.group) == 1 && board.diamond(sq) == 3-board.side {
          board.ko = board.group[0];
        }
        for i := 0; i < len(board.group); i++ {
          board.position[board.group[i]] = EMPTY;
        }
      };board.restore();
    }
  }
  board.count(sq, color);
  liberties := len(board.liberties);
  board.restore();
  if liberties == 0 {
    board.position[sq] = EMPTY;
    board.ko = oldKo;
    return false;
  };board.side = 3-board.side;
  board.lastMove = sq;
  return true;
}

func (board *Board) square(sq int) string {
  row := sq / board.size-1;
  col := sq % board.size-1;
  coord := make([]byte, 4);
  if col >= 8 { coord[0] = 'A' + byte(col) + 1;
  } else { coord[0] = 'A' + byte(col); }
  copy(coord[1:], strconv.Itoa(board.size-2-row));
  return string(coord);
}

func (board *Board) diamond(sq int) int {
  diamondColor := -1;
  otherColor := -1;
  var neighbours = []int{1, -1, board.size, -board.size};
  for i := 0; i < 4; i++ {
    if board.position[sq+neighbours[i]] == OFFBOARD { continue; }
    if board.position[sq+neighbours[i]] == EMPTY { return EMPTY; }
    if diamondColor == -1 {
      diamondColor = board.position[sq+neighbours[i]];
      otherColor = 3-diamondColor;
    } else if board.position[sq+neighbours[i]] == otherColor {
      return 0;
    }
  };diamondColor &= 3;
  return diamondColor;
}

func (board *Board) notSuicide(sq int) bool {
  liberties := 0;
  neighbours := []int{1, -1, board.size, -board.size};
  for i := 0; i < 4; i++ {
    if board.position[sq+neighbours[i]] == EMPTY { liberties++; }
  };if board.position[sq] == EMPTY &&
       liberties > 0 && sq != board.ko {
       return true;
  } else { return false; }
}

func (board *Board) copy() *Board {
  boardCopy := new(Board);
  boardCopy.size = board.size;
  boardCopy.side = board.side;
  boardCopy.ko = board.ko;
  boardCopy.lastMove = board.lastMove;
  boardCopy.position = make([]int, len(board.position));
  copy(boardCopy.position, board.position);
  return boardCopy;
}

func (board *Board) isGameOver() bool {
  return len(board.generateCandidateMoves()) == 0;
}

func (board *Board) playout() float64 {
  stones := []float64{ 0, 0, 0 };
  for {
    move := board.generatePlayoutMove(board.side);
    if !board.play(move, board.side) { break; };
  }
  for sq := 0; sq < board.size*board.size; sq++ {
    if board.position[sq] == OFFBOARD { continue; }
    if board.position[sq] != EMPTY {
      stones[board.position[sq]]++;
    }
  }
  result := stones[BLACK] - (stones[WHITE]+0.5);
  if result > 0 { return 1; } else { return -1; }
}

/*********************************************\
  ===========================================
                  HEURISTICS
  ===========================================
\*********************************************/

func (board *Board) target(color int) []int {
  var liberties []int;
  smallest := 100;
  for sq := 0; sq < board.size*board.size; sq++ {
    stone := board.position[sq];
    if stone == OFFBOARD { continue; }
    if stone & color > 0 {
      board.count(sq, color);
      if len(board.liberties) < smallest {
        smallest = len(board.liberties);
        liberties = board.liberties;
      };board.restore();
    }
  };return liberties;
}

func (board *Board) random(offset int) int {
  var moves []int;
  for i := 0; i < 1000; i++ {
    for row := 0; row < board.size; row++ {
      for col := 0; col < board.size; col++ {
        if row > offset && row < board.size - (offset+1) &&
           col > offset && col < board.size - (offset+1) {
          moves = append(moves, row * board.size + col);
        }
      }
    };randomMove := moves[rand.Intn(len(moves))];
    if board.position[randomMove] == EMPTY &&
       randomMove != board.ko &&
       board.diamond(randomMove) != WHITE &&
       board.diamond(randomMove) != BLACK {
       return randomMove;
    }
  }
  return EMPTY;
}

func (board *Board) generatePlayoutMove(color int) int {
  engine := board.target(color);
  player := board.target(3-color);
  if len(player) == 1 { /* engine captures player's group */
    if player[0] != board.ko { return player[0]; }
  }
  if len(engine) == 1 { /* engine saves own group */
    if board.notSuicide(engine[0]) {
      row := engine[0] / board.size;
      col := engine[0] % board.size;
      if row > 1 && row < (board.size-2) && col > 1 && col < (board.size-2) {
        return engine[0];
      }
    }
  }
  if len(player) > 0 { /* engine surrounds player's group */
    randomChoice := rand.Intn(len(player));
    if board.notSuicide(player[randomChoice]) {
      row := player[randomChoice] / board.size;
      col := player[randomChoice] % board.size;
      if row > 1 && row < (board.size-2) && col > 1 && col < (board.size-2) {
        return player[randomChoice];
      }
    }
  }
  if len(engine) > 0 { /* engine extends own group */
    randomChoice := rand.Intn(len(engine));
    if board.notSuicide(engine[randomChoice]) {
      row := engine[randomChoice] / board.size;
      col := engine[randomChoice] % board.size;
      if row > 1 && row < (board.size-2) && col > 1 && col < (board.size-2) {
        return engine[randomChoice];
      }
    }
  }
  for offset := 2; offset >= 0; offset-- { /* random tenuki */
    randomMove := board.random(offset);
    if randomMove != EMPTY { return randomMove; }
  };return 0; /* no legal moves found */
}

func (board *Board) generateCandidateMoves() []int {
  color := board.side;
  var moves []int;
  for i := 0; i < 10; i++ {
    candidate := board.generatePlayoutMove(color);
    duplicate := false;
    for _, move := range moves {
      if candidate == move {
        duplicate = true;
        break;
      }
    };if !duplicate {
      moves = append(moves, candidate);
    }
  };return moves;
}

/*********************************************\
  ===========================================
                     MCTS
  ===========================================
\*********************************************/

const UCTConstant = 1.41;

type Node struct { /* MCTS tree node */
  state    *Board;
  parent   *Node;
  children []*Node;
  visits   int;
  score    float64;
}

func MCTS(rootState *Board, iterations int) int {
  root := &Node { state: rootState, parent: nil, children: nil, visits: 0, score: 0.0 };
  for i := 0; i < iterations; i++ { /* selection phase */
    node := root;
    for len(node.children) > 0 {
      bestScore := -math.Inf(1);
      var bestChild *Node;
      for _, child := range node.children {
        score := 0.0;
        if child.visits == 0 {
          score = math.Inf(1);
        } else {
          score = child.score/float64(child.visits) + UCTConstant*math.Sqrt(math.Log(float64(node.visits))/float64(child.visits));
        };if score > bestScore {
          bestScore = score;
          bestChild = child;
        }
      };node = bestChild;
    }
    if !node.state.isGameOver() { /* expansion phase */
      moves := node.state.generateCandidateMoves();
      for _, move := range moves {
        childState := node.state.copy();
        if !childState.play(move, childState.side) { continue; }
        child := &Node{ state: childState, parent: node, children: nil, visits: 0, score: 0.0 };
        node.children = append(node.children, child);
      }
      if len(node.children) == 0 { continue; }
      node = node.children[rand.Intn(len(node.children))];
    }
    simState := node.state.copy(); simResult := simState.playout(); /* simulation phase */
    for node != nil { /* backpropagation phase */
      node.visits++;
      node.score += simResult;
      node = node.parent;
    }
  }

  if len(root.children) == 0 { return 0; } /* pass move */
  bestMove := EMPTY;
  bestVisits := -1;
  for _, child := range root.children {
    if child.visits > bestVisits {
      bestVisits = child.visits;
      bestMove = child.state.lastMove;
    }
  }

  fmt.Print("Best move:", root.state.square(bestMove));
  fmt.Printf(" Visits: %d", root.visits);
  if root.visits > 0 {
    winRate := root.score / float64(root.visits);
    fmt.Printf(" Win rate: %.2f\n", winRate);
  }

  if len(root.children) > 0 {
    fmt.Println("Child nodes:");
    for _, child := range root.children {
      fmt.Printf("Move: %v, Visits: %d ", child.state.square(child.state.lastMove), child.visits);
      if child.visits > 0 {
        winRate := child.score / float64(child.visits);
        fmt.Printf("Win rate: %.2f\n", winRate);
      } else { fmt.Println(); }
    }
  };return bestMove;
}

/*********************************************\
  ===========================================
                      GTP
  ===========================================
\*********************************************/

func (board *Board) gtp() {
  reader := bufio.NewReader(os.Stdin);
  writer := bufio.NewWriter(os.Stdout);
  defer writer.Flush();
  var boardSize int;
  for {
    userInput, err := reader.ReadString('\n');
    if err != nil { continue; }
    userInput = strings.TrimSpace(userInput);
    if userInput == "" { continue; }
    switch {
      case strings.HasPrefix(userInput, "quit"): return;
      case strings.HasPrefix(userInput, "name"):
        fmt.Fprintln(writer, "= GoGo\n");
      case strings.HasPrefix(userInput, "version"):
        fmt.Fprintln(writer, "= by Code Monkey King\n");
      case strings.HasPrefix(userInput, "protocol_version"):
        fmt.Fprintln(writer, "= 1\n");
      case strings.HasPrefix(userInput, "showboard"):
        fmt.Print("= current position");
        board.show();
      case strings.HasPrefix(userInput, "boardsize"):
        boardSize, _ = strconv.Atoi(userInput[10:]);
        fmt.Fprintln(writer, "=\n");
      case strings.HasPrefix(userInput, "clear_board"):
        board.init(boardSize);
        fmt.Fprintln(writer, "=\n");
      case strings.HasPrefix(userInput, "play"):
        if userInput[7:] == "pass" {
          board.side = 3-board.side;
          board.ko = EMPTY;
          fmt.Fprintln(writer, "=\n");
        } else {
          var color, col, row int;
          fmt.Sscanf(userInput, "play %c %c%d", &color, &col, &row);
          if color == 'B' { color = BLACK; }
          if color == 'W' { color = WHITE; }
          if col > 'I' { col--; }
          col = col - 'A' + 1;
          row = board.size - 1 - row;
          move := row * board.size + col;
          board.side = color;
          board.play(move, color);
          fmt.Fprintln(writer, "=\n");
        }
      case strings.HasPrefix(userInput, "genmove"):
        color := EMPTY;
        if userInput[8] == 'B' { color = BLACK; }
        if userInput[8] == 'W' { color = WHITE; }
        board.side = color;
        move := MCTS(board, 500);
        if move > 0 {
          if board.play(move, color) == false { fmt.Fprint(writer, "= pass\n\n"); }
          fmt.Fprint(writer, strings.ReplaceAll(("= " + board.square(move) + "\n\n"), "\x00", ""));
        } else { fmt.Fprint(writer, "= pass\n\n"); }
      default: fmt.Fprintln(writer, "=\n");
    };writer.Flush();
  }
}

/*********************************************\
  ===========================================
                     MAIN
  ===========================================
\*********************************************/

func main() {
  rand.Seed(time.Now().UnixNano());
  board := new(Board);
  board.init(5);
  board.gtp();
}
