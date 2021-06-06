#include <iostream>
#include <fstream>
#include <string>
#include <sstream>
#include <vector>
#include <algorithm>
#include <map>
#include <set>
#include <regex>

using namespace std;

const int STATES_COUNT = 125;
const int INPUT_SYMBOLS_COUNT = 85;
const int MAX_LENGTH_ID = 20;
const int MAX_LENGTH_INT = 11;

typedef int Matrix[STATES_COUNT][INPUT_SYMBOLS_COUNT];

std::string replaceInString(std::string haystack, std::string old, std::string newValue) {
  size_t index = 0;
  while (true) {
     index = haystack.find(old, index);
     if (index == std::string::npos) break;

     haystack.replace(index, old.length(), newValue);

     index += old.length() - 1;
  }

  return haystack;
}

string normalizeNewLines(const string & value) {
    auto newValue = value;
    return replaceInString(newValue, "\\n", "\n");
}

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
  symb.insert(pair<char, int>('A', 0));
  symb.insert(pair<char, int>('B', 1));
  symb.insert(pair<char, int>('C', 2));
  symb.insert(pair<char, int>('D', 3));
  symb.insert(pair<char, int>('E', 4));
  symb.insert(pair<char, int>('F', 5));
  symb.insert(pair<char, int>('G', 6));
  symb.insert(pair<char, int>('H', 7));
  symb.insert(pair<char, int>('I', 8));
  symb.insert(pair<char, int>('J', 9));
  symb.insert(pair<char, int>('K', 10));
  symb.insert(pair<char, int>('L', 11));
  symb.insert(pair<char, int>('M', 12));
  symb.insert(pair<char, int>('N', 13));
  symb.insert(pair<char, int>('O', 14));
  symb.insert(pair<char, int>('P', 15));
  symb.insert(pair<char, int>('Q', 16));
  symb.insert(pair<char, int>('R', 17));
  symb.insert(pair<char, int>('S', 18));
  symb.insert(pair<char, int>('T', 19));
  symb.insert(pair<char, int>('U', 20));
  symb.insert(pair<char, int>('V', 21));
  symb.insert(pair<char, int>('W', 22));
  symb.insert(pair<char, int>('X', 23));
  symb.insert(pair<char, int>('Y', 24));
  symb.insert(pair<char, int>('Z', 25));

  symb.insert(pair<char, int>('a', 26));
  symb.insert(pair<char, int>('b', 27));
  symb.insert(pair<char, int>('c', 28));
  symb.insert(pair<char, int>('d', 29));
  symb.insert(pair<char, int>('e', 30));
  symb.insert(pair<char, int>('f', 31));
  symb.insert(pair<char, int>('g', 32));
  symb.insert(pair<char, int>('h', 33));
  symb.insert(pair<char, int>('i', 34));
  symb.insert(pair<char, int>('j', 35));
  symb.insert(pair<char, int>('k', 36));
  symb.insert(pair<char, int>('l', 37));
  symb.insert(pair<char, int>('m', 38));
  symb.insert(pair<char, int>('n', 39));
  symb.insert(pair<char, int>('o', 40));
  symb.insert(pair<char, int>('p', 41));
  symb.insert(pair<char, int>('q', 42));
  symb.insert(pair<char, int>('r', 43));
  symb.insert(pair<char, int>('s', 44));
  symb.insert(pair<char, int>('t', 45));
  symb.insert(pair<char, int>('u', 46));
  symb.insert(pair<char, int>('v', 47));
  symb.insert(pair<char, int>('w', 48));
  symb.insert(pair<char, int>('x', 49));
  symb.insert(pair<char, int>('y', 50));
  symb.insert(pair<char, int>('z', 51));
  
  symb.insert(pair<char, int>('0', 52));
  symb.insert(pair<char, int>('1', 53));
  symb.insert(pair<char, int>('2', 54));
  symb.insert(pair<char, int>('3', 55));
  symb.insert(pair<char, int>('4', 56));
  symb.insert(pair<char, int>('5', 57));
  symb.insert(pair<char, int>('6', 58));
  symb.insert(pair<char, int>('7', 59));
  symb.insert(pair<char, int>('8', 60));
  symb.insert(pair<char, int>('9', 61));

  symb.insert(pair<char, int>('_', 62));
  symb.insert(pair<char, int>(' ', 63));
  symb.insert(pair<char, int>(';', 64));
  symb.insert(pair<char, int>('(', 65));
  symb.insert(pair<char, int>(')', 66));
  symb.insert(pair<char, int>('{', 67));
  symb.insert(pair<char, int>('}', 68));
  symb.insert(pair<char, int>('\n', 69));
  symb.insert(pair<char, int>(',', 70));

  symb.insert(pair<char, int>('/', 71));
  symb.insert(pair<char, int>('*', 72));
  symb.insert(pair<char, int>('+', 73));
  symb.insert(pair<char, int>('-', 74));
  
  symb.insert(pair<char, int>('=', 75));
  symb.insert(pair<char, int>('<', 76));
  symb.insert(pair<char, int>('>', 77));

  symb.insert(pair<char, int>('.', 78));
  
  symb.insert(pair<char, int>('[', 79));
  symb.insert(pair<char, int>(']', 80));

  symb.insert(pair<char, int>('"', 81));

  symb.insert(pair<char, int>('|', 82));
  
  symb.insert(pair<char, int>(':', 83));
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
  special.insert('|');
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
    3, 8, 12, 15, 17, 21, 23, 27, 32, 36, 40, 83, 87, 91, 94, 110, 111, 112, 113, 114, 115, 117
  };

  vector<int> FinalTokenStates //Вместо ключевых слов
  {
    95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109
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
      else if(j == 70) // ,
        mat[i][j] = 5;

      if(i == 116) // -> ,
        mat[i][j] = 117;
      else if (i == 119) // : -> ::
        mat[i][j] = 120;
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
  mat[0][67] = 110; // s -> {
  mat[0][66] = 56; // s -> )
  mat[0][68] = 111; // s -> }
  mat[0][79] = 112; // s -> [
  mat[0][80] = 113; // s -> ]

  mat[0][81] = 76; // s -> stringStart

  mat[0][82] = 114; // s -> |

  mat[0][83] = 118; // s -> creatingDblColon
  
  mat[0][70] = 116; // s -> preComma
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
  mat[71][82] = 7;  
  mat[71][83] = 7;

  mat[116][70] = 117; // preComma -> commaToken
  mat[116][78] = 117; // preComma -> commaToken
  mat[116][63] = 117; // preComma -> commaToken
  mat[116][64] = 117;
  mat[116][69] = 117;
  mat[116][73] = 117; 
  mat[116][74] = 117; 
  mat[116][72] = 117; 
  mat[116][71] = 117;
  mat[116][65] = 117; 
  mat[116][67] = 117; 
  mat[116][66] = 117; 
  mat[116][68] = 117;
  mat[116][79] = 117; 
  mat[116][80] = 117;
  mat[116][81] = 117; 
  mat[116][82] = 117;   
  mat[116][83] = 117;  

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

  mat[116][52] = 42; // preComma -> creatingFloat
  mat[116][53] = 42;
  mat[116][54] = 42;
  mat[116][55] = 42;
  mat[116][56] = 42;
  mat[116][57] = 42;
  mat[116][58] = 42;
  mat[116][59] = 42;
  mat[116][60] = 42;
  mat[116][61] = 42;
  
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
  mat[41][79] = 43;  
  mat[41][80] = 43; 
  mat[41][81] = 43; 
  mat[41][82] = 43;
  mat[41][83] = 43; 

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
  mat[75][82] = 44;  
  mat[75][83] = 44;   

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
  mat[45][82] = 6;  
  mat[45][83] = 6;  

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
  mat[42][82] = 44; 
  mat[42][83] = 44; 

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
  mat[61][82] = 44;    
  mat[61][83] = 44;   

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
  mat[64][82] = 65;    
  mat[64][83] = 65;  

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
    mat[preFinalStates[i]][63] = FinalTokenStates[i]; // keyword
    mat[preFinalStates[i]][64] = FinalTokenStates[i];
    mat[preFinalStates[i]][69] = FinalTokenStates[i];
    mat[preFinalStates[i]][70] = FinalTokenStates[i];
    mat[preFinalStates[i]][73] = FinalTokenStates[i]; // keyword+
    mat[preFinalStates[i]][74] = FinalTokenStates[i]; // keyword-
    mat[preFinalStates[i]][72] = FinalTokenStates[i]; // keyword/
    mat[preFinalStates[i]][71] = FinalTokenStates[i]; // keyword*  
    mat[preFinalStates[i]][77] = FinalTokenStates[i];
    mat[preFinalStates[i]][76] = FinalTokenStates[i];
    mat[preFinalStates[i]][75] = FinalTokenStates[i];
    mat[preFinalStates[i]][65] = FinalTokenStates[i]; 
    mat[preFinalStates[i]][67] = FinalTokenStates[i]; 
    mat[preFinalStates[i]][66] = FinalTokenStates[i]; 
    mat[preFinalStates[i]][68] = FinalTokenStates[i];
    mat[preFinalStates[i]][79] = FinalTokenStates[i]; 
    mat[preFinalStates[i]][80] = FinalTokenStates[i]; 
    mat[preFinalStates[i]][81] = FinalTokenStates[i];
    mat[preFinalStates[i]][82] = FinalTokenStates[i];   
    mat[preFinalStates[i]][83] = FinalTokenStates[i];        
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
  mat[52][82] = 51;  
  mat[52][83] = 51;   

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
  mat[9][82] = 5;    
  mat[9][83] = 5;

  mat[118][83] = 119; // creatingDblColon -> preDblColon    

  mat[119][70] = 120; // preComma -> commaToken
  mat[119][78] = 120; // preComma -> commaToken
  mat[119][63] = 120; // preComma -> commaToken
  mat[119][64] = 120;
  mat[119][69] = 120;
  mat[119][73] = 120; 
  mat[119][74] = 120; 
  mat[119][72] = 120; 
  mat[119][71] = 120;
  mat[119][65] = 120; 
  mat[119][67] = 120; 
  mat[119][66] = 120; 
  mat[119][68] = 120;
  mat[119][79] = 120; 
  mat[119][80] = 120;
  mat[119][81] = 120; 
  mat[119][82] = 120; 
  mat[119][83] = 45;  

  // mat[0][63] = 7;
  mat[0][64] = 115;
  // mat[0][70] = 7;
  // mat[0][78] = 7;
}

