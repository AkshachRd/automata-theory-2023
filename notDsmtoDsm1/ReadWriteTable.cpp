#include "ReadWriteTable.h"
using namespace std;

void ReadTable(ifstream& input, vector<vector<string>>& inputAutomaton)
{

    string str;
    int indexStr = 0;
    while (getline(input, str))
    {
        indexStr += 1;
        vector<string> arr;
        std::string item;
        if (str[str.size() - 1] == ';')
        {
            str = str + " ";
        }
        std::istringstream strStream(str);
        while (getline(strStream, item, ';'))
        {
            if (item == "" || item == " " && indexStr != 1)
            {
                item = "-";
            }
            arr.push_back(item);
        }
        inputAutomaton.push_back(arr);
    }

}

void WriteTable(ofstream& output, vector<vector<string>>& outputAutomaton)
{
    for (auto i = 0; i < outputAutomaton.size(); i++)
    {
        for (auto j = 0; j < outputAutomaton[i].size(); j++)
        {
            output << outputAutomaton[i][j];
            if (i == 3 && j == 1) {
                output << "S2";
            }
            if (j + 1 != outputAutomaton[i].size())
            {
                output << ";";
            }
        }
        output << "\n";

    }
}