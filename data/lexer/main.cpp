#include <iostream>
#include <fstream>
#include <string>
#include <sstream>
#include <vector>
#include <algorithm>
#include <map>
#include <set>

using namespace std;

const int STATES_COUNT = 100;
const int INPUT_SYMBOLS_COUNT = 85;
const int MAX_LENGTH_ID = 20;
const int MAX_LENGTH_INT = 11;

typedef int Matrix[STATES_COUNT][INPUT_SYMBOLS_COUNT];

string TrimBlanks(string const& str) {
  bool letterFlag = false, endFl = false;
  int start = -1, end;
  string trimedString = str;
  if (str != "")
  {
    for (int i = 0; i < str.length();++i)
    {
      if ((str[i] == ' ' || str[i] == '\t') && !letterFlag)
      {
        if (i == str.length() - 1)
        {
          trimedString = "";
          break;
        }
      }
      else if (!letterFlag)
      {
        letterFlag = true;
        start = i;
      }
      if ((str[str.length() - i - 1] == ' ' || str[str.length() - i - 1] == '\t') && !endFl)
        continue;
      else if (!endFl)
      {
        endFl = true;
        end = str.length() - i;
      }
    }
  }
  if (trimedString != "")
  {
    if (end != str.length() - 1)
      trimedString.erase(end, trimedString.length() - 1);
    if (start != -1)
      trimedString.erase(0, start);
  }

  return trimedString;
}

void FillSymbols(map<char, int>& symb)//ДОБАВИТЬ РАЗДЕЛИТЕЛИ
{
  symb.insert(pair('A', 0));
  symb.insert(pair('B', 1));
  symb.insert(pair('C', 2));
  symb.insert(pair('D', 3));
  symb.insert(pair('E', 4));
  symb.insert(pair('F', 5));
  symb.insert(pair('G', 6));
  symb.insert(pair('H', 7));
  symb.insert(pair('I', 8));
  symb.insert(pair('J', 9));
  symb.insert(pair('K', 10));
  symb.insert(pair('L', 11));
  symb.insert(pair('M', 12));
  symb.insert(pair('N', 13));
  symb.insert(pair('O', 14));
  symb.insert(pair('P', 15));
  symb.insert(pair('Q', 16));
  symb.insert(pair('R', 17));
  symb.insert(pair('S', 18));
  symb.insert(pair('T', 19));
  symb.insert(pair('U', 20));
  symb.insert(pair('V', 21));
  symb.insert(pair('W', 22));
  symb.insert(pair('X', 23));
  symb.insert(pair('Y', 24));
  symb.insert(pair('Z', 25));

  symb.insert(pair('a', 26));
  symb.insert(pair('b', 27));
  symb.insert(pair('c', 28));
  symb.insert(pair('d', 29));
  symb.insert(pair('e', 30));
  symb.insert(pair('f', 31));
  symb.insert(pair('g', 32));
  symb.insert(pair('h', 33));
  symb.insert(pair('i', 34));
  symb.insert(pair('j', 35));
  symb.insert(pair('k', 36));
  symb.insert(pair('l', 37));
  symb.insert(pair('m', 38));
  symb.insert(pair('n', 39));
  symb.insert(pair('o', 40));
  symb.insert(pair('p', 41));
  symb.insert(pair('q', 42));
  symb.insert(pair('r', 43));
  symb.insert(pair('s', 44));
  symb.insert(pair('t', 45));
  symb.insert(pair('u', 46));
  symb.insert(pair('v', 47));
  symb.insert(pair('w', 48));
  symb.insert(pair('x', 49));
  symb.insert(pair('y', 50));
  symb.insert(pair('z', 51));
  
  symb.insert(pair('0', 52));
  symb.insert(pair('1', 53));
  symb.insert(pair('2', 54));
  symb.insert(pair('3', 55));
  symb.insert(pair('4', 56));
  symb.insert(pair('5', 57));
  symb.insert(pair('6', 58));
  symb.insert(pair('7', 59));
  symb.insert(pair('8', 60));
  symb.insert(pair('9', 61));

  symb.insert(pair('_', 62));
  symb.insert(pair(' ', 63));
  symb.insert(pair(';', 64));
  symb.insert(pair('(', 65));
  symb.insert(pair(')', 66));
  symb.insert(pair('{', 67));
  symb.insert(pair('}', 68));
  symb.insert(pair('\n', 69));
  symb.insert(pair(',', 70));

  symb.insert(pair('/', 71));
  symb.insert(pair('*', 72));
  symb.insert(pair('+', 73));
  symb.insert(pair('-', 74));
  
  symb.insert(pair('=', 75));
  symb.insert(pair('<', 76));
  symb.insert(pair('>', 77));

  symb.insert(pair('.', 78));
  
  symb.insert(pair('[', 79));
  symb.insert(pair(']', 80));

  symb.insert(pair('"', 81));
}

