#pragma once
#include <iostream>
#include <vector>
#include <string>
#include <fstream>
#include <sstream>
#include <map>

void WriteTable(std::ofstream& output, std::vector<std::vector<std::string>>& outputAutomaton);

void ReadRightTable(std::ifstream& input, std::vector<std::vector<std::string>>& inputAutomaton);

void ReadLeftTable(std::ifstream& input, std::vector<std::vector<std::string>>& inputAutomaton);