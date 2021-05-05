package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexem"
	"strconv"
	"strings"
)

func runner(leftParts, rightParts [][]string) (bool){
	word, _ := lexem.FetchLexem()
  
	symbol := leftParts[0][0]
	var ptrToRight, ptrToLeft int = 0, 0
	var stack []int
  
	for word != "eof" {
	  for i := ptrToLeft; i < len(leftParts); i++ {
		if symbol == leftParts[i][0] {
		  if strings.Contains(leftParts[i][1], word) {
			ptrToRight, _ = strconv.Atoi(leftParts[i][2])
			break
		  } else if leftParts[i][3] == strconv.Itoa(1) {
			return false
		  }
		}
	  }
  
	// fmt.Println("Left", symbol, input[wordPos], ptrToLeft, ptrToRight, stack)
  
	  for i := ptrToRight; i < len(rightParts); i++ {
		if strings.Contains(rightParts[i][1], word) {
		  
		  shift, _ := strconv.Atoi(rightParts[i][3]) 
		  if shift == 1 {
			word, _ = lexem.FetchLexem()
		  }
		
		  if rightParts[i][0] == grammary.EmptySymbol{
			top := len(stack) - 1
			ptrToRight = stack[top]
			symbol = rightParts[stack[top]][0]
			stack = stack[:top]
			break
		  } 
		
		  if rightParts[i][2] == "nullptr" && !grammary.IsNonTerminalSymbol(rightParts[i][0]){
			// fmt.Println("symbol", rightParts[i][0], stack)
			if len(stack) > 0 {
			  top := len(stack) - 1
			  symbol = rightParts[stack[top]][0]
			  i = stack[top] - 1
			  stack = stack[:top]
			  continue
			} else {
			  return checkEnd(stack, word)
			}
		  } 
  
		  if grammary.IsNonTerminalSymbol(rightParts[i][0]){
			ptrToLeft, _ = strconv.Atoi(rightParts[i][2])
			symbol = rightParts[i][0]
		  
			stackPtr, _ := strconv.Atoi(rightParts[i][4]) 
			
			if stackPtr > 0 {
			  stack = append(stack, stackPtr)
			}
			break
		  } else {
			fmt.Println("symbol", rightParts[i][0])
		  }
  
		} else {
		  return false
		}
	  }
	}
  
	return checkEnd(stack, word)
  }
  
  func checkEnd(stack []int, lexem string)(answer bool){
	if len(stack) == 0 && lexem == "eof" { 
	  return true
	} else {
	  return false
	}
  }
  