void FillSpecial(set<char>& special)
{
  special.insert('{');
  special.insert('}');
  special.insert('(');
  special.insert(')');
  special.insert('[');
  special.insert(']');
  special.insert(';');
  // special.insert('.');
  // special.insert(',');
  special.insert('+');
  special.insert('-');
  special.insert('*');
  // special.insert('"');

}

void FillMatrix(Matrix& mat)
{
  vector<int> preFinalStates //Именно для ключевых слов
  {
    3, 8, 12, 15, 17, 21, 23, 27, 32, 36, 40, 83, 87, 91, 94
  };

  // ЗАПОЛНЕНИЕ МАТРИЦЫ ПЕРЕХОДОВ ДЛЯ АВТОМАТА
  for(int i = 0; i < STATES_COUNT; ++i)
  {
    for(int j = 0; j < INPUT_SYMBOLS_COUNT; ++j)
    {
      if(j != 63 && j != 64 && j != 69)
        mat[i][j] = 9;
      else
        mat[i][j] = 5;

      if(j == 73) // +
        mat[i][j] = 5;
      else if(j == 74) // -
        mat[i][j] = 5;
      else if(j == 72) // *
        mat[i][j] = 5;
      else if(j == 71) // /
        mat[i][j] = 5;
      else if(j == 75) // =
        mat[i][j] = 5;
      else if(j == 76) // <
        mat[i][j] = 5;
      else if(j == 77) // >
        mat[i][j] = 5;
    }
  }

  mat[0][34] = 1;  // s -> i
  mat[0][29] = 10; // s -> d
  mat[0][31] = 13; // s -> f
  mat[0][30] = 18; // s -> e
  mat[0][38] = 24; // s -> m
  mat[0][48] = 28; // s -> w
  mat[0][43] = 37; // s -> r
  mat[0][45] = 80; // s -> t
  mat[0][27] = 84; // s -> b
  mat[0][44] = 92; // s -> s

  mat[0][73] = 46; // s -> + 
  mat[0][74] = 47; // s -> -
  mat[0][72] = 49; // s -> *
  mat[0][71] = 57; // s -> creatingComment
  mat[0][75] = 54; // s -> =
  mat[0][76] = 52; // s -> <
  mat[0][77] = 52; // s -> >

  mat[0][65] = 55; // s -> (
  mat[0][67] = 55; // s -> {
  mat[0][66] = 56; // s -> )
  mat[0][68] = 56; // s -> }
  mat[0][79] = 55; // s -> [
  mat[0][80] = 56; // s -> ]

  mat[0][81] = 76; // s -> stringStart
  
  mat[0][70] = 71; // s -> point
  mat[0][78] = 71; // s -> point

  mat[71][70] = 7; // point -> separator
  mat[71][78] = 7; // point -> separator
  mat[71][63] = 7; // point -> separator
  mat[71][64] = 7;
  mat[71][69] = 7;
  mat[71][73] = 7; 
  mat[71][74] = 7; 
  mat[71][72] = 7; 
  mat[71][71] = 7;
  mat[71][65] = 7; 
  mat[71][67] = 7; 
  mat[71][66] = 7; 
  mat[71][68] = 7;
  mat[71][79] = 7; 
  mat[71][80] = 7;
  mat[71][81] = 7;   

  for(int i = 0; i < INPUT_SYMBOLS_COUNT; ++i)
  {
    mat[57][i] = 48; // division 
    mat[58][i] = 59; // one line comment
    mat[60][i] = 66; // open many line comment
    mat[66][i] = 67;
    mat[67][i] = 67; // finding end of comment
    mat[68][i] = 67;
    mat[69][i] = 70; // many line comment
    mat[71][i] = 7;
    mat[76][i] = 77; // creating string
    mat[77][i] = 77; // creating string
    mat[78][i] = 79; // stringEnd -> string
  }

  mat[77][81] = 78; // creating string -> stringEnd

  mat[71][52] = 42; // point -> creatingFloat
  mat[71][53] = 42;
  mat[71][54] = 42;
  mat[71][55] = 42;
  mat[71][56] = 42;
  mat[71][57] = 42;
  mat[71][58] = 42;
  mat[71][59] = 42;
  mat[71][60] = 42;
  mat[71][61] = 42;
  
  mat[57][71] = 58; // creatingComment -> creatingOneLineComment2
  mat[57][72] = 60; // creatingComment -> creatingManyLineComment2

  mat[66][72] = 68; 
  mat[67][72] = 68;
  mat[68][71] = 69;

  mat[0][52] = 61; // s -> zero 
  mat[0][53] = 41; // creatingInteger
  mat[0][54] = 41;
  mat[0][55] = 41;
  mat[0][56] = 41;
  mat[0][57] = 41;
  mat[0][58] = 41;
  mat[0][59] = 41;
  mat[0][60] = 41;
  mat[0][61] = 41;

  for(int i = 0; i < INPUT_SYMBOLS_COUNT; ++i)
  {
    mat[41][i] = 45;
    mat[42][i] = 45;
    mat[72][i] = 45;
    mat[73][i] = 45;
    mat[74][i] = 45;
    mat[75][i] = 45;
  }

  mat[41][52] = 41; // creatingInteger
  mat[41][53] = 41;
  mat[41][54] = 41;
  mat[41][55] = 41;
  mat[41][56] = 41;
  mat[41][57] = 41;
  mat[41][58] = 41;
  mat[41][59] = 41;
  mat[41][60] = 41;
  mat[41][61] = 41;
  mat[41][70] = 42;
  mat[41][78] = 42;

  mat[41][30] = 72; // integer -> exponentStart

  mat[41][73] = 43; // integer
  mat[41][74] = 43; 
  mat[41][72] = 43; 
  mat[41][71] = 43; 

  mat[41][63] = 43; // integer
  mat[41][64] = 43;
  mat[41][69] = 43;

  mat[41][65] = 43; // integer
  mat[41][67] = 43; 
  mat[41][66] = 43; 
  mat[41][68] = 43; 

  mat[41][77] = 43; 
  mat[41][76] = 43; 
  mat[41][75] = 43;  
  mat[41][81] = 43; 

  mat[42][52] = 42; // creatingFloat
  mat[42][53] = 42;
  mat[42][54] = 42;
  mat[42][55] = 42;
  mat[42][56] = 42;
  mat[42][57] = 42;
  mat[42][58] = 42;
  mat[42][59] = 42;
  mat[42][60] = 42;
  mat[42][61] = 42;

  mat[42][30] = 72; // creatingFloat -> exponentStart
  // mat[41][30] = 72; // creatingInt -> exponentStart

  mat[72][73] = 73; // exponentStart -> exponentSign + 
  mat[72][74] = 73; // exponentStart -> exponentSign -

  mat[73][52] = 74; // exponentSign -> exponent1
  mat[73][53] = 74;
  mat[73][54] = 74;
  mat[73][55] = 74;
  mat[73][56] = 74;
  mat[73][57] = 74;
  mat[73][58] = 74;
  mat[73][59] = 74;
  mat[73][60] = 74;
  mat[73][61] = 74;

  mat[74][52] = 75; // exponent1 -> exponent2
  mat[74][53] = 75;
  mat[74][54] = 75;
  mat[74][55] = 75;
  mat[74][56] = 75;
  mat[74][57] = 75;
  mat[74][58] = 75;
  mat[74][59] = 75;
  mat[74][60] = 75;
  mat[74][61] = 75;

  mat[72][63] = 6; // exponentStart -> error
  mat[72][64] = 6;
  mat[72][69] = 6;
  mat[73][63] = 6; // exponentSign + -> error
  mat[73][64] = 6;
  mat[73][69] = 6;
  mat[74][63] = 6; // exponent1 -> error
  mat[74][64] = 6;
  mat[74][69] = 6;
  mat[75][63] = 6; // exponent2 -> error
  mat[75][64] = 6;
  mat[75][69] = 6;

  mat[75][73] = 44; // exponent2 -> float
  mat[75][74] = 44; 
  mat[75][72] = 44; 
  mat[75][71] = 44; 

  mat[75][63] = 44; // exponent2 -> float
  mat[75][64] = 44;
  mat[75][69] = 44;

  mat[75][65] = 44; // exponent2 -> float
  mat[75][67] = 44; 
  mat[75][66] = 44; 
  mat[75][68] = 44;
  mat[75][79] = 44; 
  mat[75][80] = 44; 
  mat[75][81] = 44;  

  mat[75][77] = 44; 
  mat[75][76] = 44; 
  mat[75][75] = 44;  

  for(int i = 0; i < INPUT_SYMBOLS_COUNT; ++i)
  {
    if(i != 63 && i != 64 && i != 69)
      mat[45][i] = 45;
    else
      mat[45][i] = 6;
  }

  mat[45][73] = 6; // error+
  mat[45][74] = 6; // error-
  mat[45][72] = 6; // error/
  mat[45][71] = 6; // error*
  mat[45][65] = 6; 
  mat[45][67] = 6; 
  mat[45][66] = 6; 
  mat[45][68] = 6; 
  mat[45][79] = 6; 
  mat[45][80] = 6; 
  mat[45][77] = 6; 
  mat[45][76] = 6; 
  mat[45][75] = 6;  
  mat[45][81] = 6; 

  mat[42][73] = 44; // float
  mat[42][74] = 44; 
  mat[42][72] = 44; 
  mat[42][71] = 44; 

  mat[42][63] = 44; // float
  mat[42][64] = 44;
  mat[42][69] = 44;

  mat[42][65] = 44; // float
  mat[42][67] = 44; 
  mat[42][66] = 44; 
  mat[42][68] = 44;
  mat[42][79] = 44; 
  mat[42][80] = 44;  
  mat[42][81] = 44; 

  mat[42][77] = 44; 
  mat[42][76] = 44; 
  mat[42][75] = 44;  

  mat[61][49] = 62; // zero -> creatingHexadecimal1
  mat[61][23] = 45;

  mat[61][73] = 43; // zero = integer
  mat[61][74] = 43; 
  mat[61][72] = 43; 
  mat[61][71] = 43; 

  mat[61][63] = 43; // zero = integer
  mat[61][64] = 43;
  mat[61][69] = 43;

  mat[61][65] = 43; // zero = integer
  mat[61][67] = 43; 
  mat[61][66] = 43; 
  mat[61][68] = 43;
  mat[61][79] = 44; 
  mat[61][80] = 44;  
  mat[61][81] = 44;   

  mat[61][77] = 43; 
  mat[61][76] = 43; 
  mat[61][75] = 43;  

  for(int i = 0; i < INPUT_SYMBOLS_COUNT; ++i)
  {
    if(i >= 52 && i <= 61) // zero -> creatingInteger
      mat[61][i] = 41;

    if(i == 70 || i == 78) // zero -> creatingFloat
      mat[61][i] = 42; 

    if((i >= 0 && i <= 5)) // creatingHexadecimal1 -> creatingHexadecimal2 
    {
      mat[62][i] = 63; // creatingHexadecimal1 -> creatingHexadecimal2 
      mat[63][i] = 64; // creatingHexadecimal2 -> creatingHexadecimal3 
    }
    else
    {
      mat[62][i] = 45; // creatingHexadecimal1 -> creatingError
      mat[63][i] = 45; // creatingHexadecimal2 -> creatingError
    }

    mat[64][i] = 45;
  }

  mat[64][73] = 65; // hexadecimal
  mat[64][74] = 65; 
  mat[64][72] = 65; 
  mat[64][71] = 65;

  mat[64][77] = 65; 
  mat[64][76] = 65; 
  mat[64][75] = 65;  

  mat[64][63] = 65; // hexadecimal
  mat[64][64] = 65;
  mat[64][69] = 65;

  mat[64][65] = 65; // hexadecimal
  mat[64][67] = 65; 
  mat[64][66] = 65; 
  mat[64][68] = 65;
  mat[64][79] = 65; 
  mat[64][80] = 65;  
  mat[64][81] = 65;  

  mat[1][39] = 2; // i -> in
  mat[1][31] = 8; // i -> if
  
  mat[2][45] = 3; // in -> int

  mat[10][27] = 11; // d -> db
  mat[11][37] = 12; // db -> dbl

  mat[13][37] = 14; // f -> fl
  mat[14][40] = 15; // fl -> flo

  mat[13][40] = 16; // f -> fo
  mat[16][43] = 17; // fo -> for

  mat[13][26] = 88; // f -> fa
  mat[88][37] = 89; // fa -> fal
  mat[89][44] = 90; // fal -> fals
  mat[90][30] = 91; // fals -> false

  mat[18][37] = 19; // e -> el
  mat[19][44] = 20;
  mat[20][30] = 21; // els -> else

  mat[19][34] = 22; // el -> eli
  mat[22][31] = 23; // eli -> elif 

  mat[24][26] = 25; // m -> ma
  mat[25][34] = 26;
  mat[26][39] = 27; // mai -> main  

  mat[28][33] = 29; // w -> wh
  mat[29][34] = 30;
  mat[30][37] = 31;
  mat[31][30] = 32; // whil -> while

  mat[28][43] = 33; // w -> wr
  mat[33][34] = 34;
  mat[34][45] = 35;
  mat[35][30] = 36; // writ -> write  

  mat[37][30] = 38; // r -> re
  mat[38][26] = 39;
  mat[39][29] = 40; // rea -> read  

  mat[80][43] = 81; // t -> tr
  mat[81][46] = 82; // tr -> tru
  mat[82][30] = 83; // tru -> true
  
  mat[84][40] = 85; // b -> bo
  mat[85][40] = 86; // bo -> boo
  mat[86][37] = 87; // boo -> bool
  
  mat[92][45] = 93; // s -> st
  mat[93][43] = 94; // st -> str

  for(int i = 0; i < preFinalStates.size(); ++i) // word -> keyword
  {
    mat[preFinalStates[i]][63] = 4; // keyword
    mat[preFinalStates[i]][64] = 4;
    mat[preFinalStates[i]][69] = 4;
    mat[preFinalStates[i]][73] = 4; // keyword+
    mat[preFinalStates[i]][74] = 4; // keyword-
    mat[preFinalStates[i]][72] = 4; // keyword/
    mat[preFinalStates[i]][71] = 4; // keyword*  
    mat[preFinalStates[i]][77] = 4;
    mat[preFinalStates[i]][76] = 4;
    mat[preFinalStates[i]][75] = 4;
    mat[preFinalStates[i]][65] = 4; 
    mat[preFinalStates[i]][67] = 4; 
    mat[preFinalStates[i]][66] = 4; 
    mat[preFinalStates[i]][68] = 4;
    mat[preFinalStates[i]][79] = 4; 
    mat[preFinalStates[i]][80] = 4; 
    mat[preFinalStates[i]][81] = 4;   
  }

  for(int i = 0; i < INPUT_SYMBOLS_COUNT; ++i)
  {
    mat[54][i] = 50;
  }

  mat[54][63] = 50; // appropriation
  mat[54][64] = 50;
  mat[54][69] = 50;
  mat[54][75] = 53; // comp1 -> comp2

  for(int i = 0; i < INPUT_SYMBOLS_COUNT; ++i)
  {
    if(i < 75 || i > 77)
    {
      mat[52][i] = 51;
      mat[53][i] = 51;
    }
    else
    { 
      mat[52][i] = 45;
      mat[53][i] = 45;
    }
  }

  mat[52][63] = 51; // comparison
  mat[52][64] = 51;
  mat[52][69] = 51;
  mat[52][65] = 51; 
  mat[52][67] = 51; 
  mat[52][66] = 51; 
  mat[52][68] = 51;
  mat[52][79] = 51; 
  mat[52][80] = 51;  
  mat[52][81] = 51;   

  mat[52][75] = 53; // comp1 -> comp2
  mat[52][76] = 53;
  mat[52][77] = 53;
  
  mat[9][63] = 5; // id
  mat[9][64] = 5;
  mat[9][69] = 5;
  mat[9][70] = 5;
  mat[9][78] = 5;

  mat[9][73] = 5; // id
  mat[9][74] = 5; 
  mat[9][72] = 5; 
  mat[9][71] = 5; 
  mat[9][75] = 5; 
  mat[9][76] = 5; 
  mat[9][77] = 5;
  mat[9][65] = 5; 
  mat[9][67] = 5; 
  mat[9][66] = 5; 
  mat[9][68] = 5;
  mat[9][79] = 5; 
  mat[9][80] = 5;    
  mat[9][81] = 5;    

  // mat[0][63] = 7;
  mat[0][64] = 7;
  // mat[0][70] = 7;
  // mat[0][78] = 7;
}

