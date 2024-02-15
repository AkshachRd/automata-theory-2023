#include <iostream>

using namespace std;

void E(char& ch);
void A(char& ch);
void B(char& ch);
void T(char& ch);
void F(char& ch);


void E(char& ch)
{
	T(ch);
	if (ch == '+')
	{
		A(ch);
	}
	else {
		cin >> ch;
	}
	if (ch == '+')
	{
		A(ch);
	}
}

void A(char& ch)
{
	cin >> ch;
	T(ch);


	cin >> ch;
	A(ch);
}

void B(char& ch)
{
	cin >> ch;
	F(ch);

	cin >> ch;
	B(ch);
}

void T(char& ch)
{
	F(ch);
	if (ch == '*')
	{
		B(ch);
	}
	else {
		cin >> ch;
	}
	if (ch == '*')
	{
		B(ch);
	}
}

void F(char& ch)
{
	if (ch == '(')
	{
		cin >> ch;

		E(ch);
		cin >> ch;
		if (ch != ')')
		{
			throw exception("Missing ')'");
		}
	}
	else if (ch == '-')
	{
		cin >> ch;
		F(ch);
	}
	else if (!(ch == 'a' || ch == 'b' || ch == '5' || ch == '3'))
	{
		throw exception("Invalid character");
	}
}

int main() {
	char ch;

	try
	{
		while (cin >> ch)
		{
			E(ch);
		}
	}
	catch (exception& e)
	{
		cout << e.what() << endl;
	}
	cout << "Success" << endl;
}