void FillStates(map<int, string>& states)
{
  states.insert(pair<int, string>(0, "start"));

  states.insert(pair<int, string>(1, "i"));
  states.insert(pair<int, string>(2, "in"));
  states.insert(pair<int, string>(3, "int"));
  
  states.insert(pair<int, string>(8, "if"));
  states.insert(pair<int, string>(9, "creatingId"));

  states.insert(pair<int, string>(10, "d"));
  states.insert(pair<int, string>(11, "db"));
  states.insert(pair<int, string>(12, "dbl"));

  states.insert(pair<int, string>(13, "f"));
  states.insert(pair<int, string>(14, "fl"));
  states.insert(pair<int, string>(15, "flo"));

  states.insert(pair<int, string>(16, "fo"));
  states.insert(pair<int, string>(17, "for"));

  states.insert(pair<int, string>(18, "e"));
  states.insert(pair<int, string>(19, "el"));
  states.insert(pair<int, string>(20, "els"));
  states.insert(pair<int, string>(21, "else"));
  
  states.insert(pair<int, string>(22, "eli"));
  states.insert(pair<int, string>(23, "elif"));
  
  states.insert(pair<int, string>(24, "m"));
  states.insert(pair<int, string>(25, "ma"));
  states.insert(pair<int, string>(26, "mai"));
  states.insert(pair<int, string>(27, "main"));
  
  states.insert(pair<int, string>(28, "w"));
  states.insert(pair<int, string>(29, "wh"));
  states.insert(pair<int, string>(30, "whi"));
  states.insert(pair<int, string>(31, "whil"));
  states.insert(pair<int, string>(32, "while"));

  states.insert(pair<int, string>(33, "wr"));
  states.insert(pair<int, string>(34, "wri"));
  states.insert(pair<int, string>(35, "writ"));
  states.insert(pair<int, string>(36, "write"));
  
  states.insert(pair<int, string>(37, "r"));
  states.insert(pair<int, string>(38, "re"));
  states.insert(pair<int, string>(39, "rea"));
  states.insert(pair<int, string>(40, "read"));
  
  states.insert(pair<int, string>(41, "creatingInt"));
  states.insert(pair<int, string>(42, "creatingFloat"));
  
  states.insert(pair<int, string>(43, "integer"));
  states.insert(pair<int, string>(44, "float"));

  states.insert(pair<int, string>(45, "creatingError"));

  states.insert(pair<int, string>(46, "addition"));
  states.insert(pair<int, string>(47, "subtraction"));
  states.insert(pair<int, string>(48, "division"));
  states.insert(pair<int, string>(49, "multiplication"));
  states.insert(pair<int, string>(50, "appropriation"));
  states.insert(pair<int, string>(51, "comparison"));

  states.insert(pair<int, string>(52, "creatingComparison"));
  states.insert(pair<int, string>(53, "creatingComparison2"));

  states.insert(pair<int, string>(54, "creatingAppropriation"));

  states.insert(pair<int, string>(55, "openParenthesis"));
  states.insert(pair<int, string>(56, "closingParenthesis"));
  
  states.insert(pair<int, string>(57, "creatingComment"));
  states.insert(pair<int, string>(58, "creatingOneLineComment2"));
  states.insert(pair<int, string>(59, "oneLineComment"));
  states.insert(pair<int, string>(60, "creatingManyLineComment2"));

  states.insert(pair<int, string>(61, "zero"));
  states.insert(pair<int, string>(62, "creatingHexadecimal1"));
  states.insert(pair<int, string>(63, "creatingHexadecimal2"));
  states.insert(pair<int, string>(64, "creatingHexadecimal3"));
  states.insert(pair<int, string>(65, "hexadecimal"));

  states.insert(pair<int, string>(66, "openManyLineComment"));
  states.insert(pair<int, string>(67, "findingEndOfManyLineComment"));
  states.insert(pair<int, string>(68, "creatingEndManyLineComment"));
  states.insert(pair<int, string>(69, "creatingEndManyLineComment2"));
  states.insert(pair<int, string>(70, "manyLineComment"));

  states.insert(pair<int, string>(71, "point"));

  states.insert(pair<int, string>(72, "exponentStart"));
  states.insert(pair<int, string>(73, "exponentSign"));
  states.insert(pair<int, string>(74, "exponent1"));
  states.insert(pair<int, string>(75, "exponent2"));
  
  states.insert(pair<int, string>(76, "stringStart"));
  states.insert(pair<int, string>(77, "creatingString"));
  states.insert(pair<int, string>(78, "stringEnd"));
  states.insert(pair<int, string>(79, "string"));

  states.insert(pair<int, string>(80, "t"));
  states.insert(pair<int, string>(81, "tr"));
  states.insert(pair<int, string>(82, "tru"));
  states.insert(pair<int, string>(83, "true"));

  states.insert(pair<int, string>(84, "b"));
  states.insert(pair<int, string>(85, "bo"));
  states.insert(pair<int, string>(86, "boo"));
  states.insert(pair<int, string>(87, "bool"));

  states.insert(pair<int, string>(88, "fa"));
  states.insert(pair<int, string>(89, "fal"));
  states.insert(pair<int, string>(90, "fals"));
  states.insert(pair<int, string>(91, "false"));
  
  states.insert(pair<int, string>(92, "s"));
  states.insert(pair<int, string>(93, "st"));
  states.insert(pair<int, string>(94, "str"));

  states.insert(pair(95, "intToken"));
  states.insert(pair(96, "ifToken"));
  states.insert(pair(97, "dblToken"));
  states.insert(pair(98, "floToken"));
  states.insert(pair(99, "forToken"));
  states.insert(pair(100, "elseToken"));
  states.insert(pair(101, "elifToken"));
  states.insert(pair(102, "mainToken"));
  states.insert(pair(103, "whileToken"));
  states.insert(pair(104, "writeToken"));
  states.insert(pair(105, "readToken"));
  states.insert(pair(106, "trueToken"));
  states.insert(pair(107, "boolToken"));
  states.insert(pair(108, "falseToken"));
  states.insert(pair(109, "strToken"));

  states.insert(pair<int, string>(110, "openCurlyParenthesis"));
  states.insert(pair<int, string>(111, "closingCurlyParenthesis"));

  states.insert(pair<int, string>(112, "openSquareParenthesis"));
  states.insert(pair<int, string>(113, "closingSquareParenthesis"));

  states.insert(pair<int, string>(114, "straightSlashToken"));

  states.insert(pair<int, string>(115, "semicolonToken"));

  states.insert(pair<int, string>(116, "preComma"));
  states.insert(pair<int, string>(117, "commaToken"));  

  states.insert(pair<int, string>(118, "creatingDblColon"));
  states.insert(pair<int, string>(119, "perDblColon"));
  states.insert(pair<int, string>(120, "dblColonToken"));

  states.insert(pair<int, string>(4, "keyword"));
  states.insert(pair<int, string>(5, "id"));
  states.insert(pair<int, string>(6, "error"));
  states.insert(pair<int, string>(7, "separator"));
}