void FillStates(map<int, string>& states)
{
  states.insert(pair(0, "start"));

  states.insert(pair(1, "i"));
  states.insert(pair(2, "in"));
  states.insert(pair(3, "int"));
  
  states.insert(pair(8, "if"));
  states.insert(pair(9, "creatingId"));

  states.insert(pair(10, "d"));
  states.insert(pair(11, "db"));
  states.insert(pair(12, "dbl"));

  states.insert(pair(13, "f"));
  states.insert(pair(14, "fl"));
  states.insert(pair(15, "flo"));

  states.insert(pair(16, "fo"));
  states.insert(pair(17, "for"));

  states.insert(pair(18, "e"));
  states.insert(pair(19, "el"));
  states.insert(pair(20, "els"));
  states.insert(pair(21, "else"));
  
  states.insert(pair(22, "eli"));
  states.insert(pair(23, "elif"));
  
  states.insert(pair(24, "m"));
  states.insert(pair(25, "ma"));
  states.insert(pair(26, "mai"));
  states.insert(pair(27, "main"));
  
  states.insert(pair(28, "w"));
  states.insert(pair(29, "wh"));
  states.insert(pair(30, "whi"));
  states.insert(pair(31, "whil"));
  states.insert(pair(32, "while"));

  states.insert(pair(33, "wr"));
  states.insert(pair(34, "wri"));
  states.insert(pair(35, "writ"));
  states.insert(pair(36, "write"));
  
  states.insert(pair(37, "r"));
  states.insert(pair(38, "re"));
  states.insert(pair(39, "rea"));
  states.insert(pair(40, "read"));
  
  states.insert(pair(41, "creatingInt"));
  states.insert(pair(42, "creatingFloat"));
  
  states.insert(pair(43, "integer"));
  states.insert(pair(44, "float"));

  states.insert(pair(45, "creatingError"));

  states.insert(pair(46, "addition"));
  states.insert(pair(47, "subtraction"));
  states.insert(pair(48, "division"));
  states.insert(pair(49, "multiplication"));
  states.insert(pair(50, "appropriation"));
  states.insert(pair(51, "comparison"));

  states.insert(pair(52, "creatingComparison"));
  states.insert(pair(53, "creatingComparison2"));

  states.insert(pair(54, "creatingAppropriation"));

  states.insert(pair(55, "openParenthesis"));
  states.insert(pair(56, "closingParenthesis"));
  
  states.insert(pair(57, "creatingComment"));
  states.insert(pair(58, "creatingOneLineComment2"));
  states.insert(pair(59, "oneLineComment"));
  states.insert(pair(60, "creatingManyLineComment2"));

  states.insert(pair(61, "zero"));
  states.insert(pair(62, "creatingHexadecimal1"));
  states.insert(pair(63, "creatingHexadecimal2"));
  states.insert(pair(64, "creatingHexadecimal3"));
  states.insert(pair(65, "hexadecimal"));

  states.insert(pair(66, "openManyLineComment"));
  states.insert(pair(67, "findingEndOfManyLineComment"));
  states.insert(pair(68, "creatingEndManyLineComment"));
  states.insert(pair(69, "creatingEndManyLineComment2"));
  states.insert(pair(70, "manyLineComment"));

  states.insert(pair(71, "point"));

  states.insert(pair(72, "exponentStart"));
  states.insert(pair(73, "exponentSign"));
  states.insert(pair(74, "exponent1"));
  states.insert(pair(75, "exponent2"));
  
  states.insert(pair(76, "stringStart"));
  states.insert(pair(77, "creatingString"));
  states.insert(pair(78, "stringEnd"));
  states.insert(pair(79, "string"));

  states.insert(pair(80, "t"));
  states.insert(pair(81, "tr"));
  states.insert(pair(82, "tru"));
  states.insert(pair(83, "true"));

  states.insert(pair(84, "b"));
  states.insert(pair(85, "bo"));
  states.insert(pair(86, "boo"));
  states.insert(pair(87, "bool"));

  states.insert(pair(88, "fa"));
  states.insert(pair(89, "fal"));
  states.insert(pair(90, "fals"));
  states.insert(pair(91, "false"));
  
  states.insert(pair(92, "s"));
  states.insert(pair(93, "st"));
  states.insert(pair(94, "str"));

  states.insert(pair(4, "keyword"));
  states.insert(pair(5, "id"));
  states.insert(pair(6, "error"));
  states.insert(pair(7, "separator"));
}

