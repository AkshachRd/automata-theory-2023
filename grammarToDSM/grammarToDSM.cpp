#include "ReadWriteTable.h"
#include <map>
#include <algorithm>

using namespace std;

struct Eclose {
    string state;
    int column;
    bool fin;
    vector<string> eStates;
};

struct State {
    string name;
    bool fin;
    vector<string> arrOfStates;
    vector<vector<string>> transitions;
    vector<string> transitionsName;
};

int main()
{

    std::string type;
    std::string file = "5_1.txt";
    std::ifstream fileIn;
    fileIn.open(file);
    file = "output.csv";
    std::ofstream fileOut;
    fileOut.open(file);

    getline(fileIn, type);
    int eStr = 0;
    vector<vector<string>> inputAutomaton;
    map<string, Eclose> ecloses;
    if (type == "left")
    {
        ReadLeftTable(fileIn, inputAutomaton);
    }
    if (type == "right")
    {
        ReadRightTable(fileIn, inputAutomaton);
    }
    WriteTable(fileOut, inputAutomaton);

    fileIn.close();
    fileOut.close();

    return 0;
}