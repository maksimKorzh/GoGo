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
  return len(board.generateCandidateMoves()) == 0
}

func (board *Board) playout() float64 {
  moveNumber := 0;
  stones := []float64{ 0, 0, 0 };
  for moveNumber < 1500 {
    move := board.generatePlayoutMove(board.side);
    board.play(move, board.side);
    moveNumber++;
  }
  for sq := 0; sq < board.size*board.size; sq++ {
    if board.position[sq] == OFFBOARD { continue; }
    if board.position[sq] == EMPTY {
      color := board.diamond(sq);
      if color != EMPTY {
        board.position[sq] = color;
        stones[color]++;
      } else {
        randomColor := rand.Intn(2)+1;
        board.position[sq] = randomColor;
        stones[randomColor]++;
      }
    } else {
      stones[board.position[sq]]++;
    }
  }
  result := stones[BLACK] - (stones[WHITE]+7.5);
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
  for row := 0; row < board.size; row++ {
    for col := 0; col < board.size; col++ {
      if row > offset && row < board.size - (offset+1) &&
         col > offset && col < board.size - (offset+1) {
        moves = append(moves, row * board.size + col);
      }
    }
  };return moves[rand.Intn(len(moves))];
}

func (board *Board) generatePlayoutMove(color int) int {
  if board.size == 21 { /* engine takes corners and sides */
    fuseki := []int { 88,340,352,100,69,129,363,297,311,371,143,77,213,227,73,367,215,225,115,325,220 };
    randomChoice := rand.Intn(len(fuseki));
    if board.position[fuseki[randomChoice]] == EMPTY &&
       fuseki[randomChoice] != board.ko &&
       board.notSuicide(fuseki[randomChoice]) {
       return fuseki[randomChoice];
    }
  }
  engine := board.target(color);
  player := board.target(3-color);
  if len(player) == 1 { /* engine captures player's group */
    if player[0] != board.ko { return player[0]; }
  }
  if len(engine) == 1 { /* engine saves own group */
    if board.notSuicide(engine[0]) { return engine[0]; }
  }
  if len(player) > 0 && len(player) <= len(engine) { /* engine surrounds player's group */
    randomChoice := rand.Intn(len(player));
    if board.notSuicide(player[randomChoice]) {
      return player[randomChoice];
    }
  } else if len(engine) > 0 && len(engine) <= len(player) { /* engine extends own group */
    randomChoice := rand.Intn(len(engine));
    if board.notSuicide(engine[randomChoice]) {
      return engine[randomChoice];
    }
  }
  for offset := 2; offset >= 0; offset-- { /* random tenuki */
    randomMove := board.random(offset);
    if board.diamond(randomMove) != 3-color &&
       board.diamond(randomMove) != color &&
       board.position[randomMove] == EMPTY {
       return randomMove;
    }
  };return 0; /* no legal moves found */
}

func (board *Board) generateCandidateMoves() []int {
  var moves []int;
  for _, move := range []int{ 88,340,352,100 } { moves = append(moves, move); }
  for _, move := range board.target(BLACK) { moves = append(moves, move); };
  for _, move := range board.target(WHITE) { moves = append(moves, move); };
  moves = append(moves, board.random(2));
  return moves;
}
/*********************************************\
  ===========================================
                     MCTS
  ===========================================
\*********************************************/

// Node represents a node in the Monte Carlo Tree
type Node struct {
	state    *Board // State of the board at this node
	parent   *Node  // Parent node
	children []*Node
	visits   int     // Number of visits to this node
	score    float64 // Cumulative score for this node's state
}

// UCT constant for balancing exploration and exploitation
const UCTConstant = 1.41

// UCTScore calculates the UCT (Upper Confidence Bound for Trees) score for a node
func UCTScore(totalVisits int, nodeVisits int, nodeScore float64) float64 {
	if nodeVisits == 0 {
		return math.Inf(1)
	}
	return nodeScore/float64(nodeVisits) + UCTConstant*math.Sqrt(math.Log(float64(totalVisits))/float64(nodeVisits))
}

// MCTS runs the Monte Carlo Tree Search algorithm
func MCTS(rootState *Board, iterations int) int {
	root := &Node{state: rootState, parent: nil, children: nil, visits: 0, score: 0.0}

	for i := 0; i < iterations; i++ {
		node := root
		// Selection phase: navigate down the tree based on UCT score until a leaf node is reached
		for len(node.children) > 0 {
			bestScore := -math.Inf(1)
			var bestChild *Node
			for _, child := range node.children {
				score := UCTScore(node.visits, child.visits, child.score)
				if score > bestScore {
					bestScore = score
					bestChild = child
				}
			}
			node = bestChild
		}

		// Expansion phase: expand the selected node if it's not terminal
		if !node.state.isGameOver() {
			moves := node.state.generateCandidateMoves()
			for _, move := range moves {
				childState := node.state.copy() // Make a copy of the board
				if !childState.play(move, childState.side) { continue; }      // Apply the move
				child := &Node{state: childState, parent: node, children: nil, visits: 0, score: 0.0}
				node.children = append(node.children, child)
			}
			node = node.children[rand.Intn(len(node.children))] // Choose a random child to simulate
		}

		// Simulation phase: simulate a random playout from the selected node
		simState := node.state.copy()
		simResult := simulateRandomPlayout(simState)

		// Backpropagation phase: propagate the result back up the tree
		for node != nil {
			node.visits++
			node.score += simResult
			node = node.parent
		}
	}

  /*fmt.Println("Root Node:")
  fmt.Printf("Visits: %d  ", root.visits)
  if root.visits > 0 {
      winRate := root.score / float64(root.visits)
      fmt.Printf("Win Rate: %.2f\n", winRate)
  }

  if len(root.children) > 0 {
      fmt.Println("Children Nodes:")
      for _, child := range root.children {
          fmt.Printf("Move: %v, Visits: %d ", child.state.lastMove, child.visits)
          if child.visits > 0 {
              winRate := child.score / float64(child.visits)
              fmt.Printf("Win Rate: %.2f\n", winRate)
          }
      }
  }*/

	// Select the best move based on the most visited child
	bestMove := getBestMove(root)
	//fmt.Printf("Best Move: %v\n", root.state.square(bestMove))
	return bestMove
}

// simulateRandomPlayout simulates a random playout from a given board state
func simulateRandomPlayout(state *Board) float64 {
	currentState := state.copy()
	return currentState.playout() // Returns 1 if black wins, -1 if white wins
}

// getBestMove selects the move with the highest visit count
func getBestMove(root *Node) int {
	bestMove := -1
	bestVisits := -1
	for _, child := range root.children {
		if child.visits > bestVisits {
			bestVisits = child.visits
			bestMove = child.state.lastMove
		}
	}
	return bestMove
}

// Assume Board, generateCandidateMoves(), makeMove(), copy(), isGameOver(), playout(), and lastMove() methods are implemented



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
        move := 0;
        for i := 0; i < 100; i++ {
          board.side = color;
          candidate := MCTS(board, 10);
          if candidate > 0 {
            move = candidate;
            break;
          }
        };if move > 0 {
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
  board.init(19);
  board.gtp();
}