int main() {
  map<char, int> symbols;
  set<char> special;
  
  map<char, int>::iterator symbolIt;

  FillSymbols(symbols);
  FillSpecial(special);

  ifstream input;
  ofstream output;
  input.open("input.txt");
  output.open("output.txt");

  int state = 0, inputSymb;//СТРОКА И СТОЛБЕЦ В МАТРИЦЕ ПЕРЕХОДОВ 

  map<int, string> states;
  FillStates(states);

  vector<int> finalStates
  {
    4, 5, 6, 7, 43, 44, 46, 47, 48, 49, 50, 51, 55, 56, 59, 65, 66, 67, 68, 69, 70, 76, 77, 79
  };

  Matrix matrix;
  FillMatrix(matrix);

  string currStr, stringToTrim;
  set<char>::iterator specialIterator;

  int startComment, startPosition;
  stringstream commentContent;

  int stringNum = 1;
  int lengthCounter;
  while(getline(input, stringToTrim))
  { 
    currStr = TrimBlanks(stringToTrim);
    string word;
    // if(!input.eof()) ВОЗМОЖНО НАДО БУДЕТ ВЕРНУТЬ ЭТУ ПРОВЕРКУ ОБРАТНО 
      currStr += "\n";

    for(int i = 0; i != currStr.length(); ++i)//ПОСИМВОЛЬНЫЙ РАЗБОР СТРОКИ С ОПРЕДЕЛЕНИЕМ ТОКЕНОВ
    {
      symbolIt = symbols.find(currStr[i]);
      if(symbolIt != symbols.end())
      {
        inputSymb = symbolIt->second;
        state = matrix[state][inputSymb];

        // map<int, string>::iterator stateI;
        // stateI = states.find(state);
        // cout << stateI->second << " ";

        if(find(finalStates.begin(), finalStates.end(), state) == finalStates.end())
        {
          word += currStr[i];

          if(state == 9)
          {
            ++lengthCounter;

            if(lengthCounter > MAX_LENGTH_ID)
              state = 45;
          }
          
          if(state == 41)
          {
            ++lengthCounter;

            if(lengthCounter > MAX_LENGTH_INT)
              state = 45;
          }
          
        }
        else
        {
          map<int, string>::iterator stateIt;
          stateIt = states.find(state);

          if(state == 59) // one Line Comment
          {
            state = 0;
            string comment = currStr;
            comment.erase(0, i);
            int commLength = comment.length();
            comment.erase(commLength-1, commLength);

            cout << stateIt->second << "(" << stringNum << ", " << i + 1 - word.length() << ") - " << comment << endl;            
            break;
          }

          if(state == 66 || state == 76)
          {
            startComment = stringNum;
            startPosition = i + 1 - word.length();
            commentContent << currStr[i];
            state = matrix[state][inputSymb];
            stateIt = states.find(state);
            word = "";
            continue;
          }

          if(state == 67 || state == 68 || state == 69 || state == 77)
          {
            commentContent << currStr[i];
            continue;
          }

          if(state == 70) // many Line Comment
          {
            int commLength = commentContent.str().length();
            cout << stateIt->second << "(" << startComment << ", " << startPosition << "; " << stringNum << ", " << i + 1 - word.length() << ") - " << commentContent.str().erase(commLength-2, commLength) << endl;  
            commentContent.str("");
            startComment = 0;
            startPosition = 0;
          }

          if(state == 79) // string
          {
            int commLength = commentContent.str().length();
            
            if(commentContent.str()[0] == '"')
              commentContent.str(commentContent.str().erase(0,1));

            cout << stateIt->second << "(" << startComment << ", " << startPosition << "; " << stringNum << ", " << i + 1 - word.length() << ") - " << commentContent.str() << endl;  

            word = "";
            commentContent.str("");
            startComment = 0;
            startPosition = 0;
            // continue;
          }

          if(word != "")
            cout << stateIt->second << "(" << stringNum << ", " << i + 1 - word.length() << ") - " << word << endl;

          state = 0;
          specialIterator = special.find(currStr[i]);
          
          word = "";

          if(specialIterator != special.end())
          {
            state = matrix[state][inputSymb];
            stateIt = states.find(state);
            cout << stateIt->second << "(" << stringNum << ", " << i + 1 << ") - " << currStr[i] << endl;
            state = 0;
          }
          else
          {
            if(inputSymb != 69 && inputSymb != 63)      
            {
              word += currStr[i];
              state = matrix[state][inputSymb];

              if(state == 66 || state == 76)
              {
                startComment = stringNum;
                startPosition = i + 1 - word.length();
              }
            }
          }
          lengthCounter = 0;
        }
      }
      else
      {
        if(state == 77 || state == 57 || state == 67)
        {
          commentContent << currStr[i];
        }
        else
        {
          state = 45;
          word += currStr[i];
        }
      }
    }

    output << currStr;
    ++stringNum;
  }

  if(commentContent.str() != "" && state == 77)
  {
    cout << "error(" << startComment << ", " << startPosition << "; end of file" << ")" << endl;
  }
  else if(commentContent.str() != "")
  {
    cout << "eof(" << startComment << ", " << startPosition << "; end of file" << ")" << endl;
  }

  if(input.eof() && commentContent.str() == "")
  {
    cout << "eof" << endl;  
  }

  return 0;
}