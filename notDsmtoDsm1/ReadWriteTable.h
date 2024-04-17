#pragma once
#include <iostream>
#include <vector>
#include <string>
#include <fstream>
#include <sstream>

void WriteTable(std::ofstream& output, std::vector<std::vector<std::string>>& outputAutomaton);

void ReadTable(std::ifstream& input, std::vector<std::vector<std::string>>& inputAutomaton);