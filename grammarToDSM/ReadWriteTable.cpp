#include "ReadWriteTable.h"
using namespace std;

void ReadRightTable(ifstream& input, vector<vector<string>>& inputAutomaton)
{
    string str;
    map<string, int> inputSymbols;
    map<string, bool> processedSymbols;
    int countOfInputSymbols = 2;
    vector<string> emptyTransitions;

    vector<string> states;
    states.push_back("");
    inputAutomaton.push_back(states);
    inputAutomaton.push_back(states);

    while (getline(input, str))
    {
        str.erase(remove(str.begin(), str.end(), ' '), str.end());
        size_t indexState = str.find("->");
        std::string state = str.substr(0, indexState);

        emptyTransitions.push_back("-");
        inputAutomaton[0].push_back("");
        inputAutomaton[1].push_back(state);

        str.erase(0, indexState + 2);

        vector<string> arr;
        std::string item;
        std::istringstream strStream(str);
        while (getline(strStream, item, '|'))
        {
            arr.push_back(item);
        }

        for (auto i = 0; i < arr.size(); i++)
        {
            string inputSymbol(1, arr[i][0]);
            if (arr[i].size() == 1)
            {
                arr[i] = "F";
            }
            else
            {
                arr[i].erase(0, 1);
            }
            if (inputSymbols.find(inputSymbol) == inputSymbols.end())
            {
                inputSymbols.insert(make_pair(inputSymbol, countOfInputSymbols));
                processedSymbols.insert(make_pair(inputSymbol, false));
                countOfInputSymbols += 1;
                emptyTransitions[0] = inputSymbol;
                inputAutomaton.push_back(emptyTransitions);
            }
            int lineNumber = inputSymbols[inputSymbol];
            if (processedSymbols[inputSymbol] == false)
            {
                inputAutomaton[lineNumber].push_back(arr[i]);
                processedSymbols[inputSymbol] = true;
            }
            else
            {
                inputAutomaton[lineNumber][inputAutomaton[lineNumber].size() - 1] = inputAutomaton[lineNumber][inputAutomaton[lineNumber].size() - 1] + "/" + arr[i];
            }

        }

        for (auto& symbol : processedSymbols)
        {
            if (symbol.second == false)
            {
                int lineNumber = inputSymbols[symbol.first];
                inputAutomaton[lineNumber].push_back("-");
            }
            if (symbol.second == true)
            {
                symbol.second = false;
            }
        }
    }
    inputAutomaton[0].push_back("F");
    inputAutomaton[1].push_back("H");
    for (auto i = 2; i < inputAutomaton.size(); i++)
    {
        inputAutomaton[i].push_back("-");
    }
    emptyTransitions.push_back("-");
    emptyTransitions.push_back("-");
    emptyTransitions[0] = "e";
    inputAutomaton.push_back(emptyTransitions);
}

struct Position
{
    int line;
    int column;
    string state;
};

void ReadLeftTable(ifstream& input, vector<vector<string>>& inputAutomaton)
{
    string str;
    map<string, int> inputSymbols;
    map<string, int> inputStates;
    int countOfInputSymbols = 2;
    int countOfInputStates = 1;

    vector<string> states;
    states.push_back("");
    inputAutomaton.push_back(states);
    inputAutomaton.push_back(states);
    vector<Position> transitions;

    inputStates.insert(make_pair("H", countOfInputStates));
    inputAutomaton[0].push_back("");
    countOfInputStates += 1;
    inputAutomaton[1].push_back("H");

    while (getline(input, str))
    {
        str.erase(remove(str.begin(), str.end(), ' '), str.end());
        size_t indexState = str.find("->");
        std::string state = str.substr(0, indexState);


        str.erase(0, indexState + 2);

        if (inputStates.find(state) == inputStates.end())
        {
            inputStates.insert(make_pair(state, countOfInputStates));
            if (countOfInputStates == 2)
            {
                inputAutomaton[0].push_back("F");
            }
            else
            {
                inputAutomaton[0].push_back("");
            }
            countOfInputStates += 1;
            inputAutomaton[1].push_back(state);
        }

        vector<string> arr;
        std::string item;
        std::istringstream strStream(str);
        while (getline(strStream, item, '|'))
        {
            arr.push_back(item);
        }

        for (auto i = 0; i < arr.size(); i++)
        {
            string inputSymbol(1, arr[i][arr[i].size() - 1]);
            if (arr[i].size() == 1)
            {
                arr[i] = "H";
            }
            else
            {
                arr[i].erase(arr[i].size() - 1, 1);
            }

            if (inputStates.find(arr[i]) == inputStates.end())
            {
                inputStates.insert(make_pair(arr[i], countOfInputStates));
                inputAutomaton[0].push_back("");
                countOfInputStates += 1;
                inputAutomaton[1].push_back(arr[i]);
            }

            if (inputSymbols.find(inputSymbol) == inputSymbols.end())
            {
                inputSymbols.insert(make_pair(inputSymbol, countOfInputSymbols));
                countOfInputSymbols += 1;
                states[0] = inputSymbol;
                inputAutomaton.push_back(states);
            }
            int lineNumber = inputSymbols[inputSymbol];
            int columnNumber = inputStates[arr[i]];
            Position position;
            position.column = columnNumber;
            position.line = lineNumber;
            position.state = state;
            transitions.push_back(position);
        }
    }

    states[0] = "e";
    inputAutomaton.push_back(states);
    countOfInputSymbols += 1;

    for (int i = 2; i < countOfInputSymbols; i++)
    {
        for (int j = 1; j < countOfInputStates; j++)
        {
            inputAutomaton[i].push_back("-");
        }
    }

    for (auto i = 0; i < transitions.size(); i++)
    {
        if (inputAutomaton[transitions[i].line][transitions[i].column] != "-")
        {
            inputAutomaton[transitions[i].line][transitions[i].column] = inputAutomaton[transitions[i].line][transitions[i].column] + "/" + transitions[i].state;
        }
        else
        {
            inputAutomaton[transitions[i].line][transitions[i].column] = transitions[i].state;
        }

    }

}


void WriteTable(ofstream& output, vector<vector<string>>& inputAutomaton)
{
    for (auto i = 1; i < inputAutomaton.size()-1; i++)
    {

        for (auto j = 0; j < inputAutomaton[i].size(); j++)
        {
            output << inputAutomaton[i][j];
            if (j + 1 != inputAutomaton[i].size())
            {
                output << ";";
            }
        }
        output << "\n";

    }
}