string initModule() {
    cout << "PING" << std::endl;

    string content;

    getline(cin, content);

    return normalizeNewLines(content);
}

int main() {
  map<char, int> symbols;
  set<char> special;
  
  map<char, int>::iterator symbolIt;

  FillSymbols(symbols);
  FillSpecial(special);

  int state = 0, inputSymb;//СТРОКА И СТОЛБЕЦ В МАТРИЦЕ ПЕРЕХОДОВ 

  map<int, string> states;
  FillStates(states);

  vector<int> finalStates
  {
    4, 5, 6, 7, 43, 44, 46, 47, 48, 49, 50, 51, 55, 56, 59, 65, 66, 67, 68, 69, 70, 76, 77, 79,
    95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 
    114, 115, 117, 120
  };

  Matrix matrix;
  FillMatrix(matrix);

  string currStr, stringToTrim;
  set<char>::iterator specialIterator;

  int startComment, startPosition;
  stringstream commentContent;

  int stringNum = 1;
  int lengthCounter;

  // init module

  auto content = initModule();

  istringstream input;
  input.str(content);

  // used to receive signal from host process
  string flush;

  string separator = "|-|";

  // init module


  while(getline(input, stringToTrim))
  {
    getline(cin, flush);
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

            cout << stateIt->second << " " << stringNum << " " << i + 1 - word.length() << separator << comment << endl;
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
            cout << stateIt->second << " " << stringNum << " " << i + 1 - word.length() << separator << word << endl;

          state = 0;
          specialIterator = special.find(currStr[i]);
          
          word = "";

          if(specialIterator != special.end())
          {
            state = matrix[state][inputSymb];
            stateIt = states.find(state);
            cout << stateIt->second << " " << stringNum << " " << i + 1 << separator << currStr[i] << endl